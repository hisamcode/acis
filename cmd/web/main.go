package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/hisamcode/acis/internal/repository"
	"github.com/hisamcode/acis/internal/repository/postgres"
	_ "github.com/lib/pq"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type DB struct {
	User repository.UserDatabaseRepo
}

type application struct {
	config        config
	logger        *slog.Logger
	DB            DB
	templateCache templateCache
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

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := application{
		config: cfg,
		logger: logger,
		DB: DB{
			User: postgres.UserModel{DB: db},
		},
		templateCache: templateCache,
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
