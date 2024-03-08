package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"go.uber.org/zap"

)

type config struct {
	env  string
	db   struct {
		dsn string
	}
}

// application config struct
type application struct {
	config config
	logger *zap.Logger
}

// Add bash scripts pull golang image before executing executor
func main() {
	ctx := context.Background()

	var cfg config

	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to loav env vars %v", err)
	}
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()


	app := &application{
		config: cfg,
		logger: logger,
	}

	// main
	// pull a job from nats queue, try to get a lock from zookeeper,
	// if fails drop the job
	// else launch the job in a container, periodically do heartbeat to zookeeper


}
