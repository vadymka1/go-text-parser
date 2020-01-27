package controllers

import (
	"testing"
)

var filePath = "../test.txt"

func TestParseFile(t *testing.T) {
	t.Run("Test open file", func(t *testing.T) {
		if _, err := ParseFile("../test.txt"); (err != nil) != false {
			t.Errorf("Parefile() error = %v", err)
		}
	})
}

func BenchmarkParseTextToWords(b *testing.B) {
	for n := 0; n < b.N; n++  {
		ParseFile(filePath)
	}
}

func BenchmarkFindNumberOfChars(b *testing.B) {
	words, _ := ParseFile(filePath)

	for n := 0; n < b.N; n++  {
		FindNumberOfChars(words)
	}
}

func BenchmarkFindNumberOfWords(b *testing.B) {
	text, _ := ParseFile(filePath)

	words := ParseTextToWords(text)

	for n := 0; n < b.N; n++  {
		FindNumberOfWords(words)
	}

}
