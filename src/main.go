package main

import (
	"context"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"os/signal"
	"sample-golang-project/config"
	"sample-golang-project/mongodb"
	"sample-golang-project/server"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	defer stop()
	errWg, ctx := errgroup.WithContext(ctx)

	cfg := config.NewConfig()

	basicLogger, err := zap.NewProduction()
	if err != nil {
		panic("failed to create logger")
	}
	logger := basicLogger.Sugar()

	imp, err := mongodb.New(cfg, logger)
	if err != nil {
		logger.Error("failed to create mongo repository")
	}

	s := server.New(imp, mux.NewRouter(), cfg, logger)

	errWg.Go(func() error {
		return s.Start(ctx)
	})

	errWg.Go(func() error {
		<-ctx.Done()
		stop()
		return s.Stop(ctx)
	})

	err = errWg.Wait()
	if err == context.Canceled || err == nil {
		logger.Info("gracefully quit server")
	} else if err != nil {
		logger.Error("server quit with error", err)
	}

}
