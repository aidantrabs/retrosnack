package instagram

import (
	"time"

	"github.com/google/uuid"
)

type Link struct {
	ID        uuid.UUID  `json:"id"`
	ProductID uuid.UUID  `json:"product_id"`
	PostURL   string     `json:"post_url"`
	EmbedHTML string     `json:"embed_html,omitempty"`
	CachedAt  *time.Time `json:"cached_at,omitempty"`
}

type oEmbedResponse struct {
	HTML string `json:"html"`
}
