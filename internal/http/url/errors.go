package url

import (
	"errors"
	"net/http"

	_domainUrl "github.com/jesslyn-ctrl/doit-url-shortener/internal/domain/url"
)

func mapDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, _domainUrl.ErrInvalidURL):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, _domainUrl.ErrNotFound), errors.Is(err, _domainUrl.ErrExpired):
		http.Error(w, "not found", http.StatusNotFound)
	default:
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
