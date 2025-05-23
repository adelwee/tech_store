package service

import (
	"Assignment2_AdelKenesova/user_service/internal/db"
	pb "Assignment2_AdelKenesova/user_service/proto"
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	db.InitDB()
	db.DB.Exec("DELETE FROM users") // Очистка перед тестом

	s := &UserService{}
	req := &pb.RegisterRequest{
		Username: "TestUser",
		Email:    "testuser@example.com",
		Password: "securepassword",
	}

	res, err := s.RegisterUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.Email, res.Email)
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	db.InitDB()
	s := &UserService{}

	// Сначала создаём пользователя
	req := &pb.RegisterRequest{
		Username: "DuplicateUser",
		Email:    "duplicate@example.com",
		Password: "password123",
	}
	_, _ = s.RegisterUser(context.Background(), req)

	// Пытаемся создать с тем же email снова
	res, err := s.RegisterUser(context.Background(), req)

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestAuthenticateUser(t *testing.T) {
	db.InitDB()
	db.DB.Exec("DELETE FROM users") // Очистка

	s := &UserService{}
	registerReq := &pb.RegisterRequest{
		Username: "AuthUser",
		Email:    "auth@example.com",
		Password: "authpass",
	}
	_, err := s.RegisterUser(context.Background(), registerReq)
	assert.NoError(t, err)

	authReq := &pb.AuthRequest{
		Email:    "auth@example.com",
		Password: "authpass",
	}
	res, err := s.AuthenticateUser(context.Background(), authReq)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.Token)
}

func TestGetUserProfile(t *testing.T) {
	db.InitDB()

	// Очистка пользователей + сброс ID-счётчика (PostgreSQL)
	db.DB.Exec("DELETE FROM users")
	db.DB.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1") // ✅ для PostgreSQL

	s := &UserService{}

	// Регистрируем пользователя
	registerReq := &pb.RegisterRequest{
		Username: "ProfileUser",
		Email:    "profile@example.com",
		Password: "profilepass",
	}
	registerRes, err := s.RegisterUser(context.Background(), registerReq)
	assert.NoError(t, err)

	// Аутентификация
	authReq := &pb.AuthRequest{
		Email:    "profile@example.com",
		Password: "profilepass",
	}
	authRes, err := s.AuthenticateUser(context.Background(), authReq)
	assert.NoError(t, err)
	fmt.Println("🧪 Token from auth:", authRes.Token)

	// Используем токен
	md := metadata.New(map[string]string{
		"authorization": authRes.Token,
	})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	res, err := s.GetUserProfile(ctx, &pb.UserID{})
	assert.NoError(t, err)
	assert.Equal(t, registerRes.Email, res.Email)
	assert.Equal(t, registerRes.Username, res.Username)
}
