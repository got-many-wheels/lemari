package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/got-many-wheels/lemari/internal/config"
	directorynode "github.com/got-many-wheels/lemari/internal/directory_node"
)

type mpdOpts struct {
	input        string
	bitrate      string
	audioRate    string
	videoProfile string
	output       string
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	files := []string{}

	for _, target := range cfg.Target {
		dn := directorynode.New()
		_, err := dn.Scan(target)
		if err != nil {
			panic(err)
		}
		files = append(files, dn.DirFiles()...)
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	opts := mpdOpts{
		input:        files[0],
		bitrate:      "1000k",
		audioRate:    "44100",
		videoProfile: "high",
		output:       filepath.Join(wd, "output", filepath.Base(files[0]), "output.mpd"),
	}

	if err := creatempd(opts); err != nil {
		fmt.Println(err.Error())
	}
}

func creatempd(opt mpdOpts) error {
	if err := os.MkdirAll(filepath.Dir(opt.output), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// TODO: add flag to transcode lower resolutions also (god i wish have a bulky pc for this)
	cmd := exec.Command(
		"ffmpeg", "-i", opt.input,
		"-map", "0:v:0",
		"-map", "0:a:0",
		"-c:a", "aac",
		"-c:v", "libx264",
		"-b:v", opt.bitrate,
		"-profile:v", opt.videoProfile,
		"-preset", "fast",
		"-bf", "1",
		"-keyint_min", "120",
		"-g", "120",
		"-sc_threshold", "0",
		"-b_strategy", "0",
		"-ar", opt.audioRate,
		"-use_timeline", "1",
		"-use_template", "1",
		"-adaptation_sets", "id=0,streams=v id=1,streams=a",
		"-f", "dash", opt.output,
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	go streamOutput("STDOUT", stdout)
	go streamOutput("STDERR", stderr)

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg failed: %w", err)
	}

	return nil
}

func streamOutput(prefix string, r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Printf("[%s] %s\n", prefix, scanner.Text())
	}
}

func escapeArg(arg string) string {
	return strings.ReplaceAll(arg, " ", "\\ ")
}
