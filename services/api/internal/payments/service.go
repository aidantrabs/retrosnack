package payments

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	square "github.com/square/square-go-sdk"
	squareclient "github.com/square/square-go-sdk/client"
	"github.com/square/square-go-sdk/checkout"
	"github.com/square/square-go-sdk/option"
	"github.com/retrosnack-clothing/retrosnack/internal/orders"
)

type Service interface {
	CreateCheckout(ctx context.Context, req CreateCheckoutRequest, redirectURL string) (*CheckoutSession, error)
	HandleWebhook(ctx context.Context, payload []byte, signatureHeader string) error
}

type service struct {
	orders          orders.Service
	square          *squareclient.Client
	locationID      string
	webhookSigKey   string
	webhookNotifURL string
}

func NewService(ordersSvc orders.Service, accessToken, locationID, webhookSigKey, webhookNotifURL string) Service {
	c := squareclient.NewClient(
		option.WithToken(accessToken),
	)
	return &service{
		orders:          ordersSvc,
		square:          c,
		locationID:      locationID,
		webhookSigKey:   webhookSigKey,
		webhookNotifURL: webhookNotifURL,
	}
}

func (s *service) CreateCheckout(ctx context.Context, req CreateCheckoutRequest, redirectURL string) (*CheckoutSession, error) {
	order, err := s.orders.GetOrder(ctx, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("order not found: %w", err)
	}

	lineItems := make([]*square.OrderLineItem, 0, len(order.Items))
	for _, item := range order.Items {
		lineItems = append(lineItems, &square.OrderLineItem{
			Name:     square.String(item.VariantID.String()),
			Quantity: strconv.Itoa(item.Quantity),
			BasePriceMoney: &square.Money{
				Amount:   square.Int64(item.PriceCents),
				Currency: square.CurrencyCad.Ptr(),
			},
			Note: square.String("order:" + order.ID.String()),
		})
	}

	idempotencyKey := order.ID.String()

	resp, err := s.square.Checkout.PaymentLinks.Create(ctx, &checkout.CreatePaymentLinkRequest{
		IdempotencyKey: &idempotencyKey,
		Order: &square.Order{
			LocationID:  s.locationID,
			ReferenceID: square.String(order.ID.String()),
			LineItems:   lineItems,
		},
		CheckoutOptions: &square.CheckoutOptions{
			RedirectURL: &redirectURL,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create payment link: %w", err)
	}

	link := resp.PaymentLink
	if link == nil || link.ID == nil {
		return nil, fmt.Errorf("square returned empty payment link")
	}

	if err := s.orders.SetCheckoutSession(ctx, order.ID, *link.ID); err != nil {
		return nil, err
	}

	url := ""
	if link.URL != nil {
		url = *link.URL
	}

	return &CheckoutSession{
		ID:      *link.ID,
		OrderID: order.ID,
		URL:     url,
	}, nil
}

func (s *service) HandleWebhook(ctx context.Context, payload []byte, signatureHeader string) error {
	err := s.square.Webhooks.VerifySignature(ctx, &square.VerifySignatureRequest{
		RequestBody:     string(payload),
		SignatureHeader: signatureHeader,
		SignatureKey:    s.webhookSigKey,
		NotificationURL: s.webhookNotifURL,
	})
	if err != nil {
		return fmt.Errorf("webhook signature verification failed: %w", err)
	}

	var event struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to parse webhook event: %w", err)
	}

	if event.Type != "payment.updated" {
		return nil
	}

	var data struct {
		Object struct {
			Payment struct {
				Status  string `json:"status"`
				OrderID string `json:"order_id"`
			} `json:"payment"`
		} `json:"object"`
	}
	if err := json.Unmarshal(event.Data, &data); err != nil {
		return fmt.Errorf("failed to parse payment data: %w", err)
	}

	if data.Object.Payment.Status != "COMPLETED" {
		return nil
	}

	squareOrderID := data.Object.Payment.OrderID
	if squareOrderID == "" {
		return nil
	}

	// look up the square order to get our reference_id (our order UUID)
	squareOrder, err := s.square.Orders.Get(ctx, &square.GetOrdersRequest{OrderID: squareOrderID})
	if err != nil {
		return fmt.Errorf("failed to get square order: %w", err)
	}

	if squareOrder.Order == nil || squareOrder.Order.ReferenceID == nil {
		return fmt.Errorf("square order missing reference_id")
	}

	orderID, err := uuid.Parse(*squareOrder.Order.ReferenceID)
	if err != nil {
		return fmt.Errorf("invalid reference_id in square order: %w", err)
	}

	return s.orders.MarkPaid(ctx, orderID)
}
