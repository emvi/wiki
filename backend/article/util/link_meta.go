package util

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/emvi/logbuch"
	"net/http"
	"unicode/utf8"
)

const (
	metaTitleMaxLength       = 200
	metaDescriptionMaxLength = 400
)

type LinkMetaData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func GetLinkMetaData(url string) LinkMetaData {
	resp, err := http.Get(url)

	if err != nil {
		return LinkMetaData{}
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logbuch.Debug("Error closing URL response body while extracting meta data", logbuch.Fields{"err": err})
		}
	}()

	return extractMetaData(resp)
}

func extractMetaData(resp *http.Response) LinkMetaData {
	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return LinkMetaData{}
	}

	titleQuery := [][]string{
		{"meta[property~='og:title']", "content"},
		{"meta[name~='twitter:title']", "content"},
		{"meta[name~='title']", "content"},
		{"title", ""},
		{"h1", ""},
		{"h2", ""},
		{"h3", ""},
		{"h4", ""},
	}
	descriptionQuery := [][]string{
		{"meta[property~='og:description']", "content"},
		{"meta[name~='twitter:description']", "content"},
		{"meta[name~='description']", "content"},
		{"p", ""},
	}
	imageQuery := [][]string{
		{"meta[property~='og:image']", "content"},
		{"meta[name~='twitter:image']", "content"},
		{"meta[name~='image']", "content"},
		{"link[rel~='icon']", "href"},
	}
	title := extractMetaField(doc, titleQuery)
	description := extractMetaField(doc, descriptionQuery)

	if utf8.RuneCountInString(title) > metaTitleMaxLength {
		title = title[:metaTitleMaxLength]
	}

	if utf8.RuneCountInString(description) > metaDescriptionMaxLength {
		description = description[:metaDescriptionMaxLength]
	}

	return LinkMetaData{
		Title:       title,
		Description: description,
		Image:       extractMetaField(doc, imageQuery),
	}
}

func extractMetaField(doc *goquery.Document, queries [][]string) string {
	value := ""

	for _, query := range queries {
		element := doc.Find(query[0])

		if len(element.Nodes) == 1 {
			if query[1] != "" {
				attr, exists := element.Attr(query[1])

				if exists {
					value = attr
					break
				}
			} else {
				value = element.Text()
			}
		}
	}

	return value
}
