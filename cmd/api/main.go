package main

import (
	"context"
	"gofile/internal/controller"
	"gofile/internal/initialize"
	"gofile/internal/repository"
	"gofile/internal/service"
	"gofile/internal/worker"
	paws "gofile/pkg/aws"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	initialize.Run()

	// Router
	g := r.Group("/file-service/api/v1/export-csv")
	{
		// AWS configs
		cfg := paws.LoadAWSConfig()
		sqsService := paws.NewSQSService(cfg)
		s3Service := paws.NewS3Service(cfg)

		repository := repository.NewAccountRepository()
		eService := service.NewExportCsvService(sqsService, s3Service, repository)
		eController := controller.NewExportCsvController(eService)

		// Initialize and start worker
		exportWorker := worker.NewExportCsvWorker(eService)
		exportWorker.Start(context.Background())

		g.POST("/create/:type", eController.CreateExportJobHandler)
		g.GET("/get-status/:id", eController.GetExportJobStatusHandler)
		g.GET("/download/:id", eController.DownloadExportFileHandler)
	}

	r.Run()
}
