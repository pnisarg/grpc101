package api

import (
	"context"
	"log"

	"git.ouroath.com/peng/test/grpc/proto"
)

// Server represents the gRPC server.
// It implements proto.PingServer.
type Server struct {
}

// SayHello generates response to a Ping request.
func (s *Server) SayHello(ctx context.Context, request *proto.PingMessage) (*proto.PingMessage, error) {
	log.Printf("SayHello: received %s", request.Greeting)
	return &proto.PingMessage{Greeting: "pong"}, nil
}
