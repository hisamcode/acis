package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/go-playground/form/v4"
	"github.com/hisamcode/acis/internal/session"
)

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), session.SessionAuthenticatedUserID)
}

func (app *application) background(fn func()) {
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		defer func() {
			if err := recover(); err != nil {
				app.logger.Error(fmt.Sprintf("%v", err))
			}
		}()

		fn()
	}()
}

type LayoutBase byte

const (
	LayoutClean LayoutBase = iota
	LayoutStandard
)

// use on application.render()
func (l LayoutBase) String() string {
	return []string{"clean-base", "base"}[l]
}

func (app *application) renderServerError(w http.ResponseWriter, err error) {
	app.logger.Error(err.Error())
	app.render(w, http.StatusInternalServerError, LayoutClean, "server-error.html", templateData{})
}
func (app *application) renderEditConflict(w http.ResponseWriter, err error) {
	app.logger.Error(err.Error())
	app.render(w, http.StatusConflict, LayoutClean, "edit-conflict.html", templateData{})
}

func (app *application) render(w http.ResponseWriter, status int, base LayoutBase, page string, data templateData) {
	var ts *template.Template
	var ok bool

	if base == LayoutClean {
		ts, ok = app.templateCache.clean[page]
		if !ok {
			app.logger.Error("the template does not exist", "template", page)
			return
		}
	}

	if base == LayoutStandard {
		ts, ok = app.templateCache.standard[page]
		if !ok {
			app.logger.Error("the template does not exist", "template", page)
			return
		}
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, base.String(), data)
	if err != nil {
		app.logger.Error(err.Error())
		return
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			app.logger.Error(err.Error())
		}

		return err
	}

	return nil
}

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
