package article

import (
	"emviwiki/backend/article/schema"
	"emviwiki/backend/article/util"
	"emviwiki/backend/context"
	"emviwiki/backend/errs"
	"emviwiki/backend/prosemirror"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
)

const (
	schemaParagraphType = "paragraph"
)

func GetArticlePreview(ctx context.EmviContext, articleId, langId hide.ID, extractParagraph bool) (string, error) {
	article, err := util.GetArticleWithAccess(nil, ctx, articleId, true)

	if err != nil {
		return "", err
	}

	content := GetArticleContent(ctx.Organization.ID, ctx.UserId, articleId, langId, 0)

	if content.Content == "" {
		return "", nil
	}

	if ctx.IsClient() {
		util.RemoveNonPublicInformationFromContent(ctx, content)
	}

	doc, err := prosemirror.ParseDoc(content.Content)

	if err != nil {
		logbuch.Error("Error parsing document when reading article preview", logbuch.Fields{"err": err, "article_id": article.ID})
		return "", errs.ArticleNotFound
	}

	// consider just one empty paragraph to be completely empty
	if len(doc.Content) == 0 || (len(doc.Content) == 1 && doc.Content[0].Type == "paragraph" && extractTextFromContent(doc) == "") {
		return "", nil
	}

	// render preview
	contentNode := doc

	if extractParagraph {
		nodes := prosemirror.FindNodes(doc, 1, schemaParagraphType)

		if len(nodes) > 0 {
			contentNode = &nodes[0]
		}
	}

	content.Content, err = RenderDocument(ctx, ctx.Organization.ID, ctx.UserId, content.LanguageId, contentNode, schema.HTMLSchema)

	if err != nil {
		logbuch.Error("Error rendering document when reading article preview", logbuch.Fields{"err": err, "article_id": article.ID})
		return "", errs.ArticleNotFound
	}

	return content.Content, nil
}
