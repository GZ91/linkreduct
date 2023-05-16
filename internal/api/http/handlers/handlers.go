package handlers

import (
	"encoding/json"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/models"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"regexp"
)

type handlerserService interface {
	GetSmallLink(string) string
	GetURL(string) (string, bool)
	Ping() error
}

type handlers struct {
	nodeService handlerserService
	URLFilter   *regexp.Regexp
}

func New(nodeService handlerserService) *handlers {
	return &handlers{nodeService: nodeService, URLFilter: regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?(\w+\.[^:\/\n]+)`)}
}

func (h *handlers) AddLongLink(w http.ResponseWriter, r *http.Request) {
	link, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !h.URLFilter.MatchString(string(link)) {
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
	_, err = w.Write([]byte(bodyText))
	if err != nil {
		logger.Log.Error("response recording error", zap.String("error", err.Error()))
	}
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

	var data models.RequestData

	err = json.Unmarshal(textBody, &data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	link := data.URL

	if !h.URLFilter.MatchString(link) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bodyText := h.nodeService.GetSmallLink(link)
	if bodyText == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Result := models.ResultReturn{Result: bodyText}

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
	_, err = w.Write(res)
	if err != nil {
		logger.Log.Error("response recording error", zap.String("error", err.Error()))
	}
}

func (h *handlers) PingDataBase(w http.ResponseWriter, r *http.Request) {
	err := h.nodeService.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handlers) AddListLongLinkJSON(w http.ResponseWriter, r *http.Request) {
	textBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var incomingBatchURL []models.IncomingBatchURL
	var releasedBatchURL []models.ReleasedBatchURL

	err = json.Unmarshal(textBody, &incomingBatchURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	for _, data := range incomingBatchURL {
		link := data.OriginalURL

		if !h.URLFilter.MatchString(link) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		shortURL := h.nodeService.GetSmallLink(link)
		if shortURL == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		releasedBatchURL = append(releasedBatchURL, models.ReleasedBatchURL{CorrelationId: data.CorrelationId, ShortURL: shortURL})
	}

	res, err := json.Marshal(releasedBatchURL)
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
	_, err = w.Write(res)
	if err != nil {
		logger.Log.Error("response recording error", zap.String("error", err.Error()))
	}
}
