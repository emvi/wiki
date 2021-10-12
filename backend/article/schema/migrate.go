package schema

import (
	"emviwiki/backend/prosemirror"
	"emviwiki/shared/constants"
	"emviwiki/shared/model"
	"encoding/json"
	"github.com/emvi/logbuch"
)

var (
	migrationSteps = []migration{
		func(content *model.ArticleContent) error { return nil }, // dummy
		migrateVersion2,
	}
)

type migration func(*model.ArticleContent) error

// Migrate migrates the given article content to newest version if required.
func Migrate(content *model.ArticleContent) error {
	if content != nil && content.Content != "" && content.SchemaVersion < constants.LatestSchemaVersion {
		for i := content.SchemaVersion; i < len(migrationSteps); i++ {
			if err := migrationSteps[i](content); err != nil {
				logbuch.Error("Error migrating article content to new schema version", logbuch.Fields{"err": err, "schema_version": content.SchemaVersion, "latest_schema_version": constants.LatestSchemaVersion, "migration_step": i})
				return err
			}
		}

		content.SchemaVersion = constants.LatestSchemaVersion

		if err := model.SaveArticleContent(nil, content); err != nil {
			logbuch.Error("Error saving article content after migrating to newest schema version", logbuch.Fields{"err": err, "schema_version": content.SchemaVersion, "latest_schema_version": constants.LatestSchemaVersion})
			return err
		}
	}

	return nil
}

/* Add caption paragraph to images.
 * Before:
 *
 * {
 * 	"type":"image",
 * 		"attrs":{
 * 		"src":"some_url"
 * 	}
 * }
 *
 * After:
 *
 * {
 * 	"type":"image",
 * 	"attrs":{
 * 		"src":"some_url"
 * 	},
 * 	"content":[
 * 		{
 * 			"type":"paragraph"
 * 		}
 * 	]
 * }
 */
func migrateVersion2(content *model.ArticleContent) error {
	doc, err := prosemirror.ParseDoc(content.Content)

	if err != nil {
		return err
	}

	prosemirror.TransformNodes(doc, "image", func(node *prosemirror.Node) {
		if len(node.Content) == 0 {
			node.Content = append(node.Content, prosemirror.Node{Type: "paragraph"})
		}
	})

	docJson, err := json.Marshal(doc)

	if err != nil {
		return err
	}

	content.Content = string(docJson)
	return nil
}
