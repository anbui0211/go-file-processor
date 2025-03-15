package paws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSService interface {
	CreateQueue(ctx context.Context, queueName string) (string, error)
	SendMessage(ctx context.Context, message SQSMessage, sqsQueueURL string) error
}

type sqsServiceImpl struct {
	client *sqs.Client
}

type SQSMessage struct {
	JobID    string            `json:"job_id"`
	Type     string            `json:"type"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

func NewSQSService(cfg aws.Config) SQSService {
	client := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		o.BaseEndpoint = aws.String("http://localhost:4566")
		o.Credentials = aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			"test",          // Access key ID
			"test",          // Secret access key
			"session-token", // Session token (can be empty for LocalStack)
		))
	})
	return &sqsServiceImpl{client: client}
}

func (s *sqsServiceImpl) CreateQueue(ctx context.Context, queueName string) (string, error) {
	output, err := s.client.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName: &queueName,
	})
	if err != nil {
		return "", err
	}
	return *output.QueueUrl, nil
}

func (s *sqsServiceImpl) SendMessage(ctx context.Context, message SQSMessage, sqsQueueURL string) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	_, err = s.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(sqsQueueURL),
		MessageBody: aws.String(string(messageBytes)),
	})

	return err
}
