package content

import (
	"context"
	"emviwiki/shared/config"
	"encoding/hex"
	"github.com/emvi/logbuch"
	"io"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
)

type GoogleCloudStore struct {
	client *storage.Client
	bucket *storage.BucketHandle
}

func NewGoogleCloudStore() *GoogleCloudStore {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		logbuch.Fatal("Error creating Google Cloud Storage client", logbuch.Fields{"err": err})
		return nil
	}

	bucketName := strings.TrimSpace(config.Get().Storage.GCSBucket)

	if bucketName == "" {
		logbuch.Fatal("Bucket name must not be empty", logbuch.Fields{"name": bucketName})
	}

	bucket := client.Bucket(bucketName)
	logbuch.Info("Bucket name set", logbuch.Fields{"name": bucketName})

	return &GoogleCloudStore{client, bucket}
}

func (store *GoogleCloudStore) Read(path string) (io.ReadCloser, error) {
	ctx := context.Background()
	return store.bucket.Object(path).NewReader(ctx)
}

func (store *GoogleCloudStore) Info(path string) (FileInfo, error) {
	ctx := context.Background()
	attrs, err := store.bucket.Object(path).Attrs(ctx)

	if err != nil {
		return FileInfo{}, err
	}

	return FileInfo{
		attrs.Size,
		hex.EncodeToString(attrs.MD5),
	}, nil
}

func (store *GoogleCloudStore) Save(path, filename string, reader io.Reader) error {
	ctx := context.Background()
	writer := store.bucket.Object(filepath.Join(path, filename)).NewWriter(ctx)

	if _, err := io.Copy(writer, reader); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	return nil
}

func (store *GoogleCloudStore) Delete(path string) error {
	ctx := context.Background()

	if err := store.bucket.Object(path).Delete(ctx); err != nil {
		return err
	}

	return nil
}
