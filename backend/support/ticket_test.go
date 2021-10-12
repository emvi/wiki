package support

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestContactSupport(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	mailsSend := 0
	var mailSubject, mailBody, mailFrom, mailTo string
	mailer := func(subject, body, from string, to ...string) error {
		mailsSend++
		mailSubject = subject
		mailBody = body
		mailFrom = from
		mailTo = to[0]
		return nil
	}

	if err := ContactSupport(orga, user.ID, "type_question", "subject", "message", mailer); err != nil {
		t.Fatalf("Expected support ticket to be created, but was: %v", err)
	}

	if mailsSend != 1 {
		t.Fatalf("Expected one mail to be send, but was: %v", mailsSend)
	}

	if mailSubject != "[Support] subject - test@user.com - test (Expert) - Firstname Lastname - Question" {
		t.Fatalf("Mail subject wrong: %v", mailSubject)
	}

	if mailBody != "message\n" {
		t.Fatalf("Body wrong: %v", mailBody)
	}

	if mailFrom != supportMailAddress {
		t.Fatalf("Sender wrong: %v", mailFrom)
	}

	if mailTo != supportMailAddress {
		t.Fatalf("Receiver wrong: %v", mailTo)
	}

	tickets := model.FindSupportTicketByOrganizationId(orga.ID)

	if len(tickets) != 1 {
		t.Fatalf("Expected one ticket to exist in database, but was: %v", len(tickets))
	}

	if tickets[0].OrganizationId != orga.ID ||
		tickets[0].UserId != user.ID ||
		tickets[0].Type != "type_question" ||
		tickets[0].Message != "message" ||
		tickets[0].Status != statusOpen {
		t.Fatalf("Ticket not as expected: %v", tickets[0])
	}
}

func TestContactSupportInvalidInput(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	mailsSend := 0
	mailer := func(subject, body, from string, to ...string) error {
		mailsSend++
		return nil
	}

	if err := ContactSupport(orga, user.ID, "", "", "message", mailer); len(err) != 1 || err[0] != errs.SubjectTooShort {
		t.Fatalf("Subject must be too short, but was: %v", err)
	}

	subject := ""

	for i := 0; i < 101; i++ {
		subject += "x"
	}

	if err := ContactSupport(orga, user.ID, "", subject, "message", mailer); len(err) != 1 || err[0] != errs.SubjectTooLong {
		t.Fatalf("Subject must be too long, but was: %v", err)
	}

	if err := ContactSupport(orga, user.ID, "", "subject", "", mailer); len(err) != 1 || err[0] != errs.MessageTooShort {
		t.Fatalf("Message must be too short, but was: %v", err)
	}

	message := ""

	for i := 0; i < 4001; i++ {
		message += "x"
	}

	if err := ContactSupport(orga, user.ID, "01234567890123456789012345678901234567891", "subject", message, mailer); len(err) != 1 || err[0] != errs.MessageTooLong {
		t.Fatalf("Message must be too long, but was: %v", err)
	}

	if mailsSend != 0 {
		t.Fatalf("No mails must have been send, but was: %v", mailsSend)
	}

	if err := ContactSupport(orga, user.ID, "01234567890123456789012345678901234567891", "subject", "okay", mailer); err != nil {
		t.Fatalf("Ticket must have been created, but was: %v", err)
	}

	tickets := model.FindSupportTicketByOrganizationId(orga.ID)

	if len(tickets) != 1 {
		t.Fatalf("Expected one ticket to exist in database, but was: %v", len(tickets))
	}

	if tickets[0].Type != "0123456789012345678901234567890123456789" {
		t.Fatalf("Type must be cut, but was: %v", tickets[0].Type)
	}

	if mailsSend != 1 {
		t.Fatalf("One mail must have been send, but was: %v", mailsSend)
	}
}

func TestGetSupportTicketReason(t *testing.T) {
	if getSupportTicketReason("foo") != "foo" {
		t.Fatal("Must not be mapped")
	}

	if getSupportTicketReason("type_question") != "Question" {
		t.Fatal("Must be mapped")
	}
}
