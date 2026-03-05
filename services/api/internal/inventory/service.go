package inventory

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	GetStock(ctx context.Context, variantID uuid.UUID) (*StockItem, error)
	Reserve(ctx context.Context, variantID uuid.UUID, qty int) error
	Release(ctx context.Context, variantID uuid.UUID, qty int) error
	Deduct(ctx context.Context, variantID uuid.UUID, qty int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetStock(ctx context.Context, variantID uuid.UUID) (*StockItem, error) {
	return s.repo.GetStock(ctx, variantID)
}

func (s *service) Reserve(ctx context.Context, variantID uuid.UUID, qty int) error {
	return s.repo.Reserve(ctx, variantID, qty)
}

func (s *service) Release(ctx context.Context, variantID uuid.UUID, qty int) error {
	return s.repo.Release(ctx, variantID, qty)
}

func (s *service) Deduct(ctx context.Context, variantID uuid.UUID, qty int) error {
	return s.repo.Deduct(ctx, variantID, qty)
}
