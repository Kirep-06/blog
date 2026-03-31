package storage

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Storage struct {
	client       *s3.Client
	bucket       string
	publicURLBase string
}

func NewS3Storage(region, endpoint, accessKey, secretKey, bucket, publicURLBase string, forcePathStyle bool) (*S3Storage, error) {
	cfg := aws.Config{
		Region:      region,
		Credentials: credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
	}

	opts := []func(*s3.Options){
		func(o *s3.Options) {
			o.UsePathStyle = forcePathStyle
		},
	}
	if endpoint != "" {
		opts = append(opts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(endpoint)
		})
	}

	client := s3.NewFromConfig(cfg, opts...)
	return &S3Storage{client: client, bucket: bucket, publicURLBase: publicURLBase}, nil
}

func (s *S3Storage) Upload(ctx context.Context, key string, r io.Reader, size int64, mimeType string) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          r,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(mimeType),
	})
	if err != nil {
		return "", fmt.Errorf("s3 put object: %w", err)
	}
	return s.publicURLBase + "/" + key, nil
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}
