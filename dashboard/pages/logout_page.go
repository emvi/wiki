package pages

import (
	"emviwiki/shared/constants"
	"github.com/emvi/logbuch"
	"net/http"
)

func LogoutPageHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(constants.AuthHeader)

	if err != nil {
		logbuch.Warn("Cookie not found", logbuch.Fields{"err": err})
	} else {
		cookie.MaxAge = -1
		http.SetCookie(w, cookie)
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}
