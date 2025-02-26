package servers

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

type HTTPServer struct {
	server *http.Server
	logger *zap.Logger
}

func NewHTTPServer(address string, port int, logger *zap.Logger, handler http.Handler) *HTTPServer {
	serverAddress := fmt.Sprintf("%s:%d", address, port)

	server := &http.Server{
		Addr:    serverAddress,
		Handler: handler,
	}

	return &HTTPServer{server: server, logger: logger}
}

func (s *HTTPServer) Run() error {
	s.logger.Info("starting HTTP server", zap.String("address", s.server.Addr))

	if err := s.server.ListenAndServe(); err != nil {
		return fmt.Errorf("unable to start HTTP server: %v", err)
	}

	return nil
}

func (s *HTTPServer) GracefulStop() error {
	s.logger.Info("gracefully shutting down HTTP server", zap.String("address", s.server.Addr))
	if err := s.server.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("unable to shutdown HTTP server: %v", err)
	}
	return nil
}

func (s *HTTPServer) ForcefulStop() error {
	s.logger.Info("forcefully shutting down HTTP server", zap.String("address", s.server.Addr))
	if err := s.server.Close(); err != nil {
		return fmt.Errorf("unable to shutdown HTTP server: %v", err)
	}
	return nil
}
