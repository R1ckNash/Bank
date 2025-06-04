package http_server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
)

type Server struct {
	log        *slog.Logger
	httpServer *http.Server
	port       int
}

func New(log *slog.Logger, handler http.Handler, port int) *Server {
	return &Server{
		log:        log,
		httpServer: &http.Server{Handler: handler},
		port:       port,
	}
}

func (a *Server) MustRun() {
	if err := a.Run(); err != nil {
		panic("server could not be started")
	}
}

func (a *Server) Run() error {
	const op = "httpserver.Run"
	log := a.log.With(slog.String("op", op), slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("starting http server: ", slog.String("addr", l.Addr().String()))

	if err := a.httpServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *Server) Stop(ctx context.Context) error {
	const op = "http.Stop"
	a.log.With(slog.String("op", op)).
		Info("stopping http server", slog.Int("port", a.port))
	return a.httpServer.Shutdown(ctx)
}
