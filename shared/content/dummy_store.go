package content

import (
	"io"
)

type DummyReadCloser struct{}

type DummySaves struct {
	Path     string
	Filename string
	Reader   io.Reader
}

type DummyStore struct {
	Reads   []string
	Infos   []string
	Saves   []DummySaves
	Deletes []string
}

func NewDummyStore() *DummyStore {
	return &DummyStore{}
}

func (closer *DummyReadCloser) Read(b []byte) (int, error) {
	return 0, nil
}

func (closer *DummyReadCloser) Close() error {
	return nil
}

func (store *DummyStore) Read(path string) (io.ReadCloser, error) {
	store.Reads = append(store.Reads, path)
	return &DummyReadCloser{}, nil
}

func (store *DummyStore) Info(path string) (FileInfo, error) {
	store.Infos = append(store.Infos, path)
	return FileInfo{}, nil
}

func (store *DummyStore) Save(path, filename string, reader io.Reader) error {
	store.Saves = append(store.Saves, DummySaves{path, filename, reader})
	return nil
}

func (store *DummyStore) Delete(path string) error {
	store.Deletes = append(store.Deletes, path)
	return nil
}

func (store *DummyStore) Reset() {
	store.Reads = []string{}
	store.Infos = []string{}
	store.Saves = []DummySaves{}
	store.Deletes = []string{}
}
