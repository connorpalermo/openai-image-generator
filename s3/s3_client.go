package s3

import (
	"bytes"
	"encoding/base64"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// must add environment variables for your access key & secret key
var awsAccessKey = os.Getenv("AWS_ACCESS_KEY")
var awsSecretKey = os.Getenv("AWS_SECRET_KEY")

func initSession() *session.Session {
	creds := credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, "")

	_, err := creds.Get()

	if err != nil {
		log.Fatal(err)
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-2"),
		Credentials: creds,
	},
	)

	return sess
}

func initAWSConnection() *s3.S3 {
	sess := initSession()

	s3Connection := s3.New(sess)

	return s3Connection
}

func initDownloader() *s3manager.Downloader {
	sess := initSession()

	s3Downloader := s3manager.NewDownloader(sess)

	return s3Downloader
}

func Upload(base64File string, objectKey string, bucketName string) error {
	decode, err := base64.StdEncoding.DecodeString(base64File)

	if err != nil {
		return err
	}

	awsSession := initAWSConnection()

	uploadParams := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(decode),
	}

	_, err = awsSession.PutObject(uploadParams)

	return err
}

func Download(item string, bucket string, outputFile string) error {
	downloader := initDownloader()

	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})
	if err != nil {
		log.Fatalf("Unable to download item %q, %v", item, err)
	}

	log.Println("Downloaded", file.Name(), numBytes, "bytes")

	return err
}
