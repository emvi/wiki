package content

import (
	"crypto/md5"
	"emviwiki/shared/config"
	"encoding/hex"
	"github.com/emvi/logbuch"
	"io"
	"os"
	"path/filepath"
)

type FileStore struct {
	BasePath string
}

func NewFileStore() *FileStore {
	return &FileStore{config.Get().Storage.Path}
}

func (store *FileStore) Read(path string) (io.ReadCloser, error) {
	file, err := os.Open(filepath.Join(store.BasePath, path))

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (store *FileStore) Info(path string) (FileInfo, error) {
	file, err := os.Open(filepath.Join(store.BasePath, path))

	if err != nil {
		return FileInfo{}, err
	}

	defer file.Close()
	stats, err := file.Stat()

	if err != nil {
		return FileInfo{}, err
	}

	h := md5.New()
	buffer := make([]byte, 512)

	for {
		_, err := file.Read(buffer)

		if err == io.EOF {
			break
		} else if err != nil {
			return FileInfo{}, err
		}

		if _, err := h.Write(buffer); err != nil {
			return FileInfo{}, err
		}
	}

	return FileInfo{
		stats.Size(),
		hex.EncodeToString(h.Sum(nil)),
	}, nil
}

func (store *FileStore) Save(path, filename string, reader io.Reader) error {
	if err := os.MkdirAll(filepath.Join(store.BasePath, path), 0776); err != nil {
		return err
	}

	f, err := os.OpenFile(filepath.Join(store.BasePath, path, filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {
		return err
	}

	defer func() {
		if err := f.Close(); err != nil {
			logbuch.Warn("Error on save closing file in file store", logbuch.Fields{"err": err, "path": path, "filename": filename})
		}
	}()

	if _, err := io.Copy(f, reader); err != nil {
		return err
	}

	return nil
}

func (store *FileStore) Delete(path string) error {
	if err := os.Remove(filepath.Join(store.BasePath, path)); err != nil {
		return err
	}

	return nil
}
