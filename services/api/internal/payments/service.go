package payments

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/webhook"
	"github.com/MobinaToorani/retrosnack/internal/orders"
)

type Service interface {
	CreateCheckout(ctx context.Context, req CreateCheckoutRequest, successURL, cancelURL string) (*CheckoutSession, error)
	HandleWebhook(ctx context.Context, payload []byte, signature string) error
}

type service struct {
	orders        orders.Service
	webhookSecret string
}

func NewService(ordersSvc orders.Service, secretKey, webhookSecret string) Service {
	stripe.Key = secretKey
	return &service{
		orders:        ordersSvc,
		webhookSecret: webhookSecret,
	}
}

func (s *service) CreateCheckout(ctx context.Context, req CreateCheckoutRequest, successURL, cancelURL string) (*CheckoutSession, error) {
	order, err := s.orders.GetOrder(ctx, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	lineItems := make([]*stripe.CheckoutSessionLineItemParams, 0, len(order.Items))
	for _, item := range order.Items {
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String("cad"),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String(item.VariantID.String()),
				},
				UnitAmount: stripe.Int64(item.PriceCents),
			},
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	params := &stripe.CheckoutSessionParams{
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems:  lineItems,
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		Metadata: map[string]string{
			"order_id": order.ID.String(),
		},
	}

	sess, err := session.New(params)
	if err != nil {
		return nil, err
	}

	if err := s.orders.SetStripeSession(ctx, order.ID, sess.ID); err != nil {
		return nil, err
	}

	return &CheckoutSession{
		ID:        sess.ID,
		OrderID:   order.ID,
		URL:       sess.URL,
		ExpiresAt: time.Unix(sess.ExpiresAt, 0),
	}, nil
}

func (s *service) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	event, err := webhook.ConstructEvent(payload, signature, s.webhookSecret)
	if err != nil {
		return fmt.Errorf("webhook signature verification failed: %w", err)
	}

	switch event.Type {
	case "checkout.session.completed":
		var sess stripe.CheckoutSession
		if err := event.DataObject(&sess); err != nil {
			return err
		}

		orderIDStr, ok := sess.Metadata["order_id"]
		if !ok {
			return fmt.Errorf("missing order_id in session metadata")
		}

		orderID, err := uuid.Parse(orderIDStr)
		if err != nil {
			return fmt.Errorf("invalid order_id in session metadata: %w", err)
		}

		return s.orders.MarkPaid(ctx, orderID)
	}

	return nil
}
