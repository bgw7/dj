package service

import (
	"context"
	"errors"
	"log/slog"
	"path/filepath"
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

func (s *DomainService) CreateTrack(ctx context.Context, t *internal.Track) (*internal.Track, error) {
	m, err := appcontext.FromContext[*internal.Metadata](ctx, appcontext.MetadataCTXKey)
	if err != nil {
		return nil, err
	}
	t.CreatedBy = m.CreatedBy
	t.CreatedWith = m.CreatedWith
	if t, err := s.datastore.CreateTrack(ctx, t); err != nil {
		if errors.Is(err, internal.ErrUniqueConstraintViolation) {
			return t, s.datastore.CreateVote(ctx, &internal.Vote{Filename: *t.Filename, Url: t.Url, VoterID: t.CreatedBy})
		}
		return nil, err

	}
	return t, nil
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

// get top voted track
// start media player
// db update hasPlayed
// ticker.Reset based on media player info maxTime - currPosition
func (s *DomainService) RunPlayNext(ctx context.Context) error {
	slog.InfoContext(ctx, "RunPlayNext started")
	delay := time.Second * 1
	ticker := time.NewTicker(delay)
	for {
		select {
		case <-ticker.C:
			delay, err := s.playNext(ctx)
			if err != nil {
				slog.ErrorContext(ctx, "media player error", "error", err)
				return termux.Notify(ctx, err.Error())
			}
			slog.InfoContext(ctx, "delay before next track play set", "delay", delay.String())
			ticker.Reset(delay)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
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
func (s *DomainService) RunSmsPoller(ctx context.Context) error {
	slog.InfoContext(ctx, "RunSmsPoller started")
	delay := time.Second * 3
	ticker := time.NewTicker(delay)
	for {
		select {
		case <-ticker.C:
			err := s.smsPoll(ctx)
			if err != nil {
				slog.ErrorContext(ctx, "s.smsPoll() error", "error", err)
				return termux.Notify(ctx, err.Error())
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *DomainService) playNext(ctx context.Context) (time.Duration, error) {
	duration := time.Second * 1
	info, err := termux.MediaInfo(ctx)
	if err != nil {
		return duration, err
	}
	if strings.Contains("Current Position:", info) {
		sub := strings.TrimPrefix(info, "Current Position:")
		times := strings.Split(sub, "/")
		currPos, err := time.ParseDuration(strings.TrimSpace(times[0]))
		if err != nil {
			return duration, err
		}
		totalDur, err := time.ParseDuration(strings.TrimSpace(times[1]))
		if err != nil {
			return duration, err
		}
		return totalDur - currPos, nil
	}
	t, err := s.ListTracks(ctx)
	if err != nil || len(t) == 0 {
		return duration, err
	}
	slog.InfoContext(ctx, "starting next track", "filename", *t[0].Filename)
	if err := termux.MediaPlay(ctx, *t[0].Filename); err != nil {
		return duration, err
	}
	u := t[0]
	u.HasPlayed = true
	if err := s.datastore.UpdateTrack(ctx, &u); err != nil {
		return duration, err
	}
	return s.playNext(ctx)
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
			return srv.saveTrack(ctx, m.ID, m.Body, m.FromNumber)
		})
	}
	return eg.Wait()
}

func (s *DomainService) saveTrack(ctx context.Context, threadID int, body, fromNumber string) error {
	_, ok := s.readMsgs.Load(threadID)
	if !ok {
		s.readMsgs.Store(threadID, "")
		if strings.Contains(body, "https://") {
			slog.InfoContext(ctx, "msg contains https://", "body", body)
			url := strings.TrimSpace(body)
			r, err := termux.YoutubeDownload(ctx, url)
			if err != nil {
				return err
			}
			path := filepath.Join("/storage/emulated/0/Termux_Downloader/Youtube/", r.Filname)
			slog.InfoContext(ctx, "saving track to DB", "filename", r.Filname, "from", fromNumber, "path", path)
			ctx := context.WithValue(ctx, appcontext.MetadataCTXKey, &internal.Metadata{
				CreatedBy:   fromNumber,
				CreatedWith: strings.Join([]string{r.Version.Repository, r.Version.Version}, "-"),
			})

			_, err = s.CreateTrack(
				ctx,
				&internal.Track{
					Url:      url,
					Filename: &path,
				})
			return err
		}
	}
	return nil
}
