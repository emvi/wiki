package main

import (
	"emviwiki/backend/api"
	"emviwiki/backend/article"
	"emviwiki/backend/billing"
	"emviwiki/backend/content"
	"emviwiki/backend/mailtpl"
	"emviwiki/backend/member"
	"emviwiki/backend/newsletter"
	"emviwiki/backend/organization"
	"emviwiki/backend/support"
	"emviwiki/shared/auth"
	"emviwiki/shared/config"
	"emviwiki/shared/db"
	"emviwiki/shared/feed"
	"emviwiki/shared/i18n"
	"emviwiki/shared/mail"
	"emviwiki/shared/model"
	"emviwiki/shared/rest"
	"emviwiki/shared/server"
	"github.com/gorilla/mux"
	"net/http"
)

func setupRouter() *mux.Router {
	router := mux.NewRouter()

	// endpoints without context (organization) check -> website
	router.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {}).Methods(http.MethodGet)
	router.Handle("/api/v1/organization", api.AuthMiddleware(api.CreateOrganizationHandler, false, false, false)).Methods(http.MethodPost)
	router.Handle("/api/v1/organization/validation", api.AuthMiddleware(api.ValidateOrganizationHandler, false, false, false)).Methods(http.MethodPost)
	router.Handle("/api/v1/organizations", api.AuthMiddleware(api.GetOrganizationsHandler, false, false, false)).Methods(http.MethodGet)
	router.Handle("/api/v1/accession", api.AuthMiddleware(api.JoinOrganizationHandler, false, false, false)).Methods(http.MethodPut)
	router.Handle("/api/v1/invitation", api.AuthMiddleware(api.ReadInvitationsHandler, false, false, false)).Methods(http.MethodGet)
	router.Handle("/api/v1/invitation", api.AuthMiddleware(api.DeleteInvitationHandler, false, false, false)).Methods(http.MethodDelete)
	router.Handle("/api/v1/invitation/organization", api.AuthMiddleware(api.GetInvitationOrganizationHandler, false, false, false)).Methods(http.MethodGet)
	router.Handle("/api/v1/invitation/{code}", api.AuthMiddleware(api.GetInvitationHandler, false, false, false)).Methods(http.MethodGet)
	router.Handle("/api/v1/user", api.AuthMiddleware(api.GetUserHandler, false, false, false)).Methods(http.MethodGet)
	router.Handle("/api/v1/user/picture", api.AuthMiddleware(api.UploadUserPictureHandler, false, false, false)).Methods(http.MethodPost)
	router.Handle("/api/v1/user/picture", api.AuthMiddleware(api.DeleteUserPictureHandler, false, false, false)).Methods(http.MethodDelete)
	router.Handle("/api/v1/user/colormode", api.AuthMiddleware(api.UpdateUserColorModeHandler, false, false, false)).Methods(http.MethodPut)
	router.Handle("/api/v1/user/introduction", api.AuthMiddleware(api.UpdateUserIntroductionHandler, false, false, false)).Methods(http.MethodPut)
	router.Handle("/api/v1/newsletter", rest.ErrorMiddleware(api.SubscribeNewsletterHandler)).Methods(http.MethodPost)
	router.Handle("/api/v1/newsletter/onpremise", rest.ErrorMiddleware(api.SubscribeOnPremiseNewsletterHandler)).Methods(http.MethodPost)
	router.Handle("/api/v1/newsletter", rest.ErrorMiddleware(api.ConfirmNewsletterHandler)).Methods(http.MethodPut)
	router.Handle("/api/v1/newsletter", rest.ErrorMiddleware(api.UnsubscribeNewsletterHandler)).Methods(http.MethodDelete)
	router.HandleFunc("/api/v1/content/{filename}", api.GetContentHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/billing/webhook", api.StripeWebhookHandler).Methods(http.MethodPost)

	// endpoints with context (organization) check -> wiki
	addRoute(router, "/api/v1/organization", http.MethodGet, api.GetOrganizationHandler, false, false, "organization:r")
	addRoute(router, "/api/v1/organization", http.MethodPut, api.UpdateOrganizationHandler, false, true)
	addRoute(router, "/api/v1/organization", http.MethodDelete, api.DeleteOrganizationHandler, false, true)
	addRoute(router, "/api/v1/organization/permissions", http.MethodPut, api.UpdateOrganizationPermissionsHandler, true, true)
	addRoute(router, "/api/v1/organization/picture", http.MethodPost, api.UploadOrganizationPictureHandler, false, true)
	addRoute(router, "/api/v1/organization/picture", http.MethodDelete, api.DeleteOrganizationPictureHandler, false, true)
	addRoute(router, "/api/v1/organization/exit", http.MethodPost, api.LeaveOrganizationHandler, false, false)
	addRoute(router, "/api/v1/organization/statistics", http.MethodGet, api.GetOrganizationStatisticsHandler, false, false)
	addRoute(router, "/api/v1/organization/invitation", http.MethodPost, api.GenerateInvitationCodeHandler, false, false)
	addRoute(router, "/api/v1/organization/invitation", http.MethodGet, api.GetInvitationCodeHandler, false, false)
	addRoute(router, "/api/v1/organization/subscription", http.MethodPost, api.CreateSubscriptionHandler, false, false)
	addRoute(router, "/api/v1/organization/subscription", http.MethodPut, api.ResumeSubscriptionHandler, false, false)
	addRoute(router, "/api/v1/organization/subscription", http.MethodDelete, api.CancelSubscriptionHandler, false, false)
	addRoute(router, "/api/v1/organization/subscription", http.MethodGet, api.GetSubscriptionHandler, false, false)
	addRoute(router, "/api/v1/organization/subscription/invoice", http.MethodGet, api.GetInvoicesHandler, false, false)
	addRoute(router, "/api/v1/organization/subscription/customer", http.MethodPost, api.UpdateCustomerHandler, false, false)
	addRoute(router, "/api/v1/organization/subscription/plan", http.MethodPut, api.UpdatePlanHandler, false, false)
	addRoute(router, "/api/v1/organization/subscription/payment", http.MethodPost, api.UpdatePaymentMethodHandler, false, false)
	addRoute(router, "/api/v1/organization/subscription/payment", http.MethodDelete, api.RemovePaymentIntentClientSecretHandler, false, false)
	addRoute(router, "/api/v1/member", http.MethodGet, api.ReadOrganizationInvitationsHandler, false, true)
	addRoute(router, "/api/v1/member", http.MethodPost, api.InviteMemberHandler, false, true)
	addRoute(router, "/api/v1/member", http.MethodDelete, api.CancelInvitationHandler, false, true)
	addRoute(router, "/api/v1/member/{id}", http.MethodDelete, api.RemoveMemberHandler, false, true)
	addRoute(router, "/api/v1/member/{id}/moderator", http.MethodPut, api.ToggleMemberModeratorHandler, true, true)
	addRoute(router, "/api/v1/member/{id}/admin", http.MethodPut, api.ToggleMemberAdminHandler, true, true)
	addRoute(router, "/api/v1/member/{id}/readonly", http.MethodPut, api.ToggleReadOnlyHandler, false, true)
	addRoute(router, "/api/v1/auth", http.MethodGet, api.AuthenticateUserHandler, false, false)
	addRoute(router, "/api/v1/article", http.MethodPost, api.SaveArticleHandler, false, true)
	addRoute(router, "/api/v1/article/content", http.MethodPost, api.UploadArticleAttachmentHandler, false, true)
	addRoute(router, "/api/v1/article/private", http.MethodGet, api.ReadPrivateArticlesHandler, false, false)
	addRoute(router, "/api/v1/article/draft", http.MethodGet, api.ReadDraftsHandler, false, false)
	addRoute(router, "/api/v1/article/invite", http.MethodPut, api.InviteEditArticleHandler, false, true)
	addRoute(router, "/api/v1/article/history", http.MethodDelete, api.DeleteArticleHistoryEntryHandler, false, true)
	addRoute(router, "/api/v1/article/{id}", http.MethodGet, api.ReadArticleHandler, false, false, "articles:r")
	addRoute(router, "/api/v1/article/{id}", http.MethodDelete, api.DeleteArticleHandler, false, true)
	addRoute(router, "/api/v1/article/{id}/preview", http.MethodGet, api.ReadArticlePreviewHandler, false, false, "articles:r")
	addRoute(router, "/api/v1/article/{id}/history", http.MethodGet, api.ReadArticleHistoryHandler, false, false, "articles:r", "article_history:r")
	addRoute(router, "/api/v1/article/{id}/recommendation", http.MethodPost, api.RecommendArticleHandler, false, false)
	addRoute(router, "/api/v1/article/{id}/recommendation", http.MethodPut, api.ConfirmRecommendationHandler, false, false)
	addRoute(router, "/api/v1/article/{id}/invite", http.MethodPut, api.InviteEditArticleHandler, false, true)
	addRoute(router, "/api/v1/article/{id}/archive", http.MethodPut, api.ArchiveArticleHandler, false, true)
	addRoute(router, "/api/v1/article/{id}/reset", http.MethodPut, api.ResetArticleHandler, false, true)
	addRoute(router, "/api/v1/article/{id}/copy", http.MethodPut, api.CopyArticleHandler, false, true)
	addRoute(router, "/api/v1/article/{id}/list", http.MethodPost, api.AddArticleToListHandler, false, true)
	addRoute(router, "/api/v1/article/{id}/export", http.MethodGet, api.ExportArticleHandler, false, false)
	addRoute(router, "/api/v1/lang", http.MethodGet, api.GetLangsHandler, false, false, "language:r")
	addRoute(router, "/api/v1/lang", http.MethodPost, api.AddLangHandler, true, true)
	addRoute(router, "/api/v1/lang", http.MethodPut, api.SwitchDefaultLangHandler, false, true)
	addRoute(router, "/api/v1/lang/{id}", http.MethodGet, api.GetLangHandler, false, false, "language:r")
	addRoute(router, "/api/v1/usergroup", http.MethodPost, api.SaveUserGroupHandler, true, true)
	addRoute(router, "/api/v1/usergroup/{id}", http.MethodDelete, api.DeleteUserGroupHandler, true, true)
	addRoute(router, "/api/v1/usergroup/{id}", http.MethodGet, api.GetUserGroupHandler, false, false)
	addRoute(router, "/api/v1/usergroup/{id}/member", http.MethodGet, api.GetUserGroupMemberHandler, false, false)
	addRoute(router, "/api/v1/usergroup/{id}/member", http.MethodPost, api.AddUserGroupMemberHandler, true, true)
	addRoute(router, "/api/v1/usergroup/{id}/member", http.MethodDelete, api.RemoveUserGroupMemberHandler, true, true)
	addRoute(router, "/api/v1/usergroup/{id}/member", http.MethodPut, api.ToggleUserGroupModeratorHandler, true, true)
	addRoute(router, "/api/v1/search", http.MethodGet, api.SearchAllHandler, false, false, "search_all:r")
	addRoute(router, "/api/v1/search/user", http.MethodGet, api.SearchUserHandler, false, false)
	addRoute(router, "/api/v1/search/usergroup", http.MethodGet, api.SearchUsergroupHandler, false, false)
	addRoute(router, "/api/v1/search/userusergroup", http.MethodGet, api.SearchUserUsergroupHandler, false, false)
	addRoute(router, "/api/v1/search/tag", http.MethodGet, api.SearchTagHandler, false, false, "tags:r", "search_tags:r")
	addRoute(router, "/api/v1/search/article", http.MethodGet, api.SearchArticleHandler, false, false, "articles:r", "search_articles:r")
	addRoute(router, "/api/v1/search/list", http.MethodGet, api.SearchArticleListHandler, false, false, "lists:r", "search_lists:r")
	addRoute(router, "/api/v1/feed", http.MethodGet, api.GetFilteredFeedHandler, false, false)
	addRoute(router, "/api/v1/feed", http.MethodPut, api.ToggleNotificationReadHandler, false, false)
	addRoute(router, "/api/v1/observe", http.MethodPost, api.ObserveObjectHandler, false, false)
	addRoute(router, "/api/v1/observe", http.MethodGet, api.ReadObservedHandler, false, false)
	addRoute(router, "/api/v1/bookmark", http.MethodPost, api.BookmarkHandler, false, false)
	addRoute(router, "/api/v1/bookmark", http.MethodGet, api.ReadBookmarksHandler, false, false)
	addRoute(router, "/api/v1/pin", http.MethodPost, api.PinHandler, false, false)
	addRoute(router, "/api/v1/pin", http.MethodGet, api.ReadPinnedHandler, false, false, "pinned:r")
	addRoute(router, "/api/v1/articlelist", http.MethodPost, api.SaveArticleListHandler, false, true)
	addRoute(router, "/api/v1/articlelist/private", http.MethodGet, api.ReadPrivateArticleListsHandler, false, false)
	addRoute(router, "/api/v1/articlelist/{id}", http.MethodDelete, api.DeleteArticleListHandler, false, true)
	addRoute(router, "/api/v1/articlelist/{id}", http.MethodGet, api.GetArticleListHandler, false, false, "lists:r")
	addRoute(router, "/api/v1/articlelist/{id}/member", http.MethodGet, api.GetArticleListMemberHandler, false, false)
	addRoute(router, "/api/v1/articlelist/{id}/member", http.MethodPost, api.AddArticleListMemberHandler, false, true)
	addRoute(router, "/api/v1/articlelist/{id}/member", http.MethodDelete, api.RemoveArticleListMemberHandler, false, true)
	addRoute(router, "/api/v1/articlelist/{id}/member", http.MethodPut, api.ToggleArticleListModeratorHandler, false, true)
	addRoute(router, "/api/v1/articlelist/{id}/entry", http.MethodGet, api.GetArticleListEntriesHandler, false, false, "lists:r", "articles:r")
	addRoute(router, "/api/v1/articlelist/{id}/entry", http.MethodPost, api.AddArticleListEntryHandler, false, true)
	addRoute(router, "/api/v1/articlelist/{id}/entry", http.MethodDelete, api.RemoveArticleListEntryHandler, false, true)
	addRoute(router, "/api/v1/articlelist/{id}/entry", http.MethodPut, api.SortArticleListEntryHandler, false, true)
	addRoute(router, "/api/v1/tag", http.MethodPost, api.AddTagHandler, false, true)
	addRoute(router, "/api/v1/tag", http.MethodPut, api.RenameTagHandler, false, true)
	addRoute(router, "/api/v1/tag", http.MethodDelete, api.RemoveTagHandler, false, true)
	addRoute(router, "/api/v1/tag", http.MethodGet, api.ValidateTagHandler, false, true)
	addRoute(router, "/api/v1/tag/{id}", http.MethodDelete, api.DeleteTagHandler, false, true)
	addRoute(router, "/api/v1/tag/{name}", http.MethodGet, api.GetTagByNameHandler, false, false, "tags:r")
	addRoute(router, "/api/v1/user/member", http.MethodGet, api.GetMemberHandler, false, false)
	addRoute(router, "/api/v1/user/member", http.MethodPost, api.SaveSettingsHandler, false, false)
	addRoute(router, "/api/v1/profile", http.MethodGet, api.GetProfileByNameHandler, false, false)
	addRoute(router, "/api/v1/profile/{id}", http.MethodGet, api.GetProfileByIdHandler, false, false)
	addRoute(router, "/api/v1/support", http.MethodPost, api.ContactSupportHandler, false, false)
	addRoute(router, "/api/v1/client", http.MethodGet, api.ReadClientHandler, false, false)
	addRoute(router, "/api/v1/client", http.MethodPost, api.SaveClientHandler, true, false)
	addRoute(router, "/api/v1/client/{id}", http.MethodGet, api.ReadClientHandler, false, false)
	addRoute(router, "/api/v1/client/{id}", http.MethodPost, api.SaveClientHandler, true, false)
	addRoute(router, "/api/v1/client/{id}", http.MethodDelete, api.DeleteClientHandler, true, false)
	addRoute(router, "/api/v1/urlmeta", http.MethodGet, api.GetLinkMetaDataHandler, false, false)

	return router
}

