package model

import (
	"github.com/emvi/hide"
	"time"
)

type SearchArticleFilter struct {
	BaseSearch

	LanguageId       hide.ID   `json:"language_id"`
	Archived         bool      `json:"archived"`
	WIP              bool      `json:"wip"`
	ClientAccess     bool      `json:"client_access"`
	Preview          bool      `json:"preview"`           // return the full article content
	PreviewParagraph bool      `json:"preview_paragraph"` // return the first text paragraph (overwrites "Preview" if paragraph is found)
	PreviewImage     bool      `json:"preview_image"`     // return the first image URL in the article
	Title            string    `json:"title"`
	Content          string    `json:"content"`
	Tags             string    `json:"tags"`
	TagIds           []hide.ID `json:"tag_ids"`
	AuthorUserIds    []hide.ID `json:"authors"`
	UserGroupIds     []hide.ID `json:"user_group_ids"`
	Commits          string    `json:"commits"`
	PublishedStart   time.Time `json:"published_start"`
	PublishedEnd     time.Time `json:"published_end"`
	SortTitle        string    `json:"sort_title"`
	SortPublished    string    `json:"sort_published"`
	SortRelevance    string    `json:"sort_relevance"`
}

type SearchArticleListFilter struct {
	BaseSearch

	ClientAccess bool      `json:"public"`
	Name         string    `json:"name"`
	Info         string    `json:"info"`
	UserIds      []hide.ID `json:"user_ids"`
	UserGroupIds []hide.ID `json:"user_group_ids"`
	SortName     string    `json:"sort_name"`
	SortInfo     string    `json:"sort_info"`
}

type SearchArticleListEntryFilter struct {
	BaseSearch

	ClientAccess    bool      `json:"client_access"`
	Archived        bool      `json:"archived"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	Tags            string    `json:"tags"`
	AuthorUserIds   []hide.ID `json:"authors"`
	Commits         string    `json:"commits"`
	CenterArticleId hide.ID   `json:"center_article_id"` // center result set around article (with limit)
	CenterBefore    int       `json:"center_before"`     // number of results before center article
	SortPosition    string    `json:"sort_position"`
	SortTitle       string    `json:"sort_title"`
}

type SearchArticleListMemberFilter struct {
	BaseSearch

	UserIds       []hide.ID `json:"user"`
	SortUsername  string    `json:"sort_username"`
	SortEmail     string    `json:"sort_email"`
	SortFirstname string    `json:"sort_firstname"`
	SortLastname  string    `json:"sort_lastname"`
}

type SearchFeedFilter struct {
	BaseSearch

	Notifications bool      `json:"notifications"`
	Unread        bool      `json:"unread"`
	Reasons       []string  `json:"reasons"`
	UserIds       []hide.ID `json:"user"`
}

type SearchTagFilter struct {
	BaseSearch

	SortUsages string `json:"sort_usages"`
	SortName   string `json:"sort_name"`
}

type SearchUserFilter struct {
	BaseSearch

	Username      string `json:"username"`
	Email         string `json:"email"`
	Firstname     string `json:"firstname"`
	Lastname      string `json:"lastname"`
	SortUsername  string `json:"sort_username"`
	SortEmail     string `json:"sort_email"`
	SortFirstname string `json:"sort_firstname"`
	SortLastname  string `json:"sort_lastname"`
}

type SearchUserGroupFilter struct {
	BaseSearch

	Name     string    `json:"name"`
	Info     string    `json:"info"`
	UserIds  []hide.ID `json:"user_ids"`
	SortName string    `json:"sort_name"`
	SortInfo string    `json:"sort_info"`

	// used to exclude group results when organization is not expert
	// not relevant for the query, but in search function
	FindGroups bool `json:"-"`
}

type SearchUserGroupMemberFilter struct {
	BaseSearch

	UserIds       []hide.ID `json:"user_ids"`
	SortUsername  string    `json:"sort_username"`
	SortEmail     string    `json:"sort_email"`
	SortFirstname string    `json:"sort_firstname"`
	SortLastname  string    `json:"sort_lastname"`
}
