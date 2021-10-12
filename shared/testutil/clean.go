package testutil

import (
	auth "emviwiki/auth/model"
	dashboard "emviwiki/dashboard/model"
	"emviwiki/shared/model"
	"testing"
)

func CleanBackendDb(t *testing.T) {
	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "article_recommendation"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "article_visit"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "newsletter_subscription"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "client_scope"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "client"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "domain_blacklist"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "support_ticket"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "bookmark"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "file"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "invitation"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "observed_object"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "feed_access"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "feed_ref"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "feed"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "article_list_entry"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "article_list_member"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "article_list_name"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "article_list"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "article_tag"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "tag"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "article_content_author"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "article_access"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "article_content"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "article"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "user_group_member"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "user_group"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "organization_member"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "language"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "organization"`); err != nil {
		t.Fatal(err)
	}

	if _, err := model.GetConnection().Exec(nil, `DELETE FROM "user"`); err != nil {
		t.Fatal(err)
	}
}

func CleanAuthDb(t *testing.T) {
	if _, err := auth.GetConnection().Exec(nil, `DELETE FROM "login"`); err != nil {
		t.Fatal(err)
	}

	if _, err := auth.GetConnection().Exec(nil, `DELETE FROM "scope"`); err != nil {
		t.Fatal(err)
	}

	if _, err := auth.GetConnection().Exec(nil, `DELETE FROM "client"`); err != nil {
		t.Fatal(err)
	}

	if _, err := auth.GetConnection().Exec(nil, `DELETE FROM "access_grant"`); err != nil {
		t.Fatal(err)
	}

	if _, err := auth.GetConnection().Exec(nil, `DELETE FROM "user"`); err != nil {
		t.Fatal(err)
	}

	if _, err := auth.GetConnection().Exec(nil, `DELETE FROM "email_blacklist"`); err != nil {
		t.Fatal(err)
	}
}

func CleanDashboardDb(t *testing.T) {
	if _, err := dashboard.GetConnection().Exec(nil, `DELETE FROM "newsletter"`); err != nil {
		t.Fatal(err)
	}

	if _, err := dashboard.GetConnection().Exec(nil, `DELETE FROM "user"`); err != nil {
		t.Fatal(err)
	}
}
