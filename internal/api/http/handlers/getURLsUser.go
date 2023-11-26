package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/models"
	"go.uber.org/zap"
)

// GetURLsUser обрабатывает HTTP-запрос на получение списка URL пользователя.
// Извлекает идентификатор пользователя (UserID) из контекста запроса.
// Если UserID отсутствует или пуст, возвращает HTTP-статус 401 Unauthorized.
// Вызывает сервис nodeService для получения списка URL пользователя.
// В случае ошибки при получении URL, возвращает HTTP-статус 500 Internal Server Error.
// Если список URL пуст, возвращает HTTP-статус 204 No Content.
// В случае успешного получения списка URL, возвращает HTTP-статус 200 OK
// и JSON-представление списка URL в теле ответа.
func (h *Handlers) GetURLsUser(w http.ResponseWriter, r *http.Request) {
	var UserID string
	var userIDCTX models.CtxString = "userID"

	// Извлекаем идентификатор пользователя (UserID) из контекста запроса
	UserIDVal := r.Context().Value(userIDCTX)
	if UserIDVal != nil {
		UserID = UserIDVal.(string)
	}

	// Если UserID отсутствует или пуст, логируем ошибку и возвращаем HTTP-статус 401 Unauthorized
	if UserID == "" {
		logger.Log.Info("trying to execute a method to retrieve a URL by a user by an unauthorized user")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Получаем список URL пользователя с использованием сервиса nodeService
	returnedURLs, err := h.nodeService.GetURLsUser(r.Context(), UserID)
	if err != nil {
		// В случае ошибки при получении URL, логируем ошибку и возвращаем HTTP-статус 500 Internal Server Error
		logger.Log.Error("when getting URLs on the user side", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Если список URL пуст, возвращаем HTTP-статус 204 No Content
	if len(returnedURLs) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Преобразуем список URL в JSON
	jsonText, err := json.Marshal(returnedURLs)
	if err != nil {
		// В случае ошибки создания JSON, логируем ошибку и возвращаем HTTP-статус 500 Internal Server Error
		logger.Log.Error("when creating a json file in the URL return procedure by user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Устанавливаем заголовок Content-Type
	w.Header().Add("Content-Type", "application/json")
	// Устанавливаем HTTP-статус 200 OK
	w.WriteHeader(http.StatusOK)
	// Отправляем JSON-представление списка URL в теле ответа
	_, err = w.Write(jsonText)
	if err != nil {
		// В случае ошибки записи ответа, логируем ошибку
		logger.Log.Error("response recording error", zap.Error(err))
	}
}
