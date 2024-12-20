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

type Version struct {
	Version    string `json:"version"`
	Repository string `json:"repository"`
}

type YTDownloadResponse struct {
	Filname string   `json:"filename"`
	Url     string   `json:"webpage_url"`
	Version *Version `json:"_version"`
}

func YoutubeDownload(ctx context.Context, youtubeShareLink string) (*YTDownloadResponse, error) {
	cmd := exec.CommandContext(ctx, "termux-url-opener", youtubeShareLink)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("termux YoutubeDownload StdoutPipe failed: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("termux YoutubeDownload cmd.Start failed: %w", err)
	}

	var lastLine string
	scanner := bufio.NewScanner(stdout)
	buf := make([]byte, 0, 64*1024) // 64KB buffer
	scanner.Buffer(buf, 1024*1024)  // Max buffer size is 1MB
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("termux YoutubeDownload scanner.Err failed: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		slog.WarnContext(ctx, "termux YoutubeDownload cmd.Wait Stderr output:", "stderr", stderr.String())
	}

	var obj YTDownloadResponse
	err = json.Unmarshal([]byte(lastLine), &obj)
	if err != nil {
		return nil, fmt.Errorf("termux YoutubeDownload json.Unmarshal failed with youtubeShareLink %s: %w", youtubeShareLink, err)
	}

	return &obj, nil
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
		return fmt.Errorf("termux media player play failed mediaFile:%s : %s\n %w", mediaFile, string(out), err)
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
		return nil, fmt.Errorf("termux GetTextMessages cmd.StdoutPipe failed: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("termux GetTextMessages cmd.start failed: %w", err)
	}

	var msgs []TextMessage
	if err := json.NewDecoder(out).Decode(&msgs); err != nil {
		return nil, fmt.Errorf("termux json.NewDecoder cmd.start failed: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		return nil, fmt.Errorf("termux sms list cmd.wait failed. cmd.Stderr: %s\n: %w", stderr.String(), err)
	}

	return msgs, nil
}
