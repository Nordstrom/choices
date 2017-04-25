all: build/elwin build/storage build/intake

build/elwin: elwin/elwin.pb.go elwin/elwin.pb.gw.go elwin/elwin.swagger.json

build/storage: storage/storage.pb.go storage/storage.pb.gw.go storage/storage.swagger.json

build/intake: intake/intake.pb.go intake/intake.pb.gw.go intake/intake.swagger.json

elwin/elwin.pb.go: elwin/elwin.proto
	protoc -I./elwin -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I${GOPATH}/src --gogoslick_out=plugins=grpc:${GOPATH}/src elwin/elwin.proto

elwin/elwin.pb.gw.go: elwin/elwin.proto
	protoc -I./elwin -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I${GOPATH}/src --grpc-gateway_out=logtostderr=true:elwin elwin/elwin.proto

elwin/elwin.swagger.json: elwin/elwin.proto
	protoc -I./elwin -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I${GOPATH}/src --swagger_out=logtostderr=true:elwin elwin/elwin.proto

storage/storage.pb.go: storage/storage.proto
	protoc -I./storage -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I${GOPATH}/src --gogoslick_out=plugins=grpc:${GOPATH}/src storage/storage.proto

storage/storage.pb.gw.go: storage/storage.proto
	protoc -I./storage -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I${GOPATH}/src --grpc-gateway_out=logtostderr=true:storage storage/storage.proto

storage/storage.swagger.json: storage/storage.proto
	protoc -I./storage -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I${GOPATH}/src --swagger_out=logtostderr=true:storage storage/storage.proto

intake/intake.pb.go: intake/intake.proto storage/storage.proto
	protoc -I./intake -I./storage -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I${GOPATH}/src --gogoslick_out=plugins=grpc:${GOPATH}/src intake/intake.proto

intake/intake.pb.gw.go: intake/intake.proto storage/storage.proto
	protoc -I./intake -I./storage -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I${GOPATH}/src --grpc-gateway_out=logtostderr=true:intake intake/intake.proto

intake/intake.swagger.json: intake/intake.proto storage/storage.proto
	protoc -I./intake -I./storage -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I${GOPATH}/src --swagger_out=logtostderr=true:intake intake/intake.proto
