# Nix Package Data

This is a small REST server, written in Golang, to retrieve nixkgs metadata. The server is started with:

```bash
go run main.go
```

The application depends on Gin to handle HTTP requests, it provides two endpoints:

* Simple health check endpoint: `GET /health`
* Get package information by name: `GET /package/:name`

The server returns HTTP status codes in case of an error, the port can be configured via `config.json` (defaults to 8080).

## Setup

The program requires devenv.sh to be installed. In case the GIN framework is not setup, execute the following command before the start, it will be installed in the project directory.

```bash
go get github.com/gin-gonic/gin
```

## Use

The default mode starts a web server, data is retrieved via curl or through a web browser

```bash
curl http://localhost:8080/package/nginx
```

### CLI Mode

Data can also be retrieved from the command line, using the CLI mode.

```bash
go run main.go cli nginx
```
