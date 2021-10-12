package organization

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/constants"
	"emviwiki/shared/model"
	"fmt"
	"github.com/emvi/hide"
	"github.com/emvi/iso-639-1"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	nameMinLength               = 1
	nameMaxLength               = 60
	domainMinLength             = 3
	domainMaxLength             = 20
	defaultNotificationInterval = 7
	maxFreeOrganizations        = 1

	introArticleTitleDE            = "Willkommen bei Emvi"
	introArticleContentDE          = `{"type":"doc","content":[{"type":"headline","attrs":{"level":2},"content":[{"type":"text","text":"Willkommen bei Emvi. Wir sind froh dich an Board zu haben!"}]},{"type":"paragraph","content":[{"type":"text","text":"Das hier ist ein Beispielartikel. Du kannst ihn löschen indem du das Befehlsmenü öffnest ("},{"type":"text","marks":[{"type":"code"}],"text":"Shift"},{"type":"text","text":" + "},{"type":"text","marks":[{"type":"code"}],"text":"Leer"},{"type":"text","text":"), "},{"type":"text","marks":[{"type":"code"}],"text":"/löschen"},{"type":"text","text":" eintippst, Enter drückst und die Aktion bestätigst."}]},{"type":"paragraph","content":[{"type":"text","text":"Falls du die Einführung erneut lesen möchtest, wechsel zu den "},{"type":"text","marks":[{"type":"code"}],"text":"/einstellungen"},{"type":"text","text":" und wähle den Punkt "},{"type":"text","marks":[{"type":"italic"}],"text":"Einführung zurücksetzen"},{"type":"text","text":". Weitere Informationen findest du in der "},{"type":"text","marks":[{"type":"code"}],"text":"/hilfe"},{"type":"text","text":"."}]},{"type":"headline","attrs":{"level":2},"content":[{"type":"text","text":"Ein paar Tipps"}]},{"type":"bullet_list","content":[{"type":"list_item","content":[{"type":"paragraph","content":[{"type":"text","text":"verwende Listen für Artikel, die in einer festgelegten Reihenfolge gelesen werden sollen"}]}]},{"type":"list_item","content":[{"type":"paragraph","content":[{"type":"text","text":"Tags sind hilfreich um Artikel thematisch zusammenzufassen"}]}]},{"type":"list_item","content":[{"type":"paragraph","content":[{"type":"text","text":"verlinke relevante Inhalte mit Hilfe des @-Zeichens. Du kannst Artikel, Listen, Mitglieder, Gruppen und Tags verlinken."}]}]},{"type":"list_item","content":[{"type":"paragraph","content":[{"type":"text","text":"Artikelempfehlungen mit Lesebestätigung sind ein hilfreiches Werkzeug um sicherzustellen, dass eine Person oder eine Gruppe einen Artikel gelesen hat"}]}]}]},{"type":"headline","attrs":{"level":2},"content":[{"type":"text","text":"Hast du Fragen?"}]},{"type":"paragraph","content":[{"type":"text","text":"Wenn du eine Frage hast oder auf ein Problem stößt, kontaktiere uns per E-Mail an "},{"type":"text","marks":[{"type":"link","attrs":{"href":"mailto:support@emvi.com"}}],"text":"support@emvi.com"},{"type":"text","text":" oder über den "},{"type":"text","marks":[{"type":"code"}],"text":"/support"},{"type":"text","text":" Befehl."}]}]}`
	introArticleTitleEN            = "Get Started"
	introArticleContentEN          = `{"type":"doc","content":[{"type":"headline","attrs":{"level":2},"content":[{"type":"text","text":"Welcome to Emvi, we're happy to have you on board!"}]},{"type":"paragraph","content":[{"type":"text","text":"This is a sample article. You can delete it by opening the command menu ("},{"type":"text","marks":[{"type":"code"}],"text":"Shift"},{"type":"text","text":" + "},{"type":"text","marks":[{"type":"code"}],"text":"Space"},{"type":"text","text":"), type "},{"type":"text","marks":[{"type":"code"}],"text":"/delete"},{"type":"text","text":", press enter, and confirm the action."}]},{"type":"paragraph","content":[{"type":"text","text":"In case you would like to read the introduction again, go to the "},{"type":"text","marks":[{"type":"code"}],"text":"/settings"},{"type":"text","text":" menu and select "},{"type":"text","marks":[{"type":"italic"}],"text":"Repeat Introduction"},{"type":"text","text":". You can find more information in the "},{"type":"text","marks":[{"type":"code"}],"text":"/help"},{"type":"text","text":" menu."}]},{"type":"headline","attrs":{"level":2},"content":[{"type":"text","text":"Some Tips"}]},{"type":"bullet_list","content":[{"type":"list_item","content":[{"type":"paragraph","content":[{"type":"text","text":"use lists for articles that should be read in a certain order"}]}]},{"type":"list_item","content":[{"type":"paragraph","content":[{"type":"text","text":"tags can be helpful to group articles by a broader topic and improve searchability"}]}]},{"type":"list_item","content":[{"type":"paragraph","content":[{"type":"text","text":"link relevant content by using @-mentions in articles. You can link to articles, lists, members, groups, and tags"}]}]},{"type":"list_item","content":[{"type":"paragraph","content":[{"type":"text","text":"recommending an article with read confirmation is a helpful tool to make sure a person or group has read an article. You will receive a notification then"}]}]}]},{"type":"headline","attrs":{"level":2},"content":[{"type":"text","text":"Have a Question?"}]},{"type":"paragraph","content":[{"type":"text","text":"In case you have a question or encounter an issue, feel free to contact us by sending an email to "},{"type":"text","marks":[{"type":"link","attrs":{"href":"mailto:support@emvi.com"}}],"text":"support@emvi.com"},{"type":"text","text":" or by using the "},{"type":"text","marks":[{"type":"code"}],"text":"/support"},{"type":"text","text":" command."}]}]}`
	introArticleReadingTimeSeconds = 120
	introTagNameDE                 = "Einführung"
	introTagNameEN                 = "Introduction"
)

