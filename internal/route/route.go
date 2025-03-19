package route

import (
	"github.com/gin-gonic/gin"
)

func Init(e *gin.Engine) {
	r := e.Group("/file-service/api/v1")

	exportCsv(r)
}
