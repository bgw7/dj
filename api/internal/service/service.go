package service

import (
	"context"
	"sync"

	"github.com/bgw7/dj/internal"
)

type DataStorage interface {
	ListTracks(ctx context.Context) ([]internal.Track, error)
	CreateTrack(ctx context.Context, t *internal.Track) error
	UpdateTrack(ctx context.Context, t *internal.Track) error
	CreateVote(ctx context.Context, url string, userId string) error
	DeleteVote(ctx context.Context, url string, userId string) error
}
type DomainService struct {
	datastore DataStorage
	readMsgs  sync.Map
}

func NewDomainService(datastore DataStorage) *DomainService {
	return &DomainService{
		datastore: datastore,
		readMsgs:  sync.Map{},
	}
}
