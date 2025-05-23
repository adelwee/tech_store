package service

import (
	"Assignment2_AdelKenesova/pkg/email"
	"Assignment2_AdelKenesova/pkg/redis"
	"Assignment2_AdelKenesova/user_service/internal/db"
	"Assignment2_AdelKenesova/user_service/internal/models"
	pb "Assignment2_AdelKenesova/user_service/proto"
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/metadata"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var jwtSecret = []byte("adel_super_secret_key_12345")

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func (s *UserService) RegisterUser(ctx context.Context, req *pb.RegisterRequest) (*pb.UserResponse, error) {
	// Проверка на дубликат email
	var existing models.User
	if err := db.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		return nil, status.Errorf(codes.AlreadyExists, "User with email %s already exists", req.Email)
	}

	// Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return nil, status.Errorf(codes.Internal, "Failed to hash password")
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create user")
	}

	//  Отправка welcome-письма
	go func() {
		subject := " Welcome to TechStore!"
		body := fmt.Sprintf("Hello, %s!\n\nThanks for registration to our website.\n\nYour respectfully,\nteam TechStore", req.Username)
		if err := email.SendEmail(user.Email, subject, body); err != nil {
			log.Printf(" Failed to send welcome email: %v", err)
		} else {
			log.Printf(" Welcome email sent to %s", user.Email)
		}
	}()

	return &pb.UserResponse{
		Id:       uint64(user.ID),
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *UserService) AuthenticateUser(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	var user models.User

	// Поиск пользователя по email
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return &pb.AuthResponse{
			Success: false,
			Message: "User not found",
		}, nil
	}

	// Проверка пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return &pb.AuthResponse{
			Success: false,
			Message: "Invalid password",
		}, nil
	}

	//  Генерация JWT токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return &pb.AuthResponse{
			Success: false,
			Message: "Failed to generate token",
		}, nil
	}

	//  Успешный ответ с токеном
	return &pb.AuthResponse{
		Success: true,
		Message: "Authentication successful",
		UserId:  uint64(user.ID),
		Token:   tokenString,
	}, nil
}

func (s *UserService) GetUserProfile(ctx context.Context, _ *pb.UserID) (*pb.UserResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "Missing metadata")
	}
	tokens := md["authorization"]
	if len(tokens) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "Authorization token not provided")
	}
	tokenStr := strings.TrimPrefix(tokens[0], "Bearer ")

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Failed to parse claims")
	}
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, status.Errorf(codes.Internal, "user_id missing or invalid")
	}
	userID := uint64(userIDFloat)

	// Redis
	cacheKey := fmt.Sprintf("user:%d", userID)
	client := redis.GetClient()
	if cached, err := client.Get(ctx, cacheKey).Result(); err == nil {
		var cachedUser models.User
		if err := json.Unmarshal([]byte(cached), &cachedUser); err == nil && cachedUser.ID != 0 {
			log.Println("✅ Returned user from Redis cache")
			return &pb.UserResponse{
				Id:       uint64(cachedUser.ID),
				Username: cachedUser.Username,
				Email:    cachedUser.Email,
			}, nil
		}
	}

	// DB
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		log.Println("❌ DB error:", err)
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	// Кэширование
	if data, err := json.Marshal(user); err == nil {
		_ = client.Set(ctx, cacheKey, data, 0).Err()
	}

	return &pb.UserResponse{
		Id:       uint64(user.ID),
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
