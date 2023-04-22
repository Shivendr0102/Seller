package model

type URLMap struct {
	ID           int    `json:"urlmap_id"`
	Url          string `json:"url"`
	ShortenedURL string `json:"short_url"`
	DateTime     string `json:"date_time"`
}
