# =========================
# Stage 1: Build
# =========================
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /app/main ./cmd/

# =========================
# Stage 2: Run
# =========================
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/main .

# Copy .env (opsional jika ingin runtime Go membaca)
COPY .env .env

# Port akan diatur lewat env saat run
# CMD agar Go app bisa membaca PORT dari env
CMD ["/app/main"]
