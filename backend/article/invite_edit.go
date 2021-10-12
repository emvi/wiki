package article

import (
	"emviwiki/backend/errs"
	"emviwiki/backend/mailtpl"
	"emviwiki/shared/i18n"
	"emviwiki/shared/mail"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"strings"
)

const (
	inviteArticleFeed    = "invite_article"
	inviteNewArticleFeed = "invite_new_article"
	inviteMailSubject    = "invite_article"
	invitePath           = "edit"
)

var inviteMailI18n = i18n.Translation{
	"en": {
		"title-1":  "invited you to edit",
		"title-2":  "on Emvi.",
		"title-3":  "invited you to edit a new article on Emvi.",
		"action":   "Open Article",
		"link":     "Or paste this link into your browser",
		"greeting": "You have been invited to edit an article.",
		"goodbye":  "Cheers, Emvi Team",
	},
	"de": {
		"title-1":  "hat Dich eingeladen",
		"title-2":  "auf Emvi zu bearbeiten.",
		"title-3":  "hat Dich eingeladen einen neuen Artikel auf Emvi zu bearbeiten.",
		"action":   "Artikel öffnen",
		"link":     "Oder kopiere diesen Link in deinen Browser",
		"greeting": "Du wurdest eingeladen einen Artikel zu bearbeiten.",
		"goodbye":  "Dein Emvi Team",
	},
}

func InviteEditArticle(orga *model.Organization, userId, articleId, langId hide.ID, roomId, message string, notifyUserIds, notifyGroupIds []hide.ID, mailer mail.Sender) error {
	message = strings.TrimSpace(message)

	if err := checkMessageLen(message); err != nil {
		return err
	}

	roomId = strings.TrimSpace(roomId)
	var article *model.Article

	if articleId != 0 {
		var err error
		article, err = checkUserReadAccess(orga.ID, userId, articleId)

		if err != nil {
			return err
		}

		// exact match
		article.LatestArticleContent = model.GetArticleContentLatestByArticleIdAndLanguageId(articleId, langId, false)

		if article.LatestArticleContent == nil {
			return errs.LanguageNotFound
		}
	}

	user := model.GetUserById(userId)

	if user == nil {
		return errs.UserNotFound
	}

	notifyUsers, err := checkAllUserExist(orga.ID, notifyUserIds)

	if err != nil {
		return err
	}

	feedReason := inviteArticleFeed

	if article == nil {
		feedReason = inviteNewArticleFeed
	}

	notifyUsers = appendGroupMembers(orga.ID, notifyUsers, notifyGroupIds)
	createRecommendInviteArticleFeed(orga, userId, notifyUsers, article, roomId, message, feedReason)
	sendRecommendInviteArticleMail(orga, article, roomId, message, user, notifyUsers, mailtpl.InviteArticleMailTemplate, inviteMailSubject, invitePath, inviteMailI18n, mailer)
	return nil
}
