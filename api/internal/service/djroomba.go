package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/bgw7/dj/internal"
	"github.com/bgw7/dj/internal/appcontext"
	"github.com/bgw7/dj/internal/termux"
	"golang.org/x/sync/errgroup"
)

func (s *DomainService) ListTracks(ctx context.Context) ([]internal.Track, error) {
	return s.datastore.ListTracks(ctx)
}

func (s *DomainService) CreatTrack(ctx context.Context, t *internal.Track) (*internal.Track, error) {
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
//   - return err
//   - return filename
//
// err to termux notify
// filename insert to DB tracks
func (srv *DomainService) RunSmsPoller(ctx context.Context) error {
	slog.InfoContext(ctx, "RunSmsPoller started")
	delay := time.Second * 3
	ticker := time.NewTicker(delay)
	for {
		select {
		case <-ticker.C:
			err := srv.smsPoll(ctx)
			if err != nil {
				slog.ErrorContext(ctx, "sms poller error", "error", err)
				return termux.Notify(ctx, err.Error())
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (srv *DomainService) smsPoll(ctx context.Context) error {
	eg := errgroup.Group{}
	msgs, err := termux.GetTextMessages(ctx)
	if err != nil {
		return err
	}
	if len(msgs) == 0 {
		return errors.New("length of messages from inbox is 0")
	}
	for _, m := range msgs {
		eg.Go(func() error {
			return srv.saveTrack(ctx, m.ThreadID, m.Body, m.FromNumber)
		})
	}
	return eg.Wait()
}

func (srv *DomainService) saveTrack(ctx context.Context, threadID int, body, fromNumber string) error {
	_, ok := srv.readMsgs.Load(threadID)
	if !ok {
		srv.readMsgs.Store(threadID, "")
		if strings.Contains(body, "https://") {
			slog.InfoContext(ctx, "srv.readMsgs.Load(threadID)", "ok", ok, "threadID", threadID)
			slog.InfoContext(ctx, "msg contains https://", "body", body)
			url := strings.TrimSpace(body)
			r, err := termux.YoutubeDownload(ctx, url)
			if err != nil {
				return err
			}
			slog.InfoContext(ctx, "saving track to DB", "file", r.Filname, "from", fromNumber)
			return srv.datastore.CreateTrack(
				ctx,
				&internal.Track{
					Url:       url,
					Filename:  &r.Filname,
					CreatedBy: fromNumber,
				})
		}
	}
	return nil
}
