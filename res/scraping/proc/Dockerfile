# Build module
FROM golang:1.16.6-alpine AS builder

# Install librdkafka
RUN apk add --no-cache librdkafka-dev build-base

# Copy source
WORKDIR /opt
COPY go/scraping/common/ ./common/
COPY go/scraping/protobuf/ ./protobuf/
COPY go/scraping/proc/ ./proc/

WORKDIR /opt/proc
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -tags dynamic . 

# Prod module
FROM alpine:edge AS prod

RUN apk add --no-cache librdkafka

# Copy binary
WORKDIR /app
COPY --from=builder /opt/proc/proc .

ENTRYPOINT ./proc
