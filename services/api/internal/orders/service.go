package orders

import (
	"context"

	"github.com/google/uuid"
	"github.com/MobinaToorani/retrosnack/internal/inventory"
)

type Service interface {
	CreateOrder(ctx context.Context, userID *uuid.UUID, req CreateOrderRequest) (*Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*Order, error)
	GetOrderByStripeSession(ctx context.Context, sessionID string) (*Order, error)
	MarkPaid(ctx context.Context, orderID uuid.UUID) error
	SetStripeSession(ctx context.Context, orderID uuid.UUID, sessionID string) error
}

type service struct {
	repo      Repository
	inventory inventory.Service
}

func NewService(repo Repository, inv inventory.Service) Service {
	return &service{repo: repo, inventory: inv}
}

func (s *service) CreateOrder(ctx context.Context, userID *uuid.UUID, req CreateOrderRequest) (*Order, error) {
	// Reserve inventory for each item before creating the order
	for _, item := range req.Items {
		if err := s.inventory.Reserve(ctx, item.VariantID, item.Quantity); err != nil {
			return nil, err
		}
	}

	var total int64
	for _, item := range req.Items {
		total += item.PriceCents * int64(item.Quantity)
	}

	return s.repo.CreateOrder(ctx, userID, req.Items, total)
}

func (s *service) GetOrder(ctx context.Context, id uuid.UUID) (*Order, error) {
	return s.repo.GetOrderByID(ctx, id)
}

func (s *service) GetOrderByStripeSession(ctx context.Context, sessionID string) (*Order, error) {
	return s.repo.GetOrderByStripeSession(ctx, sessionID)
}

func (s *service) MarkPaid(ctx context.Context, orderID uuid.UUID) error {
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	// Deduct inventory — items are now committed
	for _, item := range order.Items {
		if err := s.inventory.Deduct(ctx, item.VariantID, item.Quantity); err != nil {
			return err
		}
	}

	return s.repo.UpdateStatus(ctx, orderID, StatusPaid)
}

func (s *service) SetStripeSession(ctx context.Context, orderID uuid.UUID, sessionID string) error {
	return s.repo.SetStripeSession(ctx, orderID, sessionID)
}
