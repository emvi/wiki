package i18n

import (
	"emviwiki/shared/util"
	"html/template"
)

// Translation is map of i18n variables for a map of languages.
type Translation map[string]map[string]template.HTML

// GetVars returns the variables for given language code or the default, if not available.
func GetVars(langCode string, vars Translation) map[string]template.HTML {
	_, ok := util.SupportedLangs[langCode]

	if ok {
		return vars[langCode]
	}

	return vars[util.DefaultSupportedLang]
}
