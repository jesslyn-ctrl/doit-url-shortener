package url

import "errors"

var (
	// ErrNotFound is returned when a short URL does not exist or expired
	ErrNotFound = errors.New("short url not found")

	// ErrExpired is returned when a short URL exists but has expired
	ErrExpired = errors.New("short url expired")

	// ErrInvalidURL is returned when the provided long URL which
	// not a valid HTTP/HTTPS URL
	ErrInvalidURL = errors.New("invalid url")

	// ErrCodeCollision is returned when the system fails to
	// generate a unique short code after retries
	ErrCodeCollision = errors.New("short code collision")
)
