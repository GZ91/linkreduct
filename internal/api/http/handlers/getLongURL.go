package handlers

import (
	"github.com/GZ91/linkreduct/internal/errorsapp"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (h *handlers) GetLongURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	link, ok, err := h.nodeService.GetURL(r.Context(), id)
	if err != nil && err != errorsapp.ErrLineURLDeleted {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err == errorsapp.ErrLineURLDeleted {
		w.WriteHeader(http.StatusGone)
		return
	}
	if ok {
		w.Header().Add("Location", link)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
