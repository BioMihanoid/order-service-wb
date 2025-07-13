package service_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"order-service-wb/internal/models"
	"order-service-wb/internal/service"
	"order-service-wb/mocks"
)

func TestGetOrderByID_CacheHit(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.OrderRepository)
	mockCache := new(mocks.Cache)

	testOrder := models.Order{OrderUID: "123"}

	mockCache.On("Get", "123").Return(testOrder, true)

	srv := service.NewOrderService(mockRepo, mockCache)

	order, err := srv.GetOrderByID(context.Background(), "123")

	assert.NoError(t, err)
	assert.Equal(t, &testOrder, order)

	mockCache.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "GetOrderByID")
}

func TestGetOrderByID_CacheMissDBHit(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.OrderRepository)
	mockCache := new(mocks.Cache)

	testOrder := &models.Order{OrderUID: "123"}

	mockCache.On("Get", "123").Return(models.Order{}, false)
	mockRepo.On("GetOrderByID", mock.Anything, "123").Return(testOrder, nil)
	mockCache.On("Set", "123", *testOrder).Return()

	srv := service.NewOrderService(mockRepo, mockCache)

	order, err := srv.GetOrderByID(context.Background(), "123")

	assert.NoError(t, err)
	assert.Equal(t, testOrder, order)

	mockCache.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestCreateOrder_Success(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.OrderRepository)
	mockCache := new(mocks.Cache)

	testOrder := generateFakeOrder("123")

	mockRepo.On("CreateOrder", mock.Anything, testOrder).Return(nil)
	mockCache.On("Set", "123", *testOrder)

	srv := service.NewOrderService(mockRepo, mockCache)

	err := srv.CreateOrder(context.Background(), testOrder)

	assert.NoError(t, err)

	mockCache.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestCreateOrder_FailedValidate(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.OrderRepository)
	mockCache := new(mocks.Cache)

	testOrder := &models.Order{OrderUID: "123"}

	srv := service.NewOrderService(mockRepo, mockCache)

	err := srv.CreateOrder(context.Background(), testOrder)

	assert.Error(t, err)

	mockRepo.AssertNotCalled(t, "CreateOrder")
	mockCache.AssertNotCalled(t, "Set")
}

func TestCreateOrder_FailedDB(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.OrderRepository)
	mockCache := new(mocks.Cache)

	testOrder := generateFakeOrder("123")

	mockRepo.On("CreateOrder", mock.Anything, testOrder).Return(fmt.Errorf("error"))

	srv := service.NewOrderService(mockRepo, mockCache)

	err := srv.CreateOrder(context.Background(), testOrder)

	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
	mockCache.AssertNotCalled(t, "Set")
}

func TestLoadCache_Success(t *testing.T) {
	t.Parallel()

	mockRepo := new(mocks.OrderRepository)
	mockCache := new(mocks.Cache)

	testOrder := &models.Order{OrderUID: "123"}

	mockRepo.On("GetAllOrders", mock.Anything, 1).Return([]*models.Order{testOrder}, nil)
	mockCache.On("Set", "123", *testOrder).Return()

	srv := service.NewOrderService(mockRepo, mockCache)

	err := srv.LoadCache(context.Background(), 1)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func generateFakeOrder(id string) *models.Order {
	return &models.Order{
		OrderUID:    id,
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
			Transaction:  id,
			RequestID:    "2234555",
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
