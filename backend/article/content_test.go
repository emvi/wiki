package article

import (
	"emviwiki/backend/prosemirror"
	"testing"
)

func TestExtractTextFromContent(t *testing.T) {
	in := `{"type":"doc","content":[{"type":"paragraph","content":[{"type":"text","text":"Das ist ein ganz neuer Artikel... und sollte funktionieren :)"}]}]}`
	doc, err := prosemirror.ParseDoc(in)

	if err != nil {
		t.Fatal(err)
	}

	out := extractTextFromContent(doc)

	if out != "Das ist ein ganz neuer Artikel... und sollte funktionieren :)" {
		t.Fatalf("Expected output to contain text of article, but was: %v", out)
	}
}

func TestCalculateReadingTime(t *testing.T) {
	input := []string{
		"1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20",
		"1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20",
		"1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20",
	}
	expected := []int{
		6,
		12,
		60,
	}

	for i, in := range input {
		if rt := calculateReadingTimeSeconds(in); rt != expected[i] {
			t.Fatalf("Expected %v seconds, but was: %v", expected[i], rt)
		}
	}
}
