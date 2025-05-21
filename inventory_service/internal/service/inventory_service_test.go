package service

import (
	"Assignment2_AdelKenesova/inventory_service/internal/db"
	pb "Assignment2_AdelKenesova/inventory_service/proto"
	"Assignment2_AdelKenesova/pkg/redis"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	ctx := context.Background()

	// Инициализация Redis
	redis.InitRedis()

	// Инициализация базы данных
	db.InitDB()
	db.Migrate()

	service := &InventoryService{}

	req := &pb.CreateProductRequest{
		Name:        "Test Redis Product",
		Brand:       "RedisBrand",
		CategoryId:  1,
		Price:       123.45,
		Stock:       10,
		Description: "For unit test",
	}

	res, err := service.CreateProduct(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, req.Name, res.Product.Name)
}

func TestGetProduct(t *testing.T) {
	ctx := context.Background()

	db.InitDB()
	db.Migrate()

	service := &InventoryService{}

	createReq := &pb.CreateProductRequest{
		Name:        "Test Get Product",
		Brand:       "BrandX",
		CategoryId:  1,
		Price:       111.11,
		Stock:       5,
		Description: "Created for GetProduct test",
	}

	createRes, err := service.CreateProduct(ctx, createReq)
	assert.NoError(t, err)
	assert.NotNil(t, createRes)


	getReq := &pb.GetProductRequest{Id: createRes.Product.Id}

	getRes, err := service.GetProduct(ctx, getReq)
	assert.NoError(t, err)
	assert.NotNil(t, getRes)
	assert.Equal(t, createReq.Name, getRes.Product.Name)
}

func TestUpdateProduct(t *testing.T) {
	ctx := context.Background()

	db.InitDB()
	db.Migrate()

	service := &InventoryService{}


	createReq := &pb.CreateProductRequest{
		Name:        "To Update",
		Brand:       "UpdateBrand",
		CategoryId:  1,
		Price:       100,
		Stock:       5,
		Description: "Before update",
	}
	createRes, err := service.CreateProduct(ctx, createReq)
	assert.NoError(t, err)
	assert.NotNil(t, createRes)


	updateReq := &pb.UpdateProductRequest{
		Id:          createRes.Product.Id,
		Name:        "Updated Name",
		Brand:       "Updated Brand",
		CategoryId:  1,
		Price:       199.99,
		Stock:       10,
		Description: "After update",
	}
	updateRes, err := service.UpdateProduct(ctx, updateReq)
	assert.NoError(t, err)
	assert.NotNil(t, updateRes)
	assert.Equal(t, updateReq.Name, updateRes.Product.Name)
	assert.Equal(t, updateReq.Price, updateRes.Product.Price)
}

func TestDeleteProduct(t *testing.T) {
	ctx := context.Background()

	service := &InventoryService{}

	req := &pb.DeleteProductRequest{Id: 1}

	_, err := service.DeleteProduct(ctx, req)

	assert.NoError(t, err)
}
