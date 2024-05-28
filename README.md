Hi Ben

I followed https://grpc-ecosystem.github.io/grpc-gateway/docs/tutorials/introduction/ to create this app.
I assume that you have golang installed on your mac. If not, homebrew can install it. 
To get this app compiled you will need to install:
```
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

and

```
brew install bufbuild/buf/buf
```

I used `buf` to deal with protobuf.

```
buf dep update
buf generate
```

and then

```
go build -ldflags="-w -s"
./protobufapitest
```
The app contains 3 elements:
- grpc server (8080)
- grpc client
- grpc-gateway (8090)

You can use curl to test it:
```
curl -X POST -k http://localhost:8090/v1/example/echo -d '{"name": " hello"}'
```

When it comes to OAS file generation out of protobuf then here is the recepie. However, protobuf needs to have annotations with API endpoints added to leverage grpc-gateway.

---

To generate an OpenAPI Specification (OAS) file from a Protocol Buffers (protobuf) file using the `protoc` utility, you need to follow these steps:

### 1. Install `protoc-gen-openapiv2` Plugin

- Make sure you have Go installed.
- Install the plugin using the following command:
  ```sh
  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
  ```

### 2. Add annotations.proto and http.proto to the project - NOT NEEDED ANYMORE LEFT FOR REFERENCE

```
# in the root of the project
mkdir -p google/api
wget https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto  -O google/api/annotations.proto
wget https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto -O google/api/http.proto
```
### 3. Generate OpenAPI Specification

Run the `protoc` command with the `--openapiv2_out` option to generate the OpenAPI specification in the root of the project:

```sh
protoc -I . \                                                                                                                    
  --openapiv2_out . \
  --openapiv2_opt logtostderr=true \
  timeservice/time_service.proto
```

OAS file will be generated in `helloworld/hello_world.swagger.json`
