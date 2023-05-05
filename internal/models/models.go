package models

type StructURL struct {
	ID          string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type RequestData struct {
	URL string `json:"url"`
}

type ResultReturn struct {
	Result string `json:"result"`
}
