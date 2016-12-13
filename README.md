[![GoDoc](https://godoc.org/github.com/Nordstrom/choices?status.svg)](https://godoc.org/github.com/Nordstrom/choices)

# choices
A way to choose things.

## Building

> In order to build you need `go` installed.

> You will need go1.7.1 if you are deploying to kubernetes.

To build choices library and all the included binaries run the following in a
terminal.

```bash
go get -u github.com/Nordstrom/choices/...
```

If you are only interested in the library you can run the following in a
terminal.

```bash
go get -u github.com/Nordstrom/choices
```

## Running locally

There are two main components of choices. The storage server (mongo-store) and
the frontend experiment server (elwin). You should first start the storage
server. You will need a local version of mongo running.  [Mongo
Installation](https://docs.mongodb.com/manual/installation/). If you are using
docker you can run the following docker cmd.

```bash
docker run -d --name mongo-storage -p 27017:27017 mongo
```

With mongo running you can start the storage server.

```bash
# from the choice/cmd/mongo-store directory run
./mongo-store
```

> This creates the storage service on port 8080

Now that the storage server is up you can start the experiment server (elwin).

```bash
# from the choices/cmd/elwin directory run
JSON_ADDRESS=":8181" GRPC_ADDRESS=":8282" MONGO_ADDRESS="localhost:8080" MONGO_DATABASE="elwin" MONGO_COLLECTION="staging" ./elwin
```

The experiment store loads some example data into the *staging* environment.
You should be able to open a web browser and navigate to
`localhost:8181?label=test&userid=1234` and see some results. You can also run
the following command from a terminal if you have `curl` installed.

```bash
curl "localhost:8181?label=test&userid=1234"
```
