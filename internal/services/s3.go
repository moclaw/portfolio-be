package services

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"portfolio-be/internal/config"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
)

type S3Service struct {
	client *s3.S3
	bucket string
	config config.S3Config
}

func NewS3Service(cfg config.S3Config) (*S3Service, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(cfg.Region),
		Endpoint:         aws.String(cfg.Endpoint),
		Credentials:      credentials.NewStaticCredentials(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		S3ForcePathStyle: aws.Bool(cfg.ForcePathStyle),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	client := s3.New(sess)

	// Check if bucket exists first
	_, err = client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(cfg.Bucket),
	})

	if err != nil {
		// Bucket doesn't exist, try to create it
		fmt.Printf("Bucket '%s' doesn't exist, attempting to create it...\n", cfg.Bucket)
		_, createErr := client.CreateBucket(&s3.CreateBucketInput{
			Bucket: aws.String(cfg.Bucket),
		})
		if createErr != nil {
			if strings.Contains(createErr.Error(), "BucketAlreadyOwnedByYou") ||
				strings.Contains(createErr.Error(), "BucketAlreadyExists") {
				fmt.Printf("Bucket '%s' already exists\n", cfg.Bucket)
			} else {
				return nil, fmt.Errorf("failed to create S3 bucket '%s': %w", cfg.Bucket, createErr)
			}
		} else {
			fmt.Printf("Successfully created bucket '%s'\n", cfg.Bucket)
		}
	} else {
		fmt.Printf("Bucket '%s' already exists\n", cfg.Bucket)
	}

	return &S3Service{
		client: client,
		bucket: cfg.Bucket,
		config: cfg,
	}, nil
}

func (s *S3Service) UploadFile(file multipart.File, header *multipart.FileHeader) (string, string, error) {
	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	key := fmt.Sprintf("uploads/%s", fileName)

	// Read file content
	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(file); err != nil {
		return "", "", fmt.Errorf("failed to read file: %w", err)
	}

	// Upload to S3
	_, err := s.client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(header.Header.Get("Content-Type")),
		Metadata: map[string]*string{
			"original-filename": aws.String(header.Filename),
			"upload-time":       aws.String(time.Now().Format(time.RFC3339)),
		},
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Generate URL
	url := fmt.Sprintf("%s/%s/%s", s.config.Endpoint, s.bucket, key)

	return key, url, nil
}

func (s *S3Service) DeleteFile(key string) error {
	_, err := s.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}

	return nil
}

func (s *S3Service) GetFileURL(key string) string {
	return fmt.Sprintf("%s/%s/%s", s.config.Endpoint, s.bucket, key)
}

// GeneratePresignedURL generates a presigned URL for file access with expiration
func (s *S3Service) GeneratePresignedURL(key string, duration time.Duration) (string, error) {
	req, _ := s.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	url, err := req.Presign(duration)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url, nil
}
