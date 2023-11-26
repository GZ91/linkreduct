package handlers

import (
	"net/http"

	"github.com/GZ91/linkreduct/internal/errorsapp"
	"github.com/go-chi/chi/v5"
)

// GetLongURL обрабатывает HTTP-запрос на получение длинного URL по идентификатору.
// Принимает идентификатор URL из параметра маршрута.
// Если успешно получает длинный URL с использованием сервиса nodeService,
// выполняет редирект с HTTP-статусом 307 Temporary Redirect и заголовком Location,
// содержащим длинный URL. В случае, если URL был удален, возвращает HTTP-статус 410 Gone.
// В случае других ошибок возвращает HTTP-статус 400 Bad Request.
// Если URL не найден, также возвращает HTTP-статус 400 Bad Request.
func (h *Handlers) GetLongURL(w http.ResponseWriter, r *http.Request) {
	// Извлекаем идентификатор URL из параметра маршрута
	id := chi.URLParam(r, "id")

	// Получаем длинный URL с использованием сервиса nodeService
	link, ok, err := h.nodeService.GetURL(r.Context(), id)

	// Проверяем наличие ошибок
	if err != nil && err != errorsapp.ErrLineURLDeleted {
		// В случае других ошибок возвращаем HTTP-статус 400 Bad Request
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// В случае, если URL был удален, возвращаем HTTP-статус 410 Gone
	if err == errorsapp.ErrLineURLDeleted {
		w.WriteHeader(http.StatusGone)
		return
	}

	// Если URL найден, выполняем редирект
	if ok {
		// Устанавливаем заголовок Location с длинным URL
		w.Header().Add("Location", link)
		// Возвращаем HTTP-статус 307 Temporary Redirect
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		// В случае, если URL не найден, возвращаем HTTP-статус 400 Bad Request
		w.WriteHeader(http.StatusBadRequest)
	}
}
