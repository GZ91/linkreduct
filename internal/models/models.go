package models

//go:generate easyjson

//easyjson:json
type StructURL struct {
	ID          string `json:"uuid" db:"uuid"`
	ShortURL    string `json:"short_url" db:"ShortURL"`
	OriginalURL string `json:"original_url" db:"OriginalURL"`
	UserID      string `json:"user_id"`
	DeletedFlag bool
}

//easyjson:json
type RequestData struct {
	URL string `json:"url"`
}

//easyjson:json
type ResultReturn struct {
	Result string `json:"result"`
}

//easyjson:json
type IncomingBatchURL struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

//easyjson:json
type ReleasedBatchURL struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

//easyjson:json
type ReturnedStructURL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type CtxString string

//easyjson:json
type StructDelURLs struct {
	URL    string
	UserID string
}
