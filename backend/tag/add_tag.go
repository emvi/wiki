package tag

import (
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"strings"
	"unicode"
	"unicode/utf8"

	"emviwiki/backend/errs"
	"emviwiki/shared/model"
)

const (
	tagMaxLen         = 60
	MaxTagsPerArticle = 50
)

type AddTagData struct {
	ArticleId hide.ID `json:"article_id"`
	Tag       string  `json:"tag"`
}

func (data *AddTagData) validate() error {
	data.Tag = strings.TrimSpace(data.Tag)

	if data.Tag == "" {
		return errs.TagEmpty
	}

	if utf8.RuneCountInString(data.Tag) > tagMaxLen {
		return errs.TagLen
	}

	if err := data.checkTagNameValid(); err != nil {
		return err
	}

	return nil
}

func (data *AddTagData) checkTagNameValid() error {
	for _, c := range data.Tag {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c != '.' && c != '-' && c != '_' && c != ' ' {
			return errs.TagInvalid
		}
	}

	return nil
}

func AddTag(organization *model.Organization, data AddTagData) error {
	if err := data.validate(); err != nil {
		return err
	}

	article := model.GetArticleByOrganizationIdAndIdIgnoreArchived(organization.ID, data.ArticleId)

	if article == nil {
		return errs.ArticleNotFound
	}

	if model.GetArticleTagByOrganizationIdAndArticleIdAndName(organization.ID, article.ID, data.Tag) != nil {
		return errs.TagExistsAlready
	}

	if model.CountArticleTagByArticleId(article.ID) >= MaxTagsPerArticle {
		return errs.MaxTagsReached
	}

	// save new tag if not exists already
	tag := model.GetTagByOrganizationIdAndName(organization.ID, data.Tag)

	if tag == nil {
		tag = &model.Tag{OrganizationId: organization.ID, Name: strings.ToLower(data.Tag)}

		if err := model.SaveTag(nil, tag); err != nil {
			logbuch.Error("Error creating new tag when adding tag to article", logbuch.Fields{"err": err})
			return errs.Saving
		}
	}

	articleTag := &model.ArticleTag{ArticleId: article.ID, TagId: tag.ID}

	if err := model.SaveArticleTag(nil, articleTag); err != nil {
		logbuch.Error("Error saving article tag when adding tag", logbuch.Fields{"err": err})
		return errs.Saving
	}

	return nil
}

func ValidateTag(data AddTagData) error {
	return data.validate()
}
