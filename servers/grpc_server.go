package servers

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type GRPCServer struct {
	Server  *grpc.Server
	address string
	logger  *zap.Logger
}

func NewGRPCServer(address string, port int, logger *zap.Logger) *GRPCServer {
	server := grpc.NewServer()

	return &GRPCServer{
		address: fmt.Sprintf("%s:%d", address, port),
		logger:  logger,
		Server:  server,
	}
}

func (s *GRPCServer) Run() error {
	s.logger.Info("starting grpc_client server", zap.String("address", s.address))
	l, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to listen grpc_client server address: %w", err)
	}

	if err := s.Server.Serve(l); err != nil {
		return fmt.Errorf("failed to serve grpc_client server: %w", err)
	}

	return nil
}

func (s *GRPCServer) GracefulStop() error {
	s.logger.Sugar().Info("gracefully stopping grpc_client server", zap.String("address", s.address))
	s.Server.GracefulStop()
	return nil
}

func (s *GRPCServer) ForcefulStop() error {
	s.logger.Sugar().Info("force stopping grpc_client server", zap.String("address", s.address))
	s.Server.Stop()
	return nil
}
