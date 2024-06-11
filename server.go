package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc/credentials"
	"log"
	"math/rand"
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

func (s *server) GiveDate(ctx context.Context, in *ts.Empty) (*ts.DateReply, error) {
	year, month, day := time.Now().Date()
	log.Printf("Received request for GiveDate() method. gRPC request: %s", in.String())
	message := fmt.Sprintf("Hi %s!\nCurrent date is %02d:%02d:%02d", day, month, year)
	return &ts.DateReply{Message: message}, nil
}

var (
	quotes []string = []string{
		"Your work is going to fill a large part of your life, and the only way to be truly satisfied is to do what you believe is great work. And the only way to do great work is to love what you do.",
		"Have the courage to follow your heart and intuition. They somehow already know what you truly want to become. Everything else is secondary.",
		"Remembering that you are going to die is the best way I know to avoid the trap of thinking you have something to lose. You are already naked. There is no reason not to follow your heart.",
		"Innovation distinguishes between a leader and a follower.",
		"Sometimes when you innovate, you make mistakes. It is best to admit them quickly, and get on with improving your other innovations.",
		"Your time is limited, so don’t waste it living someone else’s life.",
		"Design is not just what it looks like and feels like. Design is how it works.",
		"The people who are crazy enough to think they can change the world are the ones who do.",
		"Stay hungry, stay foolish.",
	}
)

func (s *server) GiveQuote(ctx context.Context, in *ts.Empty) (*ts.QuoteReply, error) {
	quoteIdx := rand.Int() % len(quotes)
	log.Printf("Received request for GiveQuote() method. gRPC request: %s", in.String())
	message := fmt.Sprintf("Hi\nHere is quote for you:  %s", quotes[quoteIdx])
	return &ts.QuoteReply{Message: message}, nil
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
