package service

import (
	"context"
	"sync"

	"github.com/bgw7/dj/internal"
)

type DataStorage interface {
	GetNextTrack(ctx context.Context) (*internal.Track, error)
	CreateTrack(ctx context.Context, t *internal.Track) (*internal.Track, error)
	UpdateTrack(ctx context.Context, t *internal.Track) error
	CreateVote(ctx context.Context, v *internal.Vote) error
	DeleteVote(ctx context.Context, v *internal.Vote) error
}
type DomainService struct {
	datastore DataStorage
	readMsgs  sync.Map
}

func NewDomainService(ctx context.Context, datastore DataStorage) *DomainService {
	ds := &DomainService{
		datastore: datastore,
		readMsgs:  sync.Map{},
	}
	go ds.listenOnTextMsgs(ctx)
	go ds.playNextLoop(ctx)
	return ds
}
