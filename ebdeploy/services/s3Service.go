package services

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
)

type S3Service struct {
	client *s3.S3
}

func NewS3Service(client *s3.S3) *S3Service {
	return &S3Service{
		client: client,
	}
}

func (svc *S3Service) BucketExists(bucket string) (bool, error) {
	if output, err := svc.client.ListBuckets(new(s3.ListBucketsInput)); err == nil {
		for _, b := range output.Buckets {
			if *b.Name == bucket {
				return true, nil
			}
		}
		return false, nil
	} else {
		return false, err
	}
}

func (svc *S3Service) CreateBucket(bucket string) error {
	_, err := svc.client.CreateBucket(&s3.CreateBucketInput{Bucket: &bucket})
	return err
}

func (svc *S3Service) UploadFile(file string, bucket string, key string) error {
	if reader, err := os.Open(file); err == nil {
		uploader := s3manager.NewUploader(&s3manager.UploadOptions{S3: svc.client})

		input := &s3manager.UploadInput{
			Bucket: &bucket,
			Body:   reader,
			Key:    &key,
		}

		_, uploadErr := uploader.Upload(input)

		return uploadErr
	} else {
		return err
	}
}
