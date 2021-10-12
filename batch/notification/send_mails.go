package notification

import (
	"bytes"
	"emviwiki/batch/errs"
	"emviwiki/shared/db"
	"emviwiki/shared/feed"
	"emviwiki/shared/i18n"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"github.com/emvi/logbuch"
	"html/template"
	"sync"
	"time"
)

const (
	sendNotificationConsumer     = 10
	sendNotificationMailTemplate = "mail_notifications.html"
	sendNotificationSubject      = "mail_notifications"
)

var notificationMailI18n = i18n.Translation{
	"en": {
		"title":    "Your unread notifications on Emvi",
		"text-1":   "Here are your unread notifications of the last",
		"text-2":   "days:",
		"text-3":   "Visit Emvi to view all notifications and mark them as read. You can always change the notification settings in the preferences.",
		"greeting": "Your unread notifications.",
		"goodbye":  "Cheers, Emvi Team",
	},
	"de": {
		"title":    "Deine ungelesenen Benachrichtigungen auf Emvi",
		"text-1":   "Hier sind deine Benachrichtigungen der letzten",
		"text-2":   "Tage:",
		"text-3":   "Besuche Emvi alle Benachrichtigungen einzusehen und um sie als gelesen zu markieren. Du kannst die Benachrichtigungseinstellungen jederzeit in den Einstellungen ändern.",
		"greeting": "Deine ungelesenen Benachrichtigungen.",
		"goodbye":  "Dein Emvi Team",
	},
}

type sendNotificationMailData struct {
	Member        *model.OrganizationMember
	User          *model.User
	Organization  *model.Organization
	Notifications []sendNotificationData
	OrgaURL       string
	EndVars       map[string]template.HTML
	Vars          map[string]template.HTML
}

type sendNotificationData struct {
	Feed *model.Feed
	Text template.HTML
	When time.Time
}

func SendNotificationMails() {
	sendChan := sendNotificationsProducer()
	var wg sync.WaitGroup
	wg.Add(sendNotificationConsumer)

	for i := 0; i < sendNotificationConsumer; i++ {
		go sendNotificationsConsumer(sendChan, &wg)
	}

	wg.Wait()
}

func sendNotificationsProducer() <-chan *model.OrganizationMember {
	logbuch.Debug("Users to send notification mails", logbuch.Fields{"count": model.CountOrganizationMemberWithNextNotificationMailReached()})
	memberRows, err := model.FindOrganizationMemberWithNextNotificationMailReachedCursor()

	if err != nil {
		logbuch.Fatal("Error reading users to notify", logbuch.Fields{"err": err})
	}

	sendChan := make(chan *model.OrganizationMember)

	go func() {
		for memberRows.Next() {
			var member model.OrganizationMember

			if err := memberRows.StructScan(&member); err != nil {
				logbuch.Fatal("Error scanning organization member", logbuch.Fields{"err": err})
				continue
			}

			sendChan <- &member
		}

		db.CloseRows(memberRows)
		close(sendChan)
	}()

	return sendChan
}

func sendNotificationsConsumer(sendChan <-chan *model.OrganizationMember, wg *sync.WaitGroup) {
	for member := range sendChan {
		if err := sendNotificationForMember(member); err != nil {
			logbuch.Error("Error sending notification mail to user", logbuch.Fields{"err": err, "user_id": member.UserId})
		}
	}

	wg.Done()
}

func sendNotificationForMember(member *model.OrganizationMember) error {
	if err := updateNextNotification(member); err != nil {
		return err
	}

	lang := util.DetermineLang(nil, member.OrganizationId, member.UserId, 0)
	mailData := getMailData(member, lang)

	if mailData == nil {
		return nil
	}

	logbuch.Debug("Notifications found for user", logbuch.Fields{"user_id": member.UserId, "notifications": len(mailData.Notifications)})
	langCode := util.DetermineSystemSupportedLangCode(member.OrganizationId, member.UserId)
	renderNotificationTexts(mailData, langCode)

	if err := renderAndSendMail(member, langCode, mailData); err != nil {
		return err
	}

	logbuch.Debug("Notification mail send for organization member", logbuch.Fields{"user_id": member.UserId})
	return nil
}

func getMailData(member *model.OrganizationMember, lang *model.Language) *sendNotificationMailData {
	orga := model.GetOrganizationById(member.OrganizationId)
	minDefTime := time.Now().Add(time.Hour * 24 * time.Duration(-member.SendNotificationsInterval))
	notifications := model.FindNotificationByOrganizationIdAndUserIdAndLanguageIdAndAfterDefTimeUnread(orga.ID, member.UserId, lang.ID, minDefTime)

	if len(notifications) == 0 {
		return nil
	}

	var notificationData []sendNotificationData

	for i := range notifications {
		notificationData = append(notificationData, sendNotificationData{
			Feed: &notifications[i],
			When: notifications[i].DefTime,
		})
	}

	langCode := util.DetermineSystemSupportedLangCode(member.OrganizationId, member.UserId)
	return &sendNotificationMailData{
		member,
		member.User,
		orga,
		notificationData,
		util.InjectSubdomain(frontendHost, orga.NameNormalized),
		i18n.GetMailEndI18n(langCode),
		i18n.GetVars(langCode, notificationMailI18n),
	}
}

func renderNotificationTexts(data *sendNotificationMailData, lang string) {
	for i := range data.Notifications {
		reason, ok := feed.Reasons[lang][data.Notifications[i].Feed.Reason]

		if ok {
			if reason.Notification == "" {
				reason.Notification = reason.Feed
			}

			data.Notifications[i].Text = template.HTML(feed.RenderFeed(data.Organization, reason.Notification, feed.NotificationText, lang, data.Notifications[i].Feed))
		}
	}
}

func renderAndSendMail(member *model.OrganizationMember, lang string, data *sendNotificationMailData) error {
	tpl := tplCache.Get()
	var buffer bytes.Buffer

	if err := tpl.ExecuteTemplate(&buffer, sendNotificationMailTemplate, data); err != nil {
		logbuch.Error("Error executing mail template", logbuch.Fields{"err": err})
		return err
	}

	subject := i18n.GetMailTitle(lang)[sendNotificationSubject]

	if err := mailProvider(subject, buffer.String(), member.User.Email); err != nil {
		return err
	}

	return nil
}

func updateNextNotification(member *model.OrganizationMember) error {
	member.NextNotificationMail = time.Now().Add(time.Hour * 24 * time.Duration(member.SendNotificationsInterval))

	if err := model.SaveOrganizationMember(nil, member); err != nil {
		return errs.Saving
	}

	return nil
}
