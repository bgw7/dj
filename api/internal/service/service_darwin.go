//go:build android

package service

import (
	"context"
	"log/slog"
	"runtime"
)

func (s *DomainService) listenOnTextMsgs(ctx context.Context) {
	slog.InfoContext(ctx, "starting listenOnTextMsgs no-op", "os", runtime.GOOS)
}
