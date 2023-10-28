package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

func uploadToS3(fileHeader *multipart.FileHeader, awsSession *session.Session) (string, error) {
	// Create an S3 service client
	svc := s3.New(awsSession)

	// Specify your S3 bucket name
	bucketName := EnvConfig.AWSBucketName

	// Specify the target location in S3
	key := fmt.Sprintf("hostels/%s%s", generateUniqueFilename(), filepath.Ext(fileHeader.Filename))

	// Open the file for reading
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Upload the file to S3
	outPut, err := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   file,
	})

	if outPut == nil || err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			return "", fmt.Errorf("failed to upload file to S3: %s - %s", awsErr.Code(), awsErr.Message())
		}
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return outPut.String(), nil
}

func generateUniqueFilename() string {
	// Generate a random string (e.g., a timestamp with a random number)
	randomString := uuid.New().String()
	return randomString
}

func uploadImageLocally(file *multipart.FileHeader) (string, error) {
	// Save to media folder with a unique filename
	mediaFolder := EnvConfig.MediaPath

	// Create the media folder if it doesn't exist
	if _, err := os.Stat(mediaFolder); os.IsNotExist(err) {
		err := os.Mkdir(mediaFolder, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	// Create the file path
	filename := fmt.Sprintf("%s%s", generateUniqueFilename(), filepath.Ext(file.Filename))
	filePath := filepath.Join(mediaFolder, filename)

	// Check if the file already exists
	if _, err := os.Stat(filePath); err == nil {
		return "", fmt.Errorf("file %s already exists", filePath)
	}

	dst, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return "", err
	}

	defer dst.Close()

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Copy the file to the destination
	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}
	// Get server address and combine it with filename
	serverAddress := EnvConfig.HTTPServerAddress
	filePath = fmt.Sprintf("%s/%s", serverAddress, filename)

	return filePath, nil
}

func ProcessUploadedImage(Image *multipart.FileHeader, awsSession *session.Session) (string, error) {
	// Decide whether to upload to AWS S3 or save locally
	if EnvConfig.Environment == "production" {
		// In production, upload to AWS S3
		filePath, err := uploadToS3(Image, awsSession)
		if err != nil {
			return "", fmt.Errorf("unable to upload image to AWS S3 %w", err)
		}
		return filePath, nil

	} else {
		// In development, save locally in the media folder
		filePath, err := uploadImageLocally(Image) // saveLocally(fileHeader)
		if err != nil {
			return "", fmt.Errorf("unable to upload image locally, %w", err)
		}
		return filePath, nil
	}
}

func ProcessUploadedImages(Images []*multipart.FileHeader, awsSession *session.Session) ([]string, error) {
	filePaths := []string{}

	for _, Image := range Images {
		// Process each uploaded file here (e.g., save to storage)
		fmt.Printf("Received file: %s\n", Image.Filename)
		filePath, err := ProcessUploadedImage(Image, awsSession)
		if err != nil {
			return nil, fmt.Errorf("unable to process images %w", err)
		}
		filePaths = append(filePaths, filePath)

	}

	return filePaths, nil
}
