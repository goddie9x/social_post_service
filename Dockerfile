# syntax=docker/dockerfile:1
# ---- builder: generate protobuf stubs + compile ----
FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git protoc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.2 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
ENV PATH=$PATH:/root/go/bin
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Regenerate gRPC stubs from proto/ (keeps them in sync with the .proto source)
RUN protoc -I proto \
    --go_out=internal/grpc/post_service/proto --go_opt=paths=source_relative \
    --go-grpc_out=internal/grpc/post_service/proto --go-grpc_opt=paths=source_relative \
    proto/post.proto
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/server ./cmd/server

# ---- runtime: minimal, non-root ----
FROM alpine:latest
RUN apk add --no-cache ca-certificates && \
    addgroup -S app && adduser -S app -G app
WORKDIR /app
COPY --from=builder /app/server .
USER app
EXPOSE 3005
CMD ["./server"]
