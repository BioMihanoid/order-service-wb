package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"order-service-wb/internal/kafka"
	"order-service-wb/internal/models"
	"order-service-wb/pkg/config"
)

func main() {
	cfg := config.NewConfig()

	prod, err := kafka.NewProducer([]string{cfg.Kafka.Broker}, cfg.Kafka.Topic)
	if err != nil {
		log.Fatalf("failed to create Kafka producer: %v", err)
	}

	defer prod.Close()

	ctx := context.Background()

	for i := 0; i < 100; i++ {
		order := generateFakeOrder()
		data, err := json.Marshal(order)
		if err != nil {
			log.Printf("failed to marshal order: %v", err)
			continue
		}

		if err = prod.Send(ctx, order.OrderUID, data); err != nil {
			log.Printf("failed to send to Kafka: %v", err)
		} else {
			log.Printf("sent order %s to Kafka", order.OrderUID)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func generateFakeOrder() models.Order {
	uid := uuid.New().String()
	return models.Order{
		OrderUID:    uid,
		TrackNumber: "WBTRACK" + randSeq(4),
		Entry:       "WBIL",
		Locale:      "en",
		InternalSig: "",
		CustomerID:  "testuser",
		DeliverySrv: "meest",
		ShardKey:    "1",
		SmID:        rand.Intn(100),
		DateCreated: time.Now().UTC(),
		OofShard:    "1",
		Delivery: models.Delivery{
			Name:   "Test User",
			Phone:  "+1234567890",
			Zip:    "123456",
			City:   "TestCity",
			Addr:   "123 Test St",
			Region: "TestRegion",
			Email:  "test@example.com",
		},
		Payment: models.Payment{
			Transaction:  uid,
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1000,
			PaymentDT:    time.Now().Unix(),
			Bank:         "alpha",
			DeliveryCost: 500,
			GoodsTotal:   500,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      rand.Intn(100000),
				TrackNumber: "WBTRACK" + randSeq(4),
				Price:       500,
				Rid:         uuid.New().String(),
				Name:        "Some Product",
				Sale:        0,
				Size:        "L",
				TotalPrice:  500,
				NmID:        rand.Intn(10000),
				Brand:       "BrandName",
				Status:      202,
			},
		},
	}
}

func randSeq(n int) string {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
