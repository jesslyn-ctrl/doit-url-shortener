package url

import "time"

// ShortURL represents a shortened URL entity
type ShortURL struct {
	Code           string    `json:"code"`
	LongURL        string    `json:"long_url"`
	ClickCount     int64     `json:"click_count"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
	LastAccessedAt time.Time `json:"last_accessed_at"`
}
