pkg := github.com/Nordstrom/choices
db := docker run -v ${CURDIR}:/go/src/${pkg}:ro -v ${CURDIR}/bin:/go/bin golang:alpine go install ${pkg}
image := quay.io/nordstrom/choices
tag := jim-pprof

all: docker/build

bin/elwin: *.go vendor/* cmd/elwin/*.go
	${db}/cmd/elwin

bin/bolt-store: *.go vendor/* cmd/bolt-store/*.go
	${db}/cmd/bolt-store

bin/elwin-grpc-gateway: *.go vendor/* cmd/elwin-grpc-gateway/*.go
	${db}/cmd/elwin-grpc-gateway

bin/houston: *.go vendor/* cmd/houston/*.go
	${db}/cmd/houston

bin/json-gateway: *.go vendor/* cmd/json-gateway/*.go
	${db}/cmd/json-gateway

bin/mongo-store: *.go vendor/* cmd/mongo-store/*.go
	${db}/cmd/mongo-store

docker/build: bin/elwin bin/bolt-store bin/elwin-grpc-gateway bin/houston bin/json-gateway bin/mongo-store
	docker build -t choices .

docker/tag: docker/build
	docker tag choices ${image}:${tag}

docker/push: docker/tag
	docker push ${image}:${tag}