func addRoute(router *mux.Router, path, method string, handler api.AuthHandler, requireExpert, requireWritePermissions bool, scopes ...string) {
	router.Handle(path, api.AuthMiddleware(handler, true, requireExpert, requireWritePermissions, scopes...)).Methods(method)
}

func connectDB() *db.Connection {
	backend := config.Get().BackendDB
	return db.NewConnection(db.ConnectionData{
		Host:               backend.Host,
		Port:               backend.Port,
		User:               backend.User,
		Password:           backend.Password,
		Schema:             backend.Schema,
		SSLMode:            backend.SSLMode,
		SSLCert:            backend.SSLCert,
		SSLKey:             backend.SSLKey,
		SSLRootCert:        backend.SSLRootCert,
		MaxOpenConnections: backend.MaxOpenConnections,
	})
}

func main() {
	config.Load()
	stdout, stderr := server.ConfigureLogging()
	defer server.CloseLogger(stdout, stderr)
	db.Migrate()
	auth.LoadConfig()
	feed.LoadConfig()
	mail.LoadConfig()
	i18n.LoadConfig()
	api.LoadConfig()
	article.LoadConfig()
	content.LoadConfig()
	member.LoadConfig()
	newsletter.LoadConfig()
	organization.LoadConfig()
	support.LoadConfig()
	billing.LoadConfig()
	article.InitTemplates()
	mailtpl.InitTemplates()
	connection := connectDB()
	model.SetConnection(connection)
	defer connection.Disconnect()
	router := setupRouter()
	cors := server.ConfigureCors(router)
	server.Start(cors, nil)
}
