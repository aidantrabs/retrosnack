package media

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateImage(ctx context.Context, productID uuid.UUID, r2Key, url string, position int) (*ProductImageRecord, error)
	ListByProduct(ctx context.Context, productID uuid.UUID) ([]ProductImageRecord, error)
	CountByProduct(ctx context.Context, productID uuid.UUID) (int, error)
	DeleteImage(ctx context.Context, id uuid.UUID) (*ProductImageRecord, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) CreateImage(ctx context.Context, productID uuid.UUID, r2Key, url string, position int) (*ProductImageRecord, error) {
	var img ProductImageRecord
	err := r.db.QueryRow(ctx,
		`INSERT INTO product_images (product_id, r2_key, url, position)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, product_id, url, position`,
		productID, r2Key, url, position,
	).Scan(&img.ID, &img.ProductID, &img.URL, &img.Position)
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func (r *repository) ListByProduct(ctx context.Context, productID uuid.UUID) ([]ProductImageRecord, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, product_id, url, position
		 FROM product_images WHERE product_id = $1 ORDER BY position`,
		productID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	images := make([]ProductImageRecord, 0)
	for rows.Next() {
		var img ProductImageRecord
		if err := rows.Scan(&img.ID, &img.ProductID, &img.URL, &img.Position); err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, rows.Err()
}

func (r *repository) CountByProduct(ctx context.Context, productID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM product_images WHERE product_id = $1`,
		productID,
	).Scan(&count)
	return count, err
}

func (r *repository) DeleteImage(ctx context.Context, id uuid.UUID) (*ProductImageRecord, error) {
	var img ProductImageRecord
	err := r.db.QueryRow(ctx,
		`DELETE FROM product_images WHERE id = $1
		 RETURNING id, product_id, url, position`,
		id,
	).Scan(&img.ID, &img.ProductID, &img.URL, &img.Position)
	if err != nil {
		return nil, err
	}
	return &img, nil
}
