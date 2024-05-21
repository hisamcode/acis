package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/hisamcode/acis/internal/data"
	"github.com/hisamcode/acis/internal/session"
	"go.uber.org/zap"
)

func (app *application) getProject(r *http.Request) (*data.Project, error) {
	projectID, err := strconv.ParseInt(r.PathValue("projectID"), 10, 64)
	if err != nil {
		return nil, err
	}

	project, err := app.DB.Project.Get(projectID)
	if err != nil {
		return nil, err
	}

	return project, nil
}

// execute event on client
func (app *application) addHXTriggerAfterSettle(w http.ResponseWriter, eventName string) {
	w.Header().Add("HX-Trigger-After-Settle", eventName)

}

type hxswap uint8

const (
	HXSWAP_INNER hxswap = iota
)

func (h hxswap) String() string {
	return []string{"innerHTML"}[h]
}

func (app *application) addHXReswap(w http.ResponseWriter, swap hxswap) {
	w.Header().Add("HX-Reswap", swap.String())

}
func (app *application) isAuthenticated(r *http.Request) int64 {
	// return app.sessionManager.Exists(r.Context(), session.SessionAuthenticatedUserID)
	if !app.sessionManager.Exists(r.Context(), session.SessionAuthenticatedUserID) {
		return -1
	}
	return app.sessionManager.GetInt64(r.Context(), session.SessionAuthenticatedUserID)
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
	LayoutPartials
	LayoutProject
)

// use on application.render()
func (l LayoutBase) String() string {
	return []string{"clean-base", "base", "partials", "project"}[l]
}

func (app *application) renderServerError(w http.ResponseWriter, err error) {
	app.logger.WithOptions(zap.AddCallerSkip(1)).Error(err.Error())
	app.render(w, http.StatusInternalServerError, LayoutClean, "server-error.html", templateData{})
}
func (app *application) renderEditConflict(w http.ResponseWriter, err error) {
	app.logger.Error(err.Error())
	app.render(w, http.StatusConflict, LayoutClean, "edit-conflict.html", templateData{})
}

func (app *application) render(w http.ResponseWriter, status int, base LayoutBase, page string, data templateData) error {
	var ts *template.Template
	var ok bool

	switch base {
	case LayoutClean:
		ts, ok = app.templateCache.clean[page]
		if !ok {
			return fmt.Errorf("the template '%s' does not exist", page)
		}
	case LayoutStandard:
		ts, ok = app.templateCache.standard[page]
		if !ok {
			return fmt.Errorf("the template '%s' does not exist", page)
		}
	case LayoutPartials:
		ts, ok = app.templateCache.partials[page]
		if !ok {
			return fmt.Errorf("the template '%s' does not exist", page)
		}
	case LayoutProject:
		ts, ok = app.templateCache.project[page]
		if !ok {
			return fmt.Errorf("the template '%s' does not exist", page)
		}
	default:
		return fmt.Errorf("the template layout is required")
	}

	buf := new(bytes.Buffer)
	err := ts.ExecuteTemplate(buf, base.String(), data)
	if err != nil {
		app.logger.Error(err.Error())
		return nil
	}
	w.WriteHeader(status)
	buf.WriteTo(w)
	return nil
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
	partials map[string]*template.Template
	project  map[string]*template.Template
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

	cachePartials, err := newTemplateCachePartials()
	if err != nil {
		return templateCache{}, err
	}

	cacheProject, err := newTemplateCacheProject()
	if err != nil {
		return templateCache{}, err
	}

	return templateCache{
		clean:    cacheClean,
		standard: cacheStandard,
		partials: cachePartials,
		project:  cacheProject,
	}, nil
}

func humanDate(t time.Time) string {
	// return t.Format("02 JAN 2006 at 15:04")
	return t.Format(time.RFC822)
}

var functions = template.FuncMap{
	"humanDate": humanDate,
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

		ts, err := template.New(name).Funcs(functions).ParseFiles(patterns...)
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

		ts, err := template.New(name).Funcs(functions).ParseFiles(patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
func newTemplateCachePartials() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/htmx/partials/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFiles(patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func newTemplateCacheProject() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/htmx/pages/project/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"./ui/htmx/bases/project.html",
			"./ui/htmx/bases/head.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFiles(patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
