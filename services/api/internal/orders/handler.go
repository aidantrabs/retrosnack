package orders

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/retrosnack-clothing/retrosnack/pkg/httputil"
	"github.com/retrosnack-clothing/retrosnack/pkg/middleware"
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

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(h.jwtSecret))
		r.Use(middleware.RequireRole("admin"))
		r.Post("/orders/{id}/ship", h.shipOrder)
		r.Post("/orders/{id}/deliver", h.deliverOrder)
		r.Post("/orders/{id}/cancel", h.cancelOrder)
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

	// extract user id from jwt if present (guest checkout passes nil)
	var userID *uuid.UUID
	if claims, ok := middleware.ClaimsFromContext(r.Context()); ok {
		m := (*claims)
		if sub, _ := m["sub"].(string); sub != "" {
			if uid, err := uuid.Parse(sub); err == nil {
				userID = &uid
			}
		}
	}

	order, err := h.svc.CreateOrder(r.Context(), userID, req)
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

	// ownership check: admins can see any order, users only their own
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		httputil.ErrorMsg(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	m := (*claims)
	role, _ := m["role"].(string)
	if role != "admin" {
		sub, _ := m["sub"].(string)
		callerID, _ := uuid.Parse(sub)
		if order.UserID == nil || *order.UserID != callerID {
			httputil.ErrorMsg(w, http.StatusForbidden, "forbidden")
			return
		}
	}

	httputil.JSON(w, http.StatusOK, order)
}

func (h *Handler) shipOrder(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid order id")
		return
	}

	if err := h.svc.MarkShipped(r.Context(), id); err != nil {
		if errors.Is(err, ErrInvalidTransition) {
			httputil.ErrorMsg(w, http.StatusConflict, "order must be paid before shipping")
			return
		}
		httputil.Error(w, http.StatusInternalServerError, err)
		return
	}
	httputil.NoContent(w)
}

func (h *Handler) deliverOrder(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid order id")
		return
	}

	if err := h.svc.MarkDelivered(r.Context(), id); err != nil {
		if errors.Is(err, ErrInvalidTransition) {
			httputil.ErrorMsg(w, http.StatusConflict, "order must be shipped before delivering")
			return
		}
		httputil.Error(w, http.StatusInternalServerError, err)
		return
	}
	httputil.NoContent(w)
}

func (h *Handler) cancelOrder(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid order id")
		return
	}

	if err := h.svc.CancelOrder(r.Context(), id); err != nil {
		if errors.Is(err, ErrInvalidTransition) {
			httputil.ErrorMsg(w, http.StatusConflict, "only pending orders can be cancelled")
			return
		}
		httputil.Error(w, http.StatusInternalServerError, err)
		return
	}
	httputil.NoContent(w)
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
