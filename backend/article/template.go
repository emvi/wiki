package article

import (
	"emviwiki/shared/config"
	"emviwiki/shared/tpl"
)

var (
	tplCache *tpl.Cache
)

func InitTemplates() {
	c := config.Get()
	tplCache = tpl.NewCache(c.Template.TemplateDir, c.Template.HotReload)
}
