package images

import (
	"context"
	"encoding/base64"
	"example/openai-image-generator/s3"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

type ImageRequestLocal struct {
	Prompt   string `json:"prompt" binding:"required"`
	FilePath string `json:"filePath" binding:"required"`
}

// must set API key in env variables first
var client = openai.NewClient(os.Getenv("OPENAI_API_KEY"))

func SampleImage(c *gin.Context) {
	respData, err := imageRequest("Man serves you a hot dog in New York, cartoon style, natural light, high detail")
	if err != nil {
		log.Printf("Image creation error: %v\n", err)
		return
	}
	processB64(respData, "test.png")

	c.IndentedJSON(http.StatusOK, gin.H{"status": "Successfully downloaded file to path!"})
}

func GenerateImageLocal(c *gin.Context) {
	var request ImageRequestLocal
	c.Bind(&request)
	respData, err := imageRequest(request.Prompt)
	if err != nil {
		log.Printf("Image creation error: %v\n", err)
		return
	}
	processB64(respData, request.FilePath)

	c.IndentedJSON(http.StatusOK, gin.H{"status": "Successfully downloaded file to path!"})
}

func GenerateImageS3(c *gin.Context) {
	var request s3.ImageRequest
	c.Bind(&request)
	respData, err := imageRequest(request.Prompt)
	if err != nil {
		log.Printf("Image creation error: %v\n", err)
		return
	}
	request.Upload(respData)

	c.IndentedJSON(http.StatusOK, gin.H{"status": "Successfully saved to S3 bucket!"})
}

func DownloadImageS3(c *gin.Context) {
	var request s3.DownloadImage
	c.Bind(&request)
	log.Print(request.FilePath)
	request.Download()

	c.IndentedJSON(http.StatusOK, gin.H{"status": "Successfully downloaded file from S3 bucket!"})
}

func processB64(b64 string, filePath string) {
	dec, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		panic(err)
	}
	if err := f.Sync(); err != nil {
		panic(err)
	}
}

func imageRequest(prompt string) (resp string, err error) {
	respUrl, err := client.CreateImage(
		context.Background(),
		openai.ImageRequest{
			Prompt:         prompt,
			Size:           openai.CreateImageSize512x512,
			ResponseFormat: openai.CreateImageResponseFormatB64JSON,
			N:              1,
		},
	)
	if err != nil {
		log.Printf("Image creation error: %v\n", err)
		return "", err
	}

	return respUrl.Data[0].B64JSON, nil
}
