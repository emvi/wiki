package article

import (
	"archive/zip"
	"bytes"
	"emviwiki/backend/article/schema"
	filecontent "emviwiki/backend/content"
	"emviwiki/backend/context"
	"emviwiki/backend/errs"
	"emviwiki/backend/prosemirror"
	"emviwiki/shared/feed"
	"emviwiki/shared/i18n"
	"emviwiki/shared/model"
	"fmt"
	"github.com/emvi/hide"
	"github.com/emvi/logbuch"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

const (
	formatHTML        = "html"
	formatMarkdown    = "markdown"
	formatHTMLExt     = "html"
	formatMarkdownExt = "md"
	zipExt            = "zip"
	exportFileDir     = "files"
	exportTemplate    = "export.html"
)

var exportHTMLI18n = i18n.Translation{
	"en": {
		// used in HTML template
		"created_by": "Authors:",
		"created_on": "Created on",
		"created":    "",
		"updated_on": "Last updated on",
		"updated":    "",

		// used in Markdown
		"tags":      "Tags",
		"authors":   "Authors",
		"published": "Published",
		"changed":   "Last changed",
	},
	"de": {
		// used in HTML template
		"created_by": "Autoren:",
		"created_on": "am",
		"created":    "erstellt",
		"updated_on": "am",
		"updated":    "zuletzt geändert",

		// used in Markdown
		"tags":      "Tags",
		"authors":   "Autoren",
		"published": "Veröffentlicht",
		"changed":   "Zuletzt geändert",
	},
}

// ExportArticle exports an article to given format inside a zip archive.
// The format can be either html or markdown. If includeFiles is set to true, all attachments will be included.
func ExportArticle(ctx context.EmviContext, articleId, langId hide.ID, format string, includeFiles bool) (io.Reader, error) {
	logbuch.Debug("Exporting article", logbuch.Fields{"article_id": articleId, "lang_id": langId, "format": format})
	format = strings.ToLower(format)

	if format != formatHTML && format != formatMarkdown {
		return nil, errs.UnknownArticleFormat
	}

	article, err := checkUserReadAccess(ctx.Organization.ID, ctx.UserId, articleId)

	if err != nil {
		return nil, err
	}

	content := GetArticleContent(ctx.Organization.ID, ctx.UserId, articleId, langId, 0)

	if content == nil || content.Content == "" {
		return nil, errs.UnpublishedArticle
	}

	doc, err := prosemirror.ParseDoc(content.Content)

	if err != nil {
		logbuch.Warn("Error parsing article content to prosemirror document on export", logbuch.Fields{
			"err":             err,
			"organization_id": ctx.Organization.ID,
			"article_id":      articleId,
			"lang_id":         langId,
			"format":          format,
			"include_files":   includeFiles,
		})
		return nil, err
	}

	var files []model.File

	if includeFiles {
		files = findAndUpdateFilePathsForFilesAndImages(doc)
	}

	var ext string

	if format == formatHTML {
		ext = formatHTMLExt
		content.Content, err = RenderDocument(ctx, ctx.Organization.ID, ctx.UserId, content.LanguageId, doc, schema.HTMLSchema)

		if err != nil {
			// return technical error here because this is an "exception" state
			return nil, err
		}

		// in case we export to HTML, render the template surrounding the actual article content (for styling)
		if err := renderArticleTemplate(ctx, article, content); err != nil {
			return nil, err
		}
	} else {
		ext = formatMarkdownExt
		content.Content, err = RenderDocument(ctx, ctx.Organization.ID, ctx.UserId, content.LanguageId, doc, schema.GetMarkdownSchema(ctx.Organization))

		if err != nil {
			// return technical error here because this is an "exception" state
			return nil, err
		}

		if err := addMarkdownMetaData(ctx, article, content); err != nil {
			return nil, err
		}
	}

	reader, writer := io.Pipe()
	go createExportZip(writer, content, ext, files)
	return reader, nil
}

func findAndUpdateFilePathsForFilesAndImages(doc *prosemirror.Node) []model.File {
	files := findAndUpdateFilePaths(doc, "image", "src")
	files = mergeExportFileMaps(files, findAndUpdateFilePaths(doc, "file", "file"))
	files = mergeExportFileMaps(files, findAndUpdateFilePaths(doc, "pdf", "src"))
	fileList := make([]model.File, 0, len(files))

	for _, file := range files {
		fileList = append(fileList, file)
	}

	return fileList
}

func findAndUpdateFilePaths(doc *prosemirror.Node, node, attr string) map[hide.ID]model.File {
	files := make(map[hide.ID]model.File)

	prosemirror.TransformNodes(doc, node, func(node *prosemirror.Node) {
		filename := getNodeAttrAsString(node, attr)

		if filename != "" {
			file := model.GetFileByUniqueName(filepath.Base(filename))

			if file != nil {
				files[file.ID] = *file
				node.Attrs[attr] = getExportFilePath(file)
			}
		}
	})

	return files
}

func getNodeAttrAsString(node *prosemirror.Node, attr string) string {
	a, ok := node.Attrs[attr]

	if ok {
		var str string
		str, ok = a.(string)

		if ok {
			return str
		}
	}

	return ""
}

func getExportFilePath(file *model.File) string {
	originalName := strings.TrimSuffix(file.OriginalName, filepath.Ext(file.OriginalName))
	uniqueName := strings.TrimSuffix(file.UniqueName, filepath.Ext(file.UniqueName))
	filename := fmt.Sprintf("%s_%s%s", originalName, uniqueName, filepath.Ext(file.UniqueName))
	return filepath.Join(exportFileDir, filename)
}

func mergeExportFileMaps(a, b map[hide.ID]model.File) map[hide.ID]model.File {
	for k, v := range b {
		a[k] = v
	}

	return a
}

func renderArticleTemplate(ctx context.EmviContext, article *model.Article, content *model.ArticleContent) error {
	lang := model.GetLanguageByOrganizationIdAndId(ctx.Organization.ID, content.LanguageId)

	if lang == nil {
		return errs.LanguageNotFound
	}

	var buffer bytes.Buffer
	data := struct {
		Vars      map[string]template.HTML
		LangCode  string
		Title     string
		Content   template.HTML
		RTL       bool
		Authors   []model.User
		Tags      []model.Tag
		Published time.Time
		Updated   time.Time
	}{
		i18n.GetVars(lang.Code, exportHTMLI18n),
		lang.Code,
		content.Title,
		template.HTML(content.Content),
		content.RTL,
		getAuthors(ctx, article.ID),
		model.FindTagByOrganizationIdAndUserIdAndArticleId(ctx.Organization.ID, ctx.UserId, article.ID),
		article.Published.Time,
		content.DefTime,
	}

	if err := tplCache.Get().ExecuteTemplate(&buffer, exportTemplate, data); err != nil {
		logbuch.Error("Error rendering export HTML template", logbuch.Fields{"err": err})
		return err
	}

	content.Content = buffer.String()
	return nil
}

func addMarkdownMetaData(ctx context.EmviContext, article *model.Article, content *model.ArticleContent) error {
	lang := model.GetLanguageByOrganizationIdAndId(ctx.Organization.ID, content.LanguageId)

	if lang == nil {
		return errs.LanguageNotFound
	}

	tags := model.FindTagByOrganizationIdAndUserIdAndArticleId(ctx.Organization.ID, ctx.UserId, article.ID)
	tagList := make([]string, 0, len(tags))

	for _, tag := range tags {
		tagList = append(tagList, tag.Name)
	}

	authors := getAuthors(ctx, article.ID)
	authorsList := make([]string, 0, len(authors))

	for _, author := range authors {
		if author.Firstname != "" {
			authorsList = append(authorsList, fmt.Sprintf("%s %s", author.Firstname, author.Lastname))
		} else {
			authorsList = append(authorsList, author.OrganizationMember.Username)
		}
	}

	vars := i18n.GetVars(lang.Code, exportHTMLI18n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n\n", content.Title))
	sb.WriteString(fmt.Sprintf("%s: %s\n", vars["tags"], strings.Join(tagList, ", ")))
	sb.WriteString(fmt.Sprintf("%s: %s\n", vars["authors"], strings.Join(authorsList, ", ")))
	sb.WriteString(fmt.Sprintf("%s: %s\n", vars["published"], article.Published.Time.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("%s: %s\n\n", vars["changed"], content.DefTime.Format("2006-01-02")))
	sb.WriteString(content.Content)
	content.Content = sb.String()
	return nil
}

func createExportZip(writer *io.PipeWriter, content *model.ArticleContent, ext string, files []model.File) {
	zipFile := zip.NewWriter(writer)

	// add article file
	articleFile, err := zipFile.Create(fmt.Sprintf("%s.%s", feed.SlugWithId(content.Title, content.ArticleId), ext))

	if err != nil {
		logbuch.Error("Error creating article file in export zip", logbuch.Fields{"err": err, "article_id": content.ArticleId})
		return
	}

	if _, err := articleFile.Write([]byte(content.Content)); err != nil {
		logbuch.Error("Error writing article to export zip", logbuch.Fields{"err": err, "article_id": content.ArticleId})
		return
	}

	// add attachments
	for _, file := range files {
		f, err := zipFile.Create(getExportFilePath(&file))

		if err != nil {
			logbuch.Error("Error creating attachment file in export zip", logbuch.Fields{"err": err, "article_id": content.ArticleId})
			break
		}

		_, reader, err := filecontent.ReadFile(file.UniqueName)

		if err != nil {
			logbuch.Error("Error reading attachment file to write it to export zip", logbuch.Fields{"err": err, "article_id": content.ArticleId, "file_id": file.ID})
			break
		}

		buf, err := ioutil.ReadAll(reader)

		if err != nil {
			logbuch.Error("Error reading attachment file bytes to write it to export zip", logbuch.Fields{"err": err, "article_id": content.ArticleId, "file_id": file.ID})
			break
		}

		if _, err := f.Write(buf); err != nil {
			logbuch.Error("Error writing attachment file to export zip", logbuch.Fields{"err": err, "article_id": content.ArticleId, "file_id": file.ID})
			break
		}

		if err := reader.Close(); err != nil {
			logbuch.Error("Error closing attachment file on export", logbuch.Fields{"err": err, "article_id": content.ArticleId, "file_id": file.ID})
			break
		}
	}

	if err := zipFile.Close(); err != nil {
		logbuch.Error("Error closing article export zip file", logbuch.Fields{"err": err, "article_id": content.ArticleId})
	}

	if err := writer.Close(); err != nil {
		logbuch.Error("Error closing writer to zip file", logbuch.Fields{"err": err, "article_id": content.ArticleId})
		return
	}
}
