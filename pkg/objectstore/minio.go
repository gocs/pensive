package objectstore

import (
	"context"
	"fmt"
	"io"
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
	mc                     *minio.Client
	PresignedURLExpiration time.Duration
}

type MakeBucketOptions minio.MakeBucketOptions

func (ostore *ObjectStore) MakeBucket(ctx context.Context, bucketName string, opts MakeBucketOptions) error {
	err := ostore.mc.MakeBucket(ctx, bucketName, minio.MakeBucketOptions(opts))
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := ostore.mc.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			return err
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	return nil
}

func (ostore *ObjectStore) ListBuckets(ctx context.Context) ([]minio.BucketInfo, error) {
	return ostore.mc.ListBuckets(ctx)
}

func (ostore *ObjectStore) ListObjects(ctx context.Context, bucketName string, opts ListObjectsOptions) <-chan minio.ObjectInfo {
	return ostore.mc.ListObjects(ctx, bucketName, minio.ListObjectsOptions(opts))
}

type ListObjectsOptions minio.ListObjectsOptions

func (ostore *ObjectStore) ListAllBucketsObjects(ctx context.Context, opts ListObjectsOptions) ([]minio.ObjectInfo, error) {
	buckets, err := ostore.mc.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	objs := []minio.ObjectInfo{}
	for _, bucket := range buckets {
		objectCh := ostore.mc.ListObjects(ctx, bucket.Name, minio.ListObjectsOptions(opts))
		for object := range objectCh {
			if object.Err != nil {
				log.Println(object.Err)
				return nil, object.Err
			}
			objs = append(objs, object)
		}
	}
	return objs, nil
}

type PresignedGetObjectOptions struct {
	ReqParams              url.Values
	PresignedURLExpiration time.Duration
}

// PresignedGetObject gets presigned URL for object
// if presignedURLExpiration is not set defaults to a day
func (ostore *ObjectStore) GetPresignedURLObject(ctx context.Context, bucketName, filename string, opts PresignedGetObjectOptions) (*url.URL, error) {
	ostore.PresignedURLExpiration = opts.PresignedURLExpiration
	if ostore.PresignedURLExpiration == 0 {
		ostore.PresignedURLExpiration = time.Second * 24 * 60 * 60
	}

	rp := make(url.Values)
	rp.Set("response-content-disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	opts.ReqParams = rp

	return ostore.mc.PresignedGetObject(ctx, bucketName, filename, ostore.PresignedURLExpiration, opts.ReqParams)
}

// GetObject gets object to be written
func (ostore *ObjectStore) GetObject(ctx context.Context, bucketName, filename string) (*minio.Object, error) {
	return ostore.mc.GetObject(ctx, bucketName, filename, minio.GetObjectOptions{})
}

type PutObjectOptions minio.PutObjectOptions

// SaveObjectFromPath saves file from local disk to cloud
// if content type is not set this will try to detect and set what it detects
func (ostore *ObjectStore) SaveObjectFromPath(ctx context.Context, bucketName, objectName, filePath string, opts PutObjectOptions) (minio.UploadInfo, error) {
	return ostore.mc.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions(opts))
}

// SaveObject saves file from a reader interface to cloud
// if content type is not set this will try to detect and set what it detects
func (ostore *ObjectStore) SaveObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, opts PutObjectOptions) (minio.UploadInfo, error) {
	return ostore.mc.PutObject(ctx, bucketName, objectName, reader, objectSize, minio.PutObjectOptions(opts))
}
