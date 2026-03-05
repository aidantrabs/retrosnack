package orders

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateOrder(ctx context.Context, userID *uuid.UUID, items []OrderItemInput, totalCents int64) (*Order, error)
	GetOrderByID(ctx context.Context, id uuid.UUID) (*Order, error)
	GetOrderByStripeSession(ctx context.Context, sessionID string) (*Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status Status) error
	SetStripeSession(ctx context.Context, id uuid.UUID, sessionID string) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) CreateOrder(ctx context.Context, userID *uuid.UUID, items []OrderItemInput, totalCents int64) (*Order, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var o Order
	err = tx.QueryRow(ctx,
		`INSERT INTO orders (user_id, status, total_cents)
		 VALUES ($1, 'pending', $2)
		 RETURNING id, user_id, status, total_cents, created_at, updated_at`,
		userID, totalCents,
	).Scan(&o.ID, &o.UserID, &o.Status, &o.TotalCents, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		var oi OrderItem
		err = tx.QueryRow(ctx,
			`INSERT INTO order_items (order_id, variant_id, quantity, price_cents)
			 VALUES ($1, $2, $3, $4)
			 RETURNING id, order_id, variant_id, quantity, price_cents`,
			o.ID, item.VariantID, item.Quantity, item.PriceCents,
		).Scan(&oi.ID, &oi.OrderID, &oi.VariantID, &oi.Quantity, &oi.PriceCents)
		if err != nil {
			return nil, err
		}
		o.Items = append(o.Items, oi)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *repository) GetOrderByID(ctx context.Context, id uuid.UUID) (*Order, error) {
	var o Order
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, status, total_cents, stripe_session_id, created_at, updated_at
		 FROM orders WHERE id = $1`,
		id,
	).Scan(&o.ID, &o.UserID, &o.Status, &o.TotalCents, &o.StripeSessionID, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *repository) GetOrderByStripeSession(ctx context.Context, sessionID string) (*Order, error) {
	var o Order
	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, status, total_cents, stripe_session_id, created_at, updated_at
		 FROM orders WHERE stripe_session_id = $1`,
		sessionID,
	).Scan(&o.ID, &o.UserID, &o.Status, &o.TotalCents, &o.StripeSessionID, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *repository) UpdateStatus(ctx context.Context, id uuid.UUID, status Status) error {
	_, err := r.db.Exec(ctx,
		`UPDATE orders SET status = $2, updated_at = NOW() WHERE id = $1`,
		id, status,
	)
	return err
}

func (r *repository) SetStripeSession(ctx context.Context, id uuid.UUID, sessionID string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE orders SET stripe_session_id = $2, updated_at = NOW() WHERE id = $1`,
		id, sessionID,
	)
	return err
}
