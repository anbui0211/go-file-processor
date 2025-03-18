package worker

import (
	"context"
	"fmt"
	"gofile/internal/service"
	"log"
	"time"
)

type ExportCsvWorker struct {
	exportService *service.ExportCsvService
}

func NewExportCsvWorker(exportService *service.ExportCsvService) *ExportCsvWorker {
	return &ExportCsvWorker{
		exportService: exportService,
	}
}

func (w *ExportCsvWorker) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Export CSV worker stopped")
				return
			case <-ticker.C:
				if err := w.exportService.ProcessMessages(ctx); err != nil {
					fmt.Printf("Error processing messages: %v\n", err)
				}
			}
		}
	}()
	log.Println("Export CSV worker started")
}
