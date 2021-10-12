package feed

// Reason is a single feed reason for a language.
type Reason struct {
	Feed         string
	Notification string
}

// Reasons is a list of all feed reasons.
// The first key is the language, the second one the reason.
// Attention! When adding a new reason activities.vue must be updated too!
var Reasons = map[string]map[string]Reason{
	"en": {
		"joined_organization": {
			Feed: `joined the organization.`,
		},
		"left_organization": {
			Feed: `left the organization.`,
		},
		"create_article": {
			Feed: `created a new article <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a>{{if (index .Content 0).Commit.Valid}}:<div class="message">{{(index .Content 0).Commit.String}}</div>{{else}}.{{end}}`,
		},
		"update_article": {
			`edited the article <a class="blue-100" href="{{.FrontendHost}}/read/{{(index .Articles 0).ID | IdToString}}">{{(index .Content 0).Title}}</a>{{if (index .Content 0).Commit.Valid}}:<div class="message">{{(index .Content 0).Commit.String}}</div>{{else}}.{{end}}`,
			`edited the article <a class="blue-100" href="{{.FrontendHost}}/read/{{(index .Articles 0).ID | IdToString}}">{{(index .Content 0).Title}}</a>{{if (index .Content 0).Commit.Valid}}:<div class="message">{{(index .Content 0).Commit.String}}</div>{{else}}.{{end}}`,
		},
		"reset_article": {
			`reset the article <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a>{{if (index .Content 0).Commit.Valid}}:<div class="message">{{(index .Content 0).Commit.String}}</div>{{else}}.{{end}}`,
			`reset the article <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a>{{if (index .Content 0).Commit.Valid}}:<div class="message">{{(index .Content 0).Commit.String}}</div>{{else}}.{{end}}`,
		},
		"delete_article_history_entry": {
			Feed: `modified the history of article <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a>.`,
		},
		"add_article_list_entry": {
			`added {{if eq (len .Articles) 1}}the article{{else}}the articles{{end}} {{$frontend := .FrontendHost}}{{range $i, $a := .Articles}}{{if $i}}, {{end}}<a class="blue-100" href="{{$frontend}}/read/{{$a.ID | IdToString}}">{{$a.LatestArticleContent.Title}}</a>{{end}} to list <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a>.`,
			`added {{if eq (len .Articles) 1}}the article{{else}}the articles{{end}} {{$frontend := .FrontendHost}}{{range $i, $a := .Articles}}{{if $i}}, {{end}}<a class="blue-100" href="{{$frontend}}/read/{{$a.ID | IdToString}}">{{$a.LatestArticleContent.Title}}</a>{{end}} to list <a class="green-100" href="{{.FrontendHost}}/list/{{(index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a>.`,
		},
		"add_protected_article_list_entry": {
			Feed: `added {{if eq (len .Articles) 0}}articles{{else}}{{if eq (len .Articles) 1}}the article{{else}}the articles{{end}} {{$frontend := .FrontendHost}}{{range $i, $a := .Articles}}{{if $i}}, {{end}}<a class="blue-100" href="{{$frontend}}/read/{{$a.ID | IdToString}}">{{$a.LatestArticleContent.Title}}</a>{{end}} and more{{end}} to list <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a>.`,
		},
		"add_article_list_member": {
			`added {{$frontend := .FrontendHost}}{{range $i, $u := .User}}{{if $i}}, {{end}}<a class="pink-100" href="{{$frontend}}/member/{{$u.OrganizationMember.Username}}">{{$u.Firstname}} {{$u.Lastname}}</a>{{end}}{{if (and (len .User) (len .Groups))}},{{end}} {{range $i, $g := .Groups}}{{if $i}}, {{end}}<a class="purple-100" href="{{$frontend}}/group/{{SlugWithId $g.Name $g.ID}}">{{$g.Name}}</a>{{end}} as {{if eq (Add (len .User) (len .Groups)) 1}}new member{{else}}new members{{end}} to list <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a>.`,
			`added you to list <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a>.`,
		},
		"remove_protected_article_list_entry": {
			Feed: `removed an article from list <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a>.`,
		},
		"remove_article_list_entry": {
			Feed: `removed {{if and (len .Articles) (len .Content)}}the article <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a>{{else}}an article{{end}} from list <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a>.`,
		},
		"remove_article_list_member": {
			`removed the member {{if .User}}<a class="pink-100" href="{{.FrontendHost}}/member/{{(index .User 0).OrganizationMember.Username}}">{{(index .User 0).Firstname}} {{(index .User 0).Lastname}}</a>{{else}}<a class="purple-100" href="{{.FrontendHost}}/group/{{SlugWithId (index .Groups 0).Name (index .Groups 0).ID}}">{{(index .Groups 0).Name}}</a>{{end}} from list <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a>.`,
			`removed you from list <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a>.`,
		},
		"create_article_list": {
			`created a new list <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a>.`,
			`created a new list <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a>.`,
		},
		"update_article_list": {
			`edited the list <a class="green-100" href="{{.FrontendHost}}/list/{{(index .Lists 0).ID | IdToString}}">{{(index .Lists 0).Name.Name}}</a>.`,
			`edited the list <a class="green-100" href="{{.FrontendHost}}/list/{{(index .Lists 0).ID | IdToString}}">{{(index .Lists 0).Name.Name}}</a>.`,
		},
		"set_article_list_moderator": {
			Feed: `granted moderator permissions to you for list <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a>.`,
		},
		"remove_article_list_moderator": {
			Feed: `removed your moderator permissions from list <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a>.`,
		},
		"add_user_group_member": {
			`added {{if eq (len .User) 1}}the member{{else}}the members{{end}} {{$frontend := .FrontendHost}}{{range $i, $u := .User}}{{if $i}}, {{end}}<a class="pink-100" href="{{$frontend}}/member/{{$u.OrganizationMember.Username}}">{{$u.Firstname}} {{$u.Lastname}}</a>{{end}} to group <a class="purple-100" href="{{.FrontendHost}}/group/{{SlugWithId (index .Groups 0).Name (index .Groups 0).ID}}">{{(index .Groups 0).Name}}</a>.`,
			`added you to group <a class="purple-100" href="{{.FrontendHost}}/group/{{(index .Groups 0).ID | IdToString}}">{{(index .Groups 0).Name}}</a>.`,
		},
		"remove_user_group_member": {
			`removed the member <a class="pink-100" href="{{.FrontendHost}}/member/{{(index .User 0).OrganizationMember.Username}}">{{(index .User 0).Firstname}} {{(index .User 0).Lastname}}</a> from group <a class="purple-100" href="{{.FrontendHost}}/group/{{SlugWithId (index .Groups 0).Name (index .Groups 0).ID}}">{{(index .Groups 0).Name}}</a>.`,
			`removed you from group <a class="purple-100" href="{{.FrontendHost}}/group/{{SlugWithId (index .Groups 0).Name (index .Groups 0).ID}}">{{(index .Groups 0).Name}}</a>.`,
		},
		"create_user_group": {
			`created a new group <a class="purple-100" href="{{.FrontendHost}}/group/{{SlugWithId (index .Groups 0).Name (index .Groups 0).ID}}">{{(index .Groups 0).Name}}</a>.`,
			`created a new group <a class="purple-100" href="{{.FrontendHost}}/group/{{SlugWithId (index .Groups 0).Name (index .Groups 0).ID}}">{{(index .Groups 0).Name}}</a>.`,
		},
		"update_user_group": {
			`edited the group <a class="purple-100" href="{{.FrontendHost}}/group/{{(index .Groups 0).ID | IdToString}}">{{(index .Groups 0).Name}}</a>.`,
			`edited the group <a class="purple-100" href="{{.FrontendHost}}/group/{{(index .Groups 0).ID | IdToString}}">{{(index .Groups 0).Name}}</a>.`,
		},
		"set_user_group_moderator": {
			Feed: `granted moderator permissions to you for group <a class="purple-100" href="{{.FrontendHost}}/group/{{SlugWithId (index .Groups 0).Name (index .Groups 0).ID}}">{{(index .Groups 0).Name}}</a>.`,
		},
		"remove_user_group_moderator": {
			Feed: `removed your moderator permissions from group <a class="purple-100" href="{{.FrontendHost}}/group/{{SlugWithId (index .Groups 0).Name (index .Groups 0).ID}}">{{(index .Groups 0).Name}}</a>.`,
		},
		"recommend_article": {
			`send you a recommendation for the article <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a>{{if (index .Vars "message")}}:<div class="message">{{index .Vars "message"}}</div>{{else}}.{{end}}`,
			`send you a recommendation for the article <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a>{{if (index .Vars "message")}}:<div class="message">{{index .Vars "message"}}</div>{{else}}.{{end}}`,
		},
		"invite_article": {
			Feed: `send you an invitation to edit the article <a class="blue-100" href="{{.FrontendHost}}/edit/{{(index .Articles 0).ID | IdToString}}?lang={{index .Vars "lang_id"}}">{{(index .Content 0).Title}}</a>{{if (index .Vars "message")}}: {{index .Vars "message"}}{{end}}.`,
		},
		"invite_new_article": {
			Feed: `send you an invitation to edit a <a class="blue-100" href="{{.FrontendHost}}/edit?room={{.Feed.RoomID.String}}">new article</a>.`,
		},
		"archived_article": {
			Feed: `{{if and (len .Articles) (len .Content)}}archived the article <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a>.{{else}}archived an article.{{end}}`,
		},
		"restored_article": {
			Feed: `{{if and (len .Articles) (len .Content)}}restored the article <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a>.{{else}}restored an article.{{end}}`,
		},
		"set_organization_moderator": {
			Feed: `granted organization moderator permissions to you.`,
		},
		"remove_organization_moderator": {
			Feed: `removed your organization moderator permissions.`,
		},
		"set_organization_admin": {
			Feed: `granted organization administrator permissions to you.`,
		},
		"remove_organization_admin": {
			Feed: `removed your organization administrator permissions.`,
		},
		"copy_article": {
			Feed: `duplicated the article <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a>.`,
		},
		"delete_article": {
			Feed: `deleted article "{{index .Vars "name"}}".`,
		},
		"delete_articlelist": {
			Feed: `deleted list "{{index .Vars "name"}}".`,
		},
		"delete_usergroup": {
			Feed: `deleted group "{{index .Vars "name"}}".`,
		},
		"delete_tag": {
			Feed: `deleted tag "{{index .Vars "name"}}".`,
		},
		"remove_member_read_only": {
			Feed: `granted write permissions to you.`,
		},
		"set_member_read_only": {
			Feed: `removed your write permissions.`,
		},
		"transfer_ownership": {
			Feed: `left the organization. The following objects have been assigned to the administrator and moderator groups: {{if (len .Articles)}}Articles: {{end}}{{$frontend := .FrontendHost}}{{range $i, $a := .Articles}}{{if $i}}, {{end}}<a class="blue-100" href="{{$frontend}}/read/{{$a.ID | IdToString}}">{{$a.LatestArticleContent.Title}}</a>{{end}} {{if (len .Lists)}}Lists: {{end}}{{range $i, $l := .Lists}}{{if $i}}, {{end}}<a class="green-100" href="{{$frontend}}/list/{{SlugWithId $l.Name.Name $l.ID}}">{{$l.Name.Name}}</a>{{end}} {{if (len .Groups)}}Groups: {{end}}{{range $i, $g := .Groups}}{{if $i}}, {{end}}<a class="purple-100" href="{{$frontend}}/group/{{SlugWithId $g.Name $g.ID}}">{{$g.Name}}</a>{{end}}`,
		},
		"mentioned": {
			Feed: `mentioned you in the article <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a>.`,
		},
		"recommendation_confirmation": {
			Feed: `has read the article <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a> you recommended.`,
		},
	},
	"de": {
		"joined_organization": {
			Feed: `ist der Organisation beigetreten.`,
		},
		"left_organization": {
			Feed: `hat die Organisation verlassen.`,
		},
		"create_article": {
			Feed: `hat einen neuen Artikel <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a> angelegt{{if (index .Content 0).Commit.Valid}}:<div class="message">{{(index .Content 0).Commit.String}}</div>{{else}}.{{end}}`,
		},
		"update_article": {
			`hat den Artikel <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a> bearbeitet{{if (index .Content 0).Commit.Valid}}:<div class="message">{{(index .Content 0).Commit.String}}</div>{{else}}.{{end}}`,
			`hat den Artikel <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a> bearbeitet{{if (index .Content 0).Commit.Valid}}:<div class="message">{{(index .Content 0).Commit.String}}</div>{{else}}.{{end}}`,
		},
		"reset_article": {
			`hat den Artikel <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a> zurückgesetzt{{if (index .Content 0).Commit.Valid}}:<div class="message">{{(index .Content 0).Commit.String}}</div>{{else}}.{{end}}`,
			`hat den Artikel <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a> zurückgesetzt{{if (index .Content 0).Commit.Valid}}:<div class="message">{{(index .Content 0).Commit.String}}</div>{{else}}.{{end}}`,
		},
		"delete_article_history_entry": {
			Feed: `hat die Historie des Artikels <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a> geändert.`,
		},
		"add_article_list_entry": {
			Feed: `hat {{if eq (len .Articles) 1}}einen Artikel{{else}}die Artikel{{end}} {{$frontend := .FrontendHost}}{{range $i, $a := .Articles}}{{if $i}}, {{end}}<a class="blue-100" href="{{$frontend}}/read/{{$a.ID | IdToString}}">{{$a.LatestArticleContent.Title}}</a>{{end}} zur Liste <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a> hinzugefügt.`,
		},
		"add_protected_article_list_entry": {
			Feed: `hat {{if eq (len .Articles) 0}}Artikel{{else}}{{if eq (len .Articles) 1}}einen Artikel{{else}}die Artikel{{end}} {{$frontend := .FrontendHost}}{{range $i, $a := .Articles}}{{if $i}}, {{end}}<a class="blue-100" href="{{$frontend}}/read/{{$a.ID | IdToString}}">{{$a.LatestArticleContent.Title}}</a>{{end}} und weitere{{end}} zur Liste <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a> hinzugefügt.`,
		},
		"add_article_list_member": {
			`hat {{$frontend := .FrontendHost}}{{range $i, $u := .User}}{{if $i}}, {{end}}<a class="pink-100" href="{{$frontend}}/member/{{$u.OrganizationMember.Username}}">{{$u.Firstname}} {{$u.Lastname}}</a>{{end}}{{if (and (len .User) (len .Groups))}},{{end}} {{range $i, $g := .Groups}}{{if $i}}, {{end}}<a class="purple-100" href="{{$frontend}}/group/{{SlugWithId $g.Name $g.ID}}">{{$g.Name}}</a>{{end}} als {{if eq (Add (len .User) (len .Groups)) 1}}neues Mitglied{{else}}neue Mitglieder{{end}} zur Liste <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a> hinzugefügt.`,
			`hat dich zur Liste <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a> hinzugefügt.`,
		},
		"remove_protected_article_list_entry": {
			Feed: `hat einen Artikel aus der Liste <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a> entfernt.`,
		},
		"remove_article_list_entry": {
			Feed: `hat {{if and (len .Articles) (len .Content)}}den Artikel <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a>{{else}}einen Artikel{{end}} aus der Liste <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a> entfernt.`,
		},
		"remove_article_list_member": {
			`hat das Mitglied {{if .User}}<a class="pink-100" href="{{.FrontendHost}}/member/{{(index .User 0).OrganizationMember.Username}}">{{(index .User 0).Firstname}} {{(index .User 0).Lastname}}</a>{{else}}<a class="purple-100" href="{{.FrontendHost}}/group/{{SlugWithId (index .Groups 0).Name (index .Groups 0).ID}}">{{(index .Groups 0).Name}}</a>{{end}} aus der Liste <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a> entfernt.`,
			`hat dich aus der Liste <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a> entfernt.`,
		},
		"create_article_list": {
			`hat eine neue Liste angelegt <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a>.`,
			`hat eine neue Liste angelegt <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a>.`,
		},
		"update_article_list": {
			`hat die Liste <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a> bearbeitet.`,
			`hat die Liste <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a> bearbeitet.`,
		},
		"set_article_list_moderator": {
			Feed: `hat dir Moderatoren Rechte in der Liste <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a> erteilt.`,
		},
		"remove_article_list_moderator": {
			Feed: `hat deine Moderatoren Rechte für die Liste <a class="green-100" href="{{.FrontendHost}}/list/{{SlugWithId (index .Lists 0).Name.Name (index .Lists 0).ID}}">{{(index .Lists 0).Name.Name}}</a> entzogen.`,
		},
		"add_user_group_member": {
			`hat {{if eq (len .User) 1}}das Mitglied{{else}}die Mitglieder{{end}} {{$frontend := .FrontendHost}}{{range $i, $u := .User}}{{if $i}}, {{end}}<a class="pink-100" href="{{$frontend}}/member/{{$u.OrganizationMember.Username}}">{{$u.Firstname}} {{$u.Lastname}}</a>{{end}} zur Gruppe <a class="purple-100" href="{{.FrontendHost}}/group/{{SlugWithId (index .Groups 0).Name (index .Groups 0).ID}}">{{(index .Groups 0).Name}}</a> hinzugefügt.`,
			`hat dich zur Gruppe <a class="purple-100" href="{{.FrontendHost}}/group/{{(index .Groups 0).ID | IdToString}}">{{(index .Groups 0).Name}}</a> hinzugefügt.`,
		},
		"remove_user_group_member": {
			`hat das Mitglied <a class="pink-100" href="{{.FrontendHost}}/member/{{(index .User 0).OrganizationMember.Username}}">{{(index .User 0).Firstname}} {{(index .User 0).Lastname}}</a> von der Gruppe <a class="purple-100" href="{{.FrontendHost}}/group/{{SlugWithId (index .Groups 0).Name (index .Groups 0).ID}}">{{(index .Groups 0).Name}}</a> entfernt.`,
			`hat dich aus der Gruppe <a class="purple-100" href="{{.FrontendHost}}/group/{{SlugWithId (index .Groups 0).Name (index .Groups 0).ID}}">{{(index .Groups 0).Name}}</a> entfernt.`,
		},
		"create_user_group": {
			`hat eine neue Gruppe <a class="purple-100" href="{{.FrontendHost}}/group/{{SlugWithId (index .Groups 0).Name (index .Groups 0).ID}}">{{(index .Groups 0).Name}}</a> angelegt.`,
			`hat eine neue Gruppe <a class="purple-100" href="{{.FrontendHost}}/group/{{SlugWithId (index .Groups 0).Name (index .Groups 0).ID}}">{{(index .Groups 0).Name}}</a> angelegt.`,
		},
		"update_user_group": {
			`hat die Gruppe <a class="purple-100" href="{{.FrontendHost}}/group/{{SlugWithId (index .Groups 0).Name (index .Groups 0).ID}}">{{(index .Groups 0).Name}}</a> bearbeitet.`,
			`hat die Gruppe <a class="purple-100" href="{{.FrontendHost}}/group/{{SlugWithId (index .Groups 0).Name (index .Groups 0).ID}}">{{(index .Groups 0).Name}}</a> bearbeitet.`,
		},
		"set_user_group_moderator": {
			Feed: `hat dir Moderatoren Rechte in der Gruppe <a class="purple-100" href="{{.FrontendHost}}/group/{{SlugWithId (index .Groups 0).Name (index .Groups 0).ID}}">{{(index .Groups 0).Name}}</a> erteilt.`,
		},
		"remove_user_group_moderator": {
			Feed: `hat deine Moderatoren Rechte für die Gruppe <a class="purple-100" href="{{.FrontendHost}}/group/{{SlugWithId (index .Groups 0).Name (index .Groups 0).ID}}">{{(index .Groups 0).Name}}</a> entzogen.`,
		},
		"recommend_article": {
			`hat dir den Artikel <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a> empfohlen{{if (index .Vars "message")}}:<div class="message">{{index .Vars "message"}}</div>{{else}}.{{end}}`,
			`hat dir den Artikel <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a> empfohlen{{if (index .Vars "message")}}:<div class="message">{{index .Vars "message"}}</div>{{else}}.{{end}}`,
		},
		"invite_article": {
			Feed: `hat dir eine Einladung gesendet den Artikel <a class="blue-100" href="{{.FrontendHost}}/edit/{{(index .Articles 0).ID | IdToString}}?lang={{index .Vars "lang_id"}}">{{(index .Content 0).Title}}</a> zu bearbeiten{{if (index .Vars "message")}}: {{index .Vars "message"}}{{end}}.`,
		},
		"invite_new_article": {
			Feed: `hat dir eine Einladung gesendet einen <a class="blue-100" href="{{.FrontendHost}}/edit?room={{.Feed.RoomID.String}}">neuen Artikel</a> zu bearbeiten.`,
		},
		"archived_article": {
			Feed: `{{if and (len .Articles) (len .Content)}}hat den Artikel <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a> archiviert.{{else}}hat einen Artikel archiviert.{{end}}`,
		},
		"restored_article": {
			Feed: `{{if and (len .Articles) (len .Content)}}hat den Artikel <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a> wiederhergestellt.{{else}}hat einen Artikel wiederhergestellt.{{end}}`,
		},
		"set_organization_moderator": {
			Feed: `hat dir Moderatoren Rechte für die Organisation erteilt.`,
		},
		"remove_organization_moderator": {
			Feed: `hat deine Moderatoren Rechte für die Organisation entzogen.`,
		},
		"set_organization_admin": {
			Feed: `hat dir Administratoren Rechte für die Organisation erteilt.`,
		},
		"remove_organization_admin": {
			Feed: `hat deine Administratoren Rechte für die Organisation entzogen.`,
		},
		"copy_article": {
			Feed: `hat den Artikel <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a> dupliziert.`,
		},
		"delete_article": {
			Feed: `hat den Artikel "{{index .Vars "name"}}" gelöscht.`,
		},
		"delete_articlelist": {
			Feed: `hat die Liste "{{index .Vars "name"}}" gelöscht.`,
		},
		"delete_usergroup": {
			Feed: `hat die Gruppe "{{index .Vars "name"}}" gelöscht.`,
		},
		"delete_tag": {
			Feed: `hat den Tag "{{index .Vars "name"}}" gelöscht.`,
		},
		"remove_member_read_only": {
			Feed: `hat dir Schreibrechte für die Organisation erteilt.`,
		},
		"set_member_read_only": {
			Feed: `hat deine Schreibrechte für die Organisation entzogen.`,
		},
		"transfer_ownership": {
			Feed: `hat die Organisation verlassen. Die folgenden Objekte wurden daher der Administrator und Moderator Gruppe zugeordnet: {{if (len .Articles)}}Artikel: {{end}}{{$frontend := .FrontendHost}}{{range $i, $a := .Articles}}{{if $i}}, {{end}}<a class="blue-100" href="{{$frontend}}/read/{{$a.ID | IdToString}}">{{$a.LatestArticleContent.Title}}</a>{{end}} {{if (len .Lists)}}Listen: {{end}}{{range $i, $l := .Lists}}{{if $i}}, {{end}}<a class="green-100" href="{{$frontend}}/list/{{SlugWithId $l.Name.Name $l.ID}}">{{$l.Name.Name}}</a>{{end}} {{if (len .Groups)}}Gruppen: {{end}}{{range $i, $g := .Groups}}{{if $i}}, {{end}}<a class="purple-100" href="{{$frontend}}/group/{{SlugWithId $g.Name $g.ID}}">{{$g.Name}}</a>{{end}}`,
		},
		"mentioned": {
			Feed: `hat dich im Artikel <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a> erwähnt.`,
		},
		"recommendation_confirmation": {
			Feed: `hat den Artikel <a class="blue-100" href="{{.FrontendHost}}/read/{{SlugWithId (index .Content 0).Title (index .Articles 0).ID}}">{{(index .Content 0).Title}}</a> gelesen, den du empfohlen hast.`,
		},
	},
}

// CheckReasonExists checks if the given reason exists for any language.
func CheckReasonExists(reason string) bool {
	_, ok := Reasons["en"][reason]
	return ok
}
