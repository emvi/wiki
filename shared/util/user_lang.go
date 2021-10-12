package util

import (
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/jmoiron/sqlx"
)

const (
	DefaultSupportedLang = "en"
)

var (
	SupportedLangs = map[string]bool{
		"en": true,
		"de": true,
	}
)

// DetermineLang returns the language to be used for given organization and user.
// If the language id is passed, the language will be returned.
// Else the users preferred language will be used or, if not available, the organizations default.
// This function is nil safe.
func DetermineLang(tx *sqlx.Tx, orgaId, userId, langId hide.ID) *model.Language {
	if langId != 0 {
		if lang := model.GetLanguageByOrganizationIdAndIdTx(tx, orgaId, langId); lang != nil {
			return lang
		}
	}

	// try to select users preferred language
	if userId != 0 {
		user := model.GetUserByOrganizationIdAndIdTx(tx, orgaId, userId)

		if user != nil && user.Language.Valid {
			lang := model.GetLanguageByOrganizationIdAndCodeTx(tx, orgaId, user.Language.String)

			if lang != nil {
				return lang
			}
		}
	}

	lang := model.GetDefaultLanguageByOrganizationIdTx(tx, orgaId)

	if lang != nil {
		return lang
	}

	return &model.Language{}
}

// DetermineSystemSupportedLangCode returns the user preferred language code when supported by the system.
// Else the default supported language code will be returned.
func DetermineSystemSupportedLangCode(orgaId, userId hide.ID) string {
	user := model.GetUserByOrganizationIdAndId(orgaId, userId)
	lang := DefaultSupportedLang

	if user != nil && user.Language.Valid {
		_, ok := SupportedLangs[user.Language.String]

		if ok {
			lang = user.Language.String
		}
	}

	return lang
}

// GetSystemSupportedLangCode returns the given code if it's supported by the system, or else the default will be returned.
func GetSystemSupportedLangCode(code string) string {
	_, ok := SupportedLangs[code]

	if ok {
		return code
	}

	return DefaultSupportedLang
}
