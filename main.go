package main

import (
	"example/openai-image-generator/images"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/sample", images.SampleImage)
	router.POST("/createImageLocal", images.GenerateImageLocal)
	router.POST("/createImageS3", images.GenerateImageS3)
	router.POST("/downloadImageS3", images.DownloadImageS3)
	router.Run("localhost:8080")
}
