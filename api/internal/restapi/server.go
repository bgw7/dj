package restapi

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"errors"

	"github.com/la-viajera/reservation-service/internal"
)

// handleError provides JSON response body and HTTP
// statusCodes for common errors
func handleError(w http.ResponseWriter, err error) {
	var parseError *time.ParseError
	var jsonError *json.SyntaxError
	code := http.StatusInternalServerError
	msg := err.Error()
	switch {
	case errors.Is(err, io.EOF):
		code = http.StatusBadRequest
	case errors.As(err, &jsonError):
		code = http.StatusBadRequest
	case errors.As(err, &parseError):
		code = http.StatusBadRequest
	case errors.Is(err, internal.RecordNotFoundErr):
		code = http.StatusNotFound
	default:
		slog.Error("http internal server error", "error", err)
		code = http.StatusInternalServerError
		msg = "InternalServerError"
	}
	httpError(w, map[string]string{"error": msg}, code)
}

func httpError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}
