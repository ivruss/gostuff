package grpc_client

import (
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCConn struct {
	Conn   *grpc.ClientConn
	logger *zap.Logger
}

func NewGRPCConn(endpoint string, logger *zap.Logger) (*GRPCConn, error) {
	clientConn, err := grpc.NewClient(
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		return nil, fmt.Errorf("error connecting to grpc server: %v", err)
	}

	return &GRPCConn{
		Conn:   clientConn,
		logger: logger,
	}, nil
}

func (c *GRPCConn) Close() error {
	c.logger.Info("closing grpc_client connection")
	if err := c.Conn.Close(); err != nil {
		return fmt.Errorf("error closing grpc_client connection: %v", err)
	}
	return nil
}
