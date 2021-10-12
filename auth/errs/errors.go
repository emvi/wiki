package errs

import (
	"emviwiki/shared/rest"
	"errors"
)

var (
	TxBegin  = errors.New("Error starting transaction")
	TxCommit = errors.New("Error committing transaction")
	Saving   = errors.New("Error on save")

	UserNotFound             = rest.NewApiError("User not found", "")
	TokenInvalid             = rest.NewApiError("Token invalid", "")
	EmailEmpty               = rest.NewApiError("Email empty", "email")
	EmailTooLong             = rest.NewApiError("Email too long", "email")
	EmailInvalid             = rest.NewApiError("Email invalid", "email")
	PasswordMatch            = rest.NewApiError("Passwords do not match", "password")
	PasswordInvalid          = rest.NewApiError("Password invalid", "password")
	TermsOfService           = rest.NewApiError("Must accept terms of service", "terms_of_service")
	Privacy                  = rest.NewApiError("Must accept privacy", "privacy")
	EmailInUse               = rest.NewApiError("Email in use", "email")
	EmailNotChanged          = rest.NewApiError("Email not changed", "email")
	Mail                     = rest.NewApiError("Mail could not be send", "")
	PasswordWrong            = rest.NewApiError("Password wrong", "old_password")
	FirstnameInvalid         = rest.NewApiError("Firstname invalid", "firstname")
	LastnameInvalid          = rest.NewApiError("Lastname invalid", "lastname")
	LanguageInvalid          = rest.NewApiError("Language invalid", "language")
	GrantTypeInvalid         = rest.NewApiError("Grant type invalid", "")
	ClientCredentialsInvalid = rest.NewApiError("Client credentials invalid", "")
	ClientNameInvalid        = rest.NewApiError("Client name invalid", "name")
	ClientNameInUse          = rest.NewApiError("Client name in use", "name")
	ReCaptchaValidationErr   = rest.NewApiError("reCAPTCHA invalid", "")
	EmailForbidden           = rest.NewApiError("Email forbidden", "email")
	MaxRegistrationMailsSend = rest.NewApiError("Maximum registration mails send", "")
	UserRegisteredAlready    = rest.NewApiError("User registered already", "")
	ScopeInvalid             = rest.NewApiError("Scope invalid", "")
	ClientNotFound           = rest.NewApiError("Client not found", "")
	StepInvalid              = rest.NewApiError("Step invalid", "")
	OperationNotAllowed      = rest.NewApiError("Operation not allowed", "")
)
