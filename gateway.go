package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"

	ts "github.com/hhruszka/grpc-gateway-demo/proto/timeservice"
)

func main() {
	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.NewClient(
		"0.0.0.0:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwmux := runtime.NewServeMux()
	// Register Greeter
	err = ts.RegisterTimeCheckHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	//gwServer := &http.Server{
	//	Addr:    ":8090",
	//	Handler: gwmux,
	//}

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
