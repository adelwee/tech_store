package nats

import (
	"Assignment2_AdelKenesova/pkg/events"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

var natsConn *nats.Conn

func InitNATS(url string) error {
	var err error
	natsConn, err = nats.Connect(url)
	return err
}

func GetConn() *nats.Conn {
	return natsConn
}

func PublishProductCreated(event *events.ProductCreatedEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	err = natsConn.Publish("product.created", data)
	if err != nil {
		return err
	}
	log.Println(" Published product.created event to NATS")
	return nil
}

func Publish(subject string, data any) error {
	conn := GetConn()

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return conn.Publish(subject, payload)
}
