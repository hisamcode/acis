package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	logger *slog.Logger
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./ui/htmx/pages/home.html", "./ui/htmx/bases/layout.html", "./ui/htmx/bases/head.html")
	if err != nil {
		app.logger.Error(err.Error())
		return
	}

	err = t.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.logger.Error(err.Error())
		return
	}
}
func (app *application) transactionCreate(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./ui/htmx/pages/transaction-create.html", "./ui/htmx/bases/clean.html", "./ui/htmx/bases/head.html")
	if err != nil {
		app.logger.Error(err.Error())
		return
	}

	err = t.ExecuteTemplate(w, "clean-base", nil)
	if err != nil {
		app.logger.Error(err.Error())
		return
	}
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

func (app *application) render(w http.ResponseWriter) {

}

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

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

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 8000, "web server port")
	flag.StringVar(&cfg.env, "env", "development", "environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://acis:acis@localhost/acis?sslmode=disable", "postgreSQL dsn")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	app := application{
		config: cfg,
		logger: logger,
	}

	server := http.Server{}
	if app.config.env == "development" {
		server.Addr = fmt.Sprintf("127.0.0.1:%d", app.config.port)
	} else {
		server.Addr = fmt.Sprintf(":%d", app.config.port)
	}
	server.Handler = app.routes()
	server.IdleTimeout = time.Minute
	server.ReadTimeout = 5 * time.Second
	server.WriteTimeout = 10 * time.Second
	server.ErrorLog = slog.NewLogLogger(logger.Handler(), slog.LevelError)

	logger.Info("starting server", "addr", server.Addr, "env", app.config.env)
	err = server.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil

}
