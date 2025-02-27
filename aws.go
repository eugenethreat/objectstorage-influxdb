package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

func init_aws_client() *s3.Client {
	// load env. variables
	err := godotenv.Load(".env")

	// Retrieve environment variables
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsRegion := os.Getenv("AWS_REGION")

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

	return client

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

		// mod to only check for list_clients
		for _, obj := range page.Contents {
			totalSize += *obj.Size
		}
	}

	return totalSize, nil
}

// // getListClientsFiles
// func getListClientsFiles(client *s3.Client, bucket, prefix string) (int64, error) {
// 	var totalSize int64

// 	// Pagination for listing objects
// 	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
// 		Bucket: aws.String(bucket),
// 		Prefix: aws.String(prefix),
// 	})

// 	for paginator.HasMorePages() {
// 		page, err := paginator.NextPage(context.TODO())
// 		if err != nil {
// 			return 0, fmt.Errorf("failed to list objects: %w", err)
// 		}

// 		// mod to only check for list_clients
// 		for _, obj := range page.Contents {
// 			obj_key_str := fmt.Sprintf("%v", *obj.Key)
// 			// fmt.Println(obj_key_str)

// 			if strings.Contains(obj_key_str, "list_clients") {
// 				// totalSize += *obj.Size
// 				// download the file if its the right file :)

// 			}
// 		}
// 	}

// 	return totalSize, nil
// }

// getListClientsFiles reads "list_clients" files from S3 for all days in 2024 without downloading them.
func getListClientsFiles(client *s3.Client, bucket string, random string) (int64, error) {
	var totalSize int64

	// Define the date range for 2024
	startDate, _ := time.Parse("2006/01/02", "2024/01/01")
	endDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC) // Start of 2025

	// Iterate over each day in 2024
	for currentDate := startDate; currentDate.Before(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		prefix := currentDate.Format("2006/01/02") // "YYYY/MM/DD" format

		// List objects with the current prefix
		paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
			Prefix: aws.String(prefix),
		})

		for paginator.HasMorePages() {
			page, err := paginator.NextPage(context.TODO())
			if err != nil {
				return 0, fmt.Errorf("failed to list objects for prefix %s: %w", prefix, err)
			}

			// Process each object in the listed results
			for _, obj := range page.Contents {
				objKey := *obj.Key

				// Check if the object contains "list_clients"
				if strings.Contains(objKey, "list_clients") {
					fmt.Println("Processing:", objKey)
					totalSize += *obj.Size

					// Read file content without downloading
					err := readObjectContent(client, bucket, objKey)
					if err != nil {
						fmt.Printf("Error reading file %s: %v\n", objKey, err)
					}
				}
			}
		}
	}

	return totalSize, nil
}

func readObjectContent(client *s3.Client, bucket, key string) error {
	resp, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to get object %s: %w", key, err)
	}
	defer resp.Body.Close()

	// Use a buffered reader to handle large files
	reader := bufio.NewReader(resp.Body)
	buffer := make([]byte, 1024*64) // 64 KB buffer size

	for {
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("error reading file %s: %w", key, err)
		}
		if n == 0 {
			break
		}
		fmt.Print(string(buffer[:n])) // Process chunk
	}

	return nil
}
