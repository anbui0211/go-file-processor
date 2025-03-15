package paws

import (
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service interface {
    CreateBucket(bucketName string) error
    UploadFile(bucketName string, key string, data []byte) error
    DownloadFile(bucketName string, key string) ([]byte, error)
}

type s3ServiceImpl struct {
    client *s3.Client
}

func NewS3Service(cfg aws.Config) S3Service {
    client := s3.NewFromConfig(cfg, func(o *s3.Options) {
        o.UsePathStyle = true
        o.BaseEndpoint = aws.String("http://localhost:4566")
    })
    return &s3ServiceImpl{client: client}
}

func (s *s3ServiceImpl) CreateBucket(bucketName string) error {
    _, err := s.client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
        Bucket: aws.String(bucketName),
    })
    return err
}

func (s *s3ServiceImpl) UploadFile(bucketName string, key string, data []byte) error {
    _, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(key),
        Body:   bytes.NewReader(data),
    })
    return err
}

func (s *s3ServiceImpl) DownloadFile(bucketName string, key string) ([]byte, error) {
    result, err := s.client.GetObject(context.TODO(), &s3.GetObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(key),
    })
    if err != nil {
        return nil, err
    }
    defer result.Body.Close()

    return io.ReadAll(result.Body)
}
