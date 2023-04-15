package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

type handlerserService interface {
	GetSmallLink(string) string
	GetURL(string) (string, bool)
}

type handlers struct {
	nodeService handlerserService
}

func New(nodeService handlerserService) *handlers {
	return &handlers{nodeService: nodeService}
}

func (h *handlers) AddLongLink(w http.ResponseWriter, r *http.Request) {
	link, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if string(link) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bodyText := h.nodeService.GetSmallLink(string(link))
	if bodyText == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "text/plain")
	w.Header().Add("Content-Length", fmt.Sprint(len(bodyText)))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(bodyText))
}

func (h *handlers) GetShortURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if link, ok := h.nodeService.GetURL(id); ok {
		w.Header().Add("Location", link)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
