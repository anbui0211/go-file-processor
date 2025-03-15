package main

import (
	"gofile/internal/controller"
	"gofile/internal/initialize"
	"gofile/internal/repository"
	"gofile/internal/service"
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

		g.POST("/create", eController.CreateExportJobHandler)
	}

	r.Run()
}
