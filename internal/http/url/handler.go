package url

import (
	"encoding/json"
	"net/http"
	"time"

	chiRouter "github.com/go-chi/chi/v5"
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

// ResolveHandler handles GET /s/{code}
func ResolveHandler(svc *_domainUrl.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chiRouter.URLParam(r, "code")
		if code == "" {
			http.Error(w, "code is required", http.StatusBadRequest)
			return
		}

		res, err := svc.Resolve(r.Context(), code)
		if err != nil {
			mapDomainError(w, err)
			return
		}

		// Requires for 302 redirect
		http.Redirect(w, r, res.LongURL, http.StatusFound)
	}
}

// GetStatsHandler handles GET /stats/{code}
func GetStatsHandler(svc *_domainUrl.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chiRouter.URLParam(r, "code")
		if code == "" {
			http.Error(w, "code is required", http.StatusBadRequest)
			return
		}

		res, err := svc.GetStats(r.Context(), code)
		if err != nil {
			mapDomainError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(res)
	}
}
