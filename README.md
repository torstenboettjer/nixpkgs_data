# Nixpkgs Search

This is a small REST server, written in Golang, that retrieves nixkgs meta data.

## Mode of Operation:

The program supports two modes of operation:

* CLI Mode: Run with go run main.go cli <package-name> for command-line usage
* API Mode: Default mode starts a web server

The server now uses Gin to handle HTTP requests and provides two endpoints:

* Simple health check endpoint: `GET /health`
* Get package information by name: `GET /package/:name`

The server returns HTTP status codes in case of an error, the port can be configured via PORT environment variable (defaults to 8080).

## Start the Server

The server depends on the GIN framework. In case it is not installed, the following command is required before the start.

```bash
go get github.com/gin-gonic/gin
```

When GIN is installed it is started with:

```bash
go run main.go
```

### HTTP Request

```bash
curl http://localhost:8080/package/nginx
```

### CLI Request

```bash
go run main.go cli nginx
```
