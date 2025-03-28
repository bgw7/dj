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
	ctxWithTimeout, cancelTimeout := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelTimeout()
	defer audio.Stop(ctxWithTimeout)
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			slog.InfoContext(ctx, "case t, ok := <-playNext: playing next")
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
			next.HasPlayed = true
			if err := s.datastore.UpdateTrack(ctx, next); err != nil {
				slog.Error("Failed to update track", "error", err, "track", next.ID)
				audio.Notify(ctx, err.Error())
				continue
			}
			slog.InfoContext(ctx, "select audio.Play", "filename", next.Filename)
			if err := audio.Play(ctx, next.Filename); err != nil {
				slog.Error("Failed to play track", "error", err, "trackFilename", next.Filename)
				audio.Notify(ctx, err.Error())
			}

		case <-ctx.Done():
			slog.Info("Shutting down playNextLoop")
			return

		}
	}
}
