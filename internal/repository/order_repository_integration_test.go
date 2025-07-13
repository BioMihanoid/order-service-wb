package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"order-service-wb/internal/models"
	"order-service-wb/internal/repository"
)

var db *sql.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	containerReq := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "pass",
			"POSTGRES_USER":     "user",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithStartupTimeout(30 * time.Second),
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
	})
	require.NoError(nil, err)

	defer postgresC.Terminate(ctx)

	host, err := postgresC.Host(ctx)
	require.NoError(nil, err)
	port, err := postgresC.MappedPort(ctx, "5432")
	require.NoError(nil, err)

	dsn := fmt.Sprintf("postgres://user:pass@%s:%s/testdb?sslmode=disable", host, port.Port())
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return
	}

	err = goose.Up(db, "./migrations")
	if err != nil {
		return
	}

	code := m.Run()

	os.Exit(code)
}

func TestCreateAndGetOrder_Success(t *testing.T) {
	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewOrderRepository(dbx)

	testOrder := &models.Order{
		OrderUID: "123",
	}

	err := repo.CreateOrder(context.Background(), testOrder)
	require.NoError(t, err)

	fetched, err := repo.GetOrderByID(context.Background(), "123")
	require.NoError(t, err)
	require.Equal(t, testOrder.OrderUID, fetched.OrderUID)
}

func TestGetAllOrders_Success(t *testing.T) {
	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewOrderRepository(dbx)

	order1 := generateFakeOrder("123")
	order2 := generateFakeOrder("321")

	err := repo.CreateOrder(context.Background(), order1)
	require.NoError(t, err)

	err = repo.CreateOrder(context.Background(), order2)
	require.NoError(t, err)

	orders, err := repo.GetAllOrders(context.Background(), 2)
	require.NoError(t, err)

	var found1, found2 bool
	for _, o := range orders {
		if o.OrderUID == order1.OrderUID {
			found1 = true
		}
		if o.OrderUID == order2.OrderUID {
			found2 = true
		}
	}
	require.True(t, found1, "order1 not found in fetched orders")
	require.True(t, found2, "order2 not found in fetched orders")
}

func TestCreateOrder_Conflict(t *testing.T) {
	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewOrderRepository(dbx)

	order := generateFakeOrder("123")

	err := repo.CreateOrder(context.Background(), order)
	require.NoError(t, err)

	err = repo.CreateOrder(context.Background(), order)
	require.Error(t, err)
}

func TestGetOrderByID_NotFound(t *testing.T) {
	dbx := sqlx.NewDb(db, "postgres")
	repo := repository.NewOrderRepository(dbx)

	_, err := repo.GetOrderByID(context.Background(), "123")
	require.Error(t, err)
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
