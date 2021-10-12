package pages

import (
	"bytes"
	"emviwiki/dashboard/auth"
	"emviwiki/dashboard/model"
	"emviwiki/shared/i18n"
	"emviwiki/shared/rest"
	"errors"
	"fmt"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strings"
	"time"
)

const (
	scheduledDateFormat     = "2006-01-02"
	newsletterStatusPlanned = "planned"
	newsletterStatusSend    = "send"
)

var newsletterMailI18n = i18n.Translation{
	"en": {
		"unsubscribe":           "You can cancel your subscription",
		"unsubscribe_link_text": "here",
	},
}

func NewsletterEditPageHandler(claims *auth.UserTokenClaims, w http.ResponseWriter, r *http.Request) {
	data := struct {
		ID        hide.ID
		Subject   string
		Scheduled string
		Content   string
		SaveError string
	}{0, "", "", "", ""}

	params := mux.Vars(r)
	var id hide.ID
	var err error

	if idStr, ok := params["id"]; ok {
		id, err = hide.FromString(idStr)

		if err != nil {
			data.SaveError = err.Error()
			RenderPage(w, newsletterEditPageTemplate, claims, &data)
			return
		}

		data.ID = id
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			data.SaveError = err.Error()
			RenderPage(w, newsletterEditPageTemplate, claims, &data)
			return
		}

		subject := strings.TrimSpace(r.Form.Get("subject"))
		scheduled := strings.TrimSpace(r.Form.Get("scheduled"))
		content := strings.TrimSpace(r.Form.Get("content"))
		id, err = saveNewsletter(id, subject, scheduled, content)

		if err != nil {
			data.Subject = subject
			data.Scheduled = scheduled
			data.Content = content
			data.SaveError = err.Error()
			RenderPage(w, newsletterEditPageTemplate, claims, &data)
			return
		}

		idStr, _ := hide.ToString(id)
		http.Redirect(w, r, fmt.Sprintf("/newsletter/edit/%s", idStr), http.StatusFound)
		return
	} else if r.Method == http.MethodPut {
		req := struct {
			Email string `json:"email"`
		}{}

		if err := rest.DecodeJSON(r, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		req.Email = strings.TrimSpace(req.Email)

		if req.Email == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := sendTestMail(id, req.Email); err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		return
	}

	if id != 0 {
		newsletter := model.GetNewsletterById(id)

		if newsletter == nil {
			data.SaveError = "Newsletter not found."
		} else {
			data.Subject = newsletter.Subject
			data.Scheduled = newsletter.Scheduled.Format(scheduledDateFormat)
			data.Content = newsletter.Content
		}
	}

	RenderPage(w, newsletterEditPageTemplate, claims, &data)
}

func saveNewsletter(id hide.ID, subject, scheduled, content string) (hide.ID, error) {
	if len(subject) == 0 || len(scheduled) == 0 || len(content) == 0 {
		return 0, errors.New("Subject, scheduled date and content must be set.")
	}

	scheduledDate, err := time.Parse(scheduledDateFormat, scheduled)

	if err != nil {
		return 0, err
	}

	if scheduledDate.Before(time.Now().Add(-time.Hour * 24)) {
		return 0, errors.New("Scheduled date must be today or in the future.")
	}

	var newsletter *model.Newsletter

	if id == 0 {
		newsletter = &model.Newsletter{Subject: subject,
			Scheduled: scheduledDate,
			Content:   content}
	} else {
		newsletter = model.GetNewsletterById(id)

		if newsletter == nil {
			return 0, errors.New("Newsletter not found.")
		}

		if newsletter.Status == newsletterStatusSend {
			return 0, errors.New("Newsletter has been send already.")
		}
	}

	newsletter.Subject = subject
	newsletter.Status = newsletterStatusPlanned
	newsletter.Scheduled = scheduledDate
	newsletter.Content = content

	if err := model.SaveNewsletter(nil, newsletter); err != nil {
		return 0, err
	}

	return newsletter.ID, nil
}

func sendTestMail(id hide.ID, email string) error {
	newsletter := model.GetNewsletterById(id)

	if newsletter == nil {
		logbuch.Error("Error finding test mail, newsletter by id not found")
		return errors.New("Newsletter not found.")
	}

	var buffer bytes.Buffer
	data := struct {
		Vars                       map[string]template.HTML
		EndVars                    map[string]template.HTML
		Subject                    string
		Content                    template.HTML
		URLunsub                   string
		NewsletterSubscriptionCode string
	}{
		i18n.GetVars("en", newsletterMailI18n),
		i18n.GetMailEndI18n("en"),
		newsletter.Subject,
		template.HTML(newsletter.Content),
		"unsub",
		"code",
	}

	if err := mailTplCache.Get().ExecuteTemplate(&buffer, newsletterMailTemplate, &data); err != nil {
		logbuch.Error("Error rendering test mail", logbuch.Fields{"err": err})
		return err
	}

	if err := mailProvider(newsletter.Subject, buffer.String(), email); err != nil {
		logbuch.Error("Error sending test mail", logbuch.Fields{"err": err})
		return err
	}

	return nil
}
