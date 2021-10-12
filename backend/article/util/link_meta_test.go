package util

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestExtractMetaDataOG(t *testing.T) {
	body := `<!DOCTYPE html>
<html lang="de">
<head>
    <base href="/" />
    <meta charset="utf-8" />
    <meta name="title" content="title" />
    <meta name="description" content="description" />
    <meta name="twitter:title" content="twitter:title" />
    <meta name="twitter:description" content="twitter:description" />
    <meta name="twitter:image" content="twitter:image" />
    <meta property="og:title" content="og:title" />
	<meta property="og:description" content="og:description" />
	<meta property="og:image" content="og:image" />
    <title>title title</title>
	<link rel="icon" type="image/png" sizes="32x32" href="image" />
</head>
<body>
	<p>Body</p>
</body>
</html>`
	reader := ioutil.NopCloser(strings.NewReader(body))
	resp := http.Response{Body: reader}
	meta := extractMetaData(&resp)

	if meta.Title != "og:title" || meta.Description != "og:description" || meta.Image != "og:image" {
		t.Fatalf("Meta data not as expected: %v", meta)
	}
}

func TestExtractMetaDataTwitter(t *testing.T) {
	body := `<!DOCTYPE html>
<html lang="de">
<head>
    <base href="/" />
    <meta charset="utf-8" />
    <meta name="title" content="title" />
    <meta name="description" content="description" />
    <meta name="twitter:title" content="twitter:title" />
    <meta name="twitter:description" content="twitter:description" />
    <meta name="twitter:image" content="twitter:image" />
    <title>title title</title>
	<link rel="icon" type="image/png" sizes="32x32" href="image" />
</head>
<body>
	<p>Body</p>
</body>
</html>`
	reader := ioutil.NopCloser(strings.NewReader(body))
	resp := http.Response{Body: reader}
	meta := extractMetaData(&resp)

	if meta.Title != "twitter:title" || meta.Description != "twitter:description" || meta.Image != "twitter:image" {
		t.Fatalf("Meta data not as expected: %v", meta)
	}
}

func TestExtractMetaDataDefault(t *testing.T) {
	body := `<!DOCTYPE html>
<html lang="de">
<head>
    <base href="/" />
    <meta charset="utf-8" />
    <meta name="title" content="title" />
    <meta name="description" content="description" />
    <title>title title</title>
	<link rel="icon" type="image/png" sizes="32x32" href="image" />
</head>
<body>
	<p>Body</p>
</body>
</html>`
	reader := ioutil.NopCloser(strings.NewReader(body))
	resp := http.Response{Body: reader}
	meta := extractMetaData(&resp)

	if meta.Title != "title" || meta.Description != "description" || meta.Image != "image" {
		t.Fatalf("Meta data not as expected: %v", meta)
	}
}

func TestExtractMetaDataTitle(t *testing.T) {
	body := `<!DOCTYPE html>
<html lang="de">
<head>
    <base href="/" />
    <meta charset="utf-8" />
    <meta name="description" content="description" />
    <title>title title</title>
	<link rel="icon" type="image/png" sizes="32x32" href="image" />
</head>
<body>
	<p>Body</p>
</body>
</html>`
	reader := ioutil.NopCloser(strings.NewReader(body))
	resp := http.Response{Body: reader}
	meta := extractMetaData(&resp)

	if meta.Title != "title title" || meta.Description != "description" || meta.Image != "image" {
		t.Fatalf("Meta data not as expected: %v", meta)
	}
}

func TestExtractMetaDataHeadlineAndParagraph(t *testing.T) {
	body := `<!DOCTYPE html>
<html lang="de">
<head>
    <base href="/" />
    <meta charset="utf-8" />
</head>
<body>
	<h1>Headline</h1>
	<p>Body</p>
</body>
</html>`
	reader := ioutil.NopCloser(strings.NewReader(body))
	resp := http.Response{Body: reader}
	meta := extractMetaData(&resp)

	if meta.Title != "Headline" || meta.Description != "Body" || meta.Image != "" {
		t.Fatalf("Meta data not as expected: %v", meta)
	}
}
