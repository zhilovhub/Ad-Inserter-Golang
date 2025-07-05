package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Profile struct {
	Resolution string
	Bitrate    string
}

func transcodeInputVideo(videoFileName string, outputDirectory string, profile Profile) string {
	if err := os.MkdirAll(outputDirectory, os.ModePerm); err != nil {
		fmt.Println("Error creating output directory:", err)
		return ""
	}

	outputPath := fmt.Sprintf("%s/output_%s_%s.mp4", outputDirectory, profile.Resolution, profile.Bitrate)

	cmd := exec.Command(
		"ffmpeg",
		"-i", videoFileName,
		"-vf", fmt.Sprintf("scale=%s", profile.Resolution),
		"-b:v", profile.Bitrate,
		"-c:a", "aac",
		"-strict", "experimental",
		outputPath,
	)

	if err := cmd.Run(); err != nil {
		fmt.Println("Error while transcoding the video:", err)
		return ""
	}

	fmt.Printf("Created %s\n", outputPath)
	return outputPath
}

func generateManifestAndSegments(videoFilePath string, outputDirectory string, manifestFileName string) string {
	if err := os.MkdirAll(outputDirectory, os.ModePerm); err != nil {
		fmt.Println("Error creating output directory:", err)
		return ""
	}

	cmd := exec.Command(
		"ffmpeg",
		"-i", videoFilePath,
		"-c:a", "copy",
		"-c:v", "copy",
		"-map", "0",
		"-f", "dash",
		fmt.Sprintf("%s/%s", outputDirectory, manifestFileName),
	)

	if err := cmd.Run(); err != nil {
		fmt.Println("Error while generating DASH manifest:", err)
		return ""
	}

	fmt.Println("DASH manifest has been created successfully")
	return fmt.Sprintf("%s/%s", outputDirectory, manifestFileName)
}

func main() {
	videoFileName := "input.mp4"
	outputDirectory := "segments"

	profiles := []Profile{
		{"1920x1080", "3000k"},
		{"1280x720", "1500k"},
		{"640x360", "500k"},
	}

	for _, profile := range profiles {
		transcodedVideoPath := transcodeInputVideo(
			videoFileName, 
			outputDirectory,
			profile,
		)

		if err := generateManifestAndSegments(
			transcodedVideoPath,
			strings.Replace(transcodedVideoPath, ".mp4", "", -1),
			"manifest.mpd",
		); err == "" {
			return
		}
	}
}
