package payments

import (
	"time"

	"github.com/google/uuid"
)

type CheckoutSession struct {
	ID        string    `json:"id"`
	OrderID   uuid.UUID `json:"order_id"`
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
}

type CreateCheckoutRequest struct {
	OrderID uuid.UUID `json:"order_id"`
}
