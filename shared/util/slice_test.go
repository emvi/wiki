package util

import (
	"github.com/emvi/hide"
	"testing"
)

func TestRemoveDuplicateIds(t *testing.T) {
	in := []hide.ID{1, 1, 2, 3, 4, 4, 4, 5, 6, 7, 7, 8}
	out := RemoveDuplicateIds(in)

	if len(out) != 8 {
		t.Fatalf("8 entries must be returned, but was: %v", len(out))
	}

	for i := 0; i < 8; i++ {
		if out[i] != hide.ID(i+1) {
			t.Fatalf("Expected entry %v to be %v, but was: %v", i, i+1, out[i])
		}
	}
}

func TestRemoveId(t *testing.T) {
	in := []hide.ID{1, 2, 3, 4, 5, 6, 7, 8}
	out := RemoveId(in, 3)

	if len(out) != 7 {
		t.Fatalf("7 entries must be returned, but was: %v", len(out))
	}

	for _, id := range out {
		if id == 3 {
			t.Fatal("ID 3 must not exist in slice")
		}
	}
}

func TestRemoveIds(t *testing.T) {
	in := []hide.ID{1, 2, 3, 4, 5, 6, 7, 8}
	remove := []hide.ID{2, 4, 7}
	out := RemoveIds(in, remove)

	if len(out) != 5 {
		t.Fatalf("5 entries must be returned, but was: %v", len(out))
	}

	for _, id := range out {
		if id == 2 || id == 4 || id == 7 {
			t.Fatal("ID 2, 4, 7 must not exist in slice")
		}
	}
}

func TestIntersectIds(t *testing.T) {
	a := []hide.ID{1, 2, 3, 4, 5, 6, 7, 8}
	b := []hide.ID{3, 4, 7, 8}
	out := IntersectIds(a, b)

	if len(out) != 4 {
		t.Fatalf("4 entries must be returned, but was: %v", len(out))
	}

	for _, id := range out {
		if id != 3 && id != 4 && id != 7 && id != 8 {
			t.Fatal("ID 2, 4, 7, 8 must exist in slice")
		}
	}
}
