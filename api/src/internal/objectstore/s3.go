package objectstore

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"com.github/w-k-s/glassdoor-hr-review-detector/pkg/services"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3 struct {
	client *s3.S3
}

func MustS3(region string) services.ObjectStoreService {
	if region == "" {
		log.Fatal("aws region not provided")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	return &S3{
		client: s3.New(sess),
	}
}

func (s *S3) Get(bucket string, object string) (io.Reader, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(object),
	}

	result, err := s.client.GetObject(input)
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return result.Body, nil
}

func (s *S3) Put(bucket string, object string, content string) error {
	input := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(object),
		Body:   bytes.NewReader([]byte(content)),
	}

	_, err := s.client.PutObject(input)
	if err != nil {
		return fmt.Errorf("failed to put object: %w", err)
	}

	return nil
}
