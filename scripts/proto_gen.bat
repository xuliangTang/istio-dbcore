protoc320 --proto_path=protos --go_out=./ dbrequest.proto
protoc320 --proto_path=protos --go-grpc_out=./ dbservice.proto

