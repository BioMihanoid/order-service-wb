package service

import (
	"context"

	"order-service-wb/internal/models"
	"order-service-wb/internal/repository"
)

type OrderService interface {
	GetOrderByID(ctx context.Context, orderID string) (*models.Order, error)
}

type Service struct {
	repo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) OrderService {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetOrderByID(ctx context.Context, orderID string) (*models.Order, error) {
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	return order, nil
}
