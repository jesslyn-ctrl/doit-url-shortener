package url

import (
	"encoding/json"
	"net/http"
	"time"

	_domainUrl "github.com/jesslyn-ctrl/doit-url-shortener/internal/domain/url"
)

// CreateShortURLHandler handles POST /shorten
func CreateShortURLHandler(svc *_domainUrl.Service, defaultTTL time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var req CreateShortURLRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if req.LongURL == "" {
			http.Error(w, "long_url is required", http.StatusBadRequest)
			return
		}

		ttl := defaultTTL
		if req.TTLSeconds > 0 {
			ttl = time.Duration(req.TTLSeconds) * time.Second
		}

		u, err := svc.CreateShortURL(r.Context(), req.LongURL, ttl)
		if err != nil {
			mapDomainError(w, err)
			return
		}

		res := CreateShortURLResponse{
			Code:      u.Code,
			LongURL:   u.LongURL,
			CreatedAt: u.CreatedAt,
			ExpiresAt: u.ExpiresAt,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(res)
	}
}
