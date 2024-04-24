package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	slog *slog.Logger
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./ui/htmx/pages/home.html", "./ui/htmx/bases/layout.html")
	if err != nil {
		app.slog.Error(err.Error())
		return
	}

	err = t.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.slog.Error(err.Error())
		return
	}
}
func (app *application) transactionPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "transactionPost")
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

func (app *application) render() {

}

func main() {

	slog := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	app := application{
		slog: slog,
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("POST /transaction", app.transactionPost)
	mux.HandleFunc("GET /categories", app.categories)
	mux.HandleFunc("GET /categories/{id}", app.categoriesView)
	mux.HandleFunc("GET /categories/create", app.categories)
	mux.HandleFunc("POST /categories", app.categoriesPost)
	mux.HandleFunc("GET /categories/{id}/edit", app.categoriesEdit)
	mux.HandleFunc("PUT /categories/{id}", app.categoriesPut)

	server := http.Server{}

	server.Addr = "127.0.0.1:8000"
	server.Handler = mux

	slog.Info("listen server", "addr", server.Addr)
	err := server.ListenAndServe()
	slog.Error(err.Error())
	os.Exit(1)

}
