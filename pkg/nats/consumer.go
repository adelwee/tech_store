package nats

import (
	pb "Assignment2_AdelKenesova/inventory_service/proto"
	"Assignment2_AdelKenesova/pkg/events"
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

var inventoryClient pb.InventoryServiceClient

func SetInventoryClient(client pb.InventoryServiceClient) {
	inventoryClient = client
}

func SubscribeToProductCreated(nc *nats.Conn) error {
	_, err := nc.Subscribe("product.created", func(m *nats.Msg) {
		var event events.ProductCreatedEvent
		if err := json.Unmarshal(m.Data, &event); err != nil {
			log.Println(" Failed to unmarshal ProductCreatedEvent:", err)
			return
		}
		log.Printf(" Received product.created event: %+v\n", event)
	})
	return err
}

func SubscribeToOrderCreated(nc *nats.Conn) error {
	_, err := nc.Subscribe("order.created", func(m *nats.Msg) {
		var event events.OrderCreatedEvent
		if err := json.Unmarshal(m.Data, &event); err != nil {
			log.Println(" Failed to unmarshal OrderCreatedEvent:", err)
			return
		}

		log.Printf("Received order.created event: %+v\n", event)

		// Уменьшение stock
		for _, item := range event.Items {
			log.Printf("🔧 Decreasing stock for product %d by %d", item.ProductID, item.Quantity)

			ctx := context.Background()
			_, err := inventoryClient.DecreaseStock(ctx, &pb.DecreaseStockRequest{
				ProductId: item.ProductID,
				Quantity:  item.Quantity,
			})
			if err != nil {
				log.Printf("Failed to decrease stock for product %d: %v", item.ProductID, err)
			} else {
				log.Printf("Stock decreased for product %d", item.ProductID)
			}
		}
	})
	return err
}
