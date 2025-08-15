# Go Application Deployment with Docker

This guide explains how to build and deploy a Go application using Docker multi-stage builds, with environment variables loaded from a `.env` file.

---

## Prerequisites

- [Docker](https://www.docker.com/) installed
- Go project with entry point in `cmd/`
- `.env` file containing environment variables, for example:


---

## Dockerfile Explanation

1. **Build Stage**: Uses `golang:1.24-alpine` to build a statically linked Go binary.
2. **Run Stage**: Uses lightweight `alpine` image, copies the binary, and runs the app.
3. **Environment Variable**: Port is read from `.env` using Go's `os.Getenv`.

---

## Build Docker Image

From project root:

```bash
$ docker build -t myapp:latest .
$ docker run --env-file .env -p ${PORT}:${PORT} myapp:latest
```

or

```bash
$ docker-compose up --build
```

## INSTALL AIR TOML

``$ go install github.com/air-verse/air@latest
``