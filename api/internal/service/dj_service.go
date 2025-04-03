package service

import (
	"context"
	"errors"

	"github.com/bgw7/dj/internal"
	"github.com/bgw7/dj/internal/appcontext"
	"github.com/bgw7/dj/internal/youtube"
)

func (s *DomainService) GetTracks(ctx context.Context) ([]internal.Track, error) {
	return s.datastore.GetTracks(ctx)
}

func (s *DomainService) GetNextTrack(ctx context.Context) (*internal.Track, error) {
	return s.datastore.GetNextTrack(ctx)
}

func (s *DomainService) CreateTrack(ctx context.Context, t *internal.Track) (*internal.Track, error) {
	m, ok := appcontext.GetMetadata(ctx)

	if !ok {
		return nil, internal.ErrCtxKeyNotFound
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
	if err != nil {
		return err
	}
	_, err = s.CreateTrack(ctx, &internal.Track{
		Url:         r.Url,
		Filename:    r.Filename,
		CreatedWith: r.CreatedWith(),
	})
	if err != nil {
		return err
	}
	return nil
}
func (s *DomainService) CreateVote(ctx context.Context) error {
	v, ok := appcontext.GetVoteRequest(ctx)
	if !ok {
		return internal.ErrCtxKeyNotFound
	}

	t, err := s.datastore.GetTrackByID(ctx, v.TrackID)
	if err != nil {
		return err
	}

	return s.datastore.CreateVote(ctx, &internal.Vote{
		Filename: t.Filename,
		Url:      t.Url,
		VoterID:  v.VoterID,
	})
}

func (s *DomainService) DeleteVote(ctx context.Context) error {
	v, ok := appcontext.GetVoteRequest(ctx)
	if !ok {
		return internal.ErrCtxKeyNotFound
	}
	return s.datastore.DeleteVote(ctx, v.ID)
}
