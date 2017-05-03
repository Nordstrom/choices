FROM golang:alpine
ARG pkg=github.com/Nordstrom/choices
RUN apk add --no-cache ca-certificates
RUN set -ex \
	&& apk add --no-cache --virtual .build-deps \
	git \
	&& mkdir -p $GOPATH/src/github.com/Nordstrom \
	&& cd $GOPATH/src/github.com/Nordstrom \
	&& git clone https://github.com/Nordstrom/choices.git \
	&& apk del .build-deps
COPY . $GOPATH/src/$pkg
RUN go install -v $(go list $pkg/... | grep -v /vendor/)
WORKDIR $GOPATH/src/$pkg
CMD echo "This the make file to build releases."; exit 1
