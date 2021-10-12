package tpl

import (
	"emviwiki/shared/config"
	"emviwiki/shared/model"
	"fmt"
	"github.com/emvi/hide"
	"html/template"
	"strings"
	"time"
)

var (
	funcMap = template.FuncMap{
		"IdToString":             hide.ToString,
		"FormatDate":             func(t time.Time, fmt string) string { return t.Format(fmt) },
		"ToHTML":                 func(str string) template.HTML { return template.HTML(str) },
		"MailGreeting":           mailGreeting,
		"MailGoodbye":            mailGoodbye,
		"MailParagraph":          mailParagraph,
		"MailTextblock":          mailTextblock,
		"MailButton":             mailButton,
		"MailNotificationsStart": mailNotificationsStart,
		"MailNotificationsEnd":   mailNotificationsEnd,
		"MailNotification":       mailNotification,
	}
)

func mailGreeting(greeting interface{}) template.HTML {
	greetingStr := strings.TrimSpace(toString(greeting))

	if greetingStr == "" {
		greetingStr = "Hello from Emvi!"
	}

	return template.HTML(fmt.Sprintf(`<h1 style="margin: 0; font-size: 24px; font-weight: 500; line-height: 36px; margin-bottom: 8px;">%s</h1>`, greetingStr))
}

func mailGoodbye(goodbye interface{}) template.HTML {
	goodbyeStr := strings.TrimSpace(toString(goodbye))

	if goodbyeStr == "" {
		goodbyeStr = "Cheers, Emvi"
	}

	return template.HTML(fmt.Sprintf(`<br><p style="margin: 0;">%s</p>`, goodbyeStr))
}

func mailParagraph(content interface{}) template.HTML {
	return template.HTML(fmt.Sprintf(`<p style="margin: 0;">%s</p>`, toString(content)))
}

func mailParagraphCentered(content interface{}) template.HTML {
	return template.HTML(fmt.Sprintf(`<p style="margin: 0;text-align: center;">%s</p>`, toString(content)))
}

func mailTextblock(content ...template.HTML) template.HTML {
	str := make([]string, 0, len(content))

	for _, c := range content {
		str = append(str, string(c))
	}

	return template.HTML(fmt.Sprintf(`<tr>
		<td align="left" style="padding: 24px; font-family: 'Inter', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 28px;">
			%s
		</td>
	</tr>`, strings.Join(str, "")))
}

func mailButton(url, label, linkText interface{}, content ...template.HTML) template.HTML {
	link := fmt.Sprintf(`%s:<br /><a href="%s" target="_blank">%s</a>`, toString(linkText), toString(url), toString(url))
	textblock := make([]template.HTML, 0, len(content)+1)
	textblock = append(textblock, mailParagraphCentered(template.HTML(link)))

	for _, c := range content {
		textblock = append(textblock, c)
	}

	return template.HTML(fmt.Sprintf(`<tr>
		<td align="left">
			<table border="0" cellpadding="0" cellspacing="0" width="100%%">
				<tr>
					<td align="center" style="padding: 12px;">
						<table border="0" cellpadding="0" cellspacing="0">
							<tr>
								<td align="center" bgcolor="#198cff" style="border-radius: 8px;">
									<a href="%s" target="_blank" style="display: inline-block; padding: 16px 36px; font-family: 'Inter', Helvetica, Arial, sans-serif; font-size: 16px; font-weight: 500; color: #ffffff; text-decoration: none; border-radius: 6px;">%s</a>
								</td>
							</tr>
						</table>
					</td>
				</tr>
			</table>
		</td>
	</tr>%s`, toString(url), toString(label), mailTextblock(textblock...)))
}

func mailNotificationsStart() template.HTML {
	return `<tr>
		<td align="left" style="padding: 0 24px; font-family: 'Inter', Helvetica, Arial, sans-serif; font-size: 16px; line-height: 28px;">
			<table border="0" cellpadding="0" cellspacing="0" width="100%">`
}

func mailNotificationsEnd() template.HTML {
	return `</table>
		</td>
	</tr>`
}

func mailNotification(user *model.User, text interface{}, when time.Time) template.HTML {
	picture := "https://emvi.com/static/img/user-placeholder.png"

	if user.Picture.Valid {
		picture = fmt.Sprintf("%s/api/v1/content/%s", config.Get().Hosts.Backend, user.Picture.String)
	}

	return template.HTML(fmt.Sprintf(`<tr>
		<td align="left" valign="center" style="padding: 16px 0;">
			<a href="#user">
				<img src="%s" style="display: block; height: 40px; width: 40px; border-radius: 50%%;">
			</a>
		</td>
		<td align="left" valign="center" style="padding: 16px 12px;">
			<p style="margin: 0;"><strong>%s %s</strong> <span>%s</span></p>
		</td>
		<td align="left" valign="center" style="padding: 16px 0;">
			<p style="margin: 0; font-size: 14px; color: #555759; white-space: nowrap;">%s</p>
		</td>
	</tr>`, picture, user.Firstname, user.Lastname, toString(text), when.Format("Mon 2. Jan 15:04")))
}

func toString(in interface{}) string {
	if str, ok := in.(string); ok {
		return str
	} else if html, ok := in.(template.HTML); ok {
		return string(html)
	}

	return ""
}
