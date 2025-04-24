# Nix Package Metadata Retrieval

This is a small REST server, written in Golang, to retrieve nixkgs metadata. The program supports two modes of operation:

* API Mode: Default mode starts a web server
* CLI Mode: At the command-line using `go run main.go`

The API mode depends on Gin to handle HTTP requests, it provides two endpoints:

* Simple health check endpoint: `GET /health`
* Get package information by name: `GET /package/:name`

The server returns HTTP status codes in case of an error, the port can be configured via PORT environment variable (defaults to 8080).

## Using the Server

In case the GIN framework is not installed, execute the following command before the start:

```bash
go get github.com/gin-gonic/gin
```

#### Start

```bash
go run main.go
```

#### API Mode

```bash
curl http://localhost:8080/package/nginx
```

#### CLI Mode

```bash
go run main.go cli nginx
```

## Open issues
* Stop the server with CTRL-C - add graceful shutdown routine
* Set port through config file
