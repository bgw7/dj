package termux

import (
	"fmt"
	"os/exec"
)

func YoutubeDownload(youtubeShareLink string) (string, error) {
	fmt.Println(youtubeShareLink)
	out, err := exec.Command("termux-url-opener", youtubeShareLink).Output()
	fmt.Println("out from youtube download", string(out))
	if err != nil {
		return "", err
	}
	return string(out), err
}

func MediaPlayer(mediaFile string) error {
	_, err := exec.Command("termux-media-player", "play", mediaFile).Output()
	return err
}

func Notify(content string) error {
	_, err := exec.Command("termux-notification", "-c", content).Output()
	return err
}
