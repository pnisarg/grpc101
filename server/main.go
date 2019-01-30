package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"git.ouroath.com/peng/test/grpc/api"
	"git.ouroath.com/peng/test/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	// FQDN of the server.
	// TODO: following values should ideally come from config.
	FQDN       = "localhost"
	port       = 4444
	certPath   = "./../tls/server.crt"
	keyPath    = "./../tls/server.key"
	caCertPath = "./../tls/ca.crt"
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

	// Create the TLS credentials for server.
	certificate, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatalf("failed to load server certificates, err: %v", err)
	}

	// Create certificate pool.
	certPool := x509.NewCertPool()

	// Read CA certificate.
	caCert, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		log.Fatalf("failed to read CA certificate, err: %v", err)
	}
	ok := certPool.AppendCertsFromPEM(caCert)
	if !ok {
		log.Fatalf("failed to parse CA cert")
	}

	// Create tls config for server.
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}

	// cert and key files should securly be deployed to the server. File path should come from config.
	// creds, err := credentials.NewServerTLSFromFile("./../tls/server.crt", "./../tls/server.key")
	// if err != nil {
	// 	log.Fatalf("could not load TLS keys: %s", err)
	// }

	creds := credentials.NewTLS(tlsConfig)

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
- create a credentials object (called creds) from your certificate and key files;
- create a grpc.ServerOption array and placed your credentials object in it;
- Create grpc server with grpc.ServerOption created in above step;

Mutual TLS Auth
- In addition to what we did for adding TLS support we have to create
a certPool by reading CA's certificate.
- Create TLS config.
	ClientCAs: We will trust the client if it presents certificate signed by any of these CA.
	Certificates: Server will present these certificates to the clients.
*/
