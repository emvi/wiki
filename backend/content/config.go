package content

import (
	"emviwiki/shared/config"
	"emviwiki/shared/content"
	"github.com/emvi/logbuch"
)

var (
	store content.ContentStore
)

func LoadConfig() {
	storeType := config.Get().Storage.Type

	if storeType == "file" {
		logbuch.Info("Using file store for content")
		store = content.NewFileStore()
	} else if storeType == "gcs" {
		logbuch.Info("Using Google Cloud Store for content")
		store = content.NewGoogleCloudStore()
	} else if storeType == "minio" {
		logbuch.Info("Using MinIO Store for content")
		store = content.NewMinioStore()
	} else {
		logbuch.Info("Using dummy store for content")
		store = content.NewDummyStore()
	}
}

// GetStore returns the configured store.
func GetStore() content.ContentStore {
	return store
}
