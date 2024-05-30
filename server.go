package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"
	"time"

	ts "github.com/hhruszka/grpc-gateway-demo/proto/timeservice"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type server struct {
	ts.UnimplementedTimeCheckServer
}

func NewServer() *server {
	return &server{}
}

func (s *server) GiveTime(ctx context.Context, in *ts.TimeRequest) (*ts.TimeReply, error) {
	hour, minute, second := time.Now().Clock()
	log.Printf("Received request for GiveTime() method. gRPC request: %s", in.String())
	message := fmt.Sprintf("Hi %s!\nCurrent time is %02d:%02d:%02d", in.Name, hour, minute, second)
	return &ts.TimeReply{Message: message}, nil
}

var (
	serverCertFile string
	serverKeyFile  string
	caCertFile     string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}
	serverKeyFile = os.Getenv("TLS_SERVER_KEY")
	serverCertFile = os.Getenv("TLS_SERVER_CERT")
	caCertFile = os.Getenv("CA_CERT")
}

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	serverCert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
	if err != nil {
		log.Fatalf("Failed to load server certificate %v", err)
	}

	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		log.Fatalf("Failed to load CA certificate %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})

	// Create a gRPC server object
	s := grpc.NewServer(grpc.Creds(creds))
	// Attach the Greeter service to the server
	ts.RegisterTimeCheckServer(s, &server{})
	// Serve gRPC server
	log.Println("Serving gRPC on 0.0.0.0:8080")
	log.Fatalln(s.Serve(lis))
}
