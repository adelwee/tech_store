package handlers

import (
	pb "Assignment2_AdelKenesova/inventory_service/proto"
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

var inventoryClient pb.InventoryServiceClient

func InitInventoryClient(conn *grpc.ClientConn) {
	inventoryClient = pb.NewInventoryServiceClient(conn)
}

func CreateProduct(c *gin.Context) {
	var req pb.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := inventoryClient.CreateProduct(ctx, &req)
	if err != nil {
		log.Println("CreateProduct gRPC error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res.Product)
}

func GetProduct(c *gin.Context) {
	idParam := c.Param("id")

	var id uint64
	_, err := fmt.Sscanf(idParam, "%d", &id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := inventoryClient.GetProduct(ctx, &pb.GetProductRequest{Id: id})
	if err != nil {
		log.Println("GetProduct gRPC error:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res.Product)
}

func ListProducts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	res, err := inventoryClient.ListProducts(ctx, &pb.Empty{})
	if err != nil {
		log.Println("ListProducts gRPC error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var reqBody struct {
		Name        string  `json:"name"`
		Brand       string  `json:"brand"`
		CategoryID  uint64  `json:"category_id"`
		Price       float64 `json:"price"`
		Stock       uint64  `json:"stock"`
		Description string  `json:"description"`
	}
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := inventoryClient.UpdateProduct(ctx, &pb.UpdateProductRequest{
		Id:          id,
		Name:        reqBody.Name,
		Brand:       reqBody.Brand,
		CategoryId:  reqBody.CategoryID,
		Price:       reqBody.Price,
		Stock:       reqBody.Stock,
		Description: reqBody.Description,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res.Product)
}

func DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = inventoryClient.DeleteProduct(ctx, &pb.DeleteProductRequest{Id: id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
