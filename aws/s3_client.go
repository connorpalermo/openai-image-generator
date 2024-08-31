package aws

import (
	"bytes"
	"encoding/base64"
	log "openai-image-generator/log"
	"openai-image-generator/model"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.uber.org/zap"
)

type S3Client struct {
	logger       *zap.Logger
	s3session    *s3.S3
	s3downloader *s3manager.Downloader
}

// must add environment variables for your access key & secret key
var (
	awsAccessKey = os.Getenv("AWS_ACCESS_KEY")
	awsSecretKey = os.Getenv("AWS_SECRET_KEY")
)

func NewS3Client() *S3Client {
	s3Client := &S3Client{
		logger: log.NewLogger(),
	}
	s3Client.Initialize()
	return s3Client
}

// allow for re-initializiation if needed
func (s3Client *S3Client) Initialize() error {
	session, err := initSession()
	if err != nil {
		s3Client.logger.Fatal("failed to initialize sesssion with error", zap.Error(err))
	}

	s3Session := s3.New(session)
	s3Downloader := s3manager.NewDownloader(session)

	s3Client.s3session = s3Session
	s3Client.s3downloader = s3Downloader

	s3Client.logger.Info("successfully initialized s3client")
	return nil
}

func (s3Client *S3Client) Upload(imageRequest model.ImageRequest, base64File string) error {
	decode, err := base64.StdEncoding.DecodeString(base64File)

	if err != nil {
		return err
	}

	uploadParams := &s3.PutObjectInput{
		Bucket: aws.String(imageRequest.BucketName),
		Key:    aws.String(imageRequest.FileName),
		Body:   bytes.NewReader(decode),
	}

	_, err = s3Client.s3session.PutObject(uploadParams)
	if err != nil {
		return err
	}

	return nil
}

func (s3Client *S3Client) Download(d model.DownloadImage) error {

	file, err := os.Create(d.FilePath)
	if err != nil {
		return err
	}

	defer file.Close()
	_, err = s3Client.s3downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(d.Bucket),
			Key:    aws.String(d.Item),
		})
	if err != nil {
		return err
	}

	return nil
}

func initSession() (*session.Session, error) {
	const defaultRegion = "us-east-2"
	creds := credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, "")

	_, err := creds.Get()

	if err != nil {
		return nil, err
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(defaultRegion),
		Credentials: creds,
	},
	)
	if err != nil {
		return nil, err
	}

	return sess, nil
}
