package main

type templateData struct {
	Form any
}

func (app *application) newTemplateData() templateData {
	return templateData{}
}
