package recaptcha

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	recaptchaVerificationURL = "https://www.google.com/recaptcha/api/siteverify"
)

type RecaptchaRequest struct {
	Secret   string `json:"secret"`
	Response string `json:"response"`
	RemoteIP string `json:"remoteip"` // optional
}

type RecaptchaResponse struct {
	Success     bool     `json:"success"`
	ChallengeTs string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

// Recaptcha is an interface to Google's Recaptcha v2.
type Recaptcha interface {
	// Validate validates the given token and returns the RecaptchaResponse object.
	Validate(string) (*RecaptchaResponse, error)
}

type RecaptchaValidator struct{}

func NewRecaptchaValidator() Recaptcha {
	return new(RecaptchaValidator)
}

func (validator *RecaptchaValidator) Validate(token string) (*RecaptchaResponse, error) {
	values := url.Values{}
	values.Add("secret", recaptchaSecret)
	values.Add("response", token)
	resp, err := http.PostForm(recaptchaVerificationURL, values)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Status not OK")
	}

	tokenResp := new(RecaptchaResponse)
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, tokenResp); err != nil {
		return nil, err
	}

	return tokenResp, nil
}
