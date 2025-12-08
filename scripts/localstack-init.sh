#!/bin/bash

# Wait for LocalStack to be ready
echo "Waiting for LocalStack to be ready..."
sleep 5

# Create S3 bucket for uploads
echo "Creating S3 bucket: portfolio-uploads"
awslocal s3 mb s3://portfolio-uploads

# Set bucket policy for public read access
echo "Setting bucket policy..."
awslocal s3api put-bucket-policy --bucket portfolio-uploads --policy '{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "PublicReadGetObject",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "arn:aws:s3:::portfolio-uploads/*"
    }
  ]
}'

# Enable CORS for the bucket
echo "Enabling CORS..."
awslocal s3api put-bucket-cors --bucket portfolio-uploads --cors-configuration '{
  "CORSRules": [
    {
      "AllowedHeaders": ["*"],
      "AllowedMethods": ["GET", "PUT", "POST", "DELETE", "HEAD"],
      "AllowedOrigins": ["*"],
      "ExposeHeaders": ["ETag"]
    }
  ]
}'

# Note: Secrets should be created via environment variables in production
# This is just for local development testing
echo "LocalStack initialization complete!"
