package catalog

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("product not found")

type Service interface {
	ListProducts(ctx context.Context) ([]Product, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*Product, error)
	CreateProduct(ctx context.Context, req CreateProductRequest) (*Product, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req UpdateProductRequest) (*Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	ListCategories(ctx context.Context) ([]Category, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ListProducts(ctx context.Context) ([]Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *service) GetProduct(ctx context.Context, id uuid.UUID) (*Product, error) {
	return s.repo.GetProductByID(ctx, id)
}

func (s *service) CreateProduct(ctx context.Context, req CreateProductRequest) (*Product, error) {
	return s.repo.CreateProduct(ctx, req)
}

func (s *service) UpdateProduct(ctx context.Context, id uuid.UUID, req UpdateProductRequest) (*Product, error) {
	return s.repo.UpdateProduct(ctx, id, req)
}

func (s *service) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteProduct(ctx, id)
}

func (s *service) ListCategories(ctx context.Context) ([]Category, error) {
	return s.repo.ListCategories(ctx)
}
