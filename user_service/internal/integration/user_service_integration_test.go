package integration

import (
	"Assignment2_AdelKenesova/user_service/internal/db"
	"Assignment2_AdelKenesova/user_service/internal/service"
	pb "Assignment2_AdelKenesova/user_service/proto"
	"context"
	"google.golang.org/grpc/metadata"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func startGRPCServer(t *testing.T) string {
	lis, err := net.Listen("tcp", ":0") // свободный порт
	assert.NoError(t, err)

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, &service.UserService{})

	go func() {
		if err := server.Serve(lis); err != nil {
			t.Fatalf("Failed to serve: %v", err)
		}
	}()
	time.Sleep(100 * time.Millisecond) // Подождать запуск сервера
	return lis.Addr().String()
}

func getClientConn(t *testing.T, addr string) pb.UserServiceClient {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	assert.NoError(t, err)
	return pb.NewUserServiceClient(conn)
}

func TestRegisterAndAuthIntegration(t *testing.T) {
	db.InitDB()

	// Очистить пользователя перед тестом
	db.DB.Exec("DELETE FROM users WHERE email = ?", "int@test.com")

	addr := startGRPCServer(t)
	client := getClientConn(t, addr)

	// Register
	res, err := client.RegisterUser(context.Background(), &pb.RegisterRequest{
		Username: "integration_user",
		Email:    "int@test.com",
		Password: "pass123",
	})
	assert.NoError(t, err)
	assert.Equal(t, "integration_user", res.Username)

	// Authenticate
	auth, err := client.AuthenticateUser(context.Background(), &pb.AuthRequest{
		Email:    "int@test.com",
		Password: "pass123",
	})
	assert.NoError(t, err)
	assert.True(t, auth.Success)
	assert.NotEmpty(t, auth.Token)

	// Add token to metadata
	ctxWithToken := metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
		"authorization", "Bearer "+auth.Token,
	))

	// Get profile
	profile, err := client.GetUserProfile(ctxWithToken, &pb.UserID{})
	assert.NoError(t, err)
	assert.Equal(t, "int@test.com", profile.Email)
}
