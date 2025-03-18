package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"gofile/global"
	"gofile/internal/constant"
	"gofile/internal/repository"
	paws "gofile/pkg/aws"
	"io"
	"log"
	"os"
	"time"

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

type ExportType string

const (
	ExportTypeAccount ExportType = "account"
	ExportTypeUser    ExportType = "user"
)

type Job struct {
	ID         string     `json:"id"`
	Status     JobStatus  `json:"status"`
	FileURL    string     `json:"file_url"`
	ExportType ExportType `json:"export_type"`
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

func (s *ExportCsvService) CreateExport(ctx context.Context, exportType string) (*Job, error) {
	jobID := uuid.New().String()
	job := Job{
		ID:         jobID,
		Status:     StatusPending,
		ExportType: ExportType(exportType),
	}

	// Save job fields separately
	err := global.Rdb.HMSet(ctx, "export_jobs:"+jobID,
		"status", string(StatusPending),
		"export_type", exportType,
	).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to save job: %w", err)
	}

	queueUrl, err := s.getQueueURl(ctx, constant.ExportCsvQueueName)
	if err != nil {
		return nil, err
	}

	sqsMessage := paws.SQSMessage{
		JobID:      jobID,
		ExportType: string(exportType),
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
	status, err := global.Rdb.HGet(ctx, "export_jobs:"+jobId, "status").Result()
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

func (s *ExportCsvService) ProcessMessages(ctx context.Context) error {
	queueUrl, err := s.getQueueURl(ctx, constant.ExportCsvQueueName)
	if err != nil {
		return err
	}

	messages, err := s.sqsService.GetQueueMessage(ctx, queueUrl)
	if err != nil {
		return err
	}

	for _, msg := range messages {
		// Update status to processing
		log.Println("Job: " + msg.JobID + " processing")
		if err = global.Rdb.HSet(ctx, "export_jobs:"+msg.JobID, "status", string(StatusProcessing)).Err(); err != nil {
			return err
		}

		// Process the export CSV logic
		exportType := ExportType(msg.ExportType)
		if err = s.processExport(ctx, exportType); err != nil {
			global.Rdb.HSet(ctx, "export_jobs:"+msg.JobID, "status", string(StatusFailed))
			log.Println("Job: " + msg.JobID + " failed to process")
			continue
		}

		// Update status to completed
		if err = global.Rdb.HSet(ctx, "export_jobs:"+msg.JobID, "status", string(StatusCompleted)).Err(); err != nil {
			return err
		}
		log.Println("Job: " + msg.JobID + " processed successfully")

		// Delete message from queue after successful processing
		if err := s.sqsService.DeleteMessage(ctx, msg, queueUrl); err != nil {
			return err
		}

	}

	return nil
}

func (s *ExportCsvService) processExport(ctx context.Context, exportType ExportType) error {
	switch exportType {
	case ExportTypeAccount:
		return s.processAccountExport(ctx)
	case ExportTypeUser:
		return fmt.Errorf("chức năng export user chưa được implement")
	default:
		return fmt.Errorf("không hỗ trợ loại export: %s", exportType)
	}
}

func (s *ExportCsvService) processAccountExport(ctx context.Context) error {
	accounts, err := s.accountRepository.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get accounts: %w", err)
	}

	tmpFile, err := os.CreateTemp("", "account-*.csv")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// create CSV writer
	writer := csv.NewWriter(tmpFile)
	defer writer.Flush()

	// write header row
	header := []string{"Code", "Name", "Type", "CreatedAt", "UpdatedAt"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write account data
	for _, account := range accounts {
		row := []string{
			account.Code,
			account.Name,
			account.Type,
			account.CreatedAt.Format(time.RFC3339), // (YYYY-MM-DDTHH:MM:SSZ)
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	// Reset file pointer to beginning
	if _, err := tmpFile.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}

	// Read file content
	data, err := io.ReadAll(tmpFile)
	if err != nil {
		return fmt.Errorf("failed to read temporary file: %w", err)
	}

	// Generate S3 key
	fileName := fmt.Sprintf("export-file/accounts-%s.csv", time.Now().Format("20060102-150405"))

	// Upload file to S3
	bucketName := "go-bucket"
	if err := s.s3Service.UploadFile(ctx, bucketName, fileName, data); err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return nil
}
