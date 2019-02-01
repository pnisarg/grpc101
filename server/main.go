package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	"git.ouroath.com/peng/test/grpc/api"
	"git.ouroath.com/peng/test/grpc/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
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

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. More: https://grpc.io/blog/coreos
func grpcHandlerFunc(ctx context.Context, grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r.WithContext(ctx))
		} else {
			otherHandler.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

// initGRPCServer initializes gRPC server with TLS credentials and registers
// proto services.
func initGRPCServer(certificate tls.Certificate, certPool *x509.CertPool) *grpc.Server {
	// create a server instance
	s := api.Server{}

	// Create tls config for gRPC server.
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		NextProtos:   []string{"h2", "http/1.1"},
		MinVersion:   tls.VersionTLS12,
	}
	creds := credentials.NewTLS(tlsConfig)

	// Create an array of gRPC options with the credentials
	opts := []grpc.ServerOption{grpc.Creds(creds)}
	grpcServer := grpc.NewServer(opts...)

	// attach the Ping service to the server
	proto.RegisterPingServer(grpcServer, &s)

	return grpcServer
}

func initHTTPServer(ctx context.Context, certificate tls.Certificate, certPool *x509.CertPool, grpcServer *grpc.Server, gwmux *runtime.ServeMux) *http.Server {
	// create a listener on TCP port
	addr := fmt.Sprintf("%s:%d", FQDN, port)

	mux := http.NewServeMux()
	// handler to check if service is up
	mux.HandleFunc("/ruok", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(w, "imok")
	})
	mux.Handle("/", gwmux)

	// Create tls config for HTTP server.
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
		ClientAuth:   tls.VerifyClientCertIfGiven,
		NextProtos:   []string{"h2", "http/1.1"},
		MinVersion:   tls.VersionTLS12,
	}
	srv := &http.Server{
		Addr:      addr,
		Handler:   grpcHandlerFunc(ctx, grpcServer, mux),
		TLSConfig: tlsConfig,
	}
	return srv
}

// main start a gRPC server and waits for connection
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	grpcServer := initGRPCServer(certificate, certPool)

	// setup gateway
	gwmux := runtime.NewServeMux()

	grpcAddr := fmt.Sprintf("%s:%d", "localhost", port)

	// create TLS config for gateway
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
		ServerName:   FQDN,
	}

	creds := credentials.NewTLS(tlsConfig)
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	if err := proto.RegisterPingHandlerFromEndpoint(ctx, gwmux, grpcAddr, dopts); err != nil {
		log.Fatalf("failed to register ping handler endpoint, err: %v", err)
	}

	server := initHTTPServer(ctx, certificate, certPool, grpcServer, gwmux)

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("starting server on %s", server.Addr)
	if err := server.Serve(tls.NewListener(listener, server.TLSConfig)); err != nil {
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

Note: Once you have TLS credentials specified, it is important to start
HTTPS server using `tls.NewListener`. If you just use http.Server.Serve(listerner)
clients will not be able to make TLS connection as your server doesn't know to server
TLS.
*/
