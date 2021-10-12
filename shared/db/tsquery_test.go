package db

import (
	"testing"
)

func TestToTSVector(t *testing.T) {
	input := []string{
		"",
		"hello",
		"	 test  ",
		"hello world",
		" 	  hello 	    world   	 ",
		"hello&world",
		"  hello	   & world   ",
		"hello|&world",
		"hello!|&world",
		"hello&(world:a|test:*)",
		"!not",
		"&test",
		"|test",
		"test!",
		"test&",
		"test|",
	}
	expected := []string{
		"",
		"hello",
		"test",
		"hello|world",
		"hello|world",
		"hello&world",
		"hello&world",
		"hello|world",
		"hello&!world",
		"hello&world|a|test",
		"!not",
		"test",
		"test",
		"test",
		"test",
		"test",
	}

	for i, in := range input {
		got := ToTSVector(in)

		if got != expected[i] {
			t.Fatalf("Expected '%v', but was: %v", expected[i], got)
		}
	}
}

func TestParseTSVectorTokens(t *testing.T) {
	input := []string{
		"",
		"hello",
		"  hello	   & world   ",
		"hello|&world",
		"hello!|&world",
		"hello&(world:a|test:*)",
	}
	expected := [][]string{
		{},
		{"hello"},
		{"hello", "&", "world"},
		{"hello", "|", "&", "world"},
		{"hello", "!", "|", "&", "world"},
		{"hello", "&", "world", "a", "|", "test"},
	}

	for i, in := range input {
		got := parseTSVectorTokens(in)

		for j, token := range got {
			if token != expected[i][j] {
				t.Fatalf("Expected '%v', but was: %v", expected[i], got)
			}
		}
	}
}
