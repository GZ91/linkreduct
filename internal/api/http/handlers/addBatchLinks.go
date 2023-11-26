package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/errorsapp"
	"github.com/GZ91/linkreduct/internal/models"
	"go.uber.org/zap"
)

// AddBatchLinks обрабатывает HTTP-запрос на добавление пакета URL.
// Принимает JSON-представление массива моделей IncomingBatchURL из тела запроса.
// В случае успешного добавления возвращает HTTP-статус 201 Created и JSON-представление
// массива освобожденных (зарезервированных) URL. В случае конфликта (URL уже существует),
// возвращает HTTP-статус 409 Conflict. При возникновении других ошибок возвращает
// HTTP-статус 400 Bad Request с соответствующим описанием ошибки в теле ответа.
// Если запрос не содержит данных (пустое тело), также возвращает HTTP-статус 400 Bad Request.
func (h *Handlers) AddBatchLinks(w http.ResponseWriter, r *http.Request) {
	StatusReturn := http.StatusCreated

	// Читаем тело запроса
	textBody, err := io.ReadAll(r.Body)
	if err != nil {
		// В случае ошибки чтения тела запроса, возвращаем HTTP-статус 400 Bad Request
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Распаковываем JSON в массив IncomingBatchURL
	var incomingBatchURL []models.IncomingBatchURL
	err = json.Unmarshal(textBody, &incomingBatchURL)
	if err != nil {
		// В случае ошибки разбора JSON, возвращаем HTTP-статус 400 Bad Request
		// и в теле ответа указываем причину ошибки
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Добавляем пакет URL с использованием сервиса nodeService
	releasedBatchURL, err := h.nodeService.AddBatchLink(r.Context(), incomingBatchURL)
	if err != nil {
		// В случае ошибки при добавлении URL
		if errors.Is(err, errorsapp.ErrLinkAlreadyExists) {
			// Если URL уже существует, возвращаем HTTP-статус 409 Conflict
			StatusReturn = http.StatusConflict
		} else {
			// В прочих случаях возвращаем HTTP-статус 400 Bad Request
			// и в теле ответа указываем причину ошибки
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}

	// Преобразуем освобожденные URL в JSON
	res, err := json.Marshal(releasedBatchURL)
	if err != nil {
		// В случае ошибки преобразования в JSON, возвращаем HTTP-статус 400 Bad Request
		// и в теле ответа указываем причину ошибки
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Если JSON-представление пусто, возвращаем HTTP-статус 400 Bad Request
	if len(res) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Устанавливаем заголовок Content-Type
	w.Header().Add("Content-Type", "application/json")
	// Устанавливаем HTTP-статус
	w.WriteHeader(StatusReturn)
	// Отправляем JSON-представление освобожденных URL в теле ответа
	_, err = w.Write(res)
	if err != nil {
		// В случае ошибки записи ответа, логируем ошибку
		logger.Log.Error("response recording error", zap.String("error", err.Error()))
	}
}
