package main

type templateData struct {
	Form          any
	TokenActivate string
}

func (app *application) newTemplateData() templateData {
	return templateData{}
}
