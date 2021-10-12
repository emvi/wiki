package content

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"io"
	"testing"
)

// FIXME add tests for protected files (organization, article)
func TestReadFile(t *testing.T) {
	testutil.CleanBackendDb(t)
	_, user := testutil.CreateOrgaAndUser(t)
	file := testutil.CreateFile(t, nil, user, nil, "")

	input := []struct {
		UniqueName string
	}{
		{""},
		{"not_found"},
		{file.UniqueName},
	}
	expected := []error{
		errs.FileNotFound,
		errs.FileNotFound,
		nil,
	}

	var result *model.File
	var reader io.Reader

	for i, in := range input {
		var err error
		result, reader, err = ReadFile(in.UniqueName)

		if err != expected[i] {
			t.Fatalf("Expected %v when reading file, but was: %v", expected[i], err)
		}
	}

	if result == nil {
		t.Fatal("Result must have been returned")
	}

	if reader == nil {
		t.Fatal("Reader must have been returned")
	}

	if result.ID != file.ID ||
		result.OrganizationId != file.OrganizationId ||
		result.UserId != file.UserId ||
		result.OriginalName != file.OriginalName ||
		result.UniqueName != file.UniqueName ||
		result.Path != file.Path ||
		result.Type != file.Type ||
		result.Size != file.Size ||
		result.MD5 != file.MD5 {
		t.Fatal("Result not as expected")
	}
}
