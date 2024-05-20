package main

import (
	"fmt"
	"net/http"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, LayoutStandard, "home.html", data)
}

func (app *application) transactionCreate(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, LayoutClean, "transaction-create.html", templateData{})
}
func (app *application) transactionPost(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *application) categoriesView(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "categories page view")
}
func (app *application) categories(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "categories page list")
}
func (app *application) categoriesPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "create categories")
}
func (app *application) categoriesEdit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "categories edit")
}
func (app *application) categoriesPut(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "categories put")
}
