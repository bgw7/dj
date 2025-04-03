package service

import (
	"context"
	"errors"

	"github.com/bgw7/dj/internal"
	"github.com/bgw7/dj/internal/appcontext"
	"github.com/bgw7/dj/internal/youtube"
)

func (s *DomainService) GetTracks(ctx context.Context) ([]internal.Track, error) {
	return nil, errors.New("not implements")
}

func (s *DomainService) GetNextTrack(ctx context.Context) (*internal.Track, error) {
	return s.datastore.GetNextTrack(ctx)
}

func (s *DomainService) CreateTrack(ctx context.Context, t *internal.Track) (*internal.Track, error) {
	m, err := appcontext.FromContext[*internal.Metadata](ctx, appcontext.MetadataCTXKey)
	if err != nil {
		return nil, err
	}
	t.CreatedBy = m.CreatedBy

	if t, err := s.datastore.CreateTrack(ctx, t); err != nil {
		if errors.Is(err, internal.ErrUniqueConstraintViolation) {
			return t, s.datastore.CreateVote(ctx, &internal.Vote{Filename: t.Filename, Url: t.Url, VoterID: t.CreatedBy})
		}
		return nil, err

	}
	return t, nil
}

func (s *DomainService) Download(ctx context.Context, req *internal.DownloadRequest) error {
	r, err := youtube.Download(ctx, s.mediaDir, req.URL)
	println(r)
	return err
}
func (s *DomainService) CreateVote(ctx context.Context) error {
	v, err := appcontext.FromContext[*internal.Vote](ctx, appcontext.DJRoombaVoteCTXKey)
	if err != nil {
		return err
	}
	return s.datastore.CreateVote(ctx, v)
}

func (s *DomainService) DeleteVote(ctx context.Context) error {
	v, err := appcontext.FromContext[*internal.Vote](ctx, appcontext.DJRoombaVoteCTXKey)
	if err != nil {
		return err
	}
	return s.datastore.DeleteVote(ctx, v)
}
