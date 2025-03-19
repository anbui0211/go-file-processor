//go:build wireinject
// +build wireinject

package wire

import (
	"gofile/internal/controller"
	"gofile/internal/service"
	"gofile/internal/wire/provider"
	"gofile/internal/worker"
	paws "gofile/pkg/aws"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/wire"
)

// Provide AWS Config
func provideAWSConfig() aws.Config {
	return paws.LoadAWSConfig()
}

type ExportHandlers struct {
	Controller *controller.ExportCsvController
	Worker     *worker.ExportCsvWorker
}

func InitExportCsvRouterHandler() (*ExportHandlers, error) {
	wire.Build(
		provideAWSConfig,
		paws.NewSQSService,
		paws.NewS3Service,
		provider.CommonRepositoryProviderSet, // Inject common repositories
		service.NewExportCsvService,
		controller.NewExportCsvController,
		worker.NewExportCsvWorker,
		wire.Struct(new(ExportHandlers), "*"),
	)
	return nil, nil
}
