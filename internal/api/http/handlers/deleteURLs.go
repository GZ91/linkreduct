package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/models"
	"go.uber.org/zap"
)

// DeleteURLs обрабатывает HTTP-запрос на удаление списка URL.
// Принимает JSON-представление массива строк (списка URL) из тела запроса.
// Если запрос не содержит данных (пустое тело) или происходит ошибка при чтении тела запроса,
// возвращает HTTP-статус 400 Bad Request. При успешном прочтении запроса, проверяет наличие
// элементов в списке URL и в случае его пустоты возвращает HTTP-статус 400 Bad Request.
// Извлекает идентификатор пользователя (UserID) из контекста запроса.
// Если UserID отсутствует или пуст, возвращает HTTP-статус 400 Bad Request.
// Запускает асинхронную операцию удаления URL с использованием сервиса nodeService.
// Возвращает HTTP-статус 202 Accepted.
func (h *Handlers) DeleteURLs(w http.ResponseWriter, r *http.Request) {
	var listURLs []string

	// Читаем тело запроса
	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		// В случае ошибки чтения тела запроса, логируем ошибку и возвращаем HTTP-статус 400 Bad Request
		logger.Log.Error("error when reading the request body", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Распаковываем JSON в массив строк (список URL)
	err = json.Unmarshal(bodyByte, &listURLs)
	if err != nil {
		// В случае ошибки разбора JSON, логируем ошибку и возвращаем HTTP-статус 400 Bad Request
		logger.Log.Error("error when reading the json conversion", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Проверяем наличие элементов в списке URL
	if len(listURLs) == 0 {
		// В случае пустого списка URL, логируем ошибку и возвращаем HTTP-статус 400 Bad Request
		logger.Log.Error("sent an empty list of links to be deleted")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var UserID string
	var userIDCTX models.CtxString = "userID"

	// Извлекаем идентификатор пользователя (UserID) из контекста запроса
	UserIDVal := r.Context().Value(userIDCTX)
	if UserIDVal != nil {
		UserID = UserIDVal.(string)
	}

	// Проверяем, что UserID не пуст
	if UserID == "" {
		// В случае отсутствия или пустоты UserID, логируем ошибку и возвращаем HTTP-статус 400 Bad Request
		logger.Log.Error("userID is not filled in")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Запускаем асинхронную операцию удаления URL с использованием сервиса nodeService
	go h.nodeService.DeletedLinks(listURLs, UserID)

	// Возвращаем HTTP-статус 202 Accepted
	w.WriteHeader(http.StatusAccepted)
}
