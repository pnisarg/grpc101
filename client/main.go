package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"git.ouroath.com/peng/test/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	// FQDN of the server.
	FQDN           = "localhost"
	clientCertPath = "./../tls/client.crt"
	clientKeyPath  = "./../tls/client.key"
	caCertPath     = "./../tls/ca.crt"
)

func main() {
	addr := fmt.Sprintf("%s:%d", FQDN, 4444)

	// Create the client TLS credentials.
	certificate, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		log.Fatalf("failed to load client certificates, err: %v", err)
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
		RootCAs:      certPool,
		ServerName:   FQDN,
	}

	creds := credentials.NewTLS(tlsConfig)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	c := proto.NewPingClient(conn)

	response, err := c.SayHello(ctx, &proto.PingMessage{Greeting: "ping"})
	if err != nil {
		log.Fatalf("error when calling SayHello: %s", err)
	}

	log.Printf("response from server: %s", response.Greeting)
}

/*
NOTES:

- function instantiates a client connection, on the TCP port the server is bound to;
- note the defer call to properly close the connection when the function returns;
- the c variable is a client for the the Ping service

Enabling TLS connection
- Create a credentials object with the CA's certificate file. Note that the client doesn't have CAs private key.
- Use grpc.Dial() to make TCP connection with the server. Provide credentials created in above step to grpc.Dial func.

Enabling Mutual TLS
- To enable mutual TLS on client we have to do all the things we did on server side to enable mutual TLS with
the following difference.
	- We will load client certificate and key. This certificate will be presented to the server while making connection.
	- TLS Config for client is bit different.
		- ServerName: FQDN of the server
		- RootCAs: We will trust certificates presented to us by server if it is signed by one of these RootCAs.
		Note: We use root CAs field instead of ClientCAs here for obvious reasons.
*/
