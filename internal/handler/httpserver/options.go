package httpserver

import (
	"go-clean-template/pkg/config"

	"go.uber.org/zap"
)

type Options func(s *Server) error

func WithConfig(cfg *config.Config) Options {
	return func(s *Server) error {
		s.Config = cfg
		return nil
	}
}

func WithLogger(l *zap.SugaredLogger) Options {
	return func(s *Server) error {
		s.Logger = l
		return nil
	}
}
