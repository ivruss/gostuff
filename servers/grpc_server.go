package servers

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type GRPCServer struct {
	server  *grpc.Server
	address string
	logger  *zap.Logger
}

func NewGRPCServer(address string, port int, logger *zap.Logger) *GRPCServer {
	server := grpc.NewServer()

	return &GRPCServer{
		address: fmt.Sprintf("%s:%d", address, port),
		logger:  logger,
		server:  server,
	}
}

func (s *GRPCServer) Run() error {
	s.logger.Info("starting gRPC server", zap.String("address", s.address))
	l, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen gRPC server address: %w", err)
	}

	if err := s.server.Serve(l); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}

func (s *GRPCServer) GracefulStop() error {
	s.logger.Sugar().Info("gracefully stopping gRPC server", zap.String("address", s.address))
	s.server.GracefulStop()
	return nil
}

func (s *GRPCServer) ForcefulStop() error {
	s.logger.Sugar().Info("force stopping gRPC server", zap.String("address", s.address))
	s.server.Stop()
	return nil
}
