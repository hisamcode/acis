package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /signin", dynamic.ThenFunc(app.signin))
	mux.Handle("POST /signin", dynamic.ThenFunc(app.signinPost))
	mux.Handle("GET /signup", dynamic.ThenFunc(app.signup))
	mux.Handle("POST /signup", dynamic.ThenFunc(app.signupPost))
	mux.Handle("GET /user/activated", dynamic.ThenFunc(app.activateAccount))

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /transaction/create", dynamic.ThenFunc(app.transactionCreate))
	mux.Handle("POST /transaction", dynamic.ThenFunc(app.transactionPost))
	mux.Handle("GET /categories", dynamic.ThenFunc(app.categories))
	mux.Handle("GET /categories/{id}", dynamic.ThenFunc(app.categoriesView))
	mux.Handle("GET /categories/create", dynamic.ThenFunc(app.categories))
	mux.Handle("POST /categories", dynamic.ThenFunc(app.categoriesPost))
	mux.Handle("GET /categories/{id}/edit", dynamic.ThenFunc(app.categoriesEdit))
	mux.Handle("PUT /categories/{id}", dynamic.ThenFunc(app.categoriesPut))
	return mux

}
