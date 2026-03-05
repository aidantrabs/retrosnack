package inventory

import (
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
	r.Get("/inventory/{variantId}", h.getStock)
}

func (h *Handler) getStock(w http.ResponseWriter, r *http.Request) {
	variantID, err := uuid.Parse(chi.URLParam(r, "variantId"))
	if err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid variant id")
		return
	}

	stock, err := h.svc.GetStock(r.Context(), variantID)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, err)
		return
	}
	httputil.JSON(w, http.StatusOK, stock)
}
