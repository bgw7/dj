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
}

func NewHandler(s Service, mediaDir string) *Handler {
	return &Handler{
		service:  s,
		mediaDir: mediaDir,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(metdataMiddleware)
	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(10 * time.Second))
	// RESTy routes for domain model
	r.Route("/tracks", func(r chi.Router) {
		r.Use(metdataMiddleware)
		r.Get("/", handleOut(h.service.GetTracks, http.StatusOK))
		r.Post("/", handleInOut(h.service.CreateTrack, http.StatusCreated))
		r.Route("/download", func(r chi.Router) {
			r.Post("/", handleIn(h.service.Download, http.StatusOK))
		})
		r.Route("/{trackId}/vote", func(r chi.Router) {
			r.Use(voteMiddleware)
			r.Post("/", handleNil(h.service.CreateVote, http.StatusCreated))
			r.Delete("/", handleNil(h.service.DeleteVote, http.StatusOK))
		})
	})
	r.ServeHTTP(w, req)
}

type targetFunc[In any, Out any] func(context.Context, In) (Out, error)
type targetInFunc[In any] func(context.Context, In) error
type targetOutFunc[Out any] func(context.Context) (Out, error)

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
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		err = json.NewEncoder(w).Encode(out)
		if err != nil {
			handleError(w, err)
			return
		}
	})
}

func handleIn[In any](f targetInFunc[In], code int) http.HandlerFunc {
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
		err = f(r.Context(), in)
		if err != nil {
			handleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
	})
}

func handleOut[Out any](f targetOutFunc[Out], code int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		out, err := f(r.Context())
		if err != nil {
			handleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		err = json.NewEncoder(w).Encode(out)
		if err != nil {
			handleError(w, err)
			return
		}
	})
}

func handleNil(f func(context.Context) error, code int) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := f(r.Context())
		if err != nil {
			handleError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
	})
}
