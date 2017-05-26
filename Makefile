pkg := github.com/Nordstrom/choices
db := docker run -v ${CURDIR}:/go/src/${pkg}:ro -v ${CURDIR}/bin:/go/bin golang:alpine go install ${pkg}
image := quay.io/nordstrom/choices
tag := v0.0.9

all: docker/build

bin/elwin: *.go vendor/* cmd/elwin/*.go
	${db}/cmd/elwin

bin/mongo-store: *.go vendor/* cmd/mongo-store/*.go
	${db}/cmd/mongo-store

bin/houston: *.go vendor/* cmd/houston/*.go
	${db}/cmd/houston

docker/build: bin/elwin bin/houston bin/mongo-store
	docker build -t choices .

docker/tag: docker/build
	docker tag choices ${image}:${tag}

docker/push: docker/tag
	docker push ${image}:${tag}

