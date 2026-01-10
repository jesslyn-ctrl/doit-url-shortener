package url

import "time"

// CreateShortURLRequest contain the request
type CreateShortURLRequest struct {
	LongURL    string `json:"long_url"`
	TTLSeconds int64  `json:"ttl_seconds,omitempty"`
}

// CreateShortURLResponse contain the response
type CreateShortURLResponse struct {
	Code      string    `json:"code"`
	LongURL   string    `json:"long_url"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
