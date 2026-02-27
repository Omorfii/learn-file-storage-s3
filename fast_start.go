package main

import "os/exec"

func processVideoForFastStart(filepath string) (string, error) {

	outputFile := filepath + ".processing"
	cmd := exec.Command("ffmpeg", "-i", filepath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", outputFile)
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return outputFile, nil
}
