package server

import (
	"context"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"sample-golang-project/api"
	"sample-golang-project/config"
	"sample-golang-project/types"
)

type Server struct {
	controllerRepo types.ControllerRepository
	router         *mux.Router
	cfg            *config.Configuration
	logger         *zap.SugaredLogger
}

func New(controllerRepo types.ControllerRepository, router *mux.Router, cfg *config.Configuration, logger *zap.SugaredLogger) *Server {
	api.NewControllerEndpoints(router, controllerRepo)
	return &Server{
		controllerRepo: controllerRepo,
		router:         router,
		cfg:            cfg,
		logger:         logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.controllerRepo.Start(ctx); err != nil {
		return err
	}
	s.logger.Info("Starting server on", s.cfg.ServerPort)
	if err := http.ListenAndServe(s.cfg.ServerPort, s.router); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.controllerRepo.Stop(ctx); err != nil {
		return err
	}
	return nil
}
