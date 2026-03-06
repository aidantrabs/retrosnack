package orders

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/MobinaToorani/retrosnack/pkg/httputil"
	"github.com/MobinaToorani/retrosnack/pkg/middleware"
)

type Handler struct {
	svc       Service
	jwtSecret string
}

func NewHandler(svc Service, jwtSecret string) *Handler {
	return &Handler{svc: svc, jwtSecret: jwtSecret}
}

func (h *Handler) Register(r chi.Router) {
	r.Post("/orders", h.createOrder)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(h.jwtSecret))
		r.Get("/orders/{id}", h.getOrder)
	})
}

func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if msg := validateCreateOrder(req); msg != "" {
		httputil.ErrorMsg(w, http.StatusBadRequest, msg)
		return
	}

	order, err := h.svc.CreateOrder(r.Context(), nil, req)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, order)
}

func (h *Handler) getOrder(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid order id")
		return
	}

	order, err := h.svc.GetOrder(r.Context(), id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, err)
		return
	}
	httputil.JSON(w, http.StatusOK, order)
}

func validateCreateOrder(req CreateOrderRequest) string {
	if len(req.Items) == 0 {
		return "order must have at least one item"
	}
	if len(req.Items) > 50 {
		return "order cannot have more than 50 items"
	}
	for i, item := range req.Items {
		if item.VariantID == uuid.Nil {
			return fmt.Sprintf("item %d: variant_id is required", i)
		}
		if item.Quantity <= 0 {
			return fmt.Sprintf("item %d: quantity must be greater than zero", i)
		}
		if item.PriceCents <= 0 {
			return fmt.Sprintf("item %d: price must be greater than zero", i)
		}
	}
	return ""
}
