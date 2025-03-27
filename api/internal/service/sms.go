package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/bgw7/dj/internal"
	"github.com/bgw7/dj/internal/audio"
	"github.com/bgw7/dj/internal/youtube"
	"golang.org/x/sync/errgroup"
)

type TextMessage struct {
	ID         int    `json:"_id"`
	FromNumber string `json:"number"`
	Body       string `json:"body"`
}

func getTextMessages(ctx context.Context) ([]TextMessage, error) {
	cmd := exec.CommandContext(ctx, "termux-sms-list")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("cmd.StdoutPipe failed: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("cmd.Start failed: %w", err)
	}

	var msgs []TextMessage
	if err := json.NewDecoder(out).Decode(&msgs); err != nil {
		return nil, fmt.Errorf("json decode failed: %w", err)
	}

	// Capture stderr if the command fails
	if err := cmd.Wait(); err != nil {
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		return nil, fmt.Errorf("cmd.Wait failed, stderr: %s: %w", stderr.String(), err)
	}

	return msgs, nil
}

func (s *DomainService) listenOnTextMsgs(ctx context.Context) {
	slog.InfoContext(ctx, "SMS Poller started: poll every 3 seconds")
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			slog.InfoContext(ctx, "listenOnTextMsgs next ticker")
			if err := s.checkSMS(ctx); err != nil {
				slog.ErrorContext(ctx, "checkSMS error", "error", err)
				audio.Notify(ctx, err.Error())
			}
		case <-ctx.Done():
			slog.InfoContext(ctx, "Shutting down SMS Poller")
			return
		}
	}
}

var processedMsgs sync.Map

func (s *DomainService) checkSMS(ctx context.Context) error {
	msgs, err := getTextMessages(ctx)
	if err != nil {
		return fmt.Errorf("failed to get text messages: %w", err)
	}

	eg, egCtx := errgroup.WithContext(ctx)
	for _, msg := range msgs {
		// Prevent processing the same message more than once
		slog.InfoContext(ctx, "check processedMsgs.LoadOrStore")
		if _, loaded := processedMsgs.LoadOrStore(msg.ID, struct{}{}); !loaded {
			slog.InfoContext(ctx, "check processedMsgs.LoadOrStore: not loaded, starting routine to save track")
			eg.Go(func(m TextMessage) func() error {
				return func() error {
					return s.saveTrack(egCtx, m.Body, m.FromNumber)
				}
			}(msg))
		}
	}

	return eg.Wait()
}

func (s *DomainService) saveTrack(ctx context.Context, body string, fromNumber string) error {
	if !strings.Contains(body, "youtube") {
		slog.InfoContext(ctx, "shared track does not contain https://y")
		return nil
	}

	slog.InfoContext(ctx, "Message contains YouTube link", "body", body)
	url := strings.TrimSpace(body)

	// Download the YouTube video
	resp, err := youtube.Download(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to download YouTube video: %w", err)
	}
	slog.InfoContext(ctx, "YT Download Response", "resp", resp)

	// Create the track object
	t := &internal.Track{
		Url:         url,
		Filename:    resp.Filename,
		CreatedBy:   fromNumber,
		CreatedWith: resp.CreatedWith(),
	}

	// Store the track in the datastore
	_, err = s.datastore.CreateTrack(ctx, t)
	if err != nil {
		return fmt.Errorf("failed to create track in datastore: %w", err)
	}

	return nil
}
