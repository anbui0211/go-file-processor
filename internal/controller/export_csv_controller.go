package controller

import (
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
