package sitemap

import (
	"testing"
	"time"
)

func TestGenerateSitemap(t *testing.T) {
	urls := []URL{
		{Loc: "https://emvi.com/", Priority: "1.0"},
		{Loc: "https://emvi.com/pricing", Priority: "0.8"},
	}
	out, err := GenerateSitemap(urls)

	if err != nil {
		t.Fatalf("Sitemap must have been generated, but was: %v", err)
	}

	today := time.Now().Format("2006-01-02")

	if string(out) != `<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url><loc>https://emvi.com/</loc><lastmod>`+today+`</lastmod><priority>1.0</priority></url><url><loc>https://emvi.com/pricing</loc><lastmod>`+today+`</lastmod><priority>0.8</priority></url></urlset>` {
		t.Fatalf("Sitemap XML not as expected: %v", string(out))
	}
}
