package models

type ShortenURLRequest struct {
	OriginalURL string `json:"original_url"`
	ExpireTime  int64  `json:"expire_time"`
}
