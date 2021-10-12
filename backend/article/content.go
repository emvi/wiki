package article

import (
	"emviwiki/backend/prosemirror"
	"math"
	"strings"
)

const (
	wordsPerMinute = float64(200)
)

func extractTextFromContent(doc *prosemirror.Node) string {
	textNodes := prosemirror.FindNodes(doc, -1, "text")
	var out strings.Builder

	for _, node := range textNodes {
		out.WriteString(node.Text)
	}

	return out.String()
}

func calculateReadingTimeSeconds(content string) int {
	wordCount := len(strings.Split(content, " "))
	readingTime := float64(wordCount) / wordsPerMinute
	readingTimeMinutes := math.Abs(readingTime)
	readingTimeSeconds := (readingTime - readingTimeMinutes) * 0.6
	return int(readingTimeMinutes*60 + readingTimeSeconds)
}
