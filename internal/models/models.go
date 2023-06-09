package models

type StructURL struct {
	ID          string `json:"uuid" db:"uuid"`
	ShortURL    string `json:"short_url" db:"ShortURL"`
	OriginalURL string `json:"original_url" db:"OriginalURL"`
}

type RequestData struct {
	URL string `json:"url"`
}

type ResultReturn struct {
	Result string `json:"result"`
}

type IncomingBatchURL struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ReleasedBatchURL struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
