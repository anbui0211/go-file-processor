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
	GetQueueMessage(ctx context.Context, queueUrl string) ([]SQSMessage, error)
	DeleteMessage(ctx context.Context, message SQSMessage, sqsQueueURL string) error
}

type sqsServiceImpl struct {
	client *sqs.Client
}

type SQSMessage struct {
	JobID         string            `json:"job_id"`
	ExportType    string            `json:"export_type"`
	ReceiptHandle string            `json:"receipt_handle"`
	Metadata      map[string]string `json:"metadata,omitempty"`
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
	// Check if queue exists first
	if !s.queueExists(ctx, sqsQueueURL) {
		return fmt.Errorf("queue does not exist: %s", sqsQueueURL)
	}

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

func (s *sqsServiceImpl) DeleteMessage(ctx context.Context, message SQSMessage, sqsQueueURL string) error {
	if message.ReceiptHandle == "" {
		return fmt.Errorf("receipt handle is required to delete message")
	}

	_, err := s.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(sqsQueueURL),
		ReceiptHandle: aws.String(message.ReceiptHandle),
	})

	return err
}

func (s *sqsServiceImpl) GetQueueMessage(ctx context.Context, queueUrl string) ([]SQSMessage, error) {
	output, err := s.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueUrl),
		MaxNumberOfMessages: *aws.Int32(10),
		WaitTimeSeconds:     *aws.Int32(10),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to receive messages: %w", err)
	}

	var messages []SQSMessage
	for _, msg := range output.Messages {
		var sqsMsg SQSMessage
		if err := json.Unmarshal([]byte(*msg.Body), &sqsMsg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message: %w", err)
		}

		sqsMsg.ReceiptHandle = *msg.ReceiptHandle
		messages = append(messages, sqsMsg)
	}

	return messages, nil
}

func (s *sqsServiceImpl) queueExists(ctx context.Context, queueUrl string) bool {
	// List all queues
	result, err := s.client.ListQueues(ctx, &sqs.ListQueuesInput{})
	if err != nil {
		return false
	}

	// Check if the queue URL exists in the list
	for _, url := range result.QueueUrls {
		if url == queueUrl {
			return true
		}
	}
	return false
}

func (s *sqsServiceImpl) queueExistsV2(ctx context.Context, queueName string) bool {
	_, err := s.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	return err == nil
}
