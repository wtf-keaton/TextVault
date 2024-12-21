package cloud

import (
	cfg "TextVault/internal/config"
	"TextVault/internal/lib/log/sl"
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

var (
	ErrEntityTooLarge = "EntityTooLarge"
)

type Storage struct {
	S3Client   *s3.Client
	log        *slog.Logger
	bucketName string
}

func New(log *slog.Logger, s3config cfg.S3Config) (*Storage, error) {
	ctx := context.Background()

	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(sdkConfig)

	return &Storage{
		S3Client:   s3Client,
		log:        log,
		bucketName: s3config.BucketName,
	}, nil
}

func (c *Storage) BucketExists(ctx context.Context) error {
	_, err := c.S3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(c.bucketName),
	})

	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				c.log.Error(fmt.Sprintf("Bucket '%s' is not available.", c.bucketName), sl.Err(err))
				return nil
			default:
				c.log.Error(fmt.Sprintf("Either you don't have access to bucket '%s' or another error occurred.", c.bucketName), sl.Err(err))
			}
		}
	}

	return err
}

func (c *Storage) UploadPaste(ctx context.Context, objectKey string, content []byte) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	contentBuffer := bytes.NewBuffer(content)

	_, err := c.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(objectKey),
		Body:   contentBuffer,
	})

	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) && apiError.ErrorCode() == ErrEntityTooLarge {
			c.log.Error(fmt.Sprintf("Error while uploading object to '%v'. The object is too large."+
				"To upload objects larger than 5GB, use the S3 console (160GB max)"+
				"or the multipart upload API (5TB max).", c.bucketName))
		} else {
			c.log.Error(fmt.Sprintf("Couldn't upload file %v to '%v':%v",
				c.bucketName, objectKey, objectKey), sl.Err(err))
		}

	} else {
		err = s3.NewObjectExistsWaiter(c.S3Client).Wait(
			ctx, &s3.HeadObjectInput{
				Bucket: aws.String(c.bucketName),
				Key:    aws.String(objectKey),
			},
			time.Minute)
		if err != nil {
			c.log.Error(fmt.Sprintf("Failed attempt to wait for object %s to exist.", objectKey), sl.Err(err))
		}
	}

	return err
}

func (c *Storage) GetPasteContent(ctx context.Context, objectKey string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var partMiBs int64 = 10
	downloader := manager.NewDownloader(c.S3Client, func(d *manager.Downloader) {
		d.PartSize = partMiBs * 1024 * 1024
	})

	buffer := manager.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(ctx, buffer, &s3.GetObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		c.log.Error(fmt.Sprintf("Couldn't download large object from %v:%v",
			c.bucketName, objectKey), sl.Err(err))
	}
	return buffer.Bytes(), err
}

func (c *Storage) DeletePaste(ctx context.Context, objectKey string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := c.S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		return err
	}

	return nil
}
