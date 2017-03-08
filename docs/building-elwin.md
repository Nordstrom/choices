# Compiling and Running Elwin from Source

# 1. Overview

Elwin is an executable that can perform a/b/c and multivariate
testing. The main library that performs the majority of the work is
called `choices`. It is written in go. Elwin contains the server logic
and evaluates experiments based on the requests it receives. There are
two binaries that make up a basic elwin deployment.

# 2. Compiling

## 2.1 Go compiler

In order to compile an Elwin executable you will need a go compiler. I
use the most up to date compiler version, currently go1.8. You can
download go at [https://golang.org/dl/](https://golang.org/dl/).

## 2.2 Getting source code

The first step will be to clone the elwin repo locally. Go requires
files to be placed in certain directories based on your `GOPATH`
environment variable. If you are using the default `GOPATH` in go1.8
It will be `$HOME/go`. There are two options for getting the elwin
code. One option is to use Go's built in `get` command. The second
option is to manually download and create the necessary directory
structure.

### 2.2.1 Using `go get`

In your terminal run the following command.

```bash
go get github.com/Nordstrom/choices
```

> This requires `go` and `git` to be installed

You should find the elwin files at
`$GOPATH/src/github.com/Nordstrom/choices` if you have `GOPATH` set or
`$HOME/go/src/github.com/Nordtsrom/choices` if you are using the
default `GOPATH` in go1.8.

### 2.2.2 Manual `git clone`

You should check if your `GOPATH` is set.

```bash
echo $GOPATH
# If nothing is displayed run the following commands.
mkdir $HOME/go
export GOPATH=$HOME/go
```

> If you set the `GOPATH` it will only be set for this terminal. If
> you close the terminal or start a new session it will need to be set
> again.

Next you need to create the directory structure that go expects to
hold the files. Then you can clone the repo.

```bash
mkdir -p $GOPATH/src/github.com/Nordstrom/
cd $GOPATH/src/github.com/Nordstrom/
git clone https://github.com/Nordstrom/choices
```

## 2.3 Compiling source components

### 2.3.1 Compiling Elwin binary

Assuming you have the downloaded the code and set your `GOPATH`, now
you can compile the elwin executable. To compile an executable that
can run on your local machine you could run the following.

```bash
cd $GOPATH/src/github.com/Nordstrom/choices/cmd/elwin
go build
```
 
### 2.3.2 Compiling Storage binary

The most up to date storage binary is the bolt-store implementation.
You compile this in a similar way to elwin.

```bash
cd $GOPATH/src/github.com/Nordstrom/choices/cmd/bolt-store
go build
```

# 3. Running Elwin

> TODO: use `spf13/viper` for configuration

# 3.1 Running Elwin Locally

To run elwin locally you will need to supply some configuration values so the ports don't collide.

```bash
cd $GOPATH/src/github.com/Nordstrom/choices/cmd/bolt-store
./bolt-store
```

In a separate terminal, run the following.

```bash
cd $GOPATH/src/github.com/Nordstrom/choices/cmd/elwin
JSON_ADDRESS=:8082 GRPC_ADDRESS=:8083 MONGO_ADDRESS=:8080 ./elwin
```
