package main

import (
	"html/template"
	"path/filepath"
)

type templateCache struct {
	clean    map[string]*template.Template
	standard map[string]*template.Template
}

func newTemplateCache() (templateCache, error) {
	cacheClean, err := newTemplateCacheClean()
	if err != nil {
		return templateCache{}, err
	}

	cacheStandard, err := newTemplateCacheStandard()
	if err != nil {
		return templateCache{}, err
	}

	return templateCache{
		clean:    cacheClean,
		standard: cacheStandard,
	}, nil
}

func newTemplateCacheClean() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/htmx/pages/clean/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"./ui/htmx/bases/clean.html",
			"./ui/htmx/bases/head.html",
			page,
		}

		ts, err := template.New(name).ParseFiles(patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
func newTemplateCacheStandard() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/htmx/pages/standard/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"./ui/htmx/bases/standard.html",
			"./ui/htmx/bases/head.html",
			page,
		}

		ts, err := template.New(name).ParseFiles(patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
