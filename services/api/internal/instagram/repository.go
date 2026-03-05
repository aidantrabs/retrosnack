package instagram

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetByProductID(ctx context.Context, productID uuid.UUID) (*Link, error)
	Upsert(ctx context.Context, productID uuid.UUID, postURL, embedHTML string) (*Link, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) GetByProductID(ctx context.Context, productID uuid.UUID) (*Link, error) {
	var l Link
	err := r.db.QueryRow(ctx,
		`SELECT id, product_id, post_url, embed_html, cached_at
		 FROM instagram_links WHERE product_id = $1`,
		productID,
	).Scan(&l.ID, &l.ProductID, &l.PostURL, &l.EmbedHTML, &l.CachedAt)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *repository) Upsert(ctx context.Context, productID uuid.UUID, postURL, embedHTML string) (*Link, error) {
	var l Link
	err := r.db.QueryRow(ctx,
		`INSERT INTO instagram_links (product_id, post_url, embed_html, cached_at)
		 VALUES ($1, $2, $3, NOW())
		 ON CONFLICT (product_id)
		 DO UPDATE SET post_url = $2, embed_html = $3, cached_at = NOW()
		 RETURNING id, product_id, post_url, embed_html, cached_at`,
		productID, postURL, embedHTML,
	).Scan(&l.ID, &l.ProductID, &l.PostURL, &l.EmbedHTML, &l.CachedAt)
	if err != nil {
		return nil, err
	}
	return &l, nil
}
