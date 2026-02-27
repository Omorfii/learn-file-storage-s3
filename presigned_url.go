package main

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func generatePresignedUrl(s3Client *s3.Client, bucket, key string, expireTime time.Duration) (string, error) {

	presignClient := s3.NewPresignClient(s3Client)

	parameters := s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	request, err := presignClient.PresignGetObject(context.Background(), &parameters, s3.WithPresignExpires(expireTime))
	if err != nil {
		return "", err
	}

	return request.URL, nil
}
