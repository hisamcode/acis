package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

type templateData struct {
	Form          any
	TokenActivate string
	CSRFToken     string
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CSRFToken: nosurf.Token(r),
	}
}
