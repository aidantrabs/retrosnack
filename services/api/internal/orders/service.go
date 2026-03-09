package orders

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/retrosnack-clothing/retrosnack/internal/inventory"
)

var ErrInvalidTransition = errors.New("invalid status transition")

type Service interface {
	CreateOrder(ctx context.Context, userID *uuid.UUID, req CreateOrderRequest) (*Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*Order, error)
	GetOrderByCheckoutSession(ctx context.Context, sessionID string) (*Order, error)
	ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Order, error)
	ListAll(ctx context.Context, limit, offset int) ([]Order, error)
	MarkPaid(ctx context.Context, orderID uuid.UUID) error
	MarkShipped(ctx context.Context, orderID uuid.UUID) error
	MarkDelivered(ctx context.Context, orderID uuid.UUID) error
	CancelOrder(ctx context.Context, orderID uuid.UUID) error
	SetCheckoutSession(ctx context.Context, orderID uuid.UUID, sessionID string) error
}

type service struct {
	repo      Repository
	inventory inventory.Service
}

func NewService(repo Repository, inv inventory.Service) Service {
	return &service{repo: repo, inventory: inv}
}

func (s *service) CreateOrder(ctx context.Context, userID *uuid.UUID, req CreateOrderRequest) (*Order, error) {
	var reserved []OrderItemInput

	for _, item := range req.Items {
		if err := s.inventory.Reserve(ctx, item.VariantID, item.Quantity); err != nil {
			s.releaseReserved(ctx, reserved)
			return nil, err
		}
		reserved = append(reserved, item)
	}

	var total int64
	for _, item := range req.Items {
		total += item.PriceCents * int64(item.Quantity)
	}

	order, err := s.repo.CreateOrder(ctx, userID, req.Items, total)
	if err != nil {
		s.releaseReserved(ctx, reserved)
		return nil, err
	}

	return order, nil
}

func (s *service) releaseReserved(ctx context.Context, items []OrderItemInput) {
	for _, item := range items {
		_ = s.inventory.Release(ctx, item.VariantID, item.Quantity)
	}
}

func (s *service) GetOrder(ctx context.Context, id uuid.UUID) (*Order, error) {
	return s.repo.GetOrderByID(ctx, id)
}

func (s *service) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]Order, error) {
	return s.repo.ListByUser(ctx, userID, limit, offset)
}

func (s *service) ListAll(ctx context.Context, limit, offset int) ([]Order, error) {
	return s.repo.ListAll(ctx, limit, offset)
}

func (s *service) GetOrderByCheckoutSession(ctx context.Context, sessionID string) (*Order, error) {
	return s.repo.GetOrderByCheckoutSession(ctx, sessionID)
}

func (s *service) MarkPaid(ctx context.Context, orderID uuid.UUID) error {
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	// idempotency: if already paid (or further along), skip
	if order.Status != StatusPending {
		return nil
	}

	if err := s.repo.UpdateStatus(ctx, orderID, StatusPaid); err != nil {
		return err
	}

	for _, item := range order.Items {
		if err := s.inventory.Deduct(ctx, item.VariantID, item.Quantity); err != nil {
			slog.Error("failed to deduct inventory after payment",
				"order_id", orderID,
				"variant_id", item.VariantID,
				"quantity", item.Quantity,
				"error", err,
			)
		}
	}

	return nil
}

func (s *service) MarkShipped(ctx context.Context, orderID uuid.UUID) error {
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}
	if order.Status != StatusPaid {
		return ErrInvalidTransition
	}
	return s.repo.UpdateStatus(ctx, orderID, StatusShipped)
}

func (s *service) MarkDelivered(ctx context.Context, orderID uuid.UUID) error {
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}
	if order.Status != StatusShipped {
		return ErrInvalidTransition
	}
	return s.repo.UpdateStatus(ctx, orderID, StatusDelivered)
}

func (s *service) CancelOrder(ctx context.Context, orderID uuid.UUID) error {
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}
	if order.Status != StatusPending {
		return ErrInvalidTransition
	}

	if err := s.repo.UpdateStatus(ctx, orderID, StatusCancelled); err != nil {
		return err
	}

	for _, item := range order.Items {
		_ = s.inventory.Release(ctx, item.VariantID, item.Quantity)
	}

	return nil
}

func (s *service) SetCheckoutSession(ctx context.Context, orderID uuid.UUID, sessionID string) error {
	return s.repo.SetCheckoutSession(ctx, orderID, sessionID)
}
