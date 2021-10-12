package tpl

import (
	"html/template"
	"sync"
)

// Cache caches a single HTML template.
type Cache struct {
	temp     template.Template
	funcMap  template.FuncMap
	dir      string
	disabled bool
	loaded   bool
	m        sync.RWMutex
}

// NewCache creates a new template cache with the default set of template functions.
func NewCache(dir string, disabled bool) *Cache {
	cache := &Cache{funcMap: funcMap, dir: dir, disabled: disabled}

	if err := cache.loadTemplate(); err != nil {
		panic(err)
	}

	return cache
}

// Get returns the HTML template or loads it in case the cache is disabled or it hasn't been loaded yet.
func (cache *Cache) Get() *template.Template {
	if cache.disabled || !cache.loaded {
		if err := cache.loadTemplate(); err != nil {
			panic(err)
		}
	}

	cache.m.RLock()
	defer cache.m.RUnlock()
	return &cache.temp
}

func (cache *Cache) loadTemplate() error {
	cache.m.Lock()
	defer cache.m.Unlock()
	t, err := template.New("template").Funcs(cache.funcMap).ParseGlob(cache.dir)

	if err != nil {
		return err
	}

	cache.temp = *t
	cache.loaded = true
	return nil
}

// Enable enables caching.
func (cache *Cache) Enable() {
	cache.m.Lock()
	defer cache.m.Unlock()
	cache.disabled = false
}

// Disable disables caching.
func (cache *Cache) Disable() {
	cache.m.Lock()
	defer cache.m.Unlock()
	cache.disabled = true
}

// Clear clears the cache.
func (cache *Cache) Clear() {
	cache.m.Lock()
	defer cache.m.Unlock()
	cache.loaded = false
}
