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

type ProcessPaymentRequest struct {
	OrderID    uuid.UUID `json:"order_id"`
	SourceID   string    `json:"source_id"` // token from square web payments sdk
}

type PaymentResult struct {
	OrderID   uuid.UUID `json:"order_id"`
	PaymentID string    `json:"payment_id"`
	Status    string    `json:"status"`
}
