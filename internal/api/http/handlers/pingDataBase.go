package handlers

import "net/http"

// PingDataBase обрабатывает HTTP-запрос на проверку доступности базы данных.
// Вызывает метод Ping узер-сервиса для проверки соединения с базой данных.
// В случае успешного соединения возвращает HTTP-статус 200 OK,
// в противном случае возвращает HTTP-статус 500 Internal Server Error.
func (h *Handlers) PingDataBase(w http.ResponseWriter, r *http.Request) {
	// Вызываем метод Ping узер-сервиса для проверки доступности базы данных
	err := h.nodeService.Ping(r.Context())
	if err != nil {
		// В случае ошибки возвращаем HTTP-статус 500 Internal Server Error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// В случае успешного соединения возвращаем HTTP-статус 200 OK
	w.WriteHeader(http.StatusOK)
}
