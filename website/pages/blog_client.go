package pages

import (
	appconfig "emviwiki/shared/config"
	emvi "github.com/emvi/api-go"
	"github.com/emvi/logbuch"
	"time"
)

const (
	defaultLangCode = "en"
)

var (
	blogClient *emvi.Client
	blogLangs  map[string]emvi.Language
)

func InitBlogClient() {
	c := appconfig.Get()
	var config *emvi.Config

	if isIntegration {
		config = &emvi.Config{
			AuthHost: c.BlogClient.AuthHost,
			ApiHost:  c.BlogClient.ApiHost,
		}
	}

	blogClient = emvi.NewClient(c.BlogClient.ID, c.BlogClient.Secret, c.BlogClient.Organization, config)
}

func GetBlogLangId(code string) string {
	langs := getBlogLangs()

	if lang, ok := langs[code]; ok {
		return lang.Id
	}

	if lang, ok := langs[defaultLangCode]; ok {
		return lang.Id
	}

	return ""
}

func getBlogLangs() map[string]emvi.Language {
	if blogLangs == nil {
		langs, err := blogClient.GetLanguages()

		for err != nil {
			logbuch.Warn("Error reading languages from blog client", logbuch.Fields{"err": err})
			time.Sleep(time.Second * 5)
			langs, err = blogClient.GetLanguages()
		}

		logbuch.Info("Loaded languages from blog client")
		blogLangs = make(map[string]emvi.Language)

		for _, lang := range langs {
			blogLangs[lang.Code] = lang
		}
	}

	return blogLangs
}
