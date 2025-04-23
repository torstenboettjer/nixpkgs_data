# Nix Package Metadata

This is a small REST server, written in Golang, that retrieves nixkgs meta data.

## Mode of Operation:

The program supports two modes of operation:

* CLI Mode: Run with go run main.go cli <package-name> for command-line usage
* API Mode: Default mode starts a web server

The server uses Gin to handle HTTP requests and provides two endpoints:

* Simple health check endpoint: `GET /health`
* Get package information by name: `GET /package/:name`

The server returns HTTP status codes in case of an error, the port can be configured via PORT environment variable (defaults to 8080).

## Start the Server

**Step 1:** The GIN framework is required, in case it is not installed, execute the following command:

```bash
go get github.com/gin-gonic/gin
```

**Step 2** Started the server with

```bash
go run main.go
```

#### HTTP Request

```bash
curl http://localhost:8080/package/nginx
```

#### CLI Request

```bash
go run main.go cli nginx
```

### Open issues
* Stop the server with CTRL-C - add graceful shutdown routine
* Set port through config file
