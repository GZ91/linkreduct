package handlers

import (
	"encoding/json"
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

func (h *handlers) AddLongLinkJSON(w http.ResponseWriter, r *http.Request) {
	textBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	type requestData struct {
		URL string `json:"url"`
	}
	var data requestData
	json.Unmarshal(textBody, &data)
	link := data.URL

	if link == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bodyText := h.nodeService.GetSmallLink(link)
	if bodyText == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	type result struct {
		Result string `json:"result"`
	}

	Result := result{Result: bodyText}

	res, err := json.Marshal(Result)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if len(res) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)
}
