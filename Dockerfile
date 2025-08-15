FROM golang:1.24-alpine3.21 AS builder
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY . .
RUN go build -o subaggr ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app .
ENTRYPOINT [ "/app/subaggr" ]