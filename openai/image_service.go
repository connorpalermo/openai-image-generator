package openai

import (
	"context"
	"encoding/base64"
	"net/http"
	aws "openai-image-generator/aws"
	log "openai-image-generator/log"
	"openai-image-generator/model"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
)

type ImageService struct {
	logger   *zap.Logger
	s3Client *aws.S3Client
	client   *openai.Client
}

const (
	status               = "status"
	localDownloadSuccess = "successfully downloaded file to path"
	processImageError    = "failed to process image response"
	localImageSuccess    = "local image created successfully"
)

func NewImageService() *ImageService {
	is := &ImageService{
		logger: log.NewLogger(),
	}
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		is.logger.Fatal("open ai API key not present")
	}
	is.client = openai.NewClient(key)
	is.s3Client = aws.NewS3Client()
	return is
}

func (is *ImageService) SampleImage(c *gin.Context) {
	const (
		samplePrompt   = "Man serves you a hot dog in New York, cartoon style, natural light, high detail"
		sampleFileName = "test.png"
	)
	respData, err := is.imageRequest(samplePrompt)
	if err != nil {
		is.logger.Error("error creating image", zap.Error(err))
		c.IndentedJSON(http.StatusOK, gin.H{status: processImageError})
		return
	}
	err = processB64(respData, sampleFileName)
	if err != nil {
		is.logger.Error(processImageError)
		c.IndentedJSON(http.StatusOK, gin.H{status: processImageError})
		return
	}

	is.logger.Info(localImageSuccess)

	c.IndentedJSON(http.StatusOK, gin.H{status: localDownloadSuccess})
}

func (is *ImageService) GenerateImageLocal(c *gin.Context) {
	var request model.ImageRequestLocal
	bindErr := c.Bind(&request)
	if bindErr != nil {
		is.logger.Error("request format incorrect", zap.Error(bindErr))
		return
	}
	respData, err := is.imageRequest(request.Prompt)
	if err != nil {
		is.logger.Error("error creating image", zap.Error(err))
		c.IndentedJSON(http.StatusOK, gin.H{status: processImageError})
		return
	}
	err = processB64(respData, request.FilePath)
	if err != nil {
		is.logger.Error(processImageError)
		c.IndentedJSON(http.StatusOK, gin.H{status: processImageError})
		return
	}

	is.logger.Info(localImageSuccess)

	c.IndentedJSON(http.StatusOK, gin.H{status: localDownloadSuccess})
}

func (is *ImageService) GenerateImageS3(c *gin.Context) {
	const (
		successfullyPostedToS3 = "successfully saved to S3 bucket"
		failedPostToS3         = "failed to save image to S3 bucket"
	)
	var request model.ImageRequest
	bindErr := c.Bind(&request)
	if bindErr != nil {
		is.logger.Error("request format incorrect", zap.Error(bindErr))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{status: failedPostToS3})
		return
	}
	respData, err := is.imageRequest(request.Prompt)
	if err != nil {
		is.logger.Error("error creating image", zap.Error(err))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{status: failedPostToS3})
		return
	}
	err = is.s3Client.Upload(request, respData)
	if err != nil {
		is.logger.Error(failedPostToS3, zap.Error(err))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{status: failedPostToS3})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{status: successfullyPostedToS3})
}

func (is *ImageService) DownloadImageS3(c *gin.Context) {
	const (
		successfullyDownloadedFromS3 = "successfully downloaded file from S3 bucket"
		failedDownloadFromS3         = "failed to download file from S3 bucket"
	)
	var request model.DownloadImage
	bindErr := c.Bind(&request)
	if bindErr != nil {
		is.logger.Error("request format incorrect", zap.Error(bindErr))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{status: failedDownloadFromS3})
		return
	}

	err := is.s3Client.Download(request)
	if err != nil {
		is.logger.Error("error downloading image", zap.Error(err))
		c.IndentedJSON(http.StatusInternalServerError, gin.H{status: failedDownloadFromS3})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{status: successfullyDownloadedFromS3})
}

func (is *ImageService) imageRequest(prompt string) (resp string, err error) {
	respUrl, err := is.client.CreateImage(
		context.Background(),
		openai.ImageRequest{
			Prompt:         prompt,
			Size:           openai.CreateImageSize512x512,
			ResponseFormat: openai.CreateImageResponseFormatB64JSON,
			N:              1,
		},
	)
	if err != nil {
		return "", err
	}

	return respUrl.Data[0].B64JSON, nil
}

func processB64(b64 string, filePath string) error {
	dec, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(err)
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}

	return nil
}
