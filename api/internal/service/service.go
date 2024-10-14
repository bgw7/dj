package service

import (
	"context"

	"github.com/la-viajera/reservation-service/internal"
)

type DataStorage interface {
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
