package pages

import (
	"emviwiki/dashboard/auth"
	"github.com/emvi/logbuch"
	"net/http"
)

type PageData struct {
	Vars   interface{}
	Claims *auth.UserTokenClaims
}

func RenderPage(w http.ResponseWriter, name string, claims *auth.UserTokenClaims, data interface{}) {
	tpl := tplCache.Get()

	if tpl == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	pageData := PageData{
		Vars:   data,
		Claims: claims,
	}

	if err := tpl.ExecuteTemplate(w, name, &pageData); err != nil {
		logbuch.Error("Error rendering page", logbuch.Fields{"err": err, "page": "login"})
		w.WriteHeader(http.StatusInternalServerError)
	}
}
