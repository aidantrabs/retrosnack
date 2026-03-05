package media

import "github.com/google/uuid"

type Upload struct {
	Key      string    `json:"key"`
	URL      string    `json:"url"`
	MimeType string    `json:"mime_type"`
	Size     int64     `json:"size"`
}

type ProductImageRecord struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"product_id"`
	URL       string    `json:"url"`
	Position  int       `json:"position"`
}
