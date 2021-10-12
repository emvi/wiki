package member

import (
	"bytes"
	"emviwiki/backend/errs"
	"emviwiki/backend/mailtpl"
	"emviwiki/backend/perm"
	"emviwiki/shared/i18n"
	"emviwiki/shared/mail"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"html/template"
	"strings"
)

const (
	codeLen           = 32
	mailInviteSubject = "invite"
)

var inviteNewUserMailI18n = i18n.Translation{
	"en": {
		"title-1":  "invited you to join",
		"title-2":  "on Emvi.",
		"text":     "You can start by following the link below.",
		"action":   "Join Now",
		"link":     "Or paste this link into your browser",
		"greeting": "You have been invited to join an organization.",
		"goodbye":  "Cheers, Emvi Team",
	},
	"de": {
		"title-1":  "hat Dich eingeladen",
		"title-2":  "auf Emvi beizutreten.",
		"text":     "Klicke auf den folgenden Link um zu Starten.",
		"action":   "Jetzt Beitreten",
		"link":     "Oder kopiere diesen Link in deinen Browser",
		"greeting": "Du wurdest in eine Organisation eingeladen.",
		"goodbye":  "Dein Emvi Team",
	},
}

type InviteMemberData struct {
	Emails   []string `json:"emails"`
	ReadOnly bool     `json:"read_only"`
}

func (data *InviteMemberData) removeDuplicates(filter string) {
	filter = strings.ToLower(filter)
	clean := make([]string, 0)

	for _, email := range data.Emails {
		email = strings.ToLower(email)

		if !containsString(clean, email) && email != filter {
			clean = append(clean, email)
		}
	}

	data.Emails = clean
}

func (data *InviteMemberData) validate() error {
	for _, email := range data.Emails {
		if !mail.EmailValid(email) {
			return errs.EmailInvalid
		}
	}

	return nil
}

func InviteMember(orga *model.Organization, userId hide.ID, data InviteMemberData, mailer mail.Sender) error {
	if _, err := perm.CheckUserIsAdmin(orga.ID, userId); err != nil {
		return err
	}

	user := model.GetUserById(userId)

	if user == nil {
		return errs.UserNotFound
	}

	data.removeDuplicates(user.Email)

	if err := data.validate(); err != nil {
		return err
	}

	// inviting members as read only is an expert feature
	if !orga.Expert {
		data.ReadOnly = false
	}

	language := util.DetermineLang(nil, orga.ID, userId, 0)

	for _, email := range data.Emails {
		sendInvitationForNewUser(orga, user, email, language, data.ReadOnly, mailer)
	}

	return nil
}

func containsString(haystack []string, needle string) bool {
	for _, entry := range haystack {
		if entry == needle {
			return true
		}
	}

	return false
}

func sendInvitationForNewUser(orga *model.Organization, user *model.User, email string, language *model.Language, readOnly bool, mailer mail.Sender) {
	code, err := createInvitation(orga.ID, email, readOnly)

	if err != nil {
		return
	}

	subject := i18n.GetMailTitle(language.Code)[mailInviteSubject]
	body, err := renderInviteNewUser(orga, user, email, code)

	if err != nil {
		return
	}

	if err := mailer(subject, body, email); err != nil {
		logbuch.Error("Error sending invitation mail for new user", logbuch.Fields{"err": err, "orga_id": orga.ID, "email": email})
	}
}

func createInvitation(orgaId hide.ID, email string, readOnly bool) (string, error) {
	if err := model.DeleteInvitationByOrganizationIdAndEmail(nil, orgaId, email); err != nil {
		logbuch.Error("Error cleaning up old organization invitations", logbuch.Fields{"err": err, "email": email})
	}

	code := util.GenRandomString(codeLen)
	invitation := &model.Invitation{OrganizationId: orgaId,
		Email:    email,
		Code:     code,
		ReadOnly: readOnly}

	if err := model.SaveInvitation(nil, invitation); err != nil {
		// nothing special to log here, model does this for us
		return "", errs.Saving
	}

	return code, nil
}

func renderInviteNewUser(orga *model.Organization, user *model.User, email, code string) (string, error) {
	tpl := mailtpl.Cache.Get()
	lang := util.DetermineSystemSupportedLangCode(orga.ID, user.ID)
	data := struct {
		EndVars      map[string]template.HTML
		Vars         map[string]template.HTML
		WebsiteHost  string
		AuthHost     string
		FrontendHost string
		Organization *model.Organization
		User         *model.User
		Email        string
		URL          string
		Code         string
	}{
		i18n.GetMailEndI18n(lang),
		inviteNewUserMailI18n[lang],
		websiteHost,
		authHost,
		util.InjectSubdomain(frontendHost, orga.NameNormalized),
		orga,
		user,
		email,
		websiteHost + "/join",
		code,
	}
	var buffer bytes.Buffer

	if err := tpl.ExecuteTemplate(&buffer, mailtpl.InviteUserMailTemplate, &data); err != nil {
		logbuch.Error("Error rendering mail template", logbuch.Fields{"err": err, "template": mailtpl.InviteUserMailTemplate})
		return "", err
	}

	return buffer.String(), nil
}
