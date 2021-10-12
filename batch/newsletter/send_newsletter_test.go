package newsletter

import (
	"emviwiki/dashboard/model"
	"emviwiki/shared/testutil"
	"strings"
	"testing"
	"time"
)

func TestRenderNewsletter(t *testing.T) {
	testutil.CleanBackendDb(t)
	testutil.CleanDashboardDb(t)
	newsletter := createNewsletter(t, time.Now().Add(-time.Minute), "planned")

	if subject, body, err := renderNewsletter(newsletter, "code"); err != nil || subject == "" || body == "" {
		t.Fatalf("Expected newsletter to be rendered, but was: %v", err)
	}
}

func TestSendNewsletterMailsFindNewsletter(t *testing.T) {
	testutil.CleanBackendDb(t)
	testutil.CleanDashboardDb(t)
	input := []struct {
		scheduled time.Time
		status    string
	}{
		{time.Now().Add(time.Minute), "planned"},
		{time.Now().Add(time.Minute), "send"},
		{time.Now().Add(-time.Minute), "send"},
		{time.Now().Add(-time.Minute), "planned"},
	}
	expected := []bool{false, false, false, true}

	for i, in := range input {
		createNewsletter(t, in.scheduled, in.status)
		found := model.GetLatestNewsletterToSend()

		if (found != nil) != expected[i] {
			t.Fatalf("Expected %v newsletters to be found, but was: %v", expected[i], found)
		}
	}
}

func TestMarkNewsletterSend(t *testing.T) {
	testutil.CleanDashboardDb(t)
	newsletter := createNewsletter(t, time.Now().Add(-time.Minute), "planned")
	markNewsletterSend(newsletter)
	newsletter = model.GetNewsletterById(newsletter.ID)

	if newsletter.Status != "send" {
		t.Fatalf("Status must have been updated, but was: %v", newsletter.Status)
	}
}

func TestRenderNewsletterMail(t *testing.T) {
	newsletter := &model.Newsletter{
		Subject: "mail-subject",
		Content: "mail-content",
	}
	subject, body, err := renderNewsletter(newsletter, "mail-subcode")

	if err != nil {
		t.Fatalf("Mail must have been rendered, but was: %v", err)
	}

	if subject != "mail-subject" {
		t.Fatalf("Subject not as expected: %v", subject)
	}

	if !strings.Contains(body, "mail-content") ||
		!strings.Contains(body, "mail-subcode") {
		t.Fatalf("Body not as expected: %v", body)
	}
}

func createNewsletter(t *testing.T, scheduled time.Time, status string) *model.Newsletter {
	newsletter := &model.Newsletter{Subject: "subject",
		Content:   "content",
		Scheduled: scheduled,
		Status:    status}

	if err := model.SaveNewsletter(nil, newsletter); err != nil {
		t.Fatal(err)
	}

	return newsletter
}
