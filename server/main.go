package main

import (
	"fmt"
	"log"
	"net"

	"git.ouroath.com/peng/test/grpc/api"
	"git.ouroath.com/peng/test/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	// FQDN of the server.
	// TODO: should ideally come from config.
	FQDN = "localhost"
	// TODO: should ideally come from config.
	port = 4444
)

// main start a gRPC server and waits for connection
func main() {
	// create a listener on TCP port
	addr := fmt.Sprintf("%s:%d", FQDN, port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a server instance
	s := api.Server{}

	// Create the TLS credentials.
	// cert and key files should securly be deployed to the server. File path should come from config.
	creds, err := credentials.NewServerTLSFromFile("./../tls/server.crt", "./../tls/server.key")
	if err != nil {
		log.Fatalf("could not load TLS keys: %s", err)
	}

	// Create an array of gRPC options with the credentials
	opts := []grpc.ServerOption{grpc.Creds(creds)}

	// create a gRPC server object
	grpcServer := grpc.NewServer(opts...)

	// attach the Ping service to the server
	proto.RegisterPingServer(grpcServer, &s)

	// start the server
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}

/*
NOTES:

- the main function starts by creating a TCP listener on the port you want to bind your gRPC server to;
- then you create instance of your `Server`
- create an instance of a gRPC server
- register the service
- start the gRPC server

Adding TLS support
- you created a credentials object (called creds) from your certificate and key files;
- you created a grpc.ServerOption array and placed your credentials object in it;
- when creating the grpc server, you provided the constructor with you array of grpc.ServerOption;
- you must have noticed that you need to precisely specify the IP you bind your server to, so that the IP matches the FQDN used in the certificate.
*/
