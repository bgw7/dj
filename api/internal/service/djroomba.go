package service

import (
	"context"

	"github.com/la-viajera/reservation-service/internal"
	"github.com/la-viajera/reservation-service/internal/appcontext"
	"github.com/la-viajera/reservation-service/internal/termux"
)

func (s *DomainService) ListTracks(ctx context.Context) ([]internal.Track, error) {
	return s.datastore.ListTracks(ctx)
}

func (s *DomainService) CreatTrack(ctx context.Context, t internal.Track) (*internal.Track, error) {
	m, err := appcontext.FromContext[*internal.Metadata](ctx, appcontext.MetadataCTXKey)
	if err != nil {
		return nil, err
	}
	t.CreatedBy = m.CreatedBy
	return nil, s.datastore.CreateTrack(ctx, t)
}

func (s *DomainService) CreateVote(ctx context.Context) error {
	v, err := appcontext.FromContext[*internal.Vote](ctx, appcontext.DJRoombaVoteCTXKey)
	if err != nil {
		return err
	}
	return s.datastore.CreateVote(ctx, v.TrackID, v.UserID)
}

func (s *DomainService) DeleteVote(ctx context.Context) error {
	v, err := appcontext.FromContext[*internal.Vote](ctx, appcontext.DJRoombaVoteCTXKey)
	if err != nil {
		return err
	}
	return s.datastore.DeleteVote(ctx, v.TrackID, v.UserID)
}

// get SMS
// for each
// ID not in map
// body is https://
// download
//   - err into chan
//   - return filename into chan
//
// err chan -> termux notify
// filename chan -> upsert DB tracks
func (s *DomainService) StartTracksPolling(ctx context.Context) error {
	_, err := termux.GetTextMessages(ctx)
	if err != nil {
		return err
	}

	return nil
}
