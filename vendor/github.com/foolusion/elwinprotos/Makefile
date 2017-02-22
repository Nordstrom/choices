all: build/elwin build/storage

build/elwin: elwin/elwin.pb.go elwin/elwin.pb.gw.go elwin/elwin.swagger.json

build/storage: storage/storage.pb.go storage/storage.pb.gw.go storage/storage.swagger.json

elwin/elwin.pb.go: protos/elwin.proto
	protoc -I./protos -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --gogoslick_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:elwin protos/elwin.proto

elwin/elwin.pb.gw.go: protos/elwin.proto
	protoc -I./protos -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:elwin protos/elwin.proto

elwin/elwin.swagger.json: protos/elwin.proto
	protoc -I./protos -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:elwin protos/elwin.proto

storage/storage.pb.go: protos/storage.proto
	protoc -I./protos -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --gogoslick_out=Mgoogle/api/annotations.proto=github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api,plugins=grpc:storage protos/storage.proto

storage/storage.pb.gw.go: protos/storage.proto
	protoc -I./protos -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out=logtostderr=true:storage protos/storage.proto

storage/storage.swagger.json: protos/storage.proto
	protoc -I./protos -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --swagger_out=logtostderr=true:storage protos/storage.proto
