package catalog

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("product not found")

type Service interface {
	ListProducts(ctx context.Context, limit, offset int) ([]Product, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*Product, error)
	CreateProduct(ctx context.Context, sellerID *uuid.UUID, req CreateProductRequest) (*Product, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req UpdateProductRequest) (*Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	ListCategories(ctx context.Context) ([]Category, error)
	ListVariants(ctx context.Context, productID uuid.UUID) ([]Variant, error)
	CreateVariant(ctx context.Context, productID uuid.UUID, req CreateVariantRequest) (*Variant, error)
	DeleteVariant(ctx context.Context, id uuid.UUID) error
	SetStock(ctx context.Context, variantID uuid.UUID, quantity int) error
	ListDrops(ctx context.Context) ([]Drop, error)
	GetDropBySlug(ctx context.Context, slug string) (*Drop, error)
	GetDropProducts(ctx context.Context, dropID uuid.UUID) ([]Product, error)
	CreateDrop(ctx context.Context, req CreateDropRequest) (*Drop, error)
	UpdateDrop(ctx context.Context, id uuid.UUID, req UpdateDropRequest) (*Drop, error)
	DeleteDrop(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ListProducts(ctx context.Context, limit, offset int) ([]Product, error) {
	return s.repo.ListProducts(ctx, limit, offset)
}

func (s *service) GetProduct(ctx context.Context, id uuid.UUID) (*Product, error) {
	return s.repo.GetProductByID(ctx, id)
}

func (s *service) CreateProduct(ctx context.Context, sellerID *uuid.UUID, req CreateProductRequest) (*Product, error) {
	return s.repo.CreateProduct(ctx, sellerID, req)
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

func (s *service) ListVariants(ctx context.Context, productID uuid.UUID) ([]Variant, error) {
	return s.repo.ListVariants(ctx, productID)
}

func (s *service) CreateVariant(ctx context.Context, productID uuid.UUID, req CreateVariantRequest) (*Variant, error) {
	return s.repo.CreateVariant(ctx, productID, req)
}

func (s *service) DeleteVariant(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteVariant(ctx, id)
}

func (s *service) SetStock(ctx context.Context, variantID uuid.UUID, quantity int) error {
	return s.repo.SetStock(ctx, variantID, quantity)
}

func (s *service) ListDrops(ctx context.Context) ([]Drop, error) {
	return s.repo.ListDrops(ctx)
}

func (s *service) GetDropBySlug(ctx context.Context, slug string) (*Drop, error) {
	return s.repo.GetDropBySlug(ctx, slug)
}

func (s *service) GetDropProducts(ctx context.Context, dropID uuid.UUID) ([]Product, error) {
	return s.repo.GetDropProducts(ctx, dropID)
}

func (s *service) CreateDrop(ctx context.Context, req CreateDropRequest) (*Drop, error) {
	return s.repo.CreateDrop(ctx, req)
}

func (s *service) UpdateDrop(ctx context.Context, id uuid.UUID, req UpdateDropRequest) (*Drop, error) {
	return s.repo.UpdateDrop(ctx, id, req)
}

func (s *service) DeleteDrop(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteDrop(ctx, id)
}
