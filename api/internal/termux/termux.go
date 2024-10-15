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
	cmd := exec.CommandContext(ctx, "termux-url-opener", youtubeShareLink).Output()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(stdout)
	var lastLine string
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	if err := cmd.Wait(); err != nil {
		panic(err)
	}
	validJSON := bytes.ReplaceAll([]byte(lastLine), []byte("'"), []byte("\""))

	var obj OpenerResponse
	err = json.Unmarshal(validJSON, &obj)

	return &obj, err
}

func MediaPlayer(ctx context.Context, mediaFile string) error {
	_, err := exec.CommandContext(ctx, "termux-media-player", "play", mediaFile).Output()
	return fmt.Errorf("termux media player failed: %w", err)
}

func Notify(ctx context.Context, content string) error {
	_, err := exec.CommandContext(ctx, "termux-notification", "-c", content).Output()
	return fmt.Errorf("termux notification failed: %w", err)
}

type TextMessage struct {
	ThreadID   int    `json:"threadid"`
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
