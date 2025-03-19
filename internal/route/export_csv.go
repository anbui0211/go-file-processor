package route

import (
	"context"
	"gofile/internal/wire"

	"github.com/gin-gonic/gin"
)

func exportCsv(r *gin.RouterGroup) {
	g := r.Group("/export-csv")

	// Dependency injection
	handlers, _ := wire.InitExportCsvRouterHandler()
	eController := handlers.Controller
	eWorker := handlers.Worker

	// start worker for export csv
	eWorker.Start(context.Background())

	// route
	g.POST("/create/:type", eController.CreateExportJobHandler)
	g.GET("/get-status/:id", eController.GetExportJobStatusHandler)
	g.GET("/download/:id", eController.DownloadExportFileHandler)
}