type CreateOrganizationData struct {
	Name        string `json:"name"`
	Domain      string `json:"domain"`
	Username    string `json:"username"`
	DefaultLang string `json:"default_lang"`
}

func (data *CreateOrganizationData) validate() []error {
	data.Trim()
	err := make([]error, 0)

	if e := data.CheckDomainValid(data.Domain); e != nil {
		err = append(err, e)
	}

	if e := data.CheckNameValid(data.Name); e != nil {
		err = append(err, e)
	}

	if e := data.CheckUsernameValid(data.Username); e != nil {
		err = append(err, e)
	}

	if e := data.CheckLanguage(); e != nil {
		err = append(err, e)
	}

	if len(err) != 0 {
		return err
	}

	return nil
}

func (data *CreateOrganizationData) Trim() {
	data.Name = strings.TrimSpace(data.Name)
	data.Domain = strings.ToLower(strings.TrimSpace(data.Domain))
	data.Username = strings.TrimSpace(data.Username)
	data.DefaultLang = strings.ToLower(strings.TrimSpace(data.DefaultLang))
}

func (data *CreateOrganizationData) CheckDomainValid(domain string) error {
	if utf8.RuneCountInString(domain) < domainMinLength {
		return errs.DomainTooShort
	}

	if utf8.RuneCountInString(domain) > domainMaxLength {
		return errs.DomainTooLong
	}

	first := true

	for i, c := range domain {
		if first && !unicode.IsLetter(c) {
			return errs.DomainNumberFirstChar
		} else if i == len(domain)-1 && !data.checkDomainChar(c) {
			return errs.DomainLastChar
		} else if !data.checkDomainChar(c) && c != '-' {
			return errs.DomainInvalidChar
		}

		first = false
	}

	if model.GetDomainBlacklistByName(data.Domain) != nil {
		return errs.DomainNotAllowed
	}

	if model.GetOrganizationByNameNormalized(data.Domain) != nil {
		return errs.DomainInUse
	}

	return nil
}

