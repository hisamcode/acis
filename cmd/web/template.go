package main

import (
	"bytes"
	"html/template"
	"net/http"
)

type LayoutBase byte

const (
	LayoutClean LayoutBase = iota
	LayoutStandard
)

func (l LayoutBase) String() string {
	return []string{"clean-base", "base"}[l]
}

func (app *application) render(w http.ResponseWriter, base LayoutBase, page string) {
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
	err := ts.ExecuteTemplate(buf, base.String(), nil)
	if err != nil {
		app.logger.Error(err.Error())
		return
	}
	buf.WriteTo(w)
}
