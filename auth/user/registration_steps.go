package user

import (
	"emviwiki/auth/errs"
	"emviwiki/auth/model"
	"emviwiki/shared/mail"
	"github.com/emvi/logbuch"
	"strings"
	"unicode/utf8"
)

const (
	StepInitial = iota
	StepPassword
	StepPersonalData
	StepCompletion
	StepDone
)

const (
	firstnameMaxLength = 40
	lastnameMaxLength  = 40
	emailMaxLength     = 255
)

type RegistrationData struct {
	Email string `json:"email"`
}

func (data *RegistrationData) validate() error {
	data.Email = strings.TrimSpace(data.Email)

	if err := data.validateEmail(); err != nil {
		return err
	}

	// consider active users only because this is also used to resend mails
	if model.GetUserByEmail(data.Email) != nil {
		return errs.EmailInUse
	}

	return nil
}

func (data *RegistrationData) validateEmail() error {
	if data.Email == "" {
		return errs.EmailEmpty
	} else if utf8.RuneCountInString(data.Email) > emailMaxLength {
		return errs.EmailTooLong
	}

	if !mail.EmailValid(data.Email) {
		return errs.EmailInvalid
	}

	parts := strings.Split(data.Email, "@")

	if model.GetEmailBlacklistByDomain(parts[1]) != nil {
		return errs.EmailForbidden
	}

	return nil
}

type RegistrationPasswordData struct {
	Code           string `json:"code"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"password_repeat"`
}

func (data *RegistrationPasswordData) validate() error {
	if data.Password == "" {
		return errs.PasswordInvalid
	}

	if data.Password != data.PasswordRepeat {
		return errs.PasswordMatch
	}

	return nil
}

type RegistrationPersonalData struct {
	Code      string `json:"code"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Language  string `json:"language"` // optional
}

func (data *RegistrationPersonalData) validate() []error {
	userData := UserData{data.Firstname, data.Lastname, data.Language, false}
	return userData.Validate()
}

type RegistrationCompletionData struct {
	Code                 string `json:"code"`
	AcceptTermsOfService bool   `json:"accept_terms_of_service"`
	AcceptPrivacy        bool   `json:"accept_privacy"`
	AcceptMarketing      bool   `json:"accept_marketing"`
	RecaptchaToken       string `json:"recaptcha_token"`
}

func (data *RegistrationCompletionData) validate() []error {
	err := make([]error, 0)

	if !data.AcceptTermsOfService {
		err = append(err, errs.TermsOfService)
	}

	if !data.AcceptPrivacy {
		err = append(err, errs.Privacy)
	}

	if e := data.validateRecaptcha(); e != nil {
		err = append(err, e)
	}

	if len(err) != 0 {
		return err
	}

	return nil
}

func (data *RegistrationCompletionData) validateRecaptcha() error {
	resp, err := recaptchaValidator.Validate(data.RecaptchaToken)

	if err != nil {
		logbuch.Debug("Error on recaptcha request", logbuch.Fields{"err": err})
		return errs.ReCaptchaValidationErr
	}

	if !resp.Success {
		logbuch.Debug("Recaptcha request result", logbuch.Fields{"success": resp.Success, "errors": resp.ErrorCodes})
		return errs.ReCaptchaValidationErr
	}

	logbuch.Debug("recaptcha OK", logbuch.Fields{"recaptcha": resp})
	return nil
}
