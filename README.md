# pong-go

A minimalistic HTTP echo server written in Go, designed for testing and debugging HTTP requests. It can log request details and respond with either a simple `200 OK` or a structured JSON response.

## Features

- Logs HTTP request method, URL, headers, and body.
- Configurable logging via environment variables.
- Two response modes:
  - **Simple Mode**: always responds with `200 OK`.
  - **Verbose Mode**: returns a JSON payload with request details and timestamp.
- Optional random response delay controlled via CLI flag.

## Environment Variables

The behavior of the server can be customized using the following environment variables:

| Variable Name             | Description                                     | Default   |
|---------------------------|-------------------------------------------------|-----------|
| `PONG_ECHO_SERVER_ADDR`   | Address and port to bind the server to          | `:8080`   |
| `PONG_LOG_METHOD_URL`     | Log the HTTP method and URL                     | `true`    |
| `PONG_LOG_HEADERS`        | Log request headers                             | `false`   |
| `PONG_LOG_BODY`           | Log request body                                | `false`   |
| `PONG_VERBOSE_RESPONSE`   | Enable verbose JSON response mode               | `false`   |

## How to build

To embed a git-describe style version string (tag, nearest tag + short commit, and a `-dirty` suffix for local changes), build from a git checkout and pass the derived value via `-ldflags`.

```
VERSION=$(git describe --tags --dirty --always --abbrev=7 2>/dev/null || echo dev)
```

### faster build
```
go build -ldflags="-X main.version=${VERSION}" -o pong-go
```

### smaller binary
```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -X main.version=${VERSION}" -o pong-go
```

## How to use

### check version
Use `-v` or `--version` to print the embedded version string and exit.
```
./pong-go -v
```

### run server

```
export PONG_VERBOSE_RESPONSE=true
./pong-go
```

### add random delay before responses

Add a random delay between `0` and the provided duration (e.g. `5000ms`, `5s`) to every request.
```
./pong-go --random-delay 5000ms
```

### send requests
```
curl http://localhost:8080
```

```
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello, world!"}'
```

### Example output
```
$ curl -s -X POST http://localhost:8080   -H "Content-Type: application/json"   -d '{"message": "Hello, world!"}' | jq .
{
  "timestamp": "2025-05-13T13:26:00.296+02:00",
  "method": "POST",
  "url": "/",
  "headers": {
    "Accept": [
      "*/*"
    ],
    "Content-Length": [
      "28"
    ],
    "Content-Type": [
      "application/json"
    ],
    "User-Agent": [
      "curl/7.81.0"
    ]
  },
  "body": "{\"message\": \"Hello, world!\"}"
}
```
