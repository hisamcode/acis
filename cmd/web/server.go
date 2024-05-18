package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func (app *application) serve() error {

	server := http.Server{}
	server.Addr = fmt.Sprintf("127.0.0.1:%d", app.config.port)
	server.Handler = app.routes()
	server.IdleTimeout = time.Minute
	server.ReadTimeout = 5 * time.Second
	server.WriteTimeout = 10 * time.Second

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.logger.Info("shutting down server", zap.String("signal", s.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.logger.Info("completing background tasks", zap.String("addr", server.Addr))

		app.wg.Wait()
		shutdownError <- nil
	}()

	app.logger.Info("starting server", zap.String("addr", server.Addr), zap.String("env", app.config.env))

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.Info("stopped server", zap.String("addr", server.Addr))
	return nil
}
