package termux

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type OpenerResponse struct {
	Filname string `json:"Filename"`
	Url     string `json:"URL"`
}

func YoutubeDownload(youtubeShareLink string) (*OpenerResponse, error) {
	fmt.Println("starting termux-url-opener")
	out, err := exec.Command("termux-url-opener", youtubeShareLink).Output()

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

func MediaPlayer(mediaFile string) error {
	_, err := exec.Command("termux-media-player", "play", mediaFile).Output()
	return err
}

func Notify(content string) error {
	_, err := exec.Command("termux-notification", "-c", content).Output()
	return err
}
