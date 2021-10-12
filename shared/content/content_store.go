package content

import (
	"io"
)

// FileInfo contains information about a file.
type FileInfo struct {
	Size int64
	MD5  string
}

// ContentStore is an interface for a file content storage.
type ContentStore interface {
	// Read reads a file for given path and returns a reader to it, which must be closed by the caller.
	Read(string) (io.ReadCloser, error)

	// Info returns information about a file.
	Info(string) (FileInfo, error)

	// Save saves a file in given path with given name.
	Save(string, string, io.Reader) error

	// Delete deletes a file for given path.
	Delete(string) error
}
