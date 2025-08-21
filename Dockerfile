FROM golang:1.24-alpine3.21 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Disable CGO for a static binary and strip symbols
ENV CGO_ENABLED=0
RUN go build -ldflags "-s -w" -tags netgo -o subaggr ./cmd/main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=builder /app/subaggr /app/subaggr
USER nonroot:nonroot
ENTRYPOINT ["/app/subaggr"]