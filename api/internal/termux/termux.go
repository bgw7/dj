package termux

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type OpenerResponse struct {
	Filname string `json:"Filename"`
	Url     string `json:"URL"`
}

func YoutubeDownload(ctx context.Context, youtubeShareLink string) (*OpenerResponse, error) {
	fmt.Println("starting termux-url-opener")
	out, err := exec.CommandContext(ctx, "termux-url-opener", youtubeShareLink).Output()

	if err != nil {
		return nil, err
	}

	sl := strings.Split(string(out), `\n`)
	invalidJSON := sl[len(sl)-2]
	validJSON := bytes.ReplaceAll([]byte(invalidJSON), []byte("'"), []byte("\""))

	var obj OpenerResponse
	err = json.Unmarshal(validJSON, &obj)

	return &obj, err
}

func MediaPlayer(ctx context.Context, mediaFile string) error {
	_, err := exec.CommandContext(ctx, "termux-media-player", "play", mediaFile).Output()
	return err
}

func Notify(ctx context.Context, content string) error {
	_, err := exec.CommandContext(ctx, "termux-notification", "-c", content).Output()
	return err
}

type TextMessage struct {
	ThreadID int    `json:"threadid"`
	Body     string `json:"body"`
}

func GetTextMessages(ctx context.Context) ([]TextMessage, error) {
	cmd := exec.CommandContext(ctx, "termux-sms-list")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	var msgs []TextMessage
	if err := json.NewDecoder(out).Decode(&msgs); err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return msgs, err
}
