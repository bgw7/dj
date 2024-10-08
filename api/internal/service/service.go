package service

import (
	"context"

	"github.com/la-viajera/reservation-service/internal"
)

type DataStorage interface {
	FindReservation(ctx context.Context, id string) (*internal.Reservation, error)
	UpdateReservation(ctx context.Context, r *internal.Reservation) (*internal.Reservation, error)
	GetReservations(ctx context.Context) ([]internal.Reservation, error)
	CreateReservation(ctx context.Context, r *internal.Reservation) (*internal.Reservation, error)
	SearchReservations(ctx context.Context, rs *internal.ReservationSearch) ([]internal.Reservation, error)
	CreateVenue(ctx context.Context, r *internal.Venue) (*internal.Venue, error)
	GetVenues(ctx context.Context) ([]internal.Venue, error)
	SearchVenues(ctx context.Context, rs *internal.VenueSearch) ([]internal.Venue, error)
	ListTracks(ctx context.Context) ([]internal.Track, error)
	CreateTrack(ctx context.Context, t internal.Track) error
	CreateVote(ctx context.Context, trackId int, userId string) error
	DeleteVote(ctx context.Context, trackId int, userId string) error
}
type DomainService struct {
	datastore DataStorage
}

func NewDomainService(datastore DataStorage) *DomainService {
	return &DomainService{
		datastore: datastore,
	}
}
