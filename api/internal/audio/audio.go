package audio

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func MediaInfo(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "termux-media-player", "info")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("termux media info failed: %w", err)
	}
	var lastLine string
	scanner := bufio.NewScanner(&stdout)
	buf := make([]byte, 0, 64*1024) // 64KB buffer
	scanner.Buffer(buf, 1024*1024)  // Max buffer size is 1MB
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		slog.ErrorContext(ctx, "failed scanner.Err()", "error", err)
		return "", fmt.Errorf("termux YoutubeDownload scanner.Err failed: %w", err)
	}
	return lastLine, nil
}

func Stop(ctx context.Context) {
	out, err := exec.CommandContext(ctx, "termux-media-player", "stop").CombinedOutput()
	if err != nil {
		slog.ErrorContext(ctx, "audio stop failed", "error", fmt.Errorf("termux media player stop failed: %s\n %w", string(out), err))
	}
}

func Play(ctx context.Context, mediaDir, mediaFile string) error {
	// fp := filepath.Join(mediaDir, mediaFile)
	fp := mediaFile
	out, err := exec.CommandContext(ctx, "termux-media-player", "play", fp).CombinedOutput()
	if err != nil {
		return fmt.Errorf("termux media player play failed. mediaFile:%s : %s\n%w", fp, string(out), err)
	}
	return blockUntilDone(ctx)
}

func Notify(ctx context.Context, content string) {
	out, err := exec.CommandContext(ctx, "termux-notification", "-c", content).CombinedOutput()
	if err != nil {
		slog.ErrorContext(ctx, "audio notify failed", "error", fmt.Errorf("termux notification failed: %s\n %w", string(out), err))
	}
}

func remainingPlayTime(ctx context.Context) (time.Duration, bool, error) {
	info, err := MediaInfo(ctx)
	if err != nil {
		return 0, false, err
	}

	if !strings.Contains(info, "Current Position:") {
		return 0, true, nil // No media playing
	}

	parts := strings.Split(strings.TrimPrefix(info, "Current Position:"), "/")
	if len(parts) != 2 {
		return 0, false, fmt.Errorf("unexpected media info format: %s", info)
	}

	currPos, err := parseMMSS(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, false, err
	}

	totalDur, err := parseMMSS(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, false, err
	}

	return totalDur - currPos, false, nil
}

func parseMMSS(durationStr string) (time.Duration, error) {
	parts := strings.Split(durationStr, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid duration format: %s", durationStr)
	}

	minutes, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid minutes value: %v", err)
	}

	seconds, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid seconds value: %v", err)
	}

	return time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second, nil
}

func blockUntilDone(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			Stop(context.Background())
			return ctx.Err()
		case <-ticker.C:
			delay, done, err := remainingPlayTime(ctx)
			if err != nil {
				return err
			}
			if done {
				return nil
			}
			slog.InfoContext(ctx, "waiting for audio play to complete", "timeRemaining", delay)
			ticker.Reset(delay)
		}
	}
}
