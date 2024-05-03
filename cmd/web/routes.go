package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)

	mux.Handle("GET /signin", dynamic.ThenFunc(app.signin))
	mux.Handle("POST /signin", dynamic.ThenFunc(app.signinPost))
	mux.Handle("GET /signup", dynamic.ThenFunc(app.signup))
	mux.Handle("POST /signup", dynamic.ThenFunc(app.signupPost))
	mux.Handle("POST /signout", dynamic.ThenFunc(app.signout))
	mux.Handle("GET /user/activated", dynamic.ThenFunc(app.activateAccount))

	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /home", protected.ThenFunc(app.home))
	mux.Handle("GET /transaction/create", protected.ThenFunc(app.transactionCreate))
	mux.Handle("POST /transaction", protected.ThenFunc(app.transactionPost))
	mux.Handle("GET /categories", protected.ThenFunc(app.categories))
	mux.Handle("GET /categories/{id}", protected.ThenFunc(app.categoriesView))
	mux.Handle("GET /categories/create", protected.ThenFunc(app.categories))
	mux.Handle("POST /categories", protected.ThenFunc(app.categoriesPost))
	mux.Handle("GET /categories/{id}/edit", protected.ThenFunc(app.categoriesEdit))
	mux.Handle("PUT /categories/{id}", protected.ThenFunc(app.categoriesPut))
	return mux

}
