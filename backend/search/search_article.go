package search

import (
	articles "emviwiki/backend/article"
	"emviwiki/backend/article/schema"
	articleutil "emviwiki/backend/article/util"
	"emviwiki/backend/context"
	"emviwiki/backend/prosemirror"
	"emviwiki/shared/model"
	"emviwiki/shared/util"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"strings"
	"sync"
)

const (
	schemaParagraphType = "paragraph"
	schemaImageType     = "image"
	schemaImageAttrSrc  = "src"
)

// Performs a fuzzy search for articles.
// Joins the latest article content in users preferred language if available.
func SearchArticle(ctx context.EmviContext, query string, filter *model.SearchArticleFilter) ([]model.Article, int) {
	query = strings.TrimSpace(query)

	// search for all fields when no filter was passed
	if filter == nil {
		filter = new(model.SearchArticleFilter)
	}

	filter.ClientAccess = ctx.IsClient()
	langId := util.DetermineLang(nil, ctx.Organization.ID, ctx.UserId, filter.LanguageId).ID
	var wg sync.WaitGroup
	wg.Add(2)
	var results []model.Article
	var resultCount int

	go func() {
		results = model.FindArticleByOrganizationIdAndUserIdAndLanguageIdAndQueryAndFilterLimit(ctx.Organization.ID, ctx.UserId, langId, query, filter)
		wg.Done()
	}()

	go func() {
		resultCount = model.CountArticleByOrganizationIdAndUserIdAndLanguageIdAndQueryAndFilterLimit(ctx.Organization.ID, ctx.UserId, query, filter)
		wg.Done()
	}()

	wg.Wait()
	articleutil.RemoveAuthorsOrAuthorMails(ctx, results)

	if filter.Preview || filter.PreviewParagraph {
		for i := range results {
			articlePreview(ctx, &results[i], langId, filter.PreviewParagraph, filter.PreviewImage)
		}
	}

	return results, resultCount
}

func articlePreview(ctx context.EmviContext, article *model.Article, langId hide.ID, extractParagraph, extractImage bool) {
	content := articles.GetArticleContent(ctx.Organization.ID, ctx.UserId, article.ID, langId, 0)

	if ctx.IsClient() {
		articleutil.RemoveNonPublicInformationFromContent(ctx, content)
	}

	doc, err := prosemirror.ParseDoc(content.Content)

	if err != nil {
		logbuch.Error("Error parsing document when building article preview", logbuch.Fields{"err": err, "article_id": article.ID})
		return
	}

	// extract first paragraph in preview mode
	contentNode := doc

	if extractParagraph {
		nodes := prosemirror.FindNodes(doc, 1, schemaParagraphType)

		if len(nodes) > 0 {
			contentNode = &nodes[0]
		}
	}

	article.LatestArticleContent.Content, err = articles.RenderDocument(ctx, ctx.Organization.ID, ctx.UserId, langId, contentNode, schema.HTMLSchema)

	if err != nil {
		logbuch.Error("Error rendering article preview", logbuch.Fields{"err": err, "article_id": article.ID})
		return
	}

	// find and set preview image
	if extractImage {
		nodes := prosemirror.FindNodes(doc, 1, schemaImageType)

		if len(nodes) > 0 {
			src, ok := nodes[0].Attrs[schemaImageAttrSrc].(string)

			if ok {
				article.PreviewImage = src
			}
		}
	}
}
