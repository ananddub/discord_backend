package util

import (
	"context"
	"discord/config"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var minioClient *minio.Client

func MinioClient() (*minio.Client, error) {
	if minioClient != nil {
		return minioClient, nil
	}
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	endpoint := cfg.S3.Endpoint
	accessKeyID := cfg.S3.AccessKey
	secretAccessKey := cfg.S3.SecretKey
	useSSL := cfg.S3.UseSSL

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	minioClient = client
	return client, nil
}

func InitMinio() {
	client, err := MinioClient()
	if err != nil {
		log.Fatalln(err)
	}

	buckets := []string{"discord", "uploads", "profile", "mesages"}

	for _, b := range buckets {
		ensureBucket(client, b)
	}
}

func GenerateUploadURL(objectName string, imageType string) (string, error) {
	client, err := MinioClient()
	if err != nil {
		return "", err
	}

	ctx := context.Background()

	reqParams := make(url.Values)
	reqParams.Set("Content-Type", imageType)

	presignedURL, err := client.PresignedPutObject(
		ctx,
		"discord",
		objectName,
		time.Minute*15,
	)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

func ensureBucket(client *minio.Client, bucketName string) error {
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exists {
		return client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	}

	return nil
}
