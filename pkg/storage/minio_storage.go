package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
	"github.com/pkg/errors"
)

var _ Storage = (*MinioStorage)(nil)

// MinioStorage is the minio implementation of Storage.
type MinioStorage struct {
	cli      *minio.Client
	endpoint string
	secure   bool
}

// NewMinioStorage prepare a new minio storage client.
func NewMinioStorage(endpoint, accessKeyID, secretAccessKey string, secure bool) (*MinioStorage, error) {
	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: secure,
	})
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to minio")
	}
	st := MinioStorage{
		cli:      cli,
		endpoint: endpoint,
		secure:   secure,
	}
	return &st, nil
}

// NewBucket implements Storage.NewBucket
func (st *MinioStorage) NewBucket(ctx context.Context, bucketName string) error {
	return st.cli.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
}

// BucketExists implements Storage.BucketExists
func (st *MinioStorage) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	return st.cli.BucketExists(ctx, bucketName)
}

// SetBucketStatus implements Storage.SetBucketStatus
func (st *MinioStorage) SetBucketStatus(ctx context.Context, bucketName string, isPublic bool) error {
	var policy string
	if isPublic {
		policy = generateDownloadPolicy(bucketName)
	}
	err := st.cli.SetBucketPolicy(ctx, bucketName, policy)
	if err != nil {
		minioError := minio.ToErrorResponse(err)
		if minioError.StatusCode == http.StatusNotFound {
			return ErrBucketNotExists
		}
	}
	return err
}

// SetBucketLifecycle implements Storage.SetBucketLifecycle
func (st *MinioStorage) SetBucketLifecycle(ctx context.Context, bucketName string, days int) error {
	conf := lifecycle.NewConfiguration()
	conf.Rules = []lifecycle.Rule{
		{
			ID:     "expire-bucket",
			Status: "Enabled",
			Expiration: lifecycle.Expiration{
				Days: lifecycle.ExpirationDays(days),
			},
		},
	}
	err := st.cli.SetBucketLifecycle(ctx, bucketName, conf)
	if err != nil {
		minioError := minio.ToErrorResponse(err)
		if minioError.StatusCode == http.StatusNotFound {
			return ErrBucketNotExists
		}
	}
	return err
}

// DeleteBucket implements Storage.DeleteBucket
func (st *MinioStorage) DeleteBucket(ctx context.Context, bucketName string) error {
	err := st.cli.RemoveBucketWithOptions(ctx, bucketName, minio.BucketOptions{ForceDelete: true})
	if err != nil {
		minioError := minio.ToErrorResponse(err)
		if minioError.StatusCode == http.StatusNotFound {
			return ErrBucketNotExists
		}
	}
	return err
}

func generateDownloadPolicy(bucketName string) string {
	tmpl := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {
					"AWS": [
						"*"
					]
				},
				"Action": [
					"s3:GetObject"
				],
				"Resource": [
					"arn:aws:s3:::%s/*"
				]
			}
		]
	}`
	return fmt.Sprintf(tmpl, bucketName)
}

// Put implements Storage.Put
func (st *MinioStorage) Put(ctx context.Context, bucketName, objName string, r io.Reader, size int64, contentType string) error {
	_, err := st.cli.PutObject(ctx, bucketName, objName, r, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// CopyObject implements Storage.CopyObject
func (st *MinioStorage) CopyObject(ctx context.Context, srcBucket, srcObject, dstBucket, dstObject string) error {
	src := minio.CopySrcOptions{
		Bucket: srcBucket,
		Object: srcObject,
	}
	dst := minio.CopyDestOptions{
		Bucket: dstBucket,
		Object: dstObject,
	}
	_, err := st.cli.CopyObject(ctx, dst, src)
	if err != nil {
		minioError := minio.ToErrorResponse(err)
		if minioError.Code == "NoSuchKey" {
			return ErrObjectNotExists
		}
	}
	return err
}

// ObjectExists implements Storage.ObjectExists
func (st *MinioStorage) ObjectExists(ctx context.Context, bucketName, objName string) (bool, error) {
	_, err := st.cli.StatObject(ctx, bucketName, objName, minio.StatObjectOptions{})
	if err != nil {
		minioError := minio.ToErrorResponse(err)
		if minioError.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetPresignedURL implements Storage.GetPresignedURL
func (st *MinioStorage) GetPresignedURL(ctx context.Context, bucketName, objName string, expires time.Duration) (string, error) {
	url, err := st.cli.PresignedGetObject(ctx, bucketName, objName, expires, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

// GetPublicURL implements Storage.GetPublicURL
func (st *MinioStorage) GetPublicURL(bucketName, objName string) (string, error) {
	scheme := "http://"
	if st.secure {
		scheme = "https://"
	}
	return scheme + path.Join(st.endpoint, bucketName, objName), nil
}

// DeleteObject implements Storage.DeleteObject
func (st *MinioStorage) DeleteObject(ctx context.Context, bucketName, objName string) error {
	return st.cli.RemoveObject(ctx, bucketName, objName, minio.RemoveObjectOptions{})
}
