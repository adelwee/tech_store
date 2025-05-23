package service

import (
	"Assignment2_AdelKenesova/order_service/internal/db"
	"Assignment2_AdelKenesova/order_service/internal/models"
	pb "Assignment2_AdelKenesova/order_service/proto"
	rdb "Assignment2_AdelKenesova/pkg/redis"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrder(t *testing.T) {
	// ✅ Инициализация БД
	db.InitDB()

	// 🔄 Очистка таблиц
	db.DB.Exec("DELETE FROM order_items")
	db.DB.Exec("DELETE FROM orders")

	service := &OrderService{}

	req := &pb.CreateOrderRequest{
		UserId: 1,
		Items: []*pb.OrderItem{
			{ProductId: 1, Quantity: 2, Price: 100},
			{ProductId: 2, Quantity: 1, Price: 50},
		},
	}

	resp, err := service.CreateOrder(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, float64(250), resp.Order.TotalPrice)
	assert.Equal(t, "pending", resp.Order.Status)
	assert.Len(t, resp.Order.Items, 2)
}

func TestGetOrder(t *testing.T) {
	db.InitDB()
	service := &OrderService{}

	// Создаем тестовый заказ
	order := models.Order{
		UserID:     2,
		TotalPrice: 200,
		Status:     "pending",
		OrderItems: []models.OrderItem{
			{ProductID: 1, Quantity: 2, Price: 50},
			{ProductID: 2, Quantity: 1, Price: 100},
		},
	}
	err := db.DB.Create(&order).Error
	assert.NoError(t, err)

	// Тестируем GetOrder
	resp, err := service.GetOrder(context.Background(), &pb.GetOrderRequest{
		Id: uint64(order.ID),
	})
	assert.NoError(t, err)
	assert.Equal(t, uint64(order.ID), resp.Order.Id)
	assert.Equal(t, uint64(order.UserID), resp.Order.UserId)
	assert.Equal(t, len(order.OrderItems), len(resp.Order.Items))
}

func TestDeleteOrder(t *testing.T) {
	db.InitDB()
	service := &OrderService{}

	// Создаем заказ
	order := models.Order{
		UserID:     3,
		TotalPrice: 100,
		Status:     "pending",
	}
	err := db.DB.Create(&order).Error
	assert.NoError(t, err)

	// Удаляем заказ
	_, err = service.DeleteOrder(context.Background(), &pb.DeleteOrderRequest{
		Id: uint64(order.ID),
	})
	assert.NoError(t, err)

	// Проверяем, что заказ удален
	var check models.Order
	err = db.DB.First(&check, order.ID).Error
	assert.Error(t, err)
}

func TestListOrders(t *testing.T) {
	db.InitDB()
	rdb.InitRedis() // ✅ Добавь это

	service := &OrderService{}

	// Очищаем базу
	db.DB.Exec("DELETE FROM order_items")
	db.DB.Exec("DELETE FROM orders")

	// Создаем тестовые заказы
	orders := []models.Order{
		{
			UserID:     4,
			TotalPrice: 200,
			Status:     "pending",
			OrderItems: []models.OrderItem{{ProductID: 1, Quantity: 1, Price: 200}},
		},
		{
			UserID:     5,
			TotalPrice: 150,
			Status:     "pending",
			OrderItems: []models.OrderItem{{ProductID: 2, Quantity: 3, Price: 50}},
		},
	}

	for _, o := range orders {
		err := db.DB.Create(&o).Error
		assert.NoError(t, err)
	}

	// Тестируем получение списка заказов
	resp, err := service.ListOrders(context.Background(), &pb.Empty{})
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(resp.Orders), 2)
}
