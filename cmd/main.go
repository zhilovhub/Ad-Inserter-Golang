package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	videoFileName := "input.mov"
	outputDirectory := "segments"

	if err := os.MkdirAll(outputDirectory, os.ModePerm); err != nil {
		fmt.Println("Error creating output directory:", err)
	}

	profiles := []struct {
		Resolution string
		Bitrate    string
	}{
		{"1920x1080", "3000k"},
		{"1280x720", "1500k"},
		{"640x360", "500k"},
	}

	for _, profile := range profiles {
		outputFileName := fmt.Sprintf("%s/output_%s_%s.mp4", outputDirectory, profile.Resolution, profile.Bitrate)
		cmd := exec.Command(
			"ffmpeg",
			"-i", videoFileName,
			"-vf", fmt.Sprintf("scale=%s", profile.Resolution),
			"-b:v", profile.Bitrate,
			"-c:a", "aac",
			"-strict", "experimental",
			outputFileName,
		)

		if err := cmd.Run(); err != nil {
			fmt.Println("Error while transcoding the video:", err)
			return
		}

		fmt.Printf("Created %s\n", outputFileName)
	}

	cmd := exec.Command(
		"ffmpeg",
		"-f", "dash",
		"-i", fmt.Sprintf("%s/output_%%v.mp4", outputDirectory),
		"-map", "0",
		"-f", "dash",
		fmt.Sprintf("%s/manifest.mpd", outputDirectory),
	)

	if err := cmd.Run(); err != nil {
		fmt.Println("Error while generating DASH manifest:", err)
		return
	}

	fmt.Println("DASH manifest has been created successfully")

}
