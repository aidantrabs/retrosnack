package payments

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/retrosnack-clothing/retrosnack/pkg/httputil"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r chi.Router) {
	r.Post("/checkout", h.createCheckout)
	r.Post("/webhooks/stripe", h.stripeWebhook)
}

func (h *Handler) createCheckout(w http.ResponseWriter, r *http.Request) {
	var req CreateCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid request body")
		return
	}

	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = "https://retrosnack.shop"
	}

	sess, err := h.svc.CreateCheckout(
		r.Context(),
		req,
		origin+"/orders/"+req.OrderID.String()+"/success",
		origin+"/orders/"+req.OrderID.String()+"/cancel",
	)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err)
		return
	}
	httputil.JSON(w, http.StatusOK, sess)
}

func (h *Handler) stripeWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(io.LimitReader(r.Body, 65536))
	if err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "failed to read body")
		return
	}

	signature := r.Header.Get("Stripe-Signature")
	if err := h.svc.HandleWebhook(r.Context(), payload, signature); err != nil {
		httputil.Error(w, http.StatusBadRequest, err)
		return
	}

	httputil.NoContent(w)
}
