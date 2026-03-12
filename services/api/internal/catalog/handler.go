package catalog

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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
	r.Get("/products", h.listProducts)
	r.Get("/products/{id}", h.getProduct)
	r.Get("/products/{id}/variants", h.listVariants)
	r.Get("/categories", h.listCategories)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(h.jwtSecret))
		r.Use(middleware.RequireRole("admin", "seller"))
		r.Post("/products", h.createProduct)
		r.Patch("/products/{id}", h.updateProduct)
		r.Delete("/products/{id}", h.deleteProduct)
		r.Post("/products/{id}/variants", h.createVariant)
		r.Delete("/variants/{id}", h.deleteVariant)
		r.Put("/variants/{id}/stock", h.setStock)
	})
}

func (h *Handler) listProducts(w http.ResponseWriter, r *http.Request) {
	limit := 20
	offset := 0

	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = n
		}
	}

	products, err := h.svc.ListProducts(r.Context(), limit, offset)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err)
		return
	}
	httputil.JSON(w, http.StatusOK, products)
}

func (h *Handler) getProduct(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid product id")
		return
	}

	product, err := h.svc.GetProduct(r.Context(), id)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, err)
		return
	}
	httputil.JSON(w, http.StatusOK, product)
}

func (h *Handler) createProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if msg := validateCreateProduct(req); msg != "" {
		httputil.ErrorMsg(w, http.StatusBadRequest, msg)
		return
	}

	var sellerID *uuid.UUID
	if claims, ok := middleware.ClaimsFromContext(r.Context()); ok {
		m := *claims
		if sub, _ := m["sub"].(string); sub != "" {
			if id, err := uuid.Parse(sub); err == nil {
				sellerID = &id
			}
		}
	}

	product, err := h.svc.CreateProduct(r.Context(), sellerID, req)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, product)
}

func (h *Handler) updateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid product id")
		return
	}

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if msg := validateUpdateProduct(req); msg != "" {
		httputil.ErrorMsg(w, http.StatusBadRequest, msg)
		return
	}

	product, err := h.svc.UpdateProduct(r.Context(), id, req)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err)
		return
	}
	httputil.JSON(w, http.StatusOK, product)
}

func (h *Handler) deleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid product id")
		return
	}

	if err := h.svc.DeleteProduct(r.Context(), id); err != nil {
		httputil.Error(w, http.StatusInternalServerError, err)
		return
	}
	httputil.NoContent(w)
}

func (h *Handler) listCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.svc.ListCategories(r.Context())
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err)
		return
	}
	httputil.JSON(w, http.StatusOK, categories)
}

func (h *Handler) listVariants(w http.ResponseWriter, r *http.Request) {
	productID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid product id")
		return
	}

	variants, err := h.svc.ListVariants(r.Context(), productID)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err)
		return
	}
	httputil.JSON(w, http.StatusOK, variants)
}

func (h *Handler) createVariant(w http.ResponseWriter, r *http.Request) {
	productID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid product id")
		return
	}

	var req CreateVariantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.SKU) == "" {
		httputil.ErrorMsg(w, http.StatusBadRequest, "sku is required")
		return
	}
	if len(req.SKU) > 100 {
		httputil.ErrorMsg(w, http.StatusBadRequest, "sku must be at most 100 characters")
		return
	}
	if len(req.Size) > 50 {
		httputil.ErrorMsg(w, http.StatusBadRequest, "size must be at most 50 characters")
		return
	}
	if len(req.Color) > 50 {
		httputil.ErrorMsg(w, http.StatusBadRequest, "color must be at most 50 characters")
		return
	}

	variant, err := h.svc.CreateVariant(r.Context(), productID, req)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err)
		return
	}
	httputil.JSON(w, http.StatusCreated, variant)
}

func (h *Handler) deleteVariant(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid variant id")
		return
	}

	if err := h.svc.DeleteVariant(r.Context(), id); err != nil {
		httputil.Error(w, http.StatusInternalServerError, err)
		return
	}
	httputil.NoContent(w)
}

func (h *Handler) setStock(w http.ResponseWriter, r *http.Request) {
	variantID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid variant id")
		return
	}

	var req SetStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Quantity < 0 {
		httputil.ErrorMsg(w, http.StatusBadRequest, "quantity cannot be negative")
		return
	}

	if err := h.svc.SetStock(r.Context(), variantID, req.Quantity); err != nil {
		httputil.Error(w, http.StatusInternalServerError, err)
		return
	}
	httputil.NoContent(w)
}

var validConditions = map[string]bool{"excellent": true, "good": true, "fair": true}

func validateCreateProduct(req CreateProductRequest) string {
	if strings.TrimSpace(req.Title) == "" {
		return "title is required"
	}
	if len(req.Title) > 200 {
		return "title must be at most 200 characters"
	}
	if len(req.Description) > 5000 {
		return "description must be at most 5000 characters"
	}
	if req.PriceCents <= 0 {
		return "price must be greater than zero"
	}
	if !validConditions[req.Condition] {
		return "condition must be excellent, good, or fair"
	}
	if len(req.Brand) > 100 {
		return "brand must be at most 100 characters"
	}
	if req.CategoryID == uuid.Nil {
		return "category_id is required"
	}
	return ""
}

func validateUpdateProduct(req UpdateProductRequest) string {
	if req.Title != nil && strings.TrimSpace(*req.Title) == "" {
		return "title cannot be empty"
	}
	if req.Title != nil && len(*req.Title) > 200 {
		return "title must be at most 200 characters"
	}
	if req.Description != nil && len(*req.Description) > 5000 {
		return "description must be at most 5000 characters"
	}
	if req.PriceCents != nil && *req.PriceCents <= 0 {
		return "price must be greater than zero"
	}
	return ""
}
