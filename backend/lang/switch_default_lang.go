package lang

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/perm"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
)

func SwitchDefaultLanguage(orga *model.Organization, userId, langId hide.ID) error {
	if !orga.Expert {
		return errs.RequiresExpertVersion
	}

	if _, err := perm.CheckUserIsAdmin(orga.ID, userId); err != nil {
		return err
	}

	defaultLang := model.GetDefaultLanguageByOrganizationId(orga.ID)

	if defaultLang.ID == langId {
		// no update required
		return nil
	}

	lang := model.GetLanguageByOrganizationIdAndId(orga.ID, langId)

	if lang == nil {
		return errs.LanguageNotFound
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction when switching default language", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "lang_id": langId})
		return errs.TxBegin
	}

	defaultLang.Default = false
	lang.Default = true

	if err := model.SaveLanguage(tx, defaultLang); err != nil {
		logbuch.Error("Error saving old default language when switching default language", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "lang_id": langId})
		return errs.Saving
	}

	if err := model.SaveLanguage(tx, lang); err != nil {
		logbuch.Error("Error saving new default language when switching default language", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "lang_id": langId})
		return errs.Saving
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when switching default language", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId, "lang_id": langId})
		return errs.TxCommit
	}

	return nil
}
