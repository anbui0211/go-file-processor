package controller

import (
	"fmt"
	"gofile/internal/service"

	"github.com/gin-gonic/gin"
)

type ExportCsvController struct {
	exportCsvService *service.ExportCsvService
}

func NewExportCsvController(exportCsvService *service.ExportCsvService) *ExportCsvController {
	return &ExportCsvController{
		exportCsvService: exportCsvService,
	}
}

func (ec *ExportCsvController) CreateExportJobHandler(c *gin.Context) {
	var (
		ctx        = c.Request.Context()
		exportType = c.Param("type")
	)

	res, err := ec.exportCsvService.CreateExport(ctx, exportType)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "create export failed",
		})
		return
	}

	c.JSON(200, gin.H{
		"data": res,
	})
	return
}

func (ec *ExportCsvController) GetExportJobStatusHandler(c *gin.Context) {
	ctx := c.Request.Context()
	jobId := c.Param("id")
	status, err := ec.exportCsvService.GetExportStatus(ctx, jobId)
	if err != nil {
		if err.Error() == "job not found" {
			c.JSON(404, gin.H{
				"message": "job not found",
			})
			return
		}
		c.JSON(500, gin.H{
			"message": "get export status failed",
		})
		return
	}

	c.JSON(200, gin.H{
		"status_job": status,
	})
}

func (ec *ExportCsvController) DownloadExportFileHandler(c *gin.Context) {
    ctx := c.Request.Context()
    jobId := c.Param("id")

    fileContent, fileName, err := ec.exportCsvService.DownloadExportFile(ctx, jobId)
    if err != nil {
        c.JSON(500, gin.H{
            "message": "failed to download file",
            "error":   err.Error(),
        })
        return
    }

    // Set headers for file download
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Transfer-Encoding", "binary")
    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
    c.Header("Content-Type", "text/csv")
    c.Header("Content-Length", fmt.Sprintf("%d", len(fileContent)))

    // Send file content
    c.Data(200, "text/csv", fileContent)
}
