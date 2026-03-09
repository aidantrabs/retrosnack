package media

import "github.com/google/uuid"

type ProductImageRecord struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	URL       string    `json:"url"`
	Position  int       `json:"position"`
}
