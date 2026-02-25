package main

import (
	"bytes"
	"encoding/json"
	"math"
	"os/exec"
)

func getVideoAspectRatio(filePath string) (string, error) {

	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filePath)

	var b bytes.Buffer
	cmd.Stdout = &b

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	type ratios struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}

	ratio := ratios{}
	err = json.Unmarshal(b.Bytes(), &ratio)
	if err != nil {
		return "", err
	}

	divised := float64(ratio.Streams[0].Width) / float64(ratio.Streams[0].Height)
	landscape := 16.0 / 9.0
	portrait := 9.0 / 16.0

	if math.Abs(divised-landscape) < 0.1 {
		return "landscape", nil
	} else if math.Abs(divised-portrait) < 0.1 {
		return "portrait", nil
	} else {
		return "other", nil
	}

}
