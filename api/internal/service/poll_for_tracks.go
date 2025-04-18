//go:build android

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

func (s *DomainService) pollForTracks(ctx context.Context) {
	slog.InfoContext(ctx, "SMS Poller started: checking messages every 3 seconds")
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.checkSMS(ctx); err != nil {
				slog.ErrorContext(ctx, "checkSMS error", "error", err)
				audio.Notify(ctx, err.Error())
			}
		case <-ctx.Done():
			slog.WarnContext(ctx, "Shutting down SMS Poller", "contextErr", ctx.Err())
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

	var eg errgroup.Group
	for _, msg := range msgs {
		// Prevent processing the same message more than once
		if _, loaded := processedMsgs.LoadOrStore(msg.ID, struct{}{}); !loaded {
			eg.Go(func(m TextMessage) func() error {
				return func() error {
					return s.saveTrack(context.Background(), m.Body, m.FromNumber)
				}
			}(msg))
		}
	}

	return eg.Wait()
}

func (s *DomainService) saveTrack(ctx context.Context, body string, fromNumber string) error {
	if !strings.Contains(body, "http") {
		return nil
	}

	url := strings.TrimSpace(body)

	// Download the YouTube video
	resp, err := youtube.Download(ctx, s.mediaDir, url)
	if err != nil {
		return fmt.Errorf("failed to download YouTube video: %w", err)
	}

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
