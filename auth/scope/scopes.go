// ******************************************************************************
// The scope package is unused right now but might become relevant in the future!
// ******************************************************************************
package scope

import (
	"emviwiki/shared/util"
	"strings"
)

var (
	scopes = []ScopeConfig{}
)

type ScopeConfig struct {
	Name       string      `yaml:"name"`
	AllowValue bool        `yaml:"allow_value"`
	Icon       string      `yaml:"icon"`
	Text       []ScopeText `yaml:"langs"`
}

type ScopeText struct {
	Lang        string `yaml:"lang"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

func (cfg *ScopeConfig) GetText(lang string) *ScopeText {
	if lang == "" {
		lang = util.DefaultSupportedLang
	}

	for _, text := range cfg.Text {
		if text.Lang == lang {
			return &text
		}
	}

	if lang == util.DefaultSupportedLang {
		return nil
	}

	return cfg.GetText(util.DefaultSupportedLang)
}

func GetScope(name string) *ScopeConfig {
	name = strings.ToLower(name)

	for _, cfg := range scopes {
		if cfg.Name == name {
			return &cfg
		}
	}

	return nil
}
