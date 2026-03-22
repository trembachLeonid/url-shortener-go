package models

type ShortenURLRequest struct {
	OriginalURL string `json:"original_url" binding:"required"`
	ExpireTime  *int64 `json:"expire_time" binding:"omitempty,gt=0"`
}
