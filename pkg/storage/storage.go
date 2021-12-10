package storage

import (
	"context"
	"errors"
	"io"
	"time"
)

var (
	ErrBucketNotExists = errors.New("the requested bucket doesn't exists")
	ErrObjectNotExists = errors.New("the requested object doesn't exists")
)

// Storage is the implementation of object storage.
type Storage interface {
	// NewBucket will create a new bucket with given name.
	NewBucket(ctx context.Context, bucketName string) error
	// BucketExists returns whether following bucket is exists or not.
	BucketExists(ctx context.Context, bucketName string) (bool, error)
	// SetBucketStatus will set whether following bucket be public or private.
	SetBucketStatus(ctx context.Context, bucketName string, isPublic bool) error
	// SetBucketLifecycle sets objects life time in the given bucket.
	SetBucketLifecycle(ctx context.Context, bucketName string, days int) error
	// DeleteBucket will delete following bucket from storage.
	// Bucket must be empty otherwise this will fail.
	DeleteBucket(ctx context.Context, bucketName string) error
	// Put uploads given object to following bucket.
	Put(ctx context.Context, bucketName, objName string, r io.Reader, size int64, contentType string) error
	// ObjectExists returns whether following object is exists or not.
	ObjectExists(ctx context.Context, bucketName, objName string) (bool, error)
	// GetPresignedURL returns a temporary and unique URL to following object.
	GetPresignedURL(ctx context.Context, bucketName, objName string, expires time.Duration) (string, error)
	// GetPublicURL returns exact URL to following object.
	// Bucket must be public otherwise URL won't work.
	GetPublicURL(bucketName, objName string) (string, error)
	// Delete will delete following object from storage.
	DeleteObject(ctx context.Context, bucketName, objName string) error
}
