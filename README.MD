# TrainReservationSysytem
# gRPC Train Ticket Service

This project is a gRPC server implemented in Go for managing train ticket purchases, receipts, and seat allocations.

## Setup and Installation

### Prerequisites

- Go (https://golang.org/dl/)
- Protocol Buffer Compiler (protoc) (https://github.com/protocolbuffers/protobuf/releases)
- gRPC and Protocol Buffers Go Plugin

### Installing Go Packages

```sh
go get -u google.golang.org/grpc
go get -u github.com/golang/protobuf/protoc-gen-go
```

### Generating Go Code from Proto File

Run the following command in the Command Prompt to generate the Go code from the proto file:

```sh
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/train.proto
```

## Running the Server

Ensure you have the generated Go code (`train.pb.go`) in the same directory as `server.go`. Adjust the import path `pb "path/to/your/proto"` to the actual path where your generated proto files are located.

Run the server using:

```sh
go run server/main.go
```

## Testing with Postman

Follow the steps to import the `train.proto` file into Postman and create gRPC requests to test the server.

## License

This project is licensed under the MIT License.