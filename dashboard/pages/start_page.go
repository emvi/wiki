package pages

import (
	"emviwiki/dashboard/auth"
	"emviwiki/dashboard/model"
	"emviwiki/shared/config"
	"emviwiki/shared/rest"
	"encoding/json"
	"fmt"
	"github.com/emvi/logbuch"
	"html/template"
	"math"
	"net/http"
	"strings"
	"time"
)

const (
	registrationDefaultDays = 30
	loginDefaultDays        = 30
	maxStatisticDays        = 365
	statisticsDateFormat    = "2006-01-02"
)

func StartPageHandler(claims *auth.UserTokenClaims, w http.ResponseWriter, r *http.Request) {
	registrationDays, _ := rest.GetIntParam(r, "registration_days")
	loginDays, _ := rest.GetIntParam(r, "login_days")
	realtime := rest.GetBoolParam(r, "realtime")
	registrationDataOnly := registrationDays != 0
	loginDataOnly := loginDays != 0

	if registrationDays <= 0 || registrationDays > maxStatisticDays {
		registrationDays = registrationDefaultDays
	}

	if loginDays <= 0 || loginDays > maxStatisticDays {
		loginDays = loginDefaultDays
	}

	if registrationDataOnly {
		resp := struct {
			Statistics []model.Statistics `json:"statistics"`
		}{model.FindRegistrationStatisticAfter(time.Now().Add(-time.Hour * 24 * time.Duration(registrationDays)))}
		rest.WriteResponse(w, &resp)
	} else if loginDataOnly {
		resp := struct {
			Statistics []model.Statistics `json:"statistics"`
		}{model.FindLoginStatisticAfter(time.Now().Add(-time.Hour * 24 * time.Duration(loginDays)))}
		rest.WriteResponse(w, &resp)
	} else if realtime {
		roomsOpen, usersConnected := getRealtimeData()
		resp := struct {
			Rooms       int `json:"rooms"`
			Connections int `json:"connections"`
		}{roomsOpen, usersConnected}
		rest.WriteResponse(w, &resp)
	} else {
		roomsOpen, usersConnected := getRealtimeData()
		registrationLabel, registrationData := getStatisticsLabelAndData(model.FindRegistrationStatisticAfter(time.Now().Add(-time.Hour * 24 * time.Duration(registrationDays))))
		loginLabel, loginData := getStatisticsLabelAndData(model.FindLoginStatisticAfter(time.Now().Add(-time.Hour * 24 * time.Duration(loginDays))))
		orgaCount := model.CountOrganizations()
		userCount := model.CountUser()
		articleCount := model.CountArticles()
		listCount := model.CountLists()
		groupCount := model.CountGroups()
		tagCount := model.CountTags()
		data := struct {
			UsersConnectedCount         int
			RoomCount                   int
			OrganizationCount           int
			UserCount                   int
			UserCountPerOrganization    float64
			ArticleCount                int
			ArticleCountPerOrganization float64
			ListCount                   int
			ListCountPerOrganization    float64
			GroupCount                  int
			GroupCountPerOrganization   float64
			TagCount                    int
			TagCountPerOrganization     float64
			NewsletterCount             int
			NewsletterOnPremiseCount    int
			LoggedInToday               int
			LoggedInOneWeek             int
			LoggedInTwoWeeks            int
			LoggedInThreeWeeks          int
			LoggedInOneMonth            int
			LoggedInTwoMonths           int
			RegistrationStatisticsLabel template.JS
			RegistrationStatisticsData  template.JS
			LoginStatisticsLabel        template.JS
			LoginStatisticsData         template.JS
		}{
			usersConnected,
			roomsOpen,
			orgaCount,
			userCount,
			math.Round((float64(userCount)/float64(orgaCount))*100) / 100,
			articleCount,
			math.Round((float64(articleCount)/float64(orgaCount))*100) / 100,
			listCount,
			math.Round((float64(listCount)/float64(orgaCount))*100) / 100,
			groupCount - orgaCount*4, // excluding default groups
			math.Round((float64(groupCount-orgaCount*4)/float64(orgaCount))*100) / 100,
			tagCount,
			math.Round((float64(tagCount)/float64(orgaCount))*100) / 100,
			model.CountNewsletterWithoutList(),
			model.CountNewsletterWithListOnPremise(),
			model.CountUserLoggedInAfter(time.Now().Add(-time.Hour * 24)),
			model.CountUserLoggedInAfter(time.Now().Add(-time.Hour * 24 * 7)),
			model.CountUserLoggedInAfter(time.Now().Add(-time.Hour * 24 * 14)),
			model.CountUserLoggedInAfter(time.Now().Add(-time.Hour * 24 * 21)),
			model.CountUserLoggedInAfter(time.Now().Add(-time.Hour * 24 * 30)),
			model.CountUserLoggedInAfter(time.Now().Add(-time.Hour * 24 * 60)),
			registrationLabel,
			registrationData,
			loginLabel,
			loginData,
		}

		RenderPage(w, startPageTemplate, claims, &data)
	}
}

func getStatisticsLabelAndData(data []model.Statistics) (template.JS, template.JS) {
	if len(data) == 0 {
		return "", ""
	}

	var labels strings.Builder
	var datapoints strings.Builder

	for _, point := range data {
		labels.WriteString(fmt.Sprintf("'%s',", point.Date.Format(statisticsDateFormat)))
		datapoints.WriteString(fmt.Sprintf("%d,", point.Count))
	}

	labelsStr := labels.String()
	dataStr := datapoints.String()
	return template.JS(labelsStr[:len(labelsStr)-1]), template.JS(dataStr[:len(dataStr)-1])
}

func getRealtimeData() (int, int) {
	resp, err := http.Get(fmt.Sprintf("%s/api/v1/collab/stats", config.Get().Hosts.Collab))
	statusCode := 0

	if resp != nil {
		statusCode = resp.StatusCode
	}

	if err != nil || statusCode != http.StatusOK {
		logbuch.Error("Error getting real time data from collab", logbuch.Fields{"err": err, "status_code": statusCode})
		return 0, 0
	}

	req := struct {
		Rooms       int `json:"rooms"`
		Connections int `json:"connections"`
	}{}
	decoder := json.NewDecoder(resp.Body)

	if err := decoder.Decode(&req); err != nil {
		logbuch.Debug("Error decoding JSON real time data response from collab", logbuch.Fields{"err": err})
		return 0, 0
	}

	return req.Rooms, req.Connections
}
