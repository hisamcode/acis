package main

import (
	"net/http"

	"github.com/hisamcode/acis/internal/data"
	"github.com/justinas/nosurf"
)

type templateData struct {
	Form          any
	TokenActivate string
	CSRFToken     string
	Projects      []data.Project
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CSRFToken: nosurf.Token(r),
	}
}
