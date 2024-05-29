package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	ts "github.com/hhruszka/grpc-gateway-demo/proto/timeservice"
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

func main() {
	// Create a listener on TCP port
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	// Create a gRPC server object
	s := grpc.NewServer()
	// Attach the Greeter service to the server
	ts.RegisterTimeCheckServer(s, &server{})
	// Serve gRPC server
	log.Println("Serving gRPC on 0.0.0.0:8080")
	log.Fatalln(s.Serve(lis))
}
