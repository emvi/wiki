package content

import (
	"emviwiki/shared/config"
	"errors"
	"github.com/emvi/logbuch"
	"github.com/minio/minio-go"
	"io"
	"path/filepath"
	"strings"
)

type MinioStore struct {
	client     *minio.Client
	bucketName string
}

func NewMinioStore() *MinioStore {
	c := config.Get().Storage.Minio
	client, err := minio.New(c.Endpoint, c.ID, c.Secret, c.Secure)

	if err != nil {
		logbuch.Fatal("Error creating MinIO Storage client", logbuch.Fields{"err": err})
		return nil
	}

	bucketName := strings.TrimSpace(c.Bucket)

	if bucketName == "" {
		logbuch.Fatal("Bucket name must not be empty", logbuch.Fields{"name": bucketName})
	}

	// the bucket location doesn't matter, since we're not using AWS S3 here but it needs to be passed anyways...
	err = client.MakeBucket(bucketName, "us-east-1")

	if err != nil {
		exists, err := client.BucketExists(bucketName)

		if err == nil && exists {
			logbuch.Info("Using existing MinIO bucket", logbuch.Fields{"name": bucketName})
		} else {
			logbuch.Fatal("Error mounting MinIO bucket", logbuch.Fields{"name": bucketName, "err": err})
		}
	} else {
		logbuch.Info("Created new MinIO bucket", logbuch.Fields{"name": bucketName})
	}

	logbuch.Info("Connected to MinIO bucket", logbuch.Fields{"name": bucketName, "use_ssl": c.Secure})
	return &MinioStore{client, bucketName}
}

func (store *MinioStore) Read(path string) (io.ReadCloser, error) {
	return store.client.GetObject(store.bucketName, path, minio.GetObjectOptions{})
}

func (store *MinioStore) Info(path string) (FileInfo, error) {
	stat, err := store.client.StatObject(store.bucketName, path, minio.StatObjectOptions{})

	if err != nil {
		return FileInfo{}, err
	}

	etag := strings.Split(stat.ETag, "-")

	if len(etag) < 1 {
		return FileInfo{}, errors.New("ETag has unexpected number of parts (splitting at '-')")
	}

	return FileInfo{
		stat.Size,
		etag[0],
	}, nil
}

func (store *MinioStore) Save(path, filename string, reader io.Reader) error {
	_, err := store.client.PutObject(store.bucketName, filepath.Join(path, filename), reader, -1, minio.PutObjectOptions{})

	if err != nil {
		return err
	}

	return nil
}

func (store *MinioStore) Delete(path string) error {
	if err := store.client.RemoveObject(store.bucketName, path); err != nil {
		return err
	}

	return nil
}
