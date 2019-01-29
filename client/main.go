package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"git.ouroath.com/peng/test/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	// FQDN of the server.
	FQDN = "localhost"
)

func main() {
	addr := fmt.Sprintf("%s:%d", FQDN, 4444)
	// Create the client TLS credentials
	creds, err := credentials.NewClientTLSFromFile("./../tls/server.crt", "")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}
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
- you created a credentials object with the certificate file. Note that the client do not use the certificate key, the key is private to the server
- you added an option to the grpc.Dial() function, using your credentials object. Note that the grpc.Dial() function is also a variadic function, so it accepts any number of options;
*/
