package server

import (
	"context"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os/signal"
	"sample-golang-project/api"
	"sample-golang-project/config"
	"sample-golang-project/kafka"
	"sample-golang-project/types"
	"sync"
	"syscall"
)

type Server struct {
	deviceController *api.DeviceController
	deviceRepo       types.DeviceRepository
	devicePub        types.DevicePublisher
	deviceListener   *kafka.DeviceListener
	router           *mux.Router
	cfg              *config.Configuration
	logger           *zap.SugaredLogger
	wg               sync.WaitGroup
}

func New(deviceController *api.DeviceController, deviceRepository types.DeviceRepository, router *mux.Router, publisher types.DevicePublisher, cfg *config.Configuration,
	logger *zap.SugaredLogger, deviceListener *kafka.DeviceListener) *Server {

	return &Server{
		deviceController: deviceController,
		deviceRepo:       deviceRepository,
		devicePub:        publisher,
		router:           router,
		deviceListener:   deviceListener,
		cfg:              cfg,
		logger:           logger,
		wg:               sync.WaitGroup{},
	}
}

func (s *Server) Start(ctx context.Context) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)
	defer stop()
	errWg, ctx := errgroup.WithContext(ctx)
	errWg.Go(func() error {
		return s.deviceRepo.Start(ctx)
	})

	errWg.Go(func() error {
		return s.devicePub.StartPublish(ctx)
	})

	errWg.Go(func() error {
		return s.deviceListener.StartListen()
	})

	errWg.Go(func() error {
		s.logger.Info("Starting server on", s.cfg.ServerPort)
		if err := http.ListenAndServe(s.cfg.ServerPort, s.router); err != nil {
			s.logger.Errorf("failed to start server on port %s", s.cfg.ServerPort)
			return err
		}
		return nil
	})

	errWg.Go(func() error {
		<-ctx.Done()
		stop()
		s.Stop(ctx)
		return nil
	})

	err := errWg.Wait()
	if err == context.Canceled || err == nil {
		s.logger.Info("gracefully quit server")
	} else if err != nil {
		s.logger.Error("server quit with error", err)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.deviceRepo.Stop(ctx); err != nil {
		return err
	}
	if err := s.deviceListener.Stop(); err != nil {
		return err
	}
	return nil
}
