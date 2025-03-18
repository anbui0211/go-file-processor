package paws

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Service interface {
	CreateBucket(ctx context.Context, bucketName string) error
	UploadFile(ctx context.Context, bucketName string, key string, data []byte) error
	DownloadFile(ctx context.Context, bucketName string, key string) ([]byte, error)
}

type s3ServiceImpl struct {
	client *s3.Client
}

func NewS3Service(cfg aws.Config) S3Service {
	// Create S3 client with dummy credentials for localstack
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String("http://localhost:4566")
		o.Credentials = aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     "dummy",
				SecretAccessKey: "dummy",
				SessionToken:    "dummy",
				Source:          "Hard-coded credentials for localstack",
			}, nil
		})
	})
	return &s3ServiceImpl{client: client}
}

func (s *s3ServiceImpl) CreateBucket(ctx context.Context, bucketName string) error {
	_, err := s.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket %s: %w", bucketName, err)
	}
	return nil
}

func (s *s3ServiceImpl) UploadFile(ctx context.Context, bucketName string, key string, data []byte) error {
	err := s.EnsureBucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to bucket %s with key %s: %w", bucketName, key, err)
	}

	return nil
}

func (s *s3ServiceImpl) DownloadFile(ctx context.Context, bucketName string, key string) ([]byte, error) {
	if err := s.EnsureBucketExists(ctx, bucketName); err != nil {
		return nil, err
	}

	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download file from bucket %s with key %s: %w", bucketName, key, err)
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content from bucket %s with key %s: %w", bucketName, key, err)
	}

	return data, nil
}

// ExistBucket checks if a bucket exists.
func (s *s3ServiceImpl) ExistBucket(ctx context.Context, bucketName string) (bool, error) {
	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		// Check if error is because bucket doesn't exist
		var notFound *types.NotFound
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check bucket %s: %w", bucketName, err)
	}

	return true, nil
}

// EnsureBucketExists checks if a bucket exists and creates it if it doesn't.
func (s *s3ServiceImpl) EnsureBucketExists(ctx context.Context, bucketName string) error {
	exists, err := s.ExistBucket(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		if err := s.CreateBucket(ctx, bucketName); err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return nil
}
