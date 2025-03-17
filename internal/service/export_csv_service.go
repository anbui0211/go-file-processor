package service

import (
	"context"
	"fmt"
	"gofile/global"
	"gofile/internal/constant"
	"gofile/internal/repository"
	paws "gofile/pkg/aws"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type ExportCsvService struct {
	sqsService        paws.SQSService
	s3Service         paws.S3Service
	accountRepository *repository.AccountRepository
}

type JobStatus string

const (
	StatusPending    JobStatus = "pending"
	StatusProcessing JobStatus = "processing"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
)

type Job struct {
	ID      string    `json:"id"`
	Status  JobStatus `json:"status"`
	FileURL string    `json:"file_url"`
}

func NewExportCsvService(
	sqsService paws.SQSService,
	s3Service paws.S3Service,
	accountRepository *repository.AccountRepository,
) *ExportCsvService {
	return &ExportCsvService{
		sqsService:        sqsService,
		s3Service:         s3Service,
		accountRepository: accountRepository,
	}
}

func (s *ExportCsvService) CreateExport(ctx context.Context) (*Job, error) {
	jobID := uuid.New().String()
	job := Job{
		ID:     jobID,
		Status: StatusPending,
	}

	// Save job status to Redis
	err := global.Rdb.HSet(ctx, "export_jobs", jobID, string(StatusPending)).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to save job status: %w", err)
	}

	queueUrl, err := s.getQueueURl(ctx, constant.ExportCsvQueueName)
	if err != nil {
		return nil, err
	}

	sqsMessage := paws.SQSMessage{
		JobID: jobID,
		Type:  constant.SqsTypeExportCsv,
	}

	// Gửi message vào SQS queue
	err = s.sqsService.SendMessage(ctx, sqsMessage, queueUrl)
	if err != nil {
		fmt.Println("err", err)
		return nil, fmt.Errorf("failed to send message to queue: %w", err)
	}

	return &job, nil
}

func (s *ExportCsvService) GetExportStatus(ctx context.Context, jobId string) (JobStatus, error) {
	status, err := global.Rdb.HGet(ctx, "export_jobs", jobId).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("job not found")
	}
	if err != nil {
		return "", fmt.Errorf("failed to get job status: %w", err)
	}

	return JobStatus(status), nil
}

func (s *ExportCsvService) getQueueURl(ctx context.Context, queueName string) (string, error) {
	queueUrl, err := global.Rdb.Get(ctx, "export-csv-queue-url").Result()
	if err == redis.Nil {
		// key does not exist, create a new queue
		queueURL, err := s.sqsService.CreateQueue(ctx, queueName)
		if err != nil {
			return "", fmt.Errorf("failed to create queue: %w", err)
		}

		// save to redis with 1 hour expiration
		err = global.Rdb.Set(ctx, "export-csv-queue-url", queueURL, 0).Err()
		return queueURL, nil
	} else if err != nil {

		return "nil", fmt.Errorf("failed to get queue url: %w", err)
	}

	// key exists, return the value retrieved from Redis
	return queueUrl, nil
}
