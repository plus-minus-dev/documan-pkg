package grpcserver

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	server *grpc.Server
	port   int
}

func New(port int, opts ...grpc.ServerOption) *Server {
	return &Server{
		server: grpc.NewServer(opts...),
		port:   port,
	}
}

func (s *Server) Server() *grpc.Server {
	return s.server
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("grpcserver listen: %w", err)
	}

	go func() {
		if err := s.server.Serve(lis); err != nil {
			panic(fmt.Sprintf("grpcserver serve: %v", err))
		}
	}()

	return nil
}

func (s *Server) GracefulStop() {
	s.server.GracefulStop()
}
