package sitemap

import (
	"encoding/xml"
	"time"
)

const (
	header        = `<?xml version="1.0" encoding="UTF-8"?>`
	xmlns         = "http://www.sitemaps.org/schemas/sitemap/0.9"
	lastmodFormat = "2006-01-02"
)

type URLSet struct {
	XMLName xml.Name `xml:"urlset"`
	XMLNS   string   `xml:"xmlns,attr"`
	URLs    []URL
}

type URL struct {
	XMLName    xml.Name `xml:"url"`
	Loc        string   `xml:"loc"`
	Lastmod    string   `xml:"lastmod"`
	Changefreq string   `xml:"changefreq,omitempty"`
	Priority   string   `xml:"priority,omitempty"`
}

func GenerateSitemap(urls []URL) ([]byte, error) {
	now := time.Now().Format(lastmodFormat)

	for i := range urls {
		urls[i].Lastmod = now
	}

	sitemap := URLSet{XMLNS: xmlns, URLs: urls}
	out, err := xml.Marshal(&sitemap)

	if err != nil {
		return nil, err
	}

	return []byte(header + string(out)), nil
}
