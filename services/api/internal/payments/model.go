package payments

import "github.com/google/uuid"

type CheckoutSession struct {
	ID      string    `json:"id"`
	OrderID uuid.UUID `json:"order_id"`
	URL     string    `json:"url"`
}

type CreateCheckoutRequest struct {
	OrderID uuid.UUID `json:"order_id"`
}
