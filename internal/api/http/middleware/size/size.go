package sizemiddleware

import (
	"net/http"
	"strconv"
)

// Size представляет собой обертку для http.ResponseWriter с возможностью расчета размера ответа.
type Size struct {
	http.ResponseWriter
}

// Write перенаправляет вызов метода Write на внутренний ResponseWriter.
// Записывает данные ответа и обновляет заголовок Content-Length с размером ответа.
// Возвращает количество записанных байт и ошибку (если есть).
func (s Size) Write(p []byte) (int, error) {
	// Вызываем метод Write у внутреннего ResponseWriter
	len, err := s.ResponseWriter.Write(p)
	// Обновляем заголовок Content-Length с размером ответа
	s.ResponseWriter.Header().Add("Content-Length", strconv.Itoa(len))
	return len, err
}

// CalculateSize возвращает обработчик HTTP, добавляющий расчет размера ответа.
// Создает объект Size для обертки ResponseWriter с возможностью расчета размера.
func CalculateSize(h http.Handler) http.Handler {
	// CalSize - функция обработки запроса с использованием обертки Size
	CalSize := func(w http.ResponseWriter, r *http.Request) {
		// Создаем объект Size для обертки ResponseWriter
		s := Size{w}
		// Передаем управление следующему обработчику в цепочке
		h.ServeHTTP(s, r)
	}
	// Возвращаем обработчик HTTP с добавленным расчетом размера ответа
	return http.HandlerFunc(CalSize)
}
