package termux

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
)

type OpenerResponse struct {
	Filname string `json:"Filename"`
	Url     string `json:"URL"`
}

func YoutubeDownload(ctx context.Context, youtubeShareLink string) (*OpenerResponse, error) {
	slog.InfoContext(ctx, "starting termux-url-opener", "youtubeShareLink", youtubeShareLink)
	cmd := exec.CommandContext(ctx, "termux-url-opener", youtubeShareLink)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(stdout)
	var lastLine string
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	validJSON := bytes.ReplaceAll([]byte(lastLine), []byte("'"), []byte("\""))

	var obj OpenerResponse
	err = json.Unmarshal(validJSON, &obj)

	return &obj, err
}

func MediaInfo(ctx context.Context) (string, error) {
	out, err := exec.CommandContext(ctx, "termux-media-player", "info").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("termux media info failed: %s\n %w", string(out), err)
	}
	return string(out), nil
}

func MediaStop(ctx context.Context) error {
	out, err := exec.CommandContext(ctx, "termux-media-player", "stop").CombinedOutput()
	if err != nil {
		return fmt.Errorf("termux media player stop failed: %s\n %w", string(out), err)
	}
	return err
}

func MediaPlay(ctx context.Context, mediaFile string) error {
	out, err := exec.CommandContext(ctx, "termux-media-player", "play", mediaFile).CombinedOutput()
	if err != nil {
		return fmt.Errorf("termux media player play failed: %s\n %w", string(out), err)
	}
	return nil
}

func Notify(ctx context.Context, content string) error {
	out, err := exec.CommandContext(ctx, "termux-notification", "-c", content).CombinedOutput()
	if err != nil {
		return fmt.Errorf("termux notification failed: %s\n %w", string(out), err)
	}
	return nil
}

type TextMessage struct {
	ID         int    `json:"_id"`
	FromNumber string `json:"number"`
	Body       string `json:"body"`
}

func GetTextMessages(ctx context.Context) ([]TextMessage, error) {
	cmd := exec.CommandContext(ctx, "termux-sms-list")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("termux sms list cmd.start failed: %w", err)
	}

	var msgs []TextMessage
	if err := json.NewDecoder(out).Decode(&msgs); err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("termux sms list cmd.wait failed: %w", err)
	}

	return msgs, err
}
