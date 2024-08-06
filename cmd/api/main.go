package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/tconnellan/macro-tracker-backend/internal/data"
	"github.com/tconnellan/macro-tracker-backend/internal/jsonlog"
	"github.com/tconnellan/macro-tracker-backend/internal/mailer"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn                string
		maxOpenConnections int
		maxIdleConns       int
		maxIdleTime        string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("MACROTRACKER_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConnections, "db-max-open-conn", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conn", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max idle connections")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	// re-enable with new defaults
	// flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	// flag.IntVar(&cfg.smtp.port, "smtp-port", 2525, "SMTP port")
	// flag.StringVar(&cfg.smtp.username, "smtp-username", "6136cf491d04e3", "SMTP username")
	// flag.StringVar(&cfg.smtp.password, "smtp-password", "faae129e4eac5b", "SMTP password")
	// flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Greenlight <no-reply@letsgo2.eb.com>", "SMTP sender")

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModel(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

// func openDB3(cfg config) (*sql.DB, error) {
// 	db, err := sql.Open("postgres", cfg.db.dsn)
// 	if err != nil {
// 		return nil, err
// 	}

// 	db.SetMaxOpenConns(cfg.db.maxOpenConnections)
// 	db.SetMaxIdleConns(cfg.db.maxIdleConns)
// 	idleTime, err := time.ParseDuration(cfg.db.maxIdleTime)
// 	if err != nil {
// 		return nil, err
// 	}
// 	db.SetConnMaxIdleTime(idleTime)

// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	err = db.PingContext(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return db, nil

// }

func openDB(cfg config) (*pgxpool.Pool, error) {
	connConfig, err := pgxpool.ParseConfig(cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	connConfig.MaxConns = int32(cfg.db.maxOpenConnections)
	// connConfig.MaxIdleConns = int32(cfg.db.maxIdleConns) # pgx doesn't implement maxidleconns

	idleTime, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	connConfig.MaxConnIdleTime = idleTime

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
