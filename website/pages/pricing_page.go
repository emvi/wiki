package pages

import (
	"emviwiki/shared/i18n"
	"emviwiki/shared/rest"
	"emviwiki/website/legal"
	"github.com/emvi/logbuch"
	"html/template"
	"net/http"
)

var pricingPageI18n = i18n.Translation{
	"en": {
		"page_title":               "Emvi — Pricing",
		"pricing_headline":         "Simple, fair pricing.",
		"pricing_text":             "Pay monthly or annually and cancel at any time. There will always be a free plan to try it for as long as you like.",
		"monthly":                  "Monthly",
		"annually":                 "Annually",
		"entry_headline":           "Entry",
		"entry_text":               "For small teams wanting to start their own internal wiki.",
		"per_user_month":           "per user / month",
		"vat":                      "depending on your country of origin, VAT may be charged",
		"unlimited_members_tags":   "Unlimited members and tags",
		"articles_list_limit":      "Up to 100 articles and 10 lists",
		"two_user_collab":          "Collaborative editing with two users",
		"article_version_limit":    "Last three article versions",
		"5gb_file_storage":         "5 GB file storage",
		"try":                      "Try it for free",
		"expert_headline":          "Expert",
		"expert_text":              "For companies and teams building a living knowledge base.",
		"unlimited_articles_lists": "Unlimited articles and lists",
		"user_groups":              "User groups for permission control",
		"article_history":          "Complete article history",
		"article_translations":     "Article translations",
		"admin_roles":              "Multiple administrators and moderator role",
		"read_only":                "Read-only members",
		"api":                      "Client API",
		"10gb_file_storage":        "10 GB file storage per user",
		"get_started":              "Get started",
		"enterprise_headline":      "Enterprise",
		"enterprise_text":          "An on-premises solution will arrive at a later date.",
		"coming_soon":              "Coming soon",
		"everything_in_expert":     "Everything in Expert",
		"self_hosted":              "Self-hosted",
		"additional_features":      "Additional enterprise features",
		"name@company.com":         "name@company.com",
		"your_email":               "Your email address",
		"get_notified":             "Get notified",
		"newsletter_hint":          "You will receive updates on our on-premises solution only.",
		"newsletter_success":       "Thank you for subscribing to our newsletter. Please check your inbox to verify your subscription.",
		"newsletter_error":         "Email address invalid.",
		"faq_headline":             "Frequently asked questions",
		"question_1":               "Can I export my data?",
		"answer_1":                 "You can export an article to HTML or Markdown at any point.",
		"question_2":               "What does \"Fair Pricing\" mean?",
		"answer_2":                 "Organizations on the \"Expert\" plan are billed per active user, meaning those who log in at least once in a 30-day period. You receive a credit for inactive users, which will be applied to new or existing users in the next billing cycle.",
		"question_3":               "What personal information do you ask for?",
		"answer_3":                 "During registration we ask for your email address, full name and an account password. No additional information is required to sign up.",
		"question_4":               "Can I upgrade or downgrade my plan?",
		"answer_4":                 "You can upgrade to \"Expert\" at anytime. If you cancel your plan, you are downgraded at the end of your billing cycle.",
		"question_5":               "Where is my data stored?",
		"answer_5":                 "Your account and organization related data is currently stored in Nürnberg, Falkenstein, Frankfurt (Germany) and Helsinki (Finland) under European privacy law. For \"Expert\" users payment information are processed by Stripe, Inc. in the United States.",
		"question_6":               "What are the payment options?",
		"answer_6":                 "We currently accept all major credit cards (e.g. Visa, MasterCard, American Express). Depending on your country you may also use services like AliPay, Giropay, iDEAL, SEPA or SOFORT.",
		"question_7":               "What happens when I reach the storage limit?",
		"answer_7":                 "In case you reach the storage limit you can either free up space by deleting files or upgrade your organization to \"Expert\" to increase the storage limit.",
		"question_8":               "What happens to my data when I downgrade?",
		"answer_8":                 "If you downgrade an organization from \"Expert\" to \"Entry\" level you will keep all data stored in Emvi. You won't be able to use Expert level features anymore such as viewing the full article history, managing groups and an extended storage limit. If you setup groups for easier permission management, group members will keep access to the articles the groups are used in, but you won't be able to edit or assign the group anymore.",
		"more_questions":           "More questions?",
		"more_questions_1":         "Write to us at",
		"more_questions_2":         "or reach out on",
		"more_questions_3":         "We will get back to you as soon as possible.",
		"or":                       "or",
	},
	"de": {
		"page_title":               "Emvi — Preise",
		"pricing_headline":         "Einfache, faire Preisgestaltung.",
		"pricing_text":             "Bezahle pro Monat oder Jahr und kündige jederzeit. Es wird immer eine kostenlose Variante geben, um Emvi beliebig lange auszuprobieren.",
		"monthly":                  "Monatlich",
		"annually":                 "Jährlich",
		"entry_headline":           "Entry",
		"entry_text":               "Für kleine Teams die mit ihrem eigenen Wiki starten möchten.",
		"per_user_month":           "pro Nutzer / Monat",
		"vat":                      "zzgl. 19% MwSt. innerhalb der EU",
		"unlimited_members_tags":   "Unendlich Mitglieder und Tags",
		"articles_list_limit":      "Bis zu 100 Artikel und 10 Listen",
		"two_user_collab":          "Kollaboratives Schreiben mit zwei Nutzern",
		"article_version_limit":    "Die letzten drei Artikelversionen",
		"5gb_file_storage":         "5 GB Speicherplatz",
		"try":                      "Kostenlos testen",
		"expert_headline":          "Expert",
		"expert_text":              "Für Teams und Unternehmen die eine lebendinge Wissensdatenbank aufbauen.",
		"unlimited_articles_lists": "Unendlich Artikel und Listen",
		"user_groups":              "Nutzergruppen für Zugriffskontrolle",
		"article_history":          "Vollständiger Artikeländerungsverlauf",
		"article_translations":     "Artikelübersetzungen",
		"admin_roles":              "Mehrere Administratoren / Moderatoren",
		"read_only":                "Nutzer mit ausschließlich Lesezugriff",
		"api":                      "Client API",
		"10gb_file_storage":        "10 GB Speicherplatz pro Nutzer",
		"get_started":              "Jetzt starten",
		"enterprise_headline":      "Enterprise",
		"enterprise_text":          "Eine Vor-Ort-Lösung wird zu einem späteren Zeitpunkt verfügbar sein.",
		"coming_soon":              "Demnächst",
		"everything_in_expert":     "Alles in Expert",
		"self_hosted":              "Selbst gehostet",
		"additional_features":      "Zusätzliche Enterprise-Funktionen",
		"name@company.com":         "name@firma.com",
		"your_email":               "Ihre E-Mail-Adresse",
		"get_notified":             "Benachrichtigen",
		"newsletter_hint":          "Sie erhalten nur Updates zu unserer lokalen Lösung.",
		"newsletter_success":       "Vielen Dank, dass Sie sich für unseren Newsletter angemeldet hast. Bitte überprüfen Sie Ihren Posteingang, um das Abonnement zu bestätigen.",
		"newsletter_error":         "Die E-Mail-Adresse ist ungültig.",
		"faq_headline":             "Häufig gestellte Fragen",
		"question_1":               "Kann ich meine Daten exportieren?",
		"answer_1":                 "Ein Artikel kann jederzeit im Format HTML oder Markdown exportiert werden.",
		"question_2":               "Was bedeutet \"Faire Preisgestaltung\"?",
		"answer_2":                 "Für \"Expert\" Organisationen werden nur aktive Nutzer abgerechnet. Als aktiv zählen nur Nutzer, die sich mindestens einmal innerhalb einer 30-Tage Periode angemeldet haben. Du erhälst ein Guthaben für jeden inaktiven Nutzer, welches bei der nächsten Abrechnung für bestehende oder neue Nutzer berücksichtigt wird.",
		"question_3":               "Nach welchen personenbezogenen Daten fragen wir?",
		"answer_3":                 "Für die Registrierung benötigen wir Ihre E-Mail-Adresse, Vor- und Nachnamen sowie ein Konto-Passwort. Keine weitere Daten sind für die Registrierung erforderlich.",
		"question_4":               "Kann ich meinen Plan upgraden oder herabstufen?",
		"answer_4":                 "Sie können jederzeit auf \"Expert\" upgraden. Wenn Sie Ihr Abo kündigen, wird die Organisation mit Ende des Abrechnungszeitraums herabgestuft.",
		"question_5":               "Wo werden meine Daten gespeichert?",
		"answer_5":                 "Ihre personenbezogenen Daten und Organisationen werden in Nürnberg, Falkenstein, Frankfurt und Helsinki (Finnland) unter europäischem Datenschutzrecht gespeichert. Die Zahlungsinformationen für \"Expert\" Nutzer werden von Stripe Inc., USA gespeichert und verarbeitet.",
		"question_6":               "Welche Zahlungsmöglichkeiten gibt es?",
		"answer_6":                 "Aktuell akzeptieren wir alle bekannten Kreditkarten (z.B. Visa, MasterCard, American Express). Abhängig von Ihrem Standort können auch Dienste wie AliPay, Giropay, iDEAL, SEPA oder SOFORT genutzt werden.",
		"question_7":               "Was passiert wenn ich das Speicherlimit erreiche?",
		"answer_7":                 "Wenn das Speicherlimit erreicht wird, kann Speicher zurückgewonnen werden, indem alte Dateien gelöschst werden oder die Organisation upgegradet wird, um das Limit zu erhöhen.",
		"question_8":               "Was passiert mit meinen Daten wenn ich herabstufe?",
		"answer_8":                 "Wenn die Organisation von \"Expert\" auf \"Entry\" herabgestuft wird, bleiben alle Daten erhalten. \"Expert\" Funktionen, wie die Einsicht des kompletten Artikelverlaufs, Gruppenverwaltung und mehr Speicherplatz, stehen dann nicht mehr zur Verfügung. Sofern Gruppen angelegt wurden um Berechtigungen einfacher verwalten zu können, werden Gruppenmitglieder weiterhin Zugriff auf Artikel haben in denen die Gruppe verwendet wird. Allerdings können die Gruppen nicht mehr bearbeitet oder neu zugewiesen werden.",
		"more_questions":           "Weitere Fragen?",
		"more_questions_1":         "Schreiben Sie uns eine E-Mail an",
		"more_questions_2":         "oder kontaktieren Sie uns auf",
		"more_questions_3":         "Wir melden uns so bald wie möglich zurück.",
		"or":                       "oder",
	},
}

func PricingPageHandler(w http.ResponseWriter, r *http.Request) {
	tpl := tplCache.Get()

	if tpl == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	langCode := rest.GetSupportedLangCode(r)
	data := struct {
		LangCode      string
		IsBlog        bool
		Vars          map[string]template.HTML
		NavbarVars    map[string]template.HTML
		FooterVars    map[string]template.HTML
		BackendHost   string
		AuthHost      string
		WebsiteHost   string
		AuthClientID  string
		Version       string
		CookiesNote   template.HTML
		IsIntegration bool
	}{
		langCode,
		false,
		i18n.GetVars(langCode, pricingPageI18n),
		i18n.GetVars(langCode, navbarComponentI18n),
		i18n.GetVars(langCode, footerComponentI18n),
		backendHost,
		authHost,
		websiteHost,
		clientId,
		version,
		template.HTML(legal.GetCookieNote(langCode)),
		isIntegration,
	}

	if err := tpl.ExecuteTemplate(w, pricingPageTemplate, &data); err != nil {
		logbuch.Error("Error rendering pricing page", logbuch.Fields{"err": err})
		w.WriteHeader(http.StatusInternalServerError)
	}
}
