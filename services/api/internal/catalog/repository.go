package catalog

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	ListProducts(ctx context.Context, limit, offset int) ([]Product, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*Product, error)
	CreateProduct(ctx context.Context, sellerID *uuid.UUID, req CreateProductRequest) (*Product, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req UpdateProductRequest) (*Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	ListCategories(ctx context.Context) ([]Category, error)
	ListVariants(ctx context.Context, productID uuid.UUID) ([]Variant, error)
	CreateVariant(ctx context.Context, productID uuid.UUID, req CreateVariantRequest) (*Variant, error)
	DeleteVariant(ctx context.Context, id uuid.UUID) error
	SetStock(ctx context.Context, variantID uuid.UUID, quantity int) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) ListProducts(ctx context.Context, limit, offset int) ([]Product, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, title, description, category_id, brand, condition,
		        price_cents, seller_id, instagram_post_url, created_at, updated_at
		 FROM products
		 ORDER BY created_at DESC
		 LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]Product, 0)
	productIDs := make([]uuid.UUID, 0)
	for rows.Next() {
		var p Product
		if err := rows.Scan(
			&p.ID, &p.Title, &p.Description, &p.CategoryID, &p.Brand, &p.Condition,
			&p.PriceCents, &p.SellerID, &p.InstagramPostURL, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		p.Images = make([]ProductImage, 0)
		products = append(products, p)
		productIDs = append(productIDs, p.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(productIDs) == 0 {
		return products, nil
	}

	imageMap, err := r.loadProductImagesBatch(ctx, productIDs)
	if err != nil {
		return nil, err
	}
	for i := range products {
		if imgs, ok := imageMap[products[i].ID]; ok {
			products[i].Images = imgs
		}
	}
	return products, nil
}

func (r *repository) GetProductByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	var p Product
	err := r.db.QueryRow(ctx,
		`SELECT id, title, description, category_id, brand, condition,
		        price_cents, seller_id, instagram_post_url, created_at, updated_at
		 FROM products WHERE id = $1`,
		id,
	).Scan(
		&p.ID, &p.Title, &p.Description, &p.CategoryID, &p.Brand, &p.Condition,
		&p.PriceCents, &p.SellerID, &p.InstagramPostURL, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	p.Images, err = r.loadProductImages(ctx, id)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) CreateProduct(ctx context.Context, sellerID *uuid.UUID, req CreateProductRequest) (*Product, error) {
	var p Product
	err := r.db.QueryRow(ctx,
		`INSERT INTO products (title, description, category_id, brand, condition, price_cents, instagram_post_url, seller_id)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, title, description, category_id, brand, condition,
		           price_cents, seller_id, instagram_post_url, created_at, updated_at`,
		req.Title, req.Description, req.CategoryID, req.Brand, req.Condition,
		req.PriceCents, req.InstagramPostURL, sellerID,
	).Scan(
		&p.ID, &p.Title, &p.Description, &p.CategoryID, &p.Brand, &p.Condition,
		&p.PriceCents, &p.SellerID, &p.InstagramPostURL, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	p.Images = make([]ProductImage, 0)
	return &p, nil
}

func (r *repository) UpdateProduct(ctx context.Context, id uuid.UUID, req UpdateProductRequest) (*Product, error) {
	// Dynamic update — only set non-nil fields
	var p Product
	err := r.db.QueryRow(ctx,
		`UPDATE products SET
		   title             = COALESCE($2, title),
		   description       = COALESCE($3, description),
		   price_cents       = COALESCE($4, price_cents),
		   instagram_post_url = COALESCE($5, instagram_post_url),
		   updated_at        = NOW()
		 WHERE id = $1
		 RETURNING id, title, description, category_id, brand, condition,
		           price_cents, seller_id, instagram_post_url, created_at, updated_at`,
		id, req.Title, req.Description, req.PriceCents, req.InstagramPostURL,
	).Scan(
		&p.ID, &p.Title, &p.Description, &p.CategoryID, &p.Brand, &p.Condition,
		&p.PriceCents, &p.SellerID, &p.InstagramPostURL, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	p.Images = make([]ProductImage, 0)
	return &p, nil
}

func (r *repository) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM products WHERE id = $1`, id)
	return err
}

func (r *repository) ListCategories(ctx context.Context) ([]Category, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, slug, parent_id FROM categories ORDER BY name`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]Category, 0)
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.ParentID); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, rows.Err()
}

func (r *repository) ListVariants(ctx context.Context, productID uuid.UUID) ([]Variant, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, product_id, size, color, sku, created_at
		 FROM variants WHERE product_id = $1 ORDER BY created_at`,
		productID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	variants := make([]Variant, 0)
	for rows.Next() {
		var v Variant
		if err := rows.Scan(&v.ID, &v.ProductID, &v.Size, &v.Color, &v.SKU, &v.CreatedAt); err != nil {
			return nil, err
		}
		variants = append(variants, v)
	}
	return variants, rows.Err()
}

func (r *repository) CreateVariant(ctx context.Context, productID uuid.UUID, req CreateVariantRequest) (*Variant, error) {
	var v Variant
	err := r.db.QueryRow(ctx,
		`INSERT INTO variants (product_id, size, color, sku)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, product_id, size, color, sku, created_at`,
		productID, req.Size, req.Color, req.SKU,
	).Scan(&v.ID, &v.ProductID, &v.Size, &v.Color, &v.SKU, &v.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *repository) DeleteVariant(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM variants WHERE id = $1`, id)
	return err
}

func (r *repository) SetStock(ctx context.Context, variantID uuid.UUID, quantity int) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO inventory (variant_id, quantity)
		 VALUES ($1, $2)
		 ON CONFLICT (variant_id) DO UPDATE SET quantity = $2`,
		variantID, quantity,
	)
	return err
}

func (r *repository) loadProductImagesBatch(ctx context.Context, productIDs []uuid.UUID) (map[uuid.UUID][]ProductImage, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, product_id, url, position
		 FROM product_images WHERE product_id = ANY($1) ORDER BY position`,
		productIDs,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	imageMap := make(map[uuid.UUID][]ProductImage)
	for rows.Next() {
		var img ProductImage
		if err := rows.Scan(&img.ID, &img.ProductID, &img.URL, &img.Position); err != nil {
			return nil, err
		}
		imageMap[img.ProductID] = append(imageMap[img.ProductID], img)
	}
	return imageMap, rows.Err()
}

func (r *repository) loadProductImages(ctx context.Context, productID uuid.UUID) ([]ProductImage, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, product_id, url, position
		 FROM product_images WHERE product_id = $1 ORDER BY position`,
		productID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	images := make([]ProductImage, 0)
	for rows.Next() {
		var img ProductImage
		if err := rows.Scan(&img.ID, &img.ProductID, &img.URL, &img.Position); err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, rows.Err()
}
