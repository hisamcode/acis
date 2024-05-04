package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/hisamcode/acis/internal/mailer"
	"github.com/hisamcode/acis/internal/repository"
	"github.com/hisamcode/acis/internal/repository/postgres"
	_ "github.com/lib/pq"
)

const (
	DEVELOPMENT = "development"
	STAGING     = "staging"
	PRODUCTION  = "production"
)

type config struct {
	host string
	port int
	env  string
	db   struct {
		dsn string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	activationAccountDuration string
}

type DB struct {
	User  repository.UserDatabaseRepoer
	Token repository.TokenDatabaseRepoer
}

type application struct {
	config         config
	logger         *slog.Logger
	DB             DB
	templateCache  templateCache
	formDecoder    *form.Decoder
	mailer         mailer.Mailer
	wg             sync.WaitGroup
	sessionManager *scs.SessionManager
}

func main() {
	var cfg config
	flag.StringVar(&cfg.host, "addr", "127.0.0.1:8000", "addr web server, ip or domain")
	flag.IntVar(&cfg.port, "port", 8000, "web server port")
	flag.StringVar(&cfg.env, "env", DEVELOPMENT, "environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://acis:acis@localhost/acis?sslmode=disable", "postgreSQL dsn")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "301cac5016d7a1", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "4a3d921250b5ad", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Acis <no-reply@acis.hisam.my.id>", "SMTP sender")

	flag.StringVar(&cfg.activationAccountDuration, "activation-account-duration", "72h", "Activation account duration when signup")

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

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	// automatically expire after 12 hours after first being created
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Persist = false

	app := application{
		config: cfg,
		logger: logger,
		DB: DB{
			User:  postgres.UserModel{DB: db},
			Token: postgres.TokenModel{DB: db},
		},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		mailer:         mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
		sessionManager: sessionManager,
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

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
