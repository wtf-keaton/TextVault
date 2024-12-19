package cloud

import (
	"bytes"
	"context"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

var (
	ErrEntityTooLarge = "EntityTooLarge"
)

type CloudStorage struct {
	S3Client *s3.Client
}

func New() (*CloudStorage, error) {
	ctx := context.Background()

	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(sdkConfig)

	return &CloudStorage{S3Client: s3Client}, nil
}

func (c *CloudStorage) BucketExists(ctx context.Context) error {
	_, err := c.S3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String("textvault"),
	})

	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				log.Printf("Bucket 'textvault' is available.\n")
				return nil
			default:
				log.Printf("Either you don't have access to bucket 'textvault' or another error occurred. "+
					"Here's what happened: %v\n", err)
			}
		}
	}

	return err
}

func (c *CloudStorage) UploadPaste(ctx context.Context, objectKey string, content []byte) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	contentBuffer := bytes.NewBuffer(content)

	_, err := c.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String("textvault"),
		Key:    aws.String(objectKey),
		Body:   contentBuffer,
	})

	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) && apiError.ErrorCode() == ErrEntityTooLarge {
			log.Printf("Error while uploading object to 'textvault'. The object is too large.\n" +
				"To upload objects larger than 5GB, use the S3 console (160GB max)\n" +
				"or the multipart upload API (5TB max).")
		} else {
			log.Printf("Couldn't upload file %v to 'textvault':%v. Here's why: %v\n",
				objectKey, objectKey, err)
		}

	} else {
		err = s3.NewObjectExistsWaiter(c.S3Client).Wait(
			ctx, &s3.HeadObjectInput{
				Bucket: aws.String("textvault"),
				Key:    aws.String(objectKey),
			},
			time.Minute)
		if err != nil {
			log.Printf("Failed attempt to wait for object %s to exist.\n", objectKey)
		}
	}

	return err
}
