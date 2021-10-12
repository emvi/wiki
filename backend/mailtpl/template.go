package mailtpl

import (
	"emviwiki/shared/config"
	"emviwiki/shared/tpl"
)

const (
	InviteArticleMailTemplate                   = "mail_invite_article.html"
	RecommendMailTemplate                       = "mail_recommend_article.html"
	InviteUserMailTemplate                      = "mail_invite_user.html"
	MemberJoinedMailTemplate                    = "mail_member_joined.html"
	NewsletterConfirmationMailTemplate          = "mail_newsletter_confirmation.html"
	NewsletterOnPremiseConfirmationMailTemplate = "mail_newsletter_onpremise_confirmation.html"
	SupportMailTemplate                         = "mail_support.html"
	SubscriptionMailTemplate                    = "mail_subscription.html"
	SubscriptionCancelledMailTemplate           = "mail_subscription_cancelled.html"
	ResumeSubscriptionMailTemplate              = "mail_resume_subscription.html"
	CancelSubscriptionMailTemplate              = "mail_cancel_subscription.html"
	DowngradeMailTemplate                       = "mail_expert_downgrade.html"
	PaymentActionRequiredMailTemplate           = "mail_payment_action_required.html"
)

var (
	Cache *tpl.Cache
)

func InitTemplates() {
	Cache = tpl.NewCache(config.Get().Template.MailTemplateDir, config.Get().Template.HotReload)
}
