package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /signup", app.signup)
	mux.HandleFunc("POST /signup", app.signupPost)

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /transaction/create", app.transactionCreate)
	mux.HandleFunc("POST /transaction", app.transactionPost)
	mux.HandleFunc("GET /categories", app.categories)
	mux.HandleFunc("GET /categories/{id}", app.categoriesView)
	mux.HandleFunc("GET /categories/create", app.categories)
	mux.HandleFunc("POST /categories", app.categoriesPost)
	mux.HandleFunc("GET /categories/{id}/edit", app.categoriesEdit)
	mux.HandleFunc("PUT /categories/{id}", app.categoriesPut)
	return mux

}
