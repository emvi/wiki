package pages

import (
	"emviwiki/shared/i18n"
	"emviwiki/shared/rest"
	"emviwiki/website/legal"
	"github.com/emvi/logbuch"
	"html/template"
	"net/http"
)

var landingPageI18n = i18n.Translation{
	"en": {
		"page_title":           "Emvi — Knowledge management platform for companies and teams.",
		"name@company.com":     "name@company.com",
		"your_email":           "Your email address",
		"get_started":          "Get started",
		"get_started_free":     "Get started for free",
		"switch_to_night_mode": "Switch to Night Mode",
		"switch_to_light_mode": "Switch to Light Mode",
		"hero_headline":        "Knowledge sorted.",
		"hero_text":            "Whether it's documentation, tutorials or your internal newsletter. Emvi lets anyone in your company find, retain and share information.",
		"status_headline":      "Goodbye, folders",
		"status_text":          "An organization's knowledge rarely fits inside a tree-like structure. Keeping documents in silos builds borders between departments, teams and employees.",
		"features_1_headline":  "Structured by content",
		"features_1_text":      "Emvi leaves old-timey folders and messy channels behind. Organize articles with tags and curate them to lists. Fine-grained rights management gives full control.",
		"features_2_headline":  "Intuitive writing",
		"features_2_text":      "Create beautiful articles with an easy-to-use editor powered by real-time collaboration, inline mentions and an automatic table of contents.",
		"features_3_headline":  "Mighty search",
		"features_3_text":      "Emvi is centered around its powerful search. Combinable filters and a complete article change history make sure nothing is ever lost.",
		"features_4_headline":  "Set up for teamwork",
		"features_4_text":      "Stay informed with notifications and always up-to-date documents. Emvi provides a toolset for every team, whether you work in the same office or around the globe.",
		"usecases_headline":    "One internal wiki, fewer emails",
		"usecases_headline_1":  "Knowledge Base",
		"usecases_text_1":      "Experience is a companies greatest asset, Emvi gives employees a place to write it down.",
		"usecases_headline_2":  "Documentation",
		"usecases_text_2":      "Keep your projects in check by making information accessible without getting lost in folder limbo.",
		"usecases_headline_3":  "Tutorials",
		"usecases_text_3":      "Improve teamwork by enabling users to help each other and themselves.",
		"usecases_headline_4":  "Note-taking",
		"usecases_text_4":      "Meetings, phone calls or a quick chat in the office — save your daily thoughts for the long term.",
	},
	"de": {
		"page_title":           "Emvi — Wissensmanagement-Plattform für Unternehmen und Teams.",
		"name@company.com":     "name@firma.de",
		"your_email":           "Ihre E-Mail-Adresse",
		"get_started":          "Jetzt starten",
		"get_started_free":     "Jetzt kostenlos starten",
		"switch_to_night_mode": "Zum Nachtmodus wechseln",
		"switch_to_light_mode": "Zum Tagmodus wechseln",
		"hero_headline":        "Wissen schaffen.",
		"hero_text":            "Ob Dokumentation, Anleitung oder ein interner Newsletter. Mit Emvi kann jeder in Ihrem Unternehmen Informationen finden, festhalten und teilen.",
		"status_headline":      "Adieu, Ordner",
		"status_text":          "Das Wissen innerhalb einer Organisation fügt sich nur selten in eine baumartige Struktur. Mit Dokumenten in Silos entstehen Mauern zwischen Abteilungen und Mitarbeitern.",
		"features_1_headline":  "Durch Inhalt strukturiert",
		"features_1_text":      "Emvi lässt Ordner hinter sich. Artikel werden mit Tags organisiert und zu Listen gebündelt. Mittels eines fein abgestuften Rechtemanagements behalten Sie die volle Kontrolle.",
		"features_2_headline":  "Intuitives Schreiben",
		"features_2_text":      "Schreiben Sie schöne Artikel einfach online mit einem benutzerfreundlichen Editor. Unterstützt durch Echtzeit-Kollaboration, Verlinkungen, einem generiertem Inhaltsverzeichnis und mehr.",
		"features_3_headline":  "Mächtige Suche",
		"features_3_text":      "Im Zentrum von Emvi steht eine leistungsstarke Suche. Kombinierbare Filter und eine vollständige Artikelhistorie sorgen dafür, dass niemals etwas verloren geht.",
		"features_4_headline":  "Mit Fokus auf Teamarbeit",
		"features_4_text":      "Bleiben Sie jederzeit auf dem Laufenden mit Hilfe von Benachrichtigungen und immer aktuellen Dokumenten. Emvi unterstützt jedes Team, egal von wo aus es arbeitet.",
		"usecases_headline":    "Ein internes Wiki, weniger E-Mails",
		"usecases_headline_1":  "Wissensdatenbank",
		"usecases_text_1":      "Wissen und Erfahrung ist das größte Kapital eines Unternehmens. Emvi gibt Mitarbeitern die Möglichkeit es aufzuschreiben.",
		"usecases_headline_2":  "Dokumentation",
		"usecases_text_2":      "Halten Sie Ihre Projekte in Schach, indem Informationen für jeden zugänglich gemacht werden. Ganz ohne Ordnerchaos.",
		"usecases_headline_3":  "Anleitungen",
		"usecases_text_3":      "Verbessern Sie die Teamarbeit, indem Kollegen ermöglicht wird, sich gegenseitig zu helfen.",
		"usecases_headline_4":  "Notizen",
		"usecases_text_4":      "Besprechungen, Telefonate oder das kurze Gespräch auf dem Flur - speichern Sie tägliche Gedanken langfristig.",
	},
}

func LandingPageHandler(w http.ResponseWriter, r *http.Request) {
	tpl := tplCache.Get()

	if tpl == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	langCode := rest.GetSupportedLangCode(r)
	data := struct {
		LangCode             string
		IsBlog               bool
		Vars                 map[string]template.HTML
		NavbarVars           map[string]template.HTML
		FooterVars           map[string]template.HTML
		NewsletterVars       map[string]template.HTML
		BackendHost          string
		AuthHost             string
		WebsiteHost          string
		AuthClientID         string
		Version              string
		GoogleSSOClientId    string
		GitHubSSOClientId    string
		SlackSSOClientId     string
		MicrosoftSSOClientId string
		CookiesNote          template.HTML
		IsIntegration        bool
	}{
		langCode,
		false,
		i18n.GetVars(langCode, landingPageI18n),
		i18n.GetVars(langCode, navbarComponentI18n),
		i18n.GetVars(langCode, footerComponentI18n),
		i18n.GetVars(langCode, newsletterComponentI18n),
		backendHost,
		authHost,
		websiteHost,
		clientId,
		version,
		googleSSOClientId,
		githubSSOClientId,
		slackSSOClientId,
		microsoftSSOClientId,
		template.HTML(legal.GetCookieNote(langCode)),
		isIntegration,
	}

	if err := tpl.ExecuteTemplate(w, landingPageTemplate, &data); err != nil {
		logbuch.Error("Error rendering landing page", logbuch.Fields{"err": err})
		w.WriteHeader(http.StatusInternalServerError)
	}
}
