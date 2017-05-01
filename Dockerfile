FROM golang:alpine
ARG pkg=github.com/Nordstrom/choices
RUN apk add --no-cache ca-certificates
COPY . $GOPATH/src/$pkg
RUN go install -v $(go list $pkg/... | grep -v /vendor/)
WORKDIR $GOPATH/src/$pkg
CMD echo "This the make file to build releases."; exit 1
