# Build module
FROM golang:1.16.6-alpine AS builder

# Install librdkafka
RUN apk add --no-cache librdkafka-dev build-base

# Copy source
WORKDIR /opt
COPY scraping/common/ ./common/
COPY scraping/protobuf/ ./protobuf/
COPY scraping/transforms/ ./transforms/
COPY scraping/cmd/proc/ ./cmd/proc/

WORKDIR /opt/cmd/proc
RUN go mod download
RUN CGO_ENABLED=1 go test . -tags dynamic
RUN CGO_ENABLED=1 go build -tags dynamic . 

# Prod module
FROM alpine:edge AS prod

RUN apk add --no-cache librdkafka

# Copy binary
WORKDIR /app
COPY --from=builder /opt/cmd/proc/proc .

ENTRYPOINT ./proc
