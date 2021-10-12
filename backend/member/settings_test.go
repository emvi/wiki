package member

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"emviwiki/shared/testutil"
	"testing"
)

func TestSaveSettings(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	phone := "0123456789012345678901234567891"
	mobile := "0123456789012345678901234567891"
	info := "01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567891"
	interval := uint(2)
	input := []Settings{
		Settings{Phone: &phone},
		Settings{Mobile: &mobile},
		Settings{Info: &info},
		Settings{SendNotificationsInterval: &interval},
	}
	expected := []error{
		errs.PhoneLen,
		errs.MobileLen,
		errs.InfoTooLong,
		errs.NotificationIntervalInvalid,
	}

	for i, in := range input {
		if err := SaveSettings(orga, user.ID, in); len(err) != 1 || err[0] != expected[i] {
			t.Fatalf("Expected error %v, but was %v", expected[i], err)
		}
	}

	if err := SaveSettings(orga, user.ID+1, Settings{}); len(err) != 1 || err[0] != errs.MemberNotFound {
		t.Fatalf("Member must not be found, but was: %v", err)
	}
}

func TestSaveSettingsSuccess(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	phone := "6543210"
	mobile := "0123456"
	info := "info"
	interval := uint(1)
	recommendation := false
	showCreateButton := false
	showNavigation := false
	showActionButtons := false
	data := Settings{
		Phone:                     &phone,
		Mobile:                    &mobile,
		Info:                      &info,
		SendNotificationsInterval: &interval,
		RecommendationMail:        &recommendation,
		ShowCreateButton:          &showCreateButton,
		ShowNavigation:            &showNavigation,
		ShowActionButtons:         &showActionButtons,
	}

	if err := SaveSettings(orga, user.ID, data); err != nil {
		t.Fatalf("Member must have been saved, but was: %v", err)
	}

	member := model.GetOrganizationMemberByOrganizationIdAndUserId(orga.ID, user.ID)

	if member == nil ||
		member.Phone.String != "6543210" ||
		member.Mobile.String != "0123456" ||
		member.Info.String != "info" ||
		member.SendNotificationsInterval != 1 ||
		member.RecommendationMail ||
		member.ShowCreateButton ||
		member.ShowNavigation ||
		member.ShowActionButtons {
		t.Fatalf("Member not as expected: %v", member)
	}
}

func TestSaveSettingsUpdateNil(t *testing.T) {
	testutil.CleanBackendDb(t)
	orga, user := testutil.CreateOrgaAndUser(t)
	user.OrganizationMember.Mobile.SetValid("75757575")

	if err := model.SaveOrganizationMember(nil, user.OrganizationMember); err != nil {
		t.Fatal(err)
	}

	phone := "42424242"
	data := Settings{
		Phone: &phone,
	}

	if err := SaveSettings(orga, user.ID, data); err != nil {
		t.Fatalf("Member must have been saved, but was: %v", err)
	}

	member := model.GetOrganizationMemberByOrganizationIdAndUserId(orga.ID, user.ID)

	if member.Phone.String != "42424242" ||
		member.Mobile.String != "75757575" {
		t.Fatalf("Phone must have been updated, but mobile must have been left alone, but was: %v %v", member.Phone, member.Mobile)
	}
}
