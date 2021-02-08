package api

import (
	"bytes"
	"os/exec"
	"path/filepath"
)

// ConvertToMP4 ConvertToMP4
func ConvertToMP4(m3u8FilePath string, mp4FilePath string) (string, error) {
	var out bytes.Buffer
	// workDir := filepath.Dir(m3u8FilePath)
	// m3u8FileName := filepath.Base(m3u8FilePath)
	ffmpegPath, _ := exec.LookPath("ffmpeg")
	ffmpegPath, _ = filepath.Abs(ffmpegPath)
	// ffmpeg -allowed_extensions ALL -protocol_whitelist file,http,https,crypto,tcp,tls -i ./output.m3u8 -c copy ./output.mp4
	cmd := exec.Command(ffmpegPath, "-allowed_extensions", "ALL", "-protocol_whitelist", "file,http,https,crypto,tcp,tls",
		"-i", m3u8FilePath,
		"-c", "copy", mp4FilePath,
	)
	// cmd.Dir = workDir
	cmd.Stdout = &out
	cmd.Stderr = &out
	var in bytes.Buffer
	in.WriteString("y\n")
	cmd.Stdin = &in
	if err := cmd.Start(); err != nil {
		return out.String(), err
	}
	if err := cmd.Wait(); err != nil {
		return out.String(), err
	}
	return out.String(), nil
}
