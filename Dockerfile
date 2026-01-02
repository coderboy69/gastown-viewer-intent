# Multi-stage build for minimal image
FROM golang:1.22-alpine AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o gvid ./cmd/gvid

# Final image - distroless for security
FROM gcr.io/distroless/static-debian12

COPY --from=builder /build/gvid /usr/local/bin/gvid

EXPOSE 7070

ENTRYPOINT ["/usr/local/bin/gvid"]
CMD ["--host", "0.0.0.0", "--port", "7070"]
