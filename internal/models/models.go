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
