package main

import (
	"openai-image-generator/openai"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	is := openai.NewImageService()
	router.GET("/sample", is.SampleImage)
	router.POST("/createImageLocal", is.GenerateImageLocal)
	router.POST("/createImageS3", is.GenerateImageS3)
	router.POST("/downloadImageS3", is.DownloadImageS3)
	router.Run("localhost:8080")
}
