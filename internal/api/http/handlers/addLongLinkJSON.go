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

// AddLongLinkJSON обрабатывает HTTP-запрос на добавление длинного URL, принимаемого в формате JSON.
// Принимает JSON-представление модели RequestData из тела запроса, содержащей поле URL.
// Проверяет URL на соответствие фильтру URLFilter. Если URL не соответствует фильтру,
// возвращает HTTP-статус 400 Bad Request. В случае успешного добавления возвращает HTTP-статус 201 Created
// и JSON-представление модели ResultReturn, содержащей сокращенную версию URL. В случае конфликта
// (URL уже существует), возвращает HTTP-статус 409 Conflict. При возникновении других ошибок возвращает
// HTTP-статус 400 Bad Request с соответствующим описанием ошибки в теле ответа.
// Если запрос не содержит данных (пустое тело), также возвращает HTTP-статус 400 Bad Request.
func (h *Handlers) AddLongLinkJSON(w http.ResponseWriter, r *http.Request) {
	StatusReturn := http.StatusCreated

	// Читаем тело запроса
	textBody, err := io.ReadAll(r.Body)
	if err != nil {
		// В случае ошибки чтения тела запроса, возвращаем HTTP-статус 400 Bad Request
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Распаковываем JSON в модель RequestData
	var data models.RequestData
	err = json.Unmarshal(textBody, &data)
	if err != nil {
		// В случае ошибки разбора JSON, возвращаем HTTP-статус 400 Bad Request
		// и в теле ответа указываем причину ошибки
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Извлекаем URL из модели RequestData
	link := data.URL

	// Проверяем URL на соответствие фильтру
	if !h.URLFilter.MatchString(link) {
		// В случае, если URL не соответствует фильтру, возвращаем HTTP-статус 400 Bad Request
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Получаем сокращенную версию URL с использованием сервиса nodeService
	bodyText, err := h.nodeService.GetSmallLink(r.Context(), link)
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

	// Формируем JSON-представление модели ResultReturn
	Result := models.ResultReturn{Result: bodyText}
	res, err := json.Marshal(Result)
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
	// Отправляем JSON-представление модели ResultReturn в теле ответа
	_, err = w.Write(res)
	if err != nil {
		// В случае ошибки записи ответа, логируем ошибку
		logger.Log.Error("response recording error", zap.String("error", err.Error()))
	}
}
