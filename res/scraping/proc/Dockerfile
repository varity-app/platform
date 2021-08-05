# Build module
FROM golang:1.16.6-alpine AS builder

# Install librdkafka
RUN apk add --no-cache librdkafka-dev build-base

# Copy source
WORKDIR /opt
COPY scraping/common/ ./common/
COPY scraping/protobuf/ ./protobuf/
COPY scraping/transforms/ ./transforms/
COPY scraping/services/proc/ ./services/proc/

WORKDIR /opt/services/proc
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -tags dynamic . 

# Prod module
FROM alpine:edge AS prod

RUN apk add --no-cache librdkafka

# Copy binary
WORKDIR /app
COPY --from=builder /opt/services/proc/proc .

ENTRYPOINT ./proc
