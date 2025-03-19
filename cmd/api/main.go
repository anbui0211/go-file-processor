package main

import (
	"gofile/internal/initialize"
	"gofile/internal/route"

	"github.com/gin-gonic/gin"
)

func main() {
	e := gin.Default()

	initialize.Run()

	route.Init(e)

	e.Run()
}
