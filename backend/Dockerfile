FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a \
    -installsuffix cgo \
    -ldflags='-w -s' \
    -o server \
    ./cmd/server

FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /app/server .
RUN chown -R appuser:appgroup /app
USER appuser
EXPOSE 8080
ENV GIN_MODE=release
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1
CMD ["./server"]
