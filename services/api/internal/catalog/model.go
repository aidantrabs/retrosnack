package catalog

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID               uuid.UUID      `json:"id"`
	Title            string         `json:"title"`
	Description      string         `json:"description"`
	CategoryID       uuid.UUID      `json:"category_id"`
	Brand            string         `json:"brand"`
	Condition        string         `json:"condition"` // "excellent" | "good" | "fair"
	PriceCents       int64          `json:"price_cents"`
	SellerID         *uuid.UUID     `json:"seller_id,omitempty"`
	InstagramPostURL string         `json:"instagram_post_url"`
	Images           []ProductImage `json:"images"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

type ProductImage struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	R2Key     string    `json:"-"`
	URL       string    `json:"url"`
	Position  int       `json:"position"`
}

type Variant struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	Size      string    `json:"size"`
	Color     string    `json:"color"`
	SKU       string    `json:"sku"`
	CreatedAt time.Time `json:"created_at"`
}

type Category struct {
	ID       uuid.UUID  `json:"id"`
	Name     string     `json:"name"`
	Slug     string     `json:"slug"`
	ParentID *uuid.UUID `json:"parent_id,omitempty"`
}

type CreateProductRequest struct {
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	CategoryID       uuid.UUID `json:"category_id"`
	Brand            string    `json:"brand"`
	Condition        string    `json:"condition"`
	PriceCents       int64     `json:"price_cents"`
	InstagramPostURL string    `json:"instagram_post_url"`
}

type UpdateProductRequest struct {
	Title            *string `json:"title,omitempty"`
	Description      *string `json:"description,omitempty"`
	PriceCents       *int64  `json:"price_cents,omitempty"`
	InstagramPostURL *string `json:"instagram_post_url,omitempty"`
}

type CreateVariantRequest struct {
	Size  string `json:"size"`
	Color string `json:"color"`
	SKU   string `json:"sku"`
}

type SetStockRequest struct {
	Quantity int `json:"quantity"`
}
