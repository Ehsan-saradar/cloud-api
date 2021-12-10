package storage

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

	"api.cloud.io/pkg/test"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var (
	endpoint        = "localhost:9002"
	accessKeyID     = "cloud"
	secretAccessKey = "cloudPass"
	secure          = false
)

func TestMinioStorage(t *testing.T) {
	ctx := context.Background()
	st, err := NewMinioStorage(endpoint, accessKeyID, secretAccessKey, secure)
	assert.NoError(t, err)

	bucketName := uuid.NewV4().String()
	t.Run("NewBucket", func(t *testing.T) {
		err = st.NewBucket(ctx, bucketName)
		assert.NoError(t, err)
	})

	t.Run("BucketExists", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			isExists, err := st.BucketExists(ctx, bucketName)
			assert.NoError(t, err)
			assert.True(t, isExists)
		})

		t.Run("NotExists", func(t *testing.T) {
			fakeBucketName := uuid.NewV4().String()
			isExists, err := st.BucketExists(ctx, fakeBucketName)
			assert.NoError(t, err)
			assert.False(t, isExists)
		})
	})

	t.Run("SetBucketStatus", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			err := st.SetBucketStatus(ctx, bucketName, true)
			assert.NoError(t, err)
		})

		t.Run("BucketNotExists", func(t *testing.T) {
			fakeBucketName := uuid.NewV4().String()
			err := st.SetBucketStatus(ctx, fakeBucketName, true)
			assert.Equal(t, ErrBucketNotExists, err)
		})
	})

	t.Run("SetBucketLifecycle", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			err := st.SetBucketLifecycle(ctx, bucketName, 1)
			assert.NoError(t, err)
		})

		t.Run("BucketNotExists", func(t *testing.T) {
			fakeBucketName := uuid.NewV4().String()
			err := st.SetBucketLifecycle(ctx, fakeBucketName, 2)
			assert.Equal(t, ErrBucketNotExists, err)
		})
	})

	objName := uuid.NewV4().String()
	t.Run("Put", func(t *testing.T) {
		r := bytes.NewReader(test.SamplePNG)
		err := st.Put(ctx, bucketName, objName, r, r.Size(), "image/png")
		assert.NoError(t, err)
	})

	t.Run("ObjectExists", func(t *testing.T) {
		t.Run("Exists", func(t *testing.T) {
			isExists, err := st.ObjectExists(ctx, bucketName, objName)
			assert.NoError(t, err)
			assert.True(t, isExists)
		})

		t.Run("ObjectNotExists", func(t *testing.T) {
			fakeObjectName := uuid.NewV4().String()
			isExists, err := st.ObjectExists(ctx, bucketName, fakeObjectName)
			assert.NoError(t, err)
			assert.False(t, isExists)
		})
	})

	t.Run("CopyObject", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			newObjectName := uuid.NewV4().String()
			err := st.CopyObject(ctx, bucketName, objName, bucketName, newObjectName)
			assert.NoError(t, err)
		})

		t.Run("ObjectNotExists", func(t *testing.T) {
			fakeObjectName := uuid.NewV4().String()
			newObjectName := uuid.NewV4().String()
			err := st.CopyObject(ctx, bucketName, fakeObjectName, bucketName, newObjectName)
			assert.Equal(t, ErrObjectNotExists, err)
		})
	})

	t.Run("GetPresignedURL", func(t *testing.T) {
		url, err := st.GetPresignedURL(ctx, bucketName, objName, time.Hour)
		assert.NoError(t, err)

		resp, err := http.Get(url)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("GetPublicURL", func(t *testing.T) {
		url, err := st.GetPublicURL(bucketName, objName)
		assert.NoError(t, err)

		resp, err := http.Get(url)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("DeleteObject", func(t *testing.T) {
		err := st.DeleteObject(ctx, bucketName, objName)
		assert.NoError(t, err)
	})

	t.Run("DeleteBucket", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			err := st.DeleteBucket(ctx, bucketName)
			assert.NoError(t, err)
		})

		t.Run("BucketNotExists", func(t *testing.T) {
			fakeBucketName := uuid.NewV4().String()
			err := st.DeleteBucket(ctx, fakeBucketName)
			assert.Equal(t, ErrBucketNotExists, err)
		})
	})
}
