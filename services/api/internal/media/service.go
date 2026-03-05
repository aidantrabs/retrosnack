package media

import (
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/MobinaToorani/retrosnack/pkg/config"
)

type Service interface {
	Upload(ctx context.Context, productID uuid.UUID, filename string, body io.Reader, size int64) (*Upload, error)
	Delete(ctx context.Context, key string) error
}

type service struct {
	client    *s3.Client
	bucket    string
	publicURL string
}

func NewService(cfg *config.Config) Service {
	r2Endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.R2AccountID)

	awsCfg, _ := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.R2AccessKeyID, cfg.R2SecretAccessKey, ""),
		),
		awsconfig.WithRegion("auto"),
	)

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(r2Endpoint)
	})

	return &service{
		client:    client,
		bucket:    cfg.R2BucketName,
		publicURL: cfg.R2PublicURL,
	}
}

func (s *service) Upload(ctx context.Context, productID uuid.UUID, filename string, body io.Reader, size int64) (*Upload, error) {
	ext := filepath.Ext(filename)
	key := fmt.Sprintf("products/%s/%s%s", productID, uuid.New(), ext)
	mimeType := mime.TypeByExtension(ext)

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          body,
		ContentLength: aws.Int64(size),
		ContentType:   aws.String(mimeType),
		CacheControl:  aws.String(fmt.Sprintf("public, max-age=%d", int(365*24*time.Hour/time.Second))),
	})
	if err != nil {
		return nil, err
	}

	return &Upload{
		Key:      key,
		URL:      fmt.Sprintf("%s/%s", s.publicURL, key),
		MimeType: mimeType,
		Size:     size,
	}, nil
}

func (s *service) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}
