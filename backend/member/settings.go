package member

import (
	"emviwiki/backend/errs"
	"emviwiki/shared/model"
	"github.com/emvi/hide"
	"github.com/emvi/null"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	phoneMaxLen  = 30
	mobileMaxLen = 30
	infoMaxLen   = 100
)

var (
	notificationInterval = []uint{0, 1, 7, 30}
)

type Settings struct {
	Phone                     *string `json:"phone"`
	Mobile                    *string `json:"mobile"`
	Info                      *string `json:"info"`
	SendNotificationsInterval *uint   `json:"send_notifications_interval"`
	RecommendationMail        *bool   `json:"recommendation_mail"`
	ShowCreateButton          *bool   `json:"show_create_button"`
	ShowNavigation            *bool   `json:"show_navigation"`
	ShowActionButtons         *bool   `json:"show_action_buttons"`
}

func (data *Settings) validate() []error {
	err := make([]error, 0)

	if data.Phone != nil && utf8.RuneCountInString(*data.Phone) > phoneMaxLen {
		err = append(err, errs.PhoneLen)
	}

	if data.Mobile != nil && utf8.RuneCountInString(*data.Mobile) > mobileMaxLen {
		err = append(err, errs.MobileLen)
	}

	if data.Info != nil && utf8.RuneCountInString(*data.Info) > infoMaxLen {
		err = append(err, errs.InfoTooLong)
	}

	if data.SendNotificationsInterval != nil {
		intervalValid := false

		for _, i := range notificationInterval {
			if *data.SendNotificationsInterval == i {
				intervalValid = true
				break
			}
		}

		if !intervalValid {
			err = append(err, errs.NotificationIntervalInvalid)
		}
	}

	if len(err) == 0 {
		return nil
	}

	return err
}

// SaveSettings saves the personal settings for given member.
func SaveSettings(organization *model.Organization, userId hide.ID, data Settings) []error {
	data.Phone = trimString(data.Phone)
	data.Mobile = trimString(data.Mobile)
	data.Info = trimString(data.Info)

	if err := data.validate(); err != nil {
		return err
	}

	member := model.GetOrganizationMemberByOrganizationIdAndUserId(organization.ID, userId)

	if member == nil {
		return []error{errs.MemberNotFound}
	}

	if data.Phone != nil {
		member.Phone = null.NewString(*data.Phone, *data.Phone != "")
	}

	if data.Mobile != nil {
		member.Mobile = null.NewString(*data.Mobile, *data.Mobile != "")
	}

	if data.Info != nil {
		member.Info = null.NewString(*data.Info, *data.Info != "")
	}

	if data.SendNotificationsInterval != nil {
		member.SendNotificationsInterval = *data.SendNotificationsInterval
		member.NextNotificationMail = time.Now().Add(time.Hour * 24 * time.Duration(*data.SendNotificationsInterval))
	}

	if data.RecommendationMail != nil {
		member.RecommendationMail = *data.RecommendationMail
	}

	if data.ShowCreateButton != nil {
		member.ShowCreateButton = *data.ShowCreateButton
	}

	if data.ShowNavigation != nil {
		member.ShowNavigation = *data.ShowNavigation
	}

	if data.ShowActionButtons != nil {
		member.ShowActionButtons = *data.ShowActionButtons
	}

	if err := model.SaveOrganizationMember(nil, member); err != nil {
		return []error{errs.Saving}
	}

	return nil
}

func trimString(str *string) *string {
	if str != nil {
		trimmed := strings.TrimSpace(*str)
		return &trimmed
	}

	return nil
}
