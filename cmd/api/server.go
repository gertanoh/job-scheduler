package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func (app *application) serve() {

	e := app.Router()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	defer stop()

	// Start server in go routine
	go func() {
		if err := e.Start(fmt.Sprintf(":%d", app.config.port)); err != nil && err != http.ErrServerClosed {
			app.logger.Fatal("Shutting down the server")
		}
	}()

	// Wait for the stop signals
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	app.logger.Info("Shutting down the server")
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		app.logger.Fatal("When gracefully shuttind down, an error happended", zap.Error(err))
	}
}
