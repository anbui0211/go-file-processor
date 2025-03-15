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
	ctx := c.Request.Context()
	res, err := ec.exportCsvService.CreateExport(ctx)
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
