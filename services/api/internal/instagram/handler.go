package instagram

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
	r.Get("/products/{productId}/instagram", h.getEmbed)
	r.Put("/products/{productId}/instagram", h.refreshEmbed)
}

func (h *Handler) getEmbed(w http.ResponseWriter, r *http.Request) {
	productID, err := uuid.Parse(chi.URLParam(r, "productId"))
	if err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid product id")
		return
	}

	link, err := h.svc.GetEmbed(r.Context(), productID)
	if err != nil {
		httputil.Error(w, http.StatusNotFound, err)
		return
	}
	httputil.JSON(w, http.StatusOK, link)
}

func (h *Handler) refreshEmbed(w http.ResponseWriter, r *http.Request) {
	productID, err := uuid.Parse(chi.URLParam(r, "productId"))
	if err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid product id")
		return
	}

	var body struct {
		PostURL string `json:"post_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httputil.ErrorMsg(w, http.StatusBadRequest, "invalid request body")
		return
	}

	link, err := h.svc.RefreshEmbed(r.Context(), productID, body.PostURL)
	if err != nil {
		httputil.Error(w, http.StatusInternalServerError, err)
		return
	}
	httputil.JSON(w, http.StatusOK, link)
}
