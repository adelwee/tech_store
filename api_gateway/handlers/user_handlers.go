package handlers

import (
	pb "Assignment2_AdelKenesova/user_service/proto"
	"context"
	"google.golang.org/grpc/metadata"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var userClient pb.UserServiceClient

func InitUserClient(conn *grpc.ClientConn) {
	userClient = pb.NewUserServiceClient(conn)
}

func RegisterUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := userClient.RegisterUser(ctx, &pb.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       res.Id,
		"username": res.Username,
		"email":    res.Email,
	})
}

func LoginUser(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := userClient.AuthenticateUser(ctx, &pb.AuthRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		log.Println("gRPC error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !res.Success {
		log.Println("Auth failed:", res.Message)
		c.JSON(http.StatusUnauthorized, gin.H{"error": res.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id": res.UserId,
		"token":   res.Token,
	})
}

func GetUserProfile(c *gin.Context) {
	// Передаём заголовок Authorization в gRPC metadata
	md := metadata.New(map[string]string{
		"authorization": c.GetHeader("Authorization"),
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	res, err := userClient.GetUserProfile(ctx, &pb.UserID{}) // ID достаётся из токена
	if err != nil {
		log.Println("GetUserProfile gRPC error:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
