package services

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"portfolio-be/internal/config"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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

	// Ensure bucket exists
	err = ensureBucketExists(client, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure bucket exists: %w", err)
	}

	// Set up CORS configuration for the bucket
	err = setupCORS(client, cfg.Bucket)
	if err != nil {
		fmt.Printf("Warning: Failed to set up CORS for bucket '%s': %v\n", cfg.Bucket, err)
	} else {
		fmt.Printf("CORS configured successfully for bucket '%s'\n", cfg.Bucket)
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
		// Check if it's a bucket-related error and try to recreate bucket
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "NoSuchBucket" {
				fmt.Printf("Bucket missing during upload, attempting to recreate...\n")
				if recreateErr := ensureBucketExists(s.client, s.bucket); recreateErr != nil {
					return "", "", fmt.Errorf("failed to recreate bucket: %w", recreateErr)
				}

				// Retry the upload after recreating bucket
				_, retryErr := s.client.PutObject(&s3.PutObjectInput{
					Bucket:      aws.String(s.bucket),
					Key:         aws.String(key),
					Body:        bytes.NewReader(buf.Bytes()),
					ContentType: aws.String(header.Header.Get("Content-Type")),
					ACL:         aws.String("public-read"),
					Metadata: map[string]*string{
						"original-filename": aws.String(header.Filename),
						"upload-time":       aws.String(time.Now().Format(time.RFC3339)),
					},
				})
				if retryErr != nil {
					return "", "", fmt.Errorf("failed to upload to S3 after bucket recreation: %w", retryErr)
				}
			} else {
				return "", "", fmt.Errorf("failed to upload to S3: %w", err)
			}
		} else {
			return "", "", fmt.Errorf("failed to upload to S3: %w", err)
		}
	}

	url := fmt.Sprintf("%s/%s/%s", "https://media.moclawr.com", s.bucket, key)

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

// setupCORS configures CORS for the S3 bucket to allow frontend access
func setupCORS(client *s3.S3, bucket string) error {
	corsConfig := &s3.PutBucketCorsInput{
		Bucket: aws.String(bucket),
		CORSConfiguration: &s3.CORSConfiguration{
			CORSRules: []*s3.CORSRule{
				{
					AllowedOrigins: []*string{
						aws.String("http://localhost:5300"),
						aws.String("https://moclawr.com"),
					},
					AllowedMethods: []*string{
						aws.String("GET"),
						aws.String("POST"),
						aws.String("PUT"),
						aws.String("DELETE"),
						aws.String("HEAD"),
					},
					AllowedHeaders: []*string{
						aws.String("*"),
					},
					MaxAgeSeconds: aws.Int64(3000),
				},
			},
		},
	}

	_, err := client.PutBucketCors(corsConfig)
	return err
}

// ensureBucketExists checks if a bucket exists and creates it if it doesn't
func ensureBucketExists(client *s3.S3, bucketName string) error {
	// Check if bucket exists
	_, err := client.HeadBucket(&s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		// Check if it's a "NoSuchBucket" error
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "NotFound" || awsErr.Code() == "NoSuchBucket" {
				// Bucket doesn't exist, create it
				fmt.Printf("Bucket '%s' doesn't exist, attempting to create it...\n", bucketName)

				createInput := &s3.CreateBucketInput{
					Bucket: aws.String(bucketName),
				}

				_, createErr := client.CreateBucket(createInput)
				if createErr != nil {
					if awsCreateErr, ok := createErr.(awserr.Error); ok {
						if awsCreateErr.Code() == "BucketAlreadyOwnedByYou" ||
							awsCreateErr.Code() == "BucketAlreadyExists" {
							fmt.Printf("Bucket '%s' already exists\n", bucketName)
							return nil
						}
					}
					return fmt.Errorf("failed to create S3 bucket '%s': %w", bucketName, createErr)
				}

				fmt.Printf("Successfully created bucket '%s'\n", bucketName)

				// Wait a moment for bucket to be ready
				time.Sleep(2 * time.Second)

				// Verify bucket was created successfully
				_, verifyErr := client.HeadBucket(&s3.HeadBucketInput{
					Bucket: aws.String(bucketName),
				})
				if verifyErr != nil {
					return fmt.Errorf("bucket was created but verification failed: %w", verifyErr)
				}

				return nil
			}
		}

		// Other error occurred
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	fmt.Printf("Bucket '%s' already exists\n", bucketName)
	return nil
}
