package objectstore

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
}

func New(cfg Config) (*ObjectStore, error) {
	// Initialize minio client object.
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		// Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &ObjectStore{mc: minioClient}, nil
}

type ObjectStore struct {
	mc *minio.Client
}

func (os *ObjectStore) MakeBucket(ctx context.Context, bucketName, location string) error {
	err := os.mc.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := os.mc.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	return nil
}

func (os *ObjectStore) ListBuckets(ctx context.Context) ([]minio.BucketInfo, error) {
	return os.mc.ListBuckets(ctx)
}

func (os *ObjectStore) ListObjects(ctx context.Context, bucketName string, opts minio.ListObjectsOptions) <-chan minio.ObjectInfo {
	return os.mc.ListObjects(ctx, bucketName, opts)
}

func (os *ObjectStore) PresignedGetObject(ctx context.Context, bucketName, filename string) (*url.URL, error) {
	reqParams := url.Values{}
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	return os.mc.PresignedGetObject(ctx, bucketName, filename, time.Second*24*60*60, reqParams)
}
