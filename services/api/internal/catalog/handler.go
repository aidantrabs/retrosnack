package catalog

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/MobinaToorani/retrosnack/pkg/httputil"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(r chi.Router) {
	r.Get("/products", h.listProducts)
	r.Get("/products/{id}", h.getProduct)
	r.Post("/products", h.createProduct)
	r.Patch("/products/{id}", h.updateProduct)
	r.Delete("/products/{id}", h.deleteProduct)
	r.Get("/categories", h.listCategories)
}

func (h *Handler) listProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.svc.ListProducts(r.Context())
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

	product, err := h.svc.CreateProduct(r.Context(), req)
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
