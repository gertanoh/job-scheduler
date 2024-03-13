package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"time"

	"gertanoh.job-scheduler/internal/authenticator"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

// application config struct
type application struct {
	config config
	auth   *authenticator.Authenticator
	logger *zap.Logger
}

func main() {

	var cfg config

	flag.IntVar(&cfg.port, "port", 8000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL DSN")

	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to loav env vars %v", err)
	}
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	authMethod, err := authenticator.New()
	if err != nil {
		logger.Fatal("Fail to setup Auth", zap.Error(err))
	}

	logger.Info("Auth setup")

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal("Fail to setup db", zap.Error(err))
	}

	defer db.Close()
	logger.Info("DB connection setup")

	app := &application{
		config: cfg,
		auth:   authMethod,
		logger: logger,
	}

	app.serve()
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	duration, err := time.ParseDuration("15m")

	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
