package service

import (
	"context"
	"log/slog"
	"runtime"
	"time"

	"github.com/bgw7/dj/internal"
	"github.com/bgw7/dj/internal/audio"
)

func (s *DomainService) playNextLoop(ctx context.Context) {
	slog.InfoContext(ctx, "starting playNextLoop", "os", runtime.GOOS)
	playNext := make(chan *internal.Track)
	ctxWithTimeout, cancelTimeout := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelTimeout()
	defer audio.Stop(ctxWithTimeout)
	for {
		select {
		case t, ok := <-playNext:
			if !ok {
				return
			}
			t.HasPlayed = true
			if err := s.datastore.UpdateTrack(ctx, t); err != nil {
				slog.Error("Failed to update track", "error", err, "track", t.ID)
				audio.Notify(ctx, err.Error())
				continue
			}

			if err := audio.Play(ctx, t.Filename); err != nil {
				slog.Error("Failed to play track", "error", err, "trackFilename", t.Filename)
				audio.Notify(ctx, err.Error())
			}

			next, err := s.datastore.GetNextTrack(ctx)
			if err == internal.ErrRecordNotFound {
				time.Sleep(3 * time.Second)
				continue
			}
			if err != nil {
				slog.Error("Failed to fetch next track", "error", err)
				audio.Notify(ctx, err.Error())
				return
			}
			if next != nil {
				playNext <- next
			}

		case <-ctx.Done():
			slog.Info("Shutting down playNextLoop")
			close(playNext)
			return

		default:
			next, err := s.datastore.GetNextTrack(ctx)
			if err == internal.ErrRecordNotFound {
				time.Sleep(3 * time.Second)
				continue
			}
			if err != nil {
				slog.Error("Failed to fetch next track", "error", err)
				audio.Notify(ctx, err.Error())
				return
			}
			if next != nil {
				playNext <- next
			}
		}
	}
}
