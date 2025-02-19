package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	// load env. variables
	err := godotenv.Load(".env")

	// Retrieve environment variables
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION")

	bucketName := os.Getenv("BUCKET_NAME")
	folderPath := os.Getenv("FOLDER_PATH")
	customEndpoint := os.Getenv("CUSTOM_ENDPOINT")

	if awsAccessKey == "" || awsSecretKey == "" || awsRegion == "" {
		log.Fatal("Missing AWS credentials or region in environment variables")
	}

	// Manually configure AWS SDK with credentials + custom Linode endpoint
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKey, awsSecretKey, "")),
		config.WithEndpointResolver(aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			if service == s3.ServiceID {
				return aws.Endpoint{
					URL:           customEndpoint, // Use the specified custom endpoint
					SigningRegion: awsRegion,      // Required for authentication
				}, nil
			}
			return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
		})),
	)
	if err != nil {
		log.Fatalf("Unable to load AWS config: %v", err)
	}

	client := s3.NewFromConfig(cfg)

	// Calculate folder size
	totalSize, err := getFolderSize(client, bucketName, folderPath)
	if err != nil {
		log.Fatalf("Error calculating folder size: %v", err)
	}

	fmt.Printf("Total size of folder '%s' in bucket '%s': %d bytes\n", folderPath, bucketName, totalSize)
}

// getFolderSize calculates the total size of objects in a given S3 folder
func getFolderSize(client *s3.Client, bucket, prefix string) (int64, error) {
	var totalSize int64

	// Pagination for listing objects
	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return 0, fmt.Errorf("failed to list objects: %w", err)
		}

		for _, obj := range page.Contents {
			totalSize += *obj.Size
		}
	}

	return totalSize, nil
}
