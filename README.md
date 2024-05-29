## Demo of gRPC-Gateway exposing gRPC through REST API + generation of OAS specification based on a proto file

This demo was prepared based on the example from https://grpc-ecosystem.github.io/grpc-gateway/docs/tutorials/introduction/.

### Prerequisites

#### 1. The following software need to be installed prior executing next steps
- golang installed - https://go.dev/doc/install
- protoc installed - https://grpc.io/docs/protoc-installation/

#### 2. To get this demo compiled, and any other grpc and grpc gateway go applications, following needs to be installed
```
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

and

#### 3. Install bufbuild 
bufbuild can be installed from https://github.com/bufbuild/buf
In case of macOS use homebrew as show below:
```
brew install bufbuild/buf/buf
```
---
### Build gRPC Gateway Demo
First, clone this repo and then use `buf` to generate server, client and gateway code with protoc plugins folling below steps.

#### 1. Generate code gRPC code
```
buf dep update
buf generate
```
#### 2. Compile gRPC server and gRPC gateway
```
go build -o grpcserver -ldflags="-w -s" ./server.go
go build -o grpcgateway -ldflags="-w -s" ./gateway.go
```

#### 3. Run demo
The demo consists of 2 apps:
- grpcserver running on the port 8080
- grpcgateway running on the port 8090

They can be run in the following way:

1st terminal
```
./grpcserver
```

2nd terminal
```
./grpcgateway
```

Then use curl or insomnia or postman to test it. Below example of curl:
```
curl -X POST -k http://localhost:8090/v1/example/echo -d '{"name": " hello"}'
```

---

### Generating OpenAPI Specification based on proto file
When it comes to OAS file generation based on proto file then here is the recipe. This is obvious, but just to make it clear. 
As in case of building/generating grpc gateway, proto files need to be updated with API annotations for services that are
to be exposed through grpc gateway. Here is an example of such annotation:
```go
service TimeCheck {
  // Sends a greeting
  rpc GiveTime (TimeRequest) returns (TimeReply) {
    option (google.api.http) = {
      post: "/api/v1/time"
      body: "*"
    };
  }
}
```
To generate an OpenAPI Specification (OAS) file from a Protocol Buffers (proto files) file using the `protoc` utility, 
follow these steps:

#### 1. Install `protoc-gen-openapiv2` Plugin

- Make sure you have Go installed.
- Install the plugin using the following command:
```sh
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```

#### 2. Generate OpenAPI Specification

Run the `protoc` command with the `--openapiv2_out` option to generate the OpenAPI specification in the root of the project:

```sh
protoc -I . --openapiv2_out . --openapiv2_opt logtostderr=true timeservice/time_service.proto
```

OAS file will be generated in `timeservice/time_service.swagger.json`
