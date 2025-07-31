package s3

import (
	"bytes"
	"context"
	"github.com/alioth-center/infrastructure/storage"
	"github.com/alioth-center/infrastructure/trace"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
)

type s3Storage struct {
	cli  *s3.Client
	conf *Config
}

func NewS3Storage(conf *Config) storage.Storage {
	awsConf, loadErr := config.LoadDefaultConfig(
		trace.Background(),
		config.WithCredentialsProvider(aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(conf.AccessKey, conf.SecretKey, ""))),
		config.WithRegion(conf.Region),
	)
	if loadErr != nil {
		panic(loadErr)
	}
	client := s3.NewFromConfig(awsConf, func(options *s3.Options) { options.BaseEndpoint = &conf.Endpoint })

	return &s3Storage{cli: client, conf: conf}
}

func (s *s3Storage) LoadFile(ctx context.Context, key string) (content *bytes.Buffer, err error) {
	getResult, getErr := s.cli.GetObject(ctx, &s3.GetObjectInput{Bucket: &s.conf.Bucket, Key: &key})
	if getErr != nil {
		return nil, getErr
	}

	content = &bytes.Buffer{}
	defer func() { _ = getResult.Body.Close() }()
	if _, copyErr := io.Copy(content, getResult.Body); copyErr != nil {
		return nil, copyErr
	}

	return content, nil
}

func (s *s3Storage) StoreFile(ctx context.Context, key string, content *bytes.Buffer) (err error) {
	uploader := manager.NewUploader(s.cli, func(u *manager.Uploader) { u.PartSize, u.Concurrency = s.conf.PartSize, s.conf.Concurrency })
	if _, uploadErr := uploader.Upload(ctx, &s3.PutObjectInput{Bucket: &s.conf.Bucket, Key: &key, Body: content}); uploadErr != nil {
		return uploadErr
	}

	return nil
}

func (s *s3Storage) DeleteFile(ctx context.Context, key string) (err error) {
	if _, deleteErr := s.cli.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: &s.conf.Bucket, Key: &key}); deleteErr != nil {
		return deleteErr
	}

	return nil
}
