package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net/http"
	"os"

	ts "github.com/hhruszka/grpc-gateway-demo/proto/timeservice"
)

var (
	clientCertFile string
	clientKeyFile  string
	rootCaCertFile string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	clientKeyFile = os.Getenv("TLS_CLIENT_KEY")
	clientCertFile = os.Getenv("TLS_CLIENT_CERT")
	rootCaCertFile = os.Getenv("CA_CERT")
}

func main() {
	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests

	clientCert, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
	if err != nil {
		log.Fatalf("Failed to load client certificate %v", err)
	}

	caCert, err := os.ReadFile(rootCaCertFile)
	if err != nil {
		log.Fatalf("Failed to load CA certificate %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	creds := credentials.NewTLS(&tls.Config{
		RootCAs:      caCertPool,
		Certificates: []tls.Certificate{clientCert},
	})

	conn, err := grpc.NewClient(
		"0.0.0.0:8080",
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()

	// Register Server
	err = ts.RegisterTimeCheckHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8090")
	//log.Fatalln(gwServer.ListenAndServe())
	log.Fatalln(http.ListenAndServe(":8090", loggingMiddleware(gwmux)))
}

// loggingMiddleware is an example middleware that logs the request method and URL
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: from %s to %s with method %s for endpoint %s", r.RemoteAddr, r.Host, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
