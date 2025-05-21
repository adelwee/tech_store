package handlers

import (
	pb "Assignment2_AdelKenesova/user_service/proto"
	"context"
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

	// gRPC подключение к UserService
	conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Println("Failed to connect to UserService:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UserService unavailable"})
		return
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := client.RegisterUser(ctx, &pb.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       res.Id,
		"username": res.Username,
		"email":    res.Email,
	})
}

func GetUserProfile(c *gin.Context) {
	var req struct {
		ID uint64 `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := userClient.GetUserProfile(ctx, &pb.UserID{Id: req.ID})
	if err != nil {
		log.Println("GetUserProfile gRPC error:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
