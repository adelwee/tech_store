package service

import (
	"Assignment2_AdelKenesova/pkg/redis"
	"Assignment2_AdelKenesova/user_service/internal/db"
	"Assignment2_AdelKenesova/user_service/internal/models"
	pb "Assignment2_AdelKenesova/user_service/proto"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

	// Успешно
	return &pb.AuthResponse{
		Success: true,
		Message: "Authentication successful",
		UserId:  uint64(user.ID),
	}, nil
}

func (s *UserService) GetUserProfile(ctx context.Context, req *pb.UserID) (*pb.UserResponse, error) {
	cacheKey := fmt.Sprintf("user:%d", req.Id)
	client := redis.GetClient()

	// 1️ Попробуй получить из Redis
	cached, err := client.Get(ctx, cacheKey).Result()
	if err == nil {
		var cachedUser models.User
		if err := json.Unmarshal([]byte(cached), &cachedUser); err == nil {
			log.Println(" Returned user from Redis cache")
			return &pb.UserResponse{
				Id:       uint64(cachedUser.ID),
				Username: cachedUser.Username,
				Email:    cachedUser.Email,
			}, nil
		}
	}

	// 2Если нет в кэше — получи из базы
	var user models.User
	if err := db.DB.First(&user, req.Id).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	userJSON, err := json.Marshal(user)
	if err == nil {
		err = client.Set(ctx, cacheKey, userJSON, 0).Err()
		if err == nil {
			log.Println("User cached in Redis")
		} else {
			log.Println(" Failed to cache user in Redis:", err)
		}
	}

	return &pb.UserResponse{
		Id:       uint64(user.ID),
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
