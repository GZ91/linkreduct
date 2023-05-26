package handlers

import (
	"context"
	"encoding/json"
	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/errorsapp"
	"github.com/GZ91/linkreduct/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"net/http"
	"regexp"
)

type handlerserService interface {
	GetSmallLink(context.Context, string) (string, error)
	GetURL(context.Context, string) (string, bool, error)
	GetURLsUser(context.Context, string) ([]models.ReturnedStructURL, error)
	Ping(ctx context.Context) error
	AddBatchLink(context.Context, []models.IncomingBatchURL) ([]models.ReleasedBatchURL, error)
}

type handlers struct {
	nodeService handlerserService
	URLFilter   *regexp.Regexp
}

func New(nodeService handlerserService) *handlers {
	return &handlers{nodeService: nodeService, URLFilter: regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?(\w+\.[^:\/\n]+)`)}
}

func (h *handlers) AddLongLink(w http.ResponseWriter, r *http.Request) {
	StatusReturn := http.StatusCreated

	link, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if !h.URLFilter.MatchString(string(link)) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	bodyText, err := h.nodeService.GetSmallLink(r.Context(), string(link))
	if err != nil {
		if errors.Is(err, errorsapp.ErrLinkAlreadyExists) {
			StatusReturn = http.StatusConflict
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}

	if bodyText == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(StatusReturn)
	_, err = w.Write([]byte(bodyText))
	if err != nil {
		logger.Log.Error("response recording error", zap.String("error", err.Error()))
	}
}

func (h *handlers) GetLongURL(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	link, ok, err := h.nodeService.GetURL(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if ok {
		w.Header().Add("Location", link)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *handlers) AddLongLinkJSON(w http.ResponseWriter, r *http.Request) {
	StatusReturn := http.StatusCreated
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

	bodyText, err := h.nodeService.GetSmallLink(r.Context(), link)
	if err != nil {
		if errors.Is(err, errorsapp.ErrLinkAlreadyExists) {
			StatusReturn = http.StatusConflict
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
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
	w.WriteHeader(StatusReturn)
	_, err = w.Write(res)
	if err != nil {
		logger.Log.Error("response recording error", zap.String("error", err.Error()))
	}
}

func (h *handlers) PingDataBase(w http.ResponseWriter, r *http.Request) {
	err := h.nodeService.Ping(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handlers) AddBatchLinks(w http.ResponseWriter, r *http.Request) {
	StatusReturn := http.StatusCreated
	textBody, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var incomingBatchURL []models.IncomingBatchURL

	err = json.Unmarshal(textBody, &incomingBatchURL)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	releasedBatchURL, err := h.nodeService.AddBatchLink(r.Context(), incomingBatchURL)
	if err != nil {
		if errors.Is(err, errorsapp.ErrLinkAlreadyExists) {
			StatusReturn = http.StatusConflict
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
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
	w.WriteHeader(StatusReturn)
	_, err = w.Write(res)
	if err != nil {
		logger.Log.Error("response recording error", zap.String("error", err.Error()))
	}
}

func (h *handlers) GetURLsUser(w http.ResponseWriter, r *http.Request) {
	var UserID string
	UserIDVal := r.Context().Value("userID")
	if UserIDVal != nil {
		UserID = UserIDVal.(string)
	} else {
		logger.Log.Info("trying to execute a method to retrieve a URL by a user by an unauthorized user")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	returnedURLs, err := h.nodeService.GetURLsUser(r.Context(), UserID)
	if err != nil {
		logger.Log.Error("when getting URLs on the user side", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(returnedURLs) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	jsonText, err := json.Marshal(returnedURLs)
	if err != nil {
		logger.Log.Error("when creating a json file in the URL return procedure by user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonText)
	if err != nil {
		logger.Log.Error("response recording error", zap.Error(err))
	}

}
