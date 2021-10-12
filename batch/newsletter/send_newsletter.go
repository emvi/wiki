package newsletter

import (
	"bytes"
	dashboardmodel "emviwiki/dashboard/model"
	"emviwiki/shared/db"
	"emviwiki/shared/i18n"
	"emviwiki/shared/model"
	"github.com/emvi/logbuch"
	"github.com/emvi/null"
	"html/template"
	"sync"
)

const (
	sendNewsletterConsumerCount = 10
	sendNewsletterMailTemplate  = "mail_newsletter.html"
	sendNewsletterLang          = "en"
)

var newsletterMailI18n = i18n.Translation{
	"en": {
		"unsubscribe":           "You can cancel your subscription",
		"unsubscribe_link_text": "here",
	},
}

type receiver struct {
	Email                      string
	NewsletterSubscriptionCode null.String `db:"newsletter_subscription_code"`
}

func SendNewsletterMails() {
	newsletter := dashboardmodel.GetLatestNewsletterToSend()

	if newsletter == nil {
		logbuch.Info("No newsletter found")
		return
	}

	sendChan := sendNewsletterProducer()
	var wg sync.WaitGroup
	wg.Add(sendNewsletterConsumerCount)

	for i := 0; i < sendNewsletterConsumerCount; i++ {
		go sendNewsletterConsumer(sendChan, newsletter, &wg)
	}

	wg.Wait()
	markNewsletterSend(newsletter)
}

func sendNewsletterProducer() <-chan receiver {
	logbuch.Debug("Users to send newsletter mails", logbuch.Fields{"count": model.CountNewsletterSubscriptions()})
	emailRows, err := model.FindNewsletterSubscriptionEmails()

	if err != nil {
		logbuch.Fatal("Error reading users to send newsletter to", logbuch.Fields{"err": err})
	}

	sendChan := make(chan receiver)

	go func() {
		for emailRows.Next() {
			var receiver receiver

			if err := emailRows.StructScan(&receiver); err != nil {
				logbuch.Fatal("Error scanning email", logbuch.Fields{"err": err})
				continue
			}

			sendChan <- receiver
		}

		db.CloseRows(emailRows)
		close(sendChan)
	}()

	return sendChan
}

func sendNewsletterConsumer(sendChan <-chan receiver, newsletter *dashboardmodel.Newsletter, wg *sync.WaitGroup) {
	for receiver := range sendChan {
		if err := sendNewsletter(newsletter, receiver.Email, receiver.NewsletterSubscriptionCode.String); err != nil {
			logbuch.Error("Error sending newsletter", logbuch.Fields{"err": err, "email": receiver.Email})
		}
	}

	wg.Done()
}

func renderNewsletter(newsletter *dashboardmodel.Newsletter, newsletterSubscriptionCode string) (string, string, error) {
	var buffer bytes.Buffer
	templateVars := struct {
		Vars                       map[string]template.HTML
		EndVars                    map[string]template.HTML
		Subject                    string
		Content                    template.HTML
		URLunsub                   string
		NewsletterSubscriptionCode string
	}{
		i18n.GetVars(sendNewsletterLang, newsletterMailI18n),
		i18n.GetMailEndI18n(sendNewsletterLang),
		newsletter.Subject,
		template.HTML(newsletter.Content),
		newsletterUnsubscribeURI,
		newsletterSubscriptionCode,
	}

	if err := tplCache.Get().ExecuteTemplate(&buffer, sendNewsletterMailTemplate, templateVars); err != nil {
		logbuch.Error("Error rendering newsletter mail template", logbuch.Fields{"err": err})
		return "", "", err
	}

	return newsletter.Subject, buffer.String(), nil
}

func sendNewsletter(newsletter *dashboardmodel.Newsletter, email, newsletterSubscriptionCode string) error {
	subject, body, err := renderNewsletter(newsletter, newsletterSubscriptionCode)

	if err != nil {
		logbuch.Fatal("Error rendering newsletter mail", logbuch.Fields{"err": err})
	}

	if err := mailProvider(subject, body, email); err != nil {
		return err
	}

	return nil
}

func markNewsletterSend(newsletter *dashboardmodel.Newsletter) {
	newsletter.Status = "send"

	if err := dashboardmodel.SaveNewsletter(nil, newsletter); err != nil {
		logbuch.Fatal("Error updating newsletter status to 'send'", logbuch.Fields{"err": err, "newsletter_id": newsletter.ID})
	}
}
