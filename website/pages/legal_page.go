package pages

import (
	"emviwiki/shared/i18n"
	"emviwiki/shared/rest"
	"emviwiki/website/legal"
	"github.com/emvi/logbuch"
	"html/template"
	"net/http"
)

var legalPageI18n = i18n.Translation{
	"en": {
		"page_title":     "Emvi — Legal Notice",
		"title":          "Legal Notice",
		"subtitle":       "According to § 5 TMG",
		"germany":        "Germany",
		"director":       "Managing Directors and responsible for contents",
		"legal":          "Legal",
		"court":          "Local Court of Gütersloh",
		"vatid":          "VAT-ID DE326639794",
		"disclaimer":     "Disclaimer",
		"disclaimer_1":   "Limitation of liability for internal content",
		"disclaimer_1_1": "The content of our website has been compiled with meticulous care and to the best of our knowledge. However, we cannot assume any liability for the up-to-dateness, completeness or accuracy of any of the pages.",
		"disclaimer_1_2": "Pursuant to section 7, para. 1 of the TMG (Telemediengesetz –  Tele Media Act by German law), we as service providers are liable for our own content on these pages in accordance with general laws. However, pursuant to sections 8 to 10 of the TMG, we as service providers are not under obligation to monitor external information provided or stored on our website.",
		"disclaimer_1_3": "Once we have become aware of a specific infringement of the law, we will immediately remove the content in question. Any liability concerning this matter can only be assumed from the point in time at which the infringement becomes known to us.",
		"disclaimer_2":   "Limitation of liability for external links",
		"disclaimer_2_1": "Our website contains links to the websites of third parties („external links“). As the content of these websites is not under our control, we cannot assume any liability for such external content. In all cases, the provider of information of the linked websites is liable for the content and accuracy of the information provided",
		"disclaimer_2_2": "At the point in time when the links were placed, no infringements of the law were recognisable to us. As soon as an infringement of the law becomes known to us, we will immediately remove the link in question.",
		"disclaimer_3":   "Copyright",
		"disclaimer_3_1": "The content and works published on this website are governed by the copyright laws of Germany. Any duplication, processing, distribution or any form of utilisation beyond the scope of copyright law shall require the prior written consent of the author or authors in question.",
		"disclaimer_4":   "Data protection",
		"disclaimer_4_1": "A visit to our website can result in the storage on our server of information about the access (date, time, page accessed). This does not represent any analysis of personal data (e.g., name, address or e-mail address).",
		"disclaimer_4_2": "If personal data are collected, this only occurs – to the extent possible – with the prior consent of the user of the website. Any forwarding of the data to third parties without the express consent of the user shall not take place.",
		"disclaimer_4_3": "We would like to expressly point out that the transmission of data via the Internet (e.g., by e-mail) can offer security vulnerabilities. It is therefore impossible to safeguard the data completely against access by third parties.",
		"disclaimer_4_4": "It is therefore impossible to safeguard the data completely against access by third parties. We cannot assume any liability for damages arising as a result of such security vulnerabilities.",
		"disclaimer_4_5": "The use by third parties of all published contact details for the purpose of advertising is expressly excluded. We reserve the right to take legal steps in the case of the unsolicited sending of advertising information; e.g., by means of spam mail.",
	},
	"de": {
		"page_title":     "Emvi — Impressum",
		"title":          "Impressum",
		"subtitle":       "Angaben gemäß § 5 TMG",
		"germany":        "Deutschland",
		"director":       "Geschäftsführende Gesellschafter und inhaltlich Verantwortliche",
		"legal":          "Unternehmen",
		"court":          "Amtsgericht Gütersloh",
		"vatid":          "USt-IdNr. DE326639794",
		"disclaimer":     "Disclaimer – rechtliche Hinweise",
		"disclaimer_1":   "1. Haftungsbeschränkung",
		"disclaimer_1_1": "Die Webseite enthält sog. „externe Links“ (Verlinkungen) zu anderen Webseiten, auf deren Inhalt der Anbieter der Webseite keinen Einfluss hat. Aus diesem Grund kann der Anbieter für diese Inhalte auch keine Gewähr übernehmen.",
		"disclaimer_1_2": "Für die Inhalte und Richtigkeit der bereitgestellten Informationen ist der jeweilige Anbieter der verlinkten Webseite verantwortlich. Zum Zeitpunkt der Verlinkung waren keine Rechtsverstöße erkennbar. Bei Bekanntwerden einer solchen Rechtsverletzung wird der Link umgehend entfernen.",
		"disclaimer_1_3": "Als Diensteanbieter ist der Anbieter dieser Webseite gemäß § 7 Abs. 1 TMG für eigene Inhalte und bereitgestellte Informationen auf diesen Seiten nach den allgemeinen Gesetzen verantwortlich; nach den §§ 8 bis 10 TMG jedoch nicht verpflichtet, die übermittelten oder gespeicherten fremden Informationen zu überwachen. Eine Entfernung oder Sperrung dieser Inhalte erfolgt umgehend ab dem Zeitpunkt der Kenntnis einer konkreten Rechtsverletzung. Eine Haftung ist erst ab dem Zeitpunkt der Kenntniserlangung möglich.",
		"disclaimer_2":   "2. Externe Links",
		"disclaimer_2_1": "Die Webseite enthält sog. „externe Links“ (Verlinkungen) zu anderen Webseiten, auf deren Inhalt der Anbieter der Webseite keinen Einfluss hat. Aus diesem Grund kann der Anbieter für diese Inhalte auch keine Gewähr übernehmen.",
		"disclaimer_2_2": "Für die Inhalte und Richtigkeit der bereitgestellten Informationen ist der jeweilige Anbieter der verlinkten Webseite verantwortlich. Zum Zeitpunkt der Verlinkung waren keine Rechtsverstöße erkennbar. Bei Bekanntwerden einer solchen Rechtsverletzung wird der Link umgehend entfernen.",
		"disclaimer_3":   "3. Urheberrecht/Leistungsschutzrecht",
		"disclaimer_3_1": "Die auf dieser Webseite veröffentlichten Inhalte, Werke und bereitgestellten Informationen unterliegen dem deutschen Urheberrecht und Leistungsschutzrecht. Jede Art der Vervielfältigung, Bearbeitung, Verbreitung, Einspeicherung und jede Art der Verwertung außerhalb der Grenzen des Urheberrechts bedarf der vorherigen schriftlichen Zustimmung des jeweiligen Rechteinhabers. Das unerlaubte Kopieren/Speichern der bereitgestellten Informationen auf diesen Webseiten ist nicht gestattet und strafbar.",
		"disclaimer_4":   "4. Datenschutz",
		"disclaimer_4_1": "Durch den Besuch des Internetauftritts können Informationen (Datum, Uhrzeit, aufgerufene Seite) über den Zugriff auf dem Server gespeichert werden. Es werden keine personenbezogenenen (z. B. Name, Anschrift oder E-Mail-Adresse) Daten, gespeichert.",
		"disclaimer_4_2": "Sofern personenbezogene Daten erhoben werden, erfolgt dies, sofern möglich, nur mit dem vorherigen Einverständnis des Nutzers der Webseite. Eine Weitergabe der Daten an Dritte findet ohne ausdrückliche Zustimmung des Nutzers nicht statt.",
		"disclaimer_4_3": "Der Anbieter weist darauf hin, dass die Übertragung von Daten im Internet (z. B. per E-Mail) Sicherheitslücken aufweisen und ein lückenloser Schutz der Daten vor dem Zugriff Dritter nicht gewährleistet werden kann. Der Anbieter übernimmt keine Haftung für die durch solche Sicherheitslücken entstandenen Schäden.",
		"disclaimer_4_4": "Der Verwendung der Kontaktdaten durch Dritte zur gewerblichen Nutzung wird ausdrücklich widersprochen. Es sei denn, der Anbieter hat zuvor seine schriftliche Einwilligung erteilt.",
		"disclaimer_4_5": "Der Anbieter behält sich rechtliche Schritte für den Fall der unverlangten Zusendung von Werbeinformationen, z. B. durch Spam-Mails, vor.",
	},
}

func LegalPageHandler(w http.ResponseWriter, r *http.Request) {
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
		i18n.GetVars(langCode, legalPageI18n),
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

	if err := tpl.ExecuteTemplate(w, legalPageTemplate, &data); err != nil {
		logbuch.Error("Error rendering legal page", logbuch.Fields{"err": err})
		w.WriteHeader(http.StatusInternalServerError)
	}
}
