package main

import (
	"errors"
	"strings"
	"time"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {

	if video.VideoURL == nil {
		return video, nil
	}

	url := strings.SplitN(*video.VideoURL, ",", 2)
	if len(url) != 2 || url[0] == "" || url[1] == "" {
		return video, errors.New("not valid url")
	}

	bucket := url[0]
	key := url[1]
	duration := time.Duration(5 * time.Minute)

	presignedUrl, err := generatePresignedUrl(cfg.s3Client, bucket, key, duration)
	if err != nil {
		return video, err
	}

	video.VideoURL = &presignedUrl

	return video, nil
}
