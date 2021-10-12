package legal

import "github.com/emvi/logbuch"

func GetCookieNote(langCode string) string {
	note, ok := cookiesNote[langCode]

	if !ok {
		logbuch.Error("Cookie not not found for given language code", logbuch.Fields{"lang_code": langCode})
		return ""
	}

	return note
}
