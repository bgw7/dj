package restapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/bgw7/dj/internal"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type DJRoombaService interface {
	Download(ctx context.Context, url *internal.DownloadRequest) error
	GetTracks(ctx context.Context) ([]internal.Track, error)
	CreateTrack(ctx context.Context, t *internal.Track) (*internal.Track, error)
	CreateVote(ctx context.Context) error
	DeleteVote(ctx context.Context) error
}

type Service interface {
	DJRoombaService
}

type Handler struct {
	service  Service
	mediaDir string
	router   *chi.Mux
}

func NewHandler(service Service, mediaDir string) *Handler {
	h := &Handler{
		service:  service,
		mediaDir: mediaDir,
		router:   chi.NewRouter(),
	}
	h.setupRoutes()
	return h
}

// setupRoutes initializes the router and middleware.
func (h *Handler) setupRoutes() {
	r := h.router
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer, middleware.Timeout(10*time.Second))
	r.Use(metadataMiddleware)

	r.Route("/tracks", func(r chi.Router) {
		r.Get("/", handleOut(h.service.GetTracks, http.StatusOK))
		r.Post("/", handleInOut(h.service.CreateTrack, http.StatusCreated))
		r.Route("/download", func(r chi.Router) {
			r.Use(timeoutHandler(20 * time.Second))
			r.Post("/", handleIn(h.service.Download, http.StatusOK))
		})
		r.Route("/{trackId}/vote", func(r chi.Router) {
			r.Use(voteMiddleware)
			r.Post("/", handleNil(h.service.CreateVote, http.StatusCreated))

			r.Route("/{voteId}", func(r chi.Router) {
				r.Use(voteIDMiddleware) // Extracts voteId for DELETE requests
				r.Delete("/", handleNil(h.service.DeleteVote, http.StatusOK))
			})
		})
	})
}

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

// Generic handler function types.
type targetFunc[In any, Out any] func(context.Context, In) (Out, error)
type targetInFunc[In any] func(context.Context, In) error
type targetOutFunc[Out any] func(context.Context) (Out, error)

// Common JSON response helper.
func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func handleInOut[In any, Out any](f targetFunc[In, Out], code int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var in In

		// Retrieve data from request.
		err := json.NewDecoder(r.Body).Decode(&in)
		if err != nil {
			// Format error response
			handleError(w, err)
			return
		}

		// Call out to target function
		out, err := f(r.Context(), in)
		if err != nil {
			// Format error response
			handleError(w, err)
			return
		}

		// Format and write response
		writeJSONResponse(w, code, out)
	})
}

func handleIn[In any](f targetInFunc[In], code int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var in In

		err := json.NewDecoder(r.Body).Decode(&in)
		if err != nil {
			// Format error response
			handleError(w, err)
			return
		}

		err = f(r.Context(), in)
		if err != nil {
			handleError(w, err)
			return
		}

		writeJSONResponse(w, code, nil)
	})
}

func handleOut[Out any](f targetOutFunc[Out], code int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		out, err := f(r.Context())
		if err != nil {
			handleError(w, err)
			return
		}

		writeJSONResponse(w, code, out)
	})
}

func handleNil(f func(context.Context) error, code int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := f(r.Context())
		if err != nil {
			handleError(w, err)
			return
		}

		writeJSONResponse(w, code, nil)
	})
}
