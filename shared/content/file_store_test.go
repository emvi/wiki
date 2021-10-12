package content

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestFileStoreRead(t *testing.T) {
	createTestFile(t)
	bucket := &FileStore{}

	file, err := bucket.Read("bucket/some/test/unknown.test")

	if file != nil || err == nil {
		t.Fatal("File must not be found")
	}

	file, err = bucket.Read("bucket/some/test/unique.test")

	if file == nil || err != nil {
		t.Fatal("File must be found")
	}

	data, err := ioutil.ReadAll(file)

	if err := file.Close(); err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatalf("File must be readable, but was: %v", err)
	}

	if string(data) != "data" {
		t.Fatalf("File content must be correct, but was: %v", string(data))
	}
}

func TestFileStoreInfo(t *testing.T) {
	createTestFile(t)
	bucket := &FileStore{}
	info, err := bucket.Info("bucket/some/test/unique.test")

	if err != nil {
		t.Fatalf("File size must be returned, but was: %v", err)
	}

	if info.Size != 4 {
		t.Fatalf("File size must be correct, but was: %v", info.Size)
	}

	if info.MD5 != "cafbbc1fc265231f1e4595d17ea5c04b" {
		t.Fatalf("File MD5 must be correct, but was: %v", info.MD5)
	}
}

func TestFileStoreSave(t *testing.T) {
	bucket := &FileStore{}
	err := bucket.Save("bucket/some/test", "newfile.test", bytes.NewBufferString("data"))

	if err != nil {
		t.Fatalf("File must have been saved in store, but was: %v", err)
	}

	data, err := ioutil.ReadFile("bucket/some/test/newfile.test")

	if err != nil {
		t.Fatalf("File must be found after save, but was: %v", err)
	}

	if string(data) != "data" {
		t.Fatalf("File content must be correct, but was: %v", string(data))
	}
}

func TestFileStoreDelete(t *testing.T) {
	createTestFile(t)
	bucket := &FileStore{}

	if err := bucket.Delete("bucket/some/test/unknown.test"); err == nil {
		t.Fatal("File must not be found")
	}

	if err := bucket.Delete("bucket/some/test/unique.test"); err != nil {
		t.Fatalf("File must be deleted, but was: %v", err)
	}
}

func createTestFile(t *testing.T) {
	if err := os.MkdirAll("bucket/some/test", 0777); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile("bucket/some/test/unique.test", []byte("data"), 0666); err != nil {
		t.Fatal(err)
	}
}
