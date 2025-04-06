//go:build darwin

package service

import (
	"context"
	"log/slog"
	"runtime"
)

func (s *DomainService) pollForTracks(ctx context.Context) {
	slog.InfoContext(ctx, "starting listenOnTextMsgs no-op", "os", runtime.GOOS)
}
