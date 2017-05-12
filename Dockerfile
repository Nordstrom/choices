FROM alpine:3.5
RUN apk add --no-cache ca-certificates
COPY bin /bin
CMD echo "Welcome to choices. Choose a binary to run."; exit 1