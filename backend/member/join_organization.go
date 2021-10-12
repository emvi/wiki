package member

import (
	"bytes"
	"emviwiki/backend/billing"
	"emviwiki/backend/errs"
	"emviwiki/backend/feed"
	"emviwiki/backend/mailtpl"
	"emviwiki/backend/usergroup"
	"emviwiki/shared/constants"
	"emviwiki/shared/db"
	"emviwiki/shared/i18n"
	"emviwiki/shared/mail"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"fmt"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/jmoiron/sqlx"
	"html/template"
	"strings"
	"time"
	"unicode"
)

const (
	usernameMinLen                      = 3
	usernameMaxLen                      = 20
	defaultSendNotificationIntervalDays = 7
	mailJoinedSubject                   = "joined"
)

var memberJoinedMailI18n = i18n.Translation{
	"en": {
		"title-1":  "joined your organization",
		"title-2":  "",
		"joined-1": "joined your organization",
		"joined-2": ".",
		"text":     "The member joined by using the public invitation link. In case this happened unintentionally, remove the member and generate a new link in the administration settings.",
		"greeting": "A member joined your organization.",
		"goodbye":  "Cheers, Emvi Team",
	},
	"de": {
		"title-1":  "ist deiner Organisation",
		"title-2":  " beigetreten",
		"joined-1": "ist deiner Organisation",
		"joined-2": " beigetreten.",
		"text":     "Das Mitglied ist über den öffentlichen Einladungslink beigetreten. Sollte dies nicht beabsichtigt gewesen sein, entferne das Mitglied und generiere einen neuen Link in den Administrationseinstellungen.",
		"greeting": "Ein Mitglied ist deiner Organisation beigetreten.",
		"goodbye":  "Dein Emvi Team",
	},
}

type JoinOrganizationData struct {
	Username       string `json:"username"`
	Code           string `json:"code"`            // for invitation via mail
	InvitationCode string `json:"invitation_code"` // for organization
}

func (data *JoinOrganizationData) validate(user *model.User, orgaId hide.ID) error {
	if model.GetOrganizationMemberByOrganizationIdAndUserId(orgaId, user.ID) != nil {
		return errs.IsMemberAlready
	}

	data.Username = strings.TrimSpace(data.Username)

	if err := checkValidUsername(orgaId, data.Username); err != nil {
		return err
	}

	return nil
}

func checkValidUsername(orgaId hide.ID, username string) error {
	usernameRunes := []rune(username)
	n := len(usernameRunes)

	if n < usernameMinLen {
		return errs.UsernameTooShort
	}

	if n > usernameMaxLen {
		return errs.UsernameTooLong
	}

	if (!unicode.IsLetter(usernameRunes[0]) && !unicode.IsNumber(usernameRunes[0])) ||
		(!unicode.IsLetter(usernameRunes[n-1]) && !unicode.IsNumber(usernameRunes[n-1])) ||
		usernameContainsInvalidCharOrSpecialCharInRow(usernameRunes) {
		return errs.UsernameInvalid
	}

	if model.GetOrganizationMemberByOrganizationIdAndUsername(orgaId, username) != nil {
		return errs.UsernameInUse
	}

	return nil
}

func JoinOrganization(userId hide.ID, data JoinOrganizationData, mailer mail.Sender) error {
	user := model.GetUserById(userId)

	if user == nil {
		return errs.UserNotFound
	}

	data.Code = strings.TrimSpace(data.Code)
	data.InvitationCode = strings.TrimSpace(data.InvitationCode)
	orgaId, readOnly, err := findOrganization(data, user.Email)

	if err != nil {
		return err
	}

	if err := data.validate(user, orgaId); err != nil {
		return err
	}

	tx, err := model.GetConnection().Beginx()

	if err != nil {
		logbuch.Error("Error starting transaction to join organization", logbuch.Fields{"err": err})
		return errs.TxBegin
	}

	if err := model.DeleteInvitationByOrganizationIdAndEmail(tx, orgaId, user.Email); err != nil {
		return errs.Saving
	}

	if err := createMember(tx, orgaId, userId, data.Username, readOnly); err != nil {
		return err
	}

	if err := joinDefaultUserGroups(tx, orgaId, userId, readOnly); err != nil {
		return err
	}

	// if the user was a member before, we don't want him to see old notifications when he had no access
	if err := markNotificationsRead(tx, orgaId, userId); err != nil {
		return err
	}

	if err := createJoinOrganizationFeed(tx, orgaId, user); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		logbuch.Error("Error committing transaction when joining organization", logbuch.Fields{"err": err})
		return errs.TxCommit
	}

	// joined by organization invitation code, inform admins
	if data.InvitationCode != "" {
		go informAdminsMemberJoined(orgaId, user, mailer)
	}

	go func() {
		if err := billing.UpdateSubscription(model.GetOrganizationById(orgaId)); err != nil {
			logbuch.Error("Error updating subscription while joining organization", logbuch.Fields{"err": err, "orga_id": orgaId})
		}
	}()

	return nil
}

func findOrganization(data JoinOrganizationData, email string) (hide.ID, bool, error) {
	var orgaId hide.ID
	readOnly := false

	if data.InvitationCode != "" {
		orga := model.GetOrganizationByInvitationCode(data.InvitationCode)

		if orga == nil {
			return 0, false, errs.OrganizationNotFound
		}

		orgaId = orga.ID
		readOnly = orga.InvitationReadOnly
		logbuch.Debug("Found invitation for invitation code", logbuch.Fields{"orga_id": orgaId, "read_only": readOnly})
	} else {
		invitation := model.GetInvitationByEmailAndCode(email, data.Code)

		if invitation == nil {
			return 0, false, errs.InvitationNotFound
		}

		orgaId = invitation.OrganizationId
		readOnly = invitation.ReadOnly
		logbuch.Debug("Found invitation for email", logbuch.Fields{"orga_id": orgaId, "read_only": readOnly})
	}

	return orgaId, readOnly, nil
}

func usernameContainsInvalidCharOrSpecialCharInRow(username []rune) bool {
	for i, c := range username {
		if (!unicode.IsLetter(c) && !unicode.IsNumber(c) && !usernameIsSpecialChar(c)) ||
			(i > 0 && usernameIsSpecialChar(username[i-1]) && usernameIsSpecialChar(c)) {
			return true
		}
	}

	return false
}

func usernameIsSpecialChar(c rune) bool {
	return c == '.' || c == '_' || c == '-'
}

func createMember(tx *sqlx.Tx, orgaId, userId hide.ID, username string, readOnly bool) error {
	defaultLang := model.GetDefaultLanguageByOrganizationIdTx(tx, orgaId)

	if defaultLang == nil {
		logbuch.Error("Default language for organization not found when joining organization", logbuch.Fields{"user_id": userId, "orga_id": orgaId})
		db.Rollback(tx)
		return errs.Saving
	}

	// look for existing user (in case he got removed before) and update username or create completely new user
	member := model.GetOrganizationMemberByOrganizationIdAndUserIdIgnoreActiveTx(tx, orgaId, userId)

	if member == nil {
		member = &model.OrganizationMember{OrganizationId: orgaId,
			UserId:                    userId,
			LanguageId:                defaultLang.ID,
			Username:                  username,
			SendNotificationsInterval: defaultSendNotificationIntervalDays,
			NextNotificationMail:      time.Now().Add(time.Hour * 24 * defaultSendNotificationIntervalDays),
			RecommendationMail:        true,
			ReadOnly:                  readOnly,
			Active:                    true,
			ShowCreateButton:          true,
			ShowActionButtons:         true,
			ShowNavigation:            true}
	} else {
		member.Username = username
		member.Active = true
		member.ReadOnly = readOnly
	}

	if err := model.SaveOrganizationMember(tx, member); err != nil {
		return errs.Saving
	}

	return nil
}

func joinDefaultUserGroups(tx *sqlx.Tx, orgaId, userId hide.ID, readOnly bool) error {
	if err := joinDefaultUserGroup(tx, orgaId, userId, constants.GroupAllName); err != nil {
		return err
	}

	if readOnly {
		if err := joinDefaultUserGroup(tx, orgaId, userId, constants.GroupReadOnlyName); err != nil {
			return err
		}
	}

	return nil
}

func joinDefaultUserGroup(tx *sqlx.Tx, orgaId, userId hide.ID, name string) error {
	group := model.GetUserGroupByOrganizationIdAndNameTx(tx, orgaId, name)

	if group == nil {
		logbuch.Error("Default group to join not found", logbuch.Fields{"orga_id": orgaId, "user_id": userId, "name": name})
		db.Rollback(tx)
		return errs.Saving
	}

	member := model.GetUserGroupMemberByGroupIdAndUserIdTx(tx, group.ID, userId)

	if member == nil {
		member = &model.UserGroupMember{UserGroupId: group.ID, UserId: userId}

		if err := model.SaveUserGroupMember(tx, member); err != nil {
			logbuch.Error("Error saving user group member when joining organization", logbuch.Fields{"orga_id": orgaId, "user_id": userId})
			return errs.Saving
		}
	}

	return nil
}

func markNotificationsRead(tx *sqlx.Tx, orgaId, userId hide.ID) error {
	orga := model.GetOrganizationByIdTx(tx, orgaId)

	if orga == nil {
		db.Rollback(tx)
		return errs.OrganizationNotFound
	}

	if err := feed.ToggleNotificationRead(tx, orga, userId, 0); err != nil {
		logbuch.Error("Error marking old notifications read", logbuch.Fields{"err": err, "orga_id": orgaId, "user_id": userId})
	}

	return nil
}

func createJoinOrganizationFeed(tx *sqlx.Tx, orgaId hide.ID, user *model.User) error {
	orga := model.GetOrganizationByIdTx(tx, orgaId)

	if orga == nil {
		db.Rollback(tx)
		return errs.OrganizationNotFound
	}

	feedData := &feed.CreateFeedData{Tx: tx,
		Organization: orga,
		UserId:       user.ID,
		Reason:       "joined_organization",
		Public:       true}

	if err := feed.CreateFeed(feedData); err != nil {
		logbuch.Error("Error creating feed when joining organization", logbuch.Fields{"err": err})
		return err
	}

	return nil
}

func informAdminsMemberJoined(orgaId hide.ID, user *model.User, mailer mail.Sender) {
	orga := model.GetOrganizationById(orgaId)
	adminGroup := usergroup.GetAdminGroup(nil, orgaId)
	admins := model.FindUserGroupMemberUserByOrganizationIdAndUserGroupId(orgaId, adminGroup.ID)

	for _, admin := range admins {
		subject, body, err := renderMemberJoinedMail(orga, user)

		if err != nil {
			return
		}

		if err := mailer(subject, body, admin.Email); err != nil {
			logbuch.Error("Error sending member joined mail", logbuch.Fields{"err": err, "orga_id": orga.ID, "email": admin.Email})
		} else {
			logbuch.Debug("Send member joined mail", logbuch.Fields{"orga_id": orga.ID, "receiver": admin.Email})
		}
	}
}

func renderMemberJoinedMail(orga *model.Organization, user *model.User) (string, string, error) {
	tpl := mailtpl.Cache.Get()
	lang := util.DetermineSystemSupportedLangCode(orga.ID, user.ID)
	subject := fmt.Sprintf(i18n.GetMailTitle(lang)[mailJoinedSubject], user.Firstname, user.Lastname)
	data := struct {
		EndVars      map[string]template.HTML
		Vars         map[string]template.HTML
		WebsiteHost  string
		AuthHost     string
		FrontendHost string
		Organization *model.Organization
		User         *model.User
	}{
		i18n.GetMailEndI18n(lang),
		i18n.GetVars(lang, memberJoinedMailI18n),
		websiteHost,
		authHost,
		util.InjectSubdomain(frontendHost, orga.NameNormalized),
		orga,
		user,
	}
	var buffer bytes.Buffer

	if err := tpl.ExecuteTemplate(&buffer, mailtpl.MemberJoinedMailTemplate, &data); err != nil {
		logbuch.Error("Error rendering mail template", logbuch.Fields{"err": err, "template": mailtpl.MemberJoinedMailTemplate})
		return "", "", err
	}

	return subject, buffer.String(), nil
}
