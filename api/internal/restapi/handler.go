package restapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/la-viajera/reservation-service/internal"
)

type ReservationService interface {
	GetReservations(ctx context.Context) ([]internal.Reservation, error)
	CreateReservation(ctx context.Context, obj *internal.Reservation) (*internal.Reservation, error)
	SearchReservations(ctx context.Context, s *internal.ReservationSearch) (*internal.PagedResponse[internal.Reservation], error)
	FindOneReservation(ctx context.Context) (*internal.Reservation, error)
	UpdateReservation(ctx context.Context, r *internal.Reservation) (*internal.Reservation, error)
}

type VenueService interface {
	GetVenues(ctx context.Context) ([]internal.Venue, error)
	CreateVenue(ctx context.Context, obj *internal.Venue) (*internal.Venue, error)
	SuggestVenues(ctx context.Context, query string) (*internal.PagedResponse[internal.Venue], error)
	SearchVenues(ctx context.Context, s *internal.VenueSearch) (*internal.PagedResponse[internal.Venue], error)
}

type DJRoombaService interface {
	ListTracks(ctx context.Context) ([]internal.Track, error)
	CreatTrack(ctx context.Context, t internal.Track) (*internal.Track, error)
	CreateVote(ctx context.Context) error
	DeleteVote(ctx context.Context) error
}

type Service interface {
	ReservationService
	VenueService
	DJRoombaService
}

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{
		service: s,
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
	r.Route("/reservations", func(r chi.Router) {
		r.Get("/", handleOut(h.service.GetReservations, http.StatusOK))
		r.Post("/", handleInOut(h.service.CreateReservation, http.StatusCreated))
		r.Post("/search", handleInOut(h.service.SearchReservations, http.StatusOK))
		// Subrouters:  /reservations/123/
		r.Route("/{reservationID}", func(r chi.Router) {
			r.Use(h.reservationIDMiddleware)
			r.Get("/", handleOut(h.service.FindOneReservation, http.StatusOK))
			r.Patch("/", handleInOut(h.service.UpdateReservation, http.StatusOK))
		})
	})

	r.Route("/tracks", func(r chi.Router) {
		r.Use(metdataMiddleware)
		r.Get("/", handleOut(h.service.ListTracks, http.StatusOK))
		r.Post("/", handleInOut(h.service.CreatTrack, http.StatusCreated))
		r.Route("/{trackId}/votes", func(r chi.Router) {
			r.Use(djRoombaVoteMiddleware)
			r.Post("/", handleNil(h.service.CreateVote, http.StatusCreated))
			r.Delete("/", handleNil(h.service.DeleteVote, http.StatusOK))
		})
	})

	r.Route("/venues", func(r chi.Router) {
		r.Get("/", handleOut(h.service.GetVenues, http.StatusOK))
		r.Post("/", handleInOut(h.service.CreateVenue, http.StatusCreated))
		r.Post("/search", handleInOut(h.service.SearchVenues, http.StatusOK))
		r.Post("/suggestions", h.suggestVenues)
		// Subrouters:  /venues/123/
		r.Route("/{reservationID}", func(r chi.Router) {
			// r.Use(venueCtx....
			// r.Get("/", ...
			// r.Patch("/", ...
		})
	})
	r.ServeHTTP(w, req)
}

func (h *Handler) suggestVenues(w http.ResponseWriter, req *http.Request) {
	q := chi.URLParam(req, "q")
	page, err := h.service.SuggestVenues(req.Context(), q)
	if err != nil {
		handleError(w, err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(page)
}

type targetFunc[In any, Out any] func(context.Context, In) (Out, error)
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
		return
	})
}
