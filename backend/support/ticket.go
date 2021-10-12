package support

import (
	"bytes"
	"emviwiki/backend/errs"
	"emviwiki/backend/mailtpl"
	"emviwiki/shared/mail"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"html"
	"html/template"
	"strings"
	"unicode/utf8"
)

const (
	statusOpen    = "open"
	typeMaxLen    = 40
	subjectMaxLen = 100
	messageMaxLen = 4000
	subjectBase   = "[Support] "
)

var (
	supportTicketReasons = map[string]string{
		"type_question": "Question",
		"type_feature":  "Feature",
		"type_feedback": "Feedback",
		"type_issue":    "Issue",
		"type_billing":  "Billing",
	}
)

func ContactSupport(orga *model.Organization, userId hide.ID, ticketType, subject, message string, mailer mail.Sender) []error {
	ticketType = strings.TrimSpace(ticketType)
	subject = strings.TrimSpace(subject)
	message = strings.TrimSpace(message)

	if len(ticketType) > typeMaxLen {
		ticketType = ticketType[:typeMaxLen]
	}

	if err := validateSupportTicketRequest(subject, message); err != nil {
		return err
	}

	user := model.GetUserById(userId)

	if user == nil {
		logbuch.Error("Error finding user by id when creating support ticket", logbuch.Fields{"orga_id": orga.ID, "user_id": userId})
		return []error{errs.Saving}
	}

	if err := sendSupportTicketMail(orga, user, ticketType, subject, message, mailer); err != nil {
		logbuch.Error("Error sending support ticket mail", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId})
		return []error{err}
	}

	if err := createSupportTicket(orga, userId, ticketType, subject, message); err != nil {
		return []error{err}
	}

	return nil
}

func validateSupportTicketRequest(subject, message string) []error {
	err := make([]error, 0)

	if subject == "" {
		err = append(err, errs.SubjectTooShort)
	} else if utf8.RuneCountInString(subject) > subjectMaxLen {
		err = append(err, errs.SubjectTooLong)
	}

	if message == "" {
		err = append(err, errs.MessageTooShort)
	} else if utf8.RuneCountInString(message) > messageMaxLen {
		err = append(err, errs.MessageTooLong)
	}

	if len(err) != 0 {
		return err
	}

	return nil
}

func sendSupportTicketMail(orga *model.Organization, user *model.User, ticketType, subject, message string, mailer mail.Sender) error {
	tpl := mailtpl.Cache.Get()
	messageHTML := template.HTML(strings.Replace(html.EscapeString(message), "\n", "<br />", -1))
	data := struct {
		Message template.HTML
	}{
		messageHTML,
	}
	var buffer bytes.Buffer

	if err := tpl.ExecuteTemplate(&buffer, mailtpl.SupportMailTemplate, &data); err != nil {
		return errs.Saving
	}

	isExpert := ""

	if orga.Expert {
		isExpert = " (Expert)"
	}

	var sb strings.Builder
	sb.WriteString(subjectBase)
	sb.WriteString(subject)
	sb.WriteString(" - ")
	sb.WriteString(user.Email)
	sb.WriteString(" - ")
	sb.WriteString(orga.Name)
	sb.WriteString(isExpert)
	sb.WriteString(" - ")
	sb.WriteString(user.Firstname)
	sb.WriteString(" ")
	sb.WriteString(user.Lastname)
	sb.WriteString(" - ")
	sb.WriteString(getSupportTicketReason(ticketType))

	if err := mailer(sb.String(), buffer.String(), supportMailAddress, supportMailAddress); err != nil {
		return errs.Saving
	}

	return nil
}

func createSupportTicket(orga *model.Organization, userId hide.ID, ticketType, subject, message string) error {
	ticket := &model.SupportTicket{OrganizationId: orga.ID,
		UserId:  userId,
		Type:    ticketType,
		Subject: subject,
		Message: message,
		Status:  statusOpen}

	if err := model.SaveSupportTicket(nil, ticket); err != nil {
		logbuch.Error("Error saving support ticket", logbuch.Fields{"err": err, "orga_id": orga.ID, "user_id": userId})
		return errs.Saving
	}

	return nil
}

func getSupportTicketReason(reason string) string {
	mapped := supportTicketReasons[reason]

	if mapped == "" {
		return reason
	}

	return mapped
}
