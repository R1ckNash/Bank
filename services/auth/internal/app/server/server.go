package server

import (
	pkgerrors "Bank/pkg/errors"
	"context"
	"fmt"
	"go.uber.org/zap"
	"net"
	"net/http"
)

type Server struct {
	log        *zap.Logger
	httpServer *http.Server
	port       int
}

func New(log *zap.Logger, handler http.Handler, port int) *Server {
	return &Server{
		log:  log,
		port: port,
		httpServer: &http.Server{
			Handler: handler,
		},
	}
}

func (s *Server) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return pkgerrors.Wrap("http listen error", err)
	}

	s.log.Info("HTTP server started", zap.String("addr", l.Addr().String()))
	return s.httpServer.Serve(l)
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("Shutting down HTTP server gracefully...")
	return s.httpServer.Shutdown(ctx)
}
