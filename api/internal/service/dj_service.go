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

func (s *DomainService) Download(ctx context.Context, req internal.DownloadRequest) error {
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

// func (s *DomainService) RunPlayNext(ctx context.Context) error {
// 	slog.InfoContext(ctx, "RunPlayNext started")
// 	delay := time.Second * 1
// 	ticker := time.NewTicker(delay)
// 	for {
// 		select {
// 		case <-ticker.C:
// 			delay, err := s.playNext(ctx)
// 			if err != nil {
// 				slog.ErrorContext(ctx, "media player error", "error", err)
// 				audio.Notify(ctx, err.Error())
// 			}
// 			slog.InfoContext(ctx, "delay before next track play set", "delay", delay.String())
// 			ticker.Reset(delay)
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		}
// 	}
// }

// func (s *DomainService) RunSmsPoller(ctx context.Context) error {
// 	slog.InfoContext(ctx, "RunSmsPoller started")
// 	delay := time.Second * 3
// 	ticker := time.NewTicker(delay)
// 	for {
// 		select {
// 		case <-ticker.C:
// 			err := s.smsPoll(ctx)
// 			if err != nil {
// 				slog.ErrorContext(ctx, "s.smsPoll() error", "error", err)
// 				audio.Notify(ctx, err.Error())
// 			}
// 		case <-ctx.Done():
// 			return ctx.Err()
// 		}
// 	}
// }

// func (s *DomainService) playNext(ctx context.Context) (time.Duration, error) {
// 	duration := time.Second * 1
// 	info, err := termux.MediaInfo(ctx)
// 	if err != nil {
// 		return duration, err
// 	}
// 	if strings.Contains("Current Position:", info) {
// 		sub := strings.TrimPrefix(info, "Current Position:")
// 		times := strings.Split(sub, "/")
// 		currPos, err := time.ParseDuration(strings.TrimSpace(times[0]))
// 		if err != nil {
// 			return duration, err
// 		}
// 		totalDur, err := time.ParseDuration(strings.TrimSpace(times[1]))
// 		if err != nil {
// 			return duration, err
// 		}
// 		return totalDur - currPos, nil
// 	}
// 	t, err := s.GetTracks(ctx)
// 	if err != nil || len(t) == 0 {
// 		//TODO: no tracks to play, block until SMSpoller gets a track
// 		return duration, err
// 	}

// 	slog.InfoContext(ctx, "starting termux.MediaPlay", "filename", *t[0].Filename)
// 	if err := termux.MediaPlay(ctx, *t[0].Filename); err != nil {
// 		return duration, err
// 	}
// 	u := t[0]
// 	u.HasPlayed = true
// 	if err := s.datastore.UpdateTrack(ctx, &u); err != nil {
// 		return duration, err
// 	}
// 	return s.playNext(ctx)
// }

// func (srv *DomainService) smsPoll(ctx context.Context) error {
// 	eg := errgroup.Group{}
// 	msgs, err := termux.GetTextMessages(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	if len(msgs) == 0 {
// 		return errors.New("length of messages from inbox is 0")
// 	}
// 	for _, m := range msgs {
// 		eg.Go(func() error {
// 			return srv.saveTrack(ctx, m.ID, m.Body, m.FromNumber)
// 		})
// 	}
// 	return eg.Wait()
// }

// func (s *DomainService) saveTrack(ctx context.Context, threadID int, body, fromNumber string) error {
// 	_, ok := s.readMsgs.Load(threadID)
// 	if !ok {
// 		s.readMsgs.Store(threadID, "")
// 		if strings.Contains(body, "https://y") {
// 			slog.InfoContext(ctx, "msg contains https://", "body", body)
// 			url := strings.TrimSpace(body)
// 			ctx := context.WithValue(
// 				ctx,
// 				appcontext.MetadataCTXKey,
// 				&internal.Metadata{
// 					CreatedBy: fromNumber,
// 				},
// 			)
// 			t := &internal.Track{
// 				Url: url,
// 			}
// 			r, err := youtube.YoutubeDownload(ctx, t.Url)
// 			if err != nil {
// 				return err
// 			}
// 			t.CreatedWith = strings.Join([]string{"termux", r.Version.Repository, r.Version.Version}, "-")
// 			t.Filename = &r.Filename

// 			_, err = s.CreateTrack(
// 				ctx,
// 				t,
// 			)
// 			return err
// 		}
// 	}
// 	return nil
// }
