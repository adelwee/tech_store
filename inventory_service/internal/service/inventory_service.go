package service

import (
	"Assignment2_AdelKenesova/inventory_service/internal/db"
	"Assignment2_AdelKenesova/inventory_service/internal/models"
	pb "Assignment2_AdelKenesova/inventory_service/proto"
	redis "Assignment2_AdelKenesova/pkg/redis"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type InventoryService struct {
	pb.UnimplementedInventoryServiceServer
}

func (s *InventoryService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	product := models.Product{
		Name:        req.Name,
		Brand:       req.Brand,
		CategoryID:  uint(req.CategoryId),
		Price:       req.Price,
		Stock:       uint(req.Stock),
		Description: req.Description,
	}

	if err := db.DB.Create(&product).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create product: %v", err)
	}

	// Кэширование в Redis
	client := redis.GetClient()
	cacheCtx := context.Background()

	productJSON, err := json.Marshal(product)
	if err == nil {
		err = client.Set(cacheCtx, fmt.Sprintf("product:%d", product.ID), productJSON, 0).Err()
		if err != nil {
			log.Println("Failed to cache product in Redis:", err)
		} else {
			log.Println("Cached product in Redis")
		}
	} else {
		log.Println(" Failed to marshal product for Redis:", err)
	}

	return &pb.ProductResponse{
		Product: &pb.Product{
			Id:          uint64(product.ID),
			Name:        product.Name,
			Brand:       product.Brand,
			CategoryId:  uint64(product.CategoryID),
			Price:       product.Price,
			Stock:       uint64(product.Stock),
			Description: product.Description,
		},
	}, nil
}

func (s *InventoryService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	cacheKey := fmt.Sprintf("product:%d", req.Id)
	client := redis.GetClient()
	cacheCtx := context.Background()

	//  Попытка получить продукт из Redis
	cached, err := client.Get(cacheCtx, cacheKey).Result()
	if err == nil {
		var cachedProduct models.Product
		if err := json.Unmarshal([]byte(cached), &cachedProduct); err == nil {
			log.Println(" Product returned from Redis cache")
			return &pb.ProductResponse{
				Product: &pb.Product{
					Id:          uint64(cachedProduct.ID),
					Name:        cachedProduct.Name,
					Brand:       cachedProduct.Brand,
					CategoryId:  uint64(cachedProduct.CategoryID),
					Price:       cachedProduct.Price,
					Stock:       uint64(cachedProduct.Stock),
					Description: cachedProduct.Description,
				},
			}, nil
		}
	}

	// Если не найден — получаем из БД
	var product models.Product
	if err := db.DB.First(&product, req.Id).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "Product not found")
	}

	// Кэшируем в Redis
	productJSON, err := json.Marshal(product)
	if err == nil {
		err = client.Set(cacheCtx, cacheKey, productJSON, 0).Err()
		if err != nil {
			log.Println(" Failed to cache product in Redis:", err)
		} else {
			log.Println(" Product cached in Redis")
		}
	}

	return &pb.ProductResponse{
		Product: &pb.Product{
			Id:          uint64(product.ID),
			Name:        product.Name,
			Brand:       product.Brand,
			CategoryId:  uint64(product.CategoryID),
			Price:       product.Price,
			Stock:       uint64(product.Stock),
			Description: product.Description,
		},
	}, nil
}

func (s *InventoryService) ListProducts(ctx context.Context, _ *pb.Empty) (*pb.ProductListResponse, error) {
	var products []models.Product
	if err := db.DB.Find(&products).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to list products")
	}

	var protoProducts []*pb.Product
	for _, p := range products {
		protoProducts = append(protoProducts, &pb.Product{
			Id:          uint64(p.ID),
			Name:        p.Name,
			Brand:       p.Brand,
			CategoryId:  uint64(p.CategoryID),
			Price:       p.Price,
			Stock:       uint64(p.Stock),
			Description: p.Description,
		})
	}

	return &pb.ProductListResponse{Products: protoProducts}, nil
}

func (s *InventoryService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.ProductResponse, error) {
	var product models.Product

	// Проверка существования
	if err := db.DB.First(&product, req.Id).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "Product not found")
	}

	// Обновление полей
	product.Name = req.Name
	product.Brand = req.Brand
	product.CategoryID = uint(req.CategoryId)
	product.Price = req.Price
	product.Stock = uint(req.Stock)
	product.Description = req.Description

	// Сохранение
	if err := db.DB.Save(&product).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update product")
	}

	//  Обновление кэша в Redis
	productJSON, err := json.Marshal(product)
	if err == nil {
		err = redis.GetClient().Set(context.Background(), fmt.Sprintf("product:%d", product.ID), productJSON, 0).Err()
		if err != nil {
			log.Println(" Failed to update product cache in Redis:", err)
		} else {
			log.Println(" Product cache updated in Redis")
		}
	}

	// Ответ
	return &pb.ProductResponse{
		Product: &pb.Product{
			Id:          uint64(product.ID),
			Name:        product.Name,
			Brand:       product.Brand,
			CategoryId:  uint64(product.CategoryID),
			Price:       product.Price,
			Stock:       uint64(product.Stock),
			Description: product.Description,
		},
	}, nil
}

func (s *InventoryService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.Empty, error) {
	// Удаление из БД
	if err := db.DB.Delete(&models.Product{}, req.Id).Error; err != nil {
		return nil, err
	}

	//  Удаление из кэша
	err := redis.GetClient().Del(context.Background(), fmt.Sprintf("product:%d", req.Id)).Err()
	if err != nil {
		log.Println(" Failed to delete product cache from Redis:", err)
	} else {
		log.Println("🗑Product cache deleted from Redis")
	}

	return &pb.Empty{}, nil
}

func (s *InventoryService) DecreaseStock(ctx context.Context, req *pb.DecreaseStockRequest) (*pb.Empty, error) {
	var product models.Product
	if err := db.DB.First(&product, req.ProductId).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "Product not found")
	}

	if product.Stock < uint(req.Quantity) {
		return nil, status.Errorf(codes.InvalidArgument, "Not enough stock")
	}

	product.Stock -= uint(req.Quantity)

	if err := db.DB.Save(&product).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update stock")
	}

	return &pb.Empty{}, nil
}
