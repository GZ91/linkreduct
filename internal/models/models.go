package models

// StructURL представляет структуру данных для хранения информации о URL.
type StructURL struct {
	ID          string `json:"uuid" db:"uuid"`
	ShortURL    string `json:"short_url" db:"ShortURL"`
	OriginalURL string `json:"original_url" db:"OriginalURL"`
	UserID      string `json:"user_id"`
	DeletedFlag bool
}

// RequestData представляет структуру данных для запроса с одним полем URL.
type RequestData struct {
	URL string `json:"url"`
}

// ResultReturn представляет структуру данных для возвращения результата с одним полем Result.
type ResultReturn struct {
	Result string `json:"result"`
}

// IncomingBatchURL представляет структуру данных для получения пакета URL.
type IncomingBatchURL struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// ReleasedBatchURL представляет структуру данных для возвращения результата обработки пакета URL.
type ReleasedBatchURL struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// ReturnedStructURL представляет структуру данных для возвращения информации о URL.
type ReturnedStructURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// CtxString представляет собой тип для ключа в контексте с типом строки.
type CtxString string

// StructDelURLs представляет структуру данных для удаления URL по URL и UserID.
type StructDelURLs struct {
	URL    string
	UserID string
}
