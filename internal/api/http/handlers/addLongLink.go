package handlers

import (
	"errors"
	"io"
	"net/http"

	"github.com/GZ91/linkreduct/internal/app/logger"
	"github.com/GZ91/linkreduct/internal/errorsapp"
	"go.uber.org/zap"
)

// AddLongLink обрабатывает HTTP-запрос на добавление длинного URL.
// Принимает URL в виде текстового тела запроса и проверяет его на соответствие
// фильтру URLFilter. Если URL не соответствует фильтру, возвращает HTTP-статус 400 Bad Request.
// В случае успешного добавления возвращает HTTP-статус 201 Created и текстовое представление
// сокращенной версии URL. В случае конфликта (URL уже существует),
// возвращает HTTP-статус 409 Conflict. При возникновении других ошибок возвращает
// HTTP-статус 400 Bad Request с соответствующим описанием ошибки в теле ответа.
// Если запрос не содержит данных (пустое тело), также возвращает HTTP-статус 400 Bad Request.
func (h *Handlers) AddLongLink(w http.ResponseWriter, r *http.Request) {
	StatusReturn := http.StatusCreated

	// Читаем тело запроса
	link, err := io.ReadAll(r.Body)
	if err != nil {
		// В случае ошибки чтения тела запроса, возвращаем HTTP-статус 400 Bad Request
		// и в теле ответа указываем причину ошибки
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Проверяем URL на соответствие фильтру
	if !h.URLFilter.MatchString(string(link)) {
		// В случае, если URL не соответствует фильтру, возвращаем HTTP-статус 400 Bad Request
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Получаем сокращенную версию URL с использованием сервиса nodeService
	bodyText, err := h.nodeService.GetSmallLink(r.Context(), string(link))
	if err != nil {
		// В случае ошибки при получении сокращенной версии URL
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

	// Если сокращенная версия URL пуста, возвращаем HTTP-статус 400 Bad Request
	if bodyText == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Устанавливаем заголовок Content-Type
	w.Header().Add("Content-Type", "text/plain")
	// Устанавливаем HTTP-статус
	w.WriteHeader(StatusReturn)
	// Отправляем текстовое представление сокращенной версии URL в теле ответа
	_, err = w.Write([]byte(bodyText))
	if err != nil {
		// В случае ошибки записи ответа, логируем ошибку
		logger.Log.Error("response recording error", zap.String("error", err.Error()))
	}
}
