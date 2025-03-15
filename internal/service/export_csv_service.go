package service

import (
	"context"
	"fmt"
	"gofile/internal/repository"
	paws "gofile/pkg/aws"

	"github.com/google/uuid"
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
	// Tạo một job mới với ID ngẫu nhiên
	jobID := uuid.New().String()
	job := Job{
		ID:     jobID,
		Status: StatusPending,
	}

	queueName := "export-csv-queue"
	queueURL, err := s.sqsService.CreateQueue(ctx, queueName)
	if err != nil {
		return nil, fmt.Errorf("failed to create queue: %w", err)
	}

	sqsMessage := paws.SQSMessage{
		JobID: jobID,
		Type:  "export_csv",
		Metadata: map[string]string{
			"file_type": "csv",
			"entity":    "accounts",
		},
	}

	// Gửi message vào SQS queue
	err = s.sqsService.SendMessage(ctx, sqsMessage, queueURL)
	if err != nil {
		return nil, fmt.Errorf("failed to send message to queue: %w", err)
	}

	return &job, nil
}
