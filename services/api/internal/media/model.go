package media

import "github.com/google/uuid"

type ProductImageRecord struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	R2Key     string    `json:"-"`
	URL       string    `json:"url"`
	Position  int       `json:"position"`
}
