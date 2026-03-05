package inventory

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInsufficientStock = errors.New("insufficient stock")

type Repository interface {
	GetStock(ctx context.Context, variantID uuid.UUID) (*StockItem, error)
	Reserve(ctx context.Context, variantID uuid.UUID, qty int) error
	Release(ctx context.Context, variantID uuid.UUID, qty int) error
	Deduct(ctx context.Context, variantID uuid.UUID, qty int) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) GetStock(ctx context.Context, variantID uuid.UUID) (*StockItem, error) {
	var s StockItem
	err := r.db.QueryRow(ctx,
		`SELECT variant_id, quantity, reserved, quantity - reserved AS available
		 FROM inventory WHERE variant_id = $1`,
		variantID,
	).Scan(&s.VariantID, &s.Quantity, &s.Reserved, &s.Available)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) Reserve(ctx context.Context, variantID uuid.UUID, qty int) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE inventory
		 SET reserved = reserved + $2
		 WHERE variant_id = $1 AND (quantity - reserved) >= $2`,
		variantID, qty,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrInsufficientStock
	}
	return nil
}

func (r *repository) Release(ctx context.Context, variantID uuid.UUID, qty int) error {
	_, err := r.db.Exec(ctx,
		`UPDATE inventory
		 SET reserved = GREATEST(0, reserved - $2)
		 WHERE variant_id = $1`,
		variantID, qty,
	)
	return err
}

func (r *repository) Deduct(ctx context.Context, variantID uuid.UUID, qty int) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE inventory
		 SET quantity = quantity - $2, reserved = GREATEST(0, reserved - $2)
		 WHERE variant_id = $1 AND quantity >= $2`,
		variantID, qty,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrInsufficientStock
	}
	return nil
}
