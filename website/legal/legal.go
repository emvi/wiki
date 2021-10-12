package legal

import (
	"encoding/json"
	"github.com/emvi/logbuch"
	"net/http"
	"time"
)

const (
	contentCacheTTLSeconds = 3600
)

var (
	contentCache = make(map[string]LegalContent)
)

type LegalContent struct {
	Content string
	TTL     time.Time
}

func GetPrivacyPolicy(langCode string) string {
	content, err := getJsonCached(privacyPolicyURL[langCode])

	if err != nil {
		logbuch.Error("Error getting privacy policy from iubenda", logbuch.Fields{"err": err})
		return ""
	}

	return content
}

func GetCookiePolicy(langCode string) string {
	content, err := getJsonCached(cookiePolicyURL[langCode])

	if err != nil {
		logbuch.Error("Error getting cookie policy from iubenda", logbuch.Fields{"err": err})
		return ""
	}

	return content
}

func GetTermsAndConditions(langCode string) string {
	content, err := getJsonCached(termsAndConditionsURL[langCode])

	if err != nil {
		logbuch.Error("Error getting terms and conditions from iubenda", logbuch.Fields{"err": err})
		return ""
	}

	return content
}

func getJsonCached(url string) (string, error) {
	cached, ok := contentCache[url]

	if ok && cached.TTL.After(time.Now()) {
		return cached.Content, nil
	}

	logbuch.Debug("Getting legal content", logbuch.Fields{"url": url})
	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			logbuch.Error("Error closing body after reading JSON response", logbuch.Fields{"err": err, "url": url})
		}
	}()

	respJson := struct {
		Content string `json:"content"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&respJson); err != nil {
		return "", err
	}

	contentCache[url] = LegalContent{respJson.Content, time.Now().Add(time.Second * contentCacheTTLSeconds)}
	return respJson.Content, nil
}