func (data *CreateOrganizationData) checkDomainChar(r rune) bool {
	return unicode.IsNumber(r) ||
		(r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z')
}

func (data *CreateOrganizationData) CheckNameValid(name string) error {
	if utf8.RuneCountInString(name) < nameMinLength {
		return errs.NameTooShort
	}

	if utf8.RuneCountInString(name) > nameMaxLength {
		return errs.NameTooLong
	}

	return nil
}

func (data *CreateOrganizationData) CheckUsernameValid(username string) error {
	if len(username) < nameMinLength || utf8.RuneCountInString(username) > nameMaxLength {
		return errs.UsernameInvalid
	}

	first := true

	for i, c := range username {
		if first && !unicode.IsLetter(c) && !unicode.IsNumber(c) {
			return errs.UsernameFirstChar
		} else if i == len(username)-1 && !unicode.IsLetter(c) && !unicode.IsNumber(c) {
			return errs.UsernameLastChar
		} else if !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '-' && c != '_' && c != ' ' {
			return errs.UsernameInvalidChar
		}

		first = false
	}

	return nil
}

func (data *CreateOrganizationData) CheckLanguage() error {
	if !iso6391.ValidCode(data.DefaultLang) {
		return errs.DefaultLanguageInvalid
	}

	return nil
}

// CreateOrganization creates a new organization with given name, (sub)domain, languages and current user as owner.
// The languages are passed as two character strings and mapped to the available languages.
func CreateOrganization(userId hide.ID, data CreateOrganizationData) []error {
	if len(model.FindOrganizationByUserIdAndNotExpert(userId)) >= maxFreeOrganizations {
		return []error{errs.MaxFreeOrganizationsReached}
	}

	if err := data.validate(); err != nil {
		return err
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction to create organization", logbuch.Fields{"err": err})
		return []error{errs.TxBegin}
	}

	org := &model.Organization{
		Name:           data.Name,
		NameNormalized: strings.ToLower(data.Domain),
		MaxStorageGB:   constants.DefaultMaxStorageGb,
		OwnerUserId:    userId,
	}

	if err := model.SaveOrganization(tx, org); err != nil {
		return []error{errs.Saving}
	}

	memberLang, err := createAndSaveLanguage(tx, org, data.DefaultLang)

	if err != nil {
		return []error{errs.Saving}
	}

	member := &model.OrganizationMember{Username: data.Username,
		IsModerator:               true,
		IsAdmin:                   true,
		SendNotificationsInterval: defaultNotificationInterval,
		NextNotificationMail:      time.Now().Add(time.Hour * 24 * defaultNotificationInterval),
		RecommendationMail:        true,
		OrganizationId:            org.ID,
		UserId:                    userId,
		LanguageId:                memberLang.ID,
		ShowActionButtons:         true,
		ShowNavigation:            true,
		ShowCreateButton:          true,
		Active:                    true}

	if err := model.SaveOrganizationMember(tx, member); err != nil {
		return []error{errs.Saving}
	}

	if err := createDefaultUserGroups(tx, org, userId); err != nil {
		return []error{err}
	}

	if err := createIntroduction(tx, org, userId, memberLang); err != nil {
		return []error{err}
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction to create organization", logbuch.Fields{"err": err})
		return []error{errs.TxCommit}
	}

	return nil
}

func createAndSaveLanguage(tx *sqlx.Tx, org *model.Organization, defaultLangCode string) (*model.Language, error) {
	lang := iso6391.FromCode(defaultLangCode)
	language := &model.Language{Name: lang.NativeName,
		Code:           lang.Code,
		Default:        true,
		OrganizationId: org.ID}

	if err := model.SaveLanguage(tx, language); err != nil {
		logbuch.Error("Error saving default language when creating new organization", logbuch.Fields{"err": err})
		return nil, err
	}

	// return owners default language
	return language, nil
}

func createDefaultUserGroups(tx *sqlx.Tx, orga *model.Organization, userId hide.ID) error {
	if err := createDefaultUserGroup(tx, orga, userId, constants.GroupAllName, constants.GroupAllInfo); err != nil {
		return err
	}

	if err := createDefaultUserGroup(tx, orga, userId, constants.GroupAdminName, constants.GroupAdminInfo); err != nil {
		return err
	}

	if err := createDefaultUserGroup(tx, orga, userId, constants.GroupModName, constants.GroupModInfo); err != nil {
		return err
	}

	// create read only group without initial member
	readOnly := &model.UserGroup{OrganizationId: orga.ID,
		Name:      constants.GroupReadOnlyName,
		Info:      null.NewString(constants.GroupReadOnlyInfo, constants.GroupReadOnlyInfo != ""), // always true, that's fine
		Immutable: true}

	if err := model.SaveUserGroup(tx, readOnly); err != nil {
		logbuch.Error("Error creating default user group", logbuch.Fields{"err": err, "name": constants.GroupReadOnlyName})
		return errs.Saving
	}

	return nil
}

func createDefaultUserGroup(tx *sqlx.Tx, orga *model.Organization, userId hide.ID, name, info string) error {
	all := &model.UserGroup{OrganizationId: orga.ID,
		Name:      name,
		Info:      null.NewString(info, info != ""),
		Immutable: true}

	if err := model.SaveUserGroup(tx, all); err != nil {
		logbuch.Error("Error creating default user group", logbuch.Fields{"err": err, "name": name})
		return errs.Saving
	}

	allMember := &model.UserGroupMember{UserGroupId: all.ID, UserId: userId}

	if err := model.SaveUserGroupMember(tx, allMember); err != nil {
		logbuch.Error("Error creating default user group member", logbuch.Fields{"err": err, "name": name})
		return errs.Saving
	}

	return nil
}

func createIntroduction(tx *sqlx.Tx, orga *model.Organization, userId hide.ID, lang *model.Language) error {
	// create article
	article := &model.Article{OrganizationId: orga.ID,
		ReadEveryone:  true,
		WriteEveryone: true,
		Published:     null.NewTime(time.Now(), true),
		Pinned:        true,
		WIP:           -1}

	if err := model.SaveArticle(tx, article); err != nil {
		logbuch.Error("Error saving introduction article", logbuch.Fields{"err": err, "user_id": userId})
		return errs.Saving
	}

	// create article access
	access := &model.ArticleAccess{UserId: userId,
		ArticleId: article.ID,
		Write:     true}

	if err := model.SaveArticleAccess(tx, access); err != nil {
		logbuch.Error("Error saving introduction article access", logbuch.Fields{"err": err, "user_id": userId})
		return errs.Saving
	}

	// create article content
	// default = English
	title := introArticleTitleEN
	content := introArticleContentEN
	tagName := introTagNameEN

	if lang.Code == "de" {
		title = introArticleTitleDE
		content = introArticleContentDE
		tagName = introTagNameDE
	}

	content = strings.ReplaceAll(content, "https://emvi.emvi.com", fmt.Sprintf("%s://%s.%s", frontendHostProtocol, orga.NameNormalized, frontendHostWithoutProtocol))
	articleContent := &model.ArticleContent{Title: title,
		Content:     content,
		ReadingTime: introArticleReadingTimeSeconds,
		ArticleId:   article.ID,
		LanguageId:  lang.ID,
		UserId:      userId}

	if err := model.SaveArticleContent(tx, articleContent); err != nil {
		logbuch.Error("Error saving introduction article content", logbuch.Fields{"err": err, "user_id": userId, "version": 0})
		return errs.Saving
	}

	// set organization owner as author
	articleContentAuthor := &model.ArticleContentAuthor{ArticleContentId: articleContent.ID, UserId: userId}

	if err := model.SaveArticleContentAuthor(tx, articleContentAuthor); err != nil {
		logbuch.Error("Error saving introduction article content author", logbuch.Fields{"err": err, "user_id": userId, "version": 0})
		return errs.Saving
	}

	articleContent.ID = 0
	articleContent.Version = 1

	if err := model.SaveArticleContent(tx, articleContent); err != nil {
		logbuch.Error("Error saving introduction article content", logbuch.Fields{"err": err, "user_id": userId, "version": 1})
		return errs.Saving
	}

	articleContentAuthor.ArticleContentId = articleContent.ID

	if err := model.SaveArticleContentAuthor(tx, articleContentAuthor); err != nil {
		logbuch.Error("Error saving introduction article content author", logbuch.Fields{"err": err, "user_id": userId, "version": 1})
		return errs.Saving
	}

	// create tag
	tag := &model.Tag{OrganizationId: orga.ID,
		Name:   tagName,
		Usages: 1}

	if err := model.SaveTag(tx, tag); err != nil {
		logbuch.Error("Error saving introduction tag", logbuch.Fields{"err": err, "user_id": userId})
		return errs.Saving
	}

	articleTag := &model.ArticleTag{ArticleId: article.ID, TagId: tag.ID}

	if err := model.SaveArticleTag(tx, articleTag); err != nil {
		logbuch.Error("Error saving introduction article tag", logbuch.Fields{"err": err, "user_id": userId})
		return errs.Saving
	}

	return nil
}
