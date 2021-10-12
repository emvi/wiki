package pages

import (
	"emviwiki/dashboard/auth"
	"emviwiki/dashboard/model"
	"emviwiki/shared/constants"
	"emviwiki/shared/util"
	"fmt"
	"net/http"
)

const (
	loginErrMessage = "User not found for that email and password."
	startPageURL    = "/"
)

type loginData struct {
	Error string `json:"error"`
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	data := new(loginData)

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			handleLoginError(w, data)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")
		user := model.GetUserByEmail(email)

		if user == nil {
			handleLoginError(w, data)
			return
		}

		passwordHash := util.Sha256Base64(password + user.PasswordSalt)

		if passwordHash != user.Password {
			handleLoginError(w, data)
			return
		}

		if err := setUserToken(w, user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, startPageURL, http.StatusFound)
		return
	}

	RenderPage(w, loginPageTemplate, nil, data)
}

func handleLoginError(w http.ResponseWriter, data *loginData) {
	w.WriteHeader(http.StatusBadRequest)
	data.Error = loginErrMessage
	RenderPage(w, loginPageTemplate, nil, &data)
}

func setUserToken(w http.ResponseWriter, user *model.User) error {
	claims := auth.UserTokenClaims{
		UserId:    user.ID,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
	}
	token, expires, err := auth.NewUserToken(&claims)

	if err != nil {
		return err
	}

	cookie := http.Cookie{Name: constants.AuthHeader,
		Value:    fmt.Sprintf("%s %s", constants.AuthTokenType, token),
		Expires:  expires,
		Secure:   SecureCookies,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode}
	http.SetCookie(w, &cookie)
	return nil
}
