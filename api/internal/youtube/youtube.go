package youtube

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os/exec"
	"path/filepath"
	"strings"
)

type Version struct {
	Version    string `json:"version"`
	Repository string `json:"repository"`
}

type YTDownloadResponse struct {
	Filename string   `json:"filename"`
	Url      string   `json:"webpage_url"`
	Version  *Version `json:"_version"`
}

func (y *YTDownloadResponse) CreatedWith() string {
	return strings.Join([]string{"", y.Version.Repository, y.Version.Version}, "-")
}

func Download(ctx context.Context, youtubeShareLink string) (*YTDownloadResponse, error) {
	cmd := exec.CommandContext(ctx, "termux-url-opener", youtubeShareLink)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		slog.ErrorContext(ctx, "failed cmd.StdoutPipe()", "error", err)
		return nil, fmt.Errorf("termux YoutubeDownload StdoutPipe failed: %w", err)
	}

	if err := cmd.Start(); err != nil {
		slog.ErrorContext(ctx, "failed cmd.Start()", "error", err)
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
		slog.ErrorContext(ctx, "failed scanner.Err()", "error", err)
		return nil, fmt.Errorf("termux YoutubeDownload scanner.Err failed: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		slog.ErrorContext(ctx, "youtube Download cmd.Wait() error", "error", err, "stderr", stderr.String())
		return nil, err
	}

	var obj YTDownloadResponse
	err = json.Unmarshal([]byte(lastLine), &obj)
	if err != nil {
		slog.ErrorContext(ctx, "json.Unmarshal() error", "error", err)
		return nil, fmt.Errorf("termux YoutubeDownload json.Unmarshal failed with youtubeShareLink %s: %w", youtubeShareLink, err)
	}
	obj.Filename = changeFileExtension(obj.Filename)

	return &obj, nil
}

func changeFileExtension(filePath string) string {
	oldExtension := filepath.Ext(filePath)
	return strings.TrimSuffix(filePath, oldExtension) + ".mp3"
}
