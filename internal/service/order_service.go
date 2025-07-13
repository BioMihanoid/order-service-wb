package service

import (
	"context"

	"github.com/go-playground/validator/v10"

	"order-service-wb/internal/cache"
	"order-service-wb/internal/models"
	"order-service-wb/internal/repository"
)

type OrderService interface {
	GetOrderByID(ctx context.Context, orderID string) (*models.Order, error)
	LoadCache(ctx context.Context, limit int) error
	CreateOrder(ctx context.Context, order *models.Order) error
}

type Service struct {
	repo      repository.OrderRepository
	cache     cache.Cache
	validator *validator.Validate
}

func NewOrderService(repo repository.OrderRepository, cache cache.Cache) OrderService {
	return &Service{
		repo:      repo,
		cache:     cache,
		validator: validator.New(),
	}
}

func (s *Service) GetOrderByID(ctx context.Context, orderID string) (*models.Order, error) {
	if order, ok := s.cache.Get(orderID); ok {
		return &order, nil
	}

	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err == nil && order != nil {
		s.cache.Set(orderID, *order)
	}

	return order, err
}

func (s *Service) LoadCache(ctx context.Context, limit int) error {
	orders, err := s.repo.GetAllOrders(ctx, limit)
	if err != nil {
		return err
	}
	for _, order := range orders {
		s.cache.Set(order.OrderUID, *order)
	}
	return nil
}

func (s *Service) CreateOrder(ctx context.Context, order *models.Order) error {
	if err := s.validator.Struct(order); err != nil {
		return err
	}
	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return err
	}
	s.cache.Set(order.OrderUID, *order)
	return nil
}
