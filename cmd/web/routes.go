package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "oke")
	})

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)

	mux.Handle("GET /signin", dynamic.ThenFunc(app.signin))
	mux.Handle("POST /signin", dynamic.ThenFunc(app.signinPost))
	mux.Handle("GET /signup", dynamic.ThenFunc(app.signup))
	mux.Handle("POST /signup", dynamic.ThenFunc(app.signupPost))
	mux.Handle("POST /signout", dynamic.ThenFunc(app.signout))
	mux.Handle("GET /user/activated", dynamic.ThenFunc(app.activateAccount))

	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /home", protected.ThenFunc(app.home))
	mux.Handle("GET /projects", protected.ThenFunc(app.projects))
	mux.Handle("POST /projects", protected.ThenFunc(app.projectPost))
	mux.Handle("GET /projects/latest", protected.ThenFunc(app.latestProjects))

	mux.Handle("GET /projects/{projectID}/home", protected.ThenFunc(app.project))
	mux.Handle("GET /projects/{projectID}/daily", protected.ThenFunc(app.latestTransactionDaily))
	mux.Handle("GET /projects/{projectID}/monthly", protected.ThenFunc(app.latestTransactionMonthly))
	mux.Handle("GET /projects/{projectID}/settings", protected.ThenFunc(app.projectSetting))
	mux.Handle("POST /projects/{projectID}/transaction", protected.ThenFunc(app.projectTransactionPost))
	mux.Handle("PUT /projects/{projectID}", protected.ThenFunc(app.projectSettingPut))

	mux.Handle("POST /projects/{projectID}/emoji", protected.ThenFunc(app.projectEmojiPost))
	mux.Handle("PUT /projects/{projectID}/emoji", protected.ThenFunc(app.projectEmojiPut))
	mux.Handle("POST /projects/{projectID}/emoji/delete", protected.ThenFunc(app.projectEmojiDelete))

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
