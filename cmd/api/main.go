package main

import (
	"flag"
	"log"

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

	app := &application{
		config: cfg,
		auth:   authMethod,
		logger: logger,
	}

	app.serve()
}
