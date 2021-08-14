# Build module
FROM golang:1.16.6-alpine AS builder

RUN apk add --no-cache build-base

# Copy source
WORKDIR /opt
COPY scraping/common/ ./common/
COPY scraping/protobuf/ ./protobuf/
COPY scraping/transforms/ ./transforms/
COPY scraping/data/kafka/ ./data/kafka/
COPY scraping/cmd/proc/ ./cmd/proc/

WORKDIR /opt/cmd/proc
RUN go mod download
RUN go test . 
RUN go build . 

# Prod module
FROM alpine:edge AS prod

# Copy binary
WORKDIR /app
COPY --from=builder /opt/cmd/proc/proc .

ENTRYPOINT ./proc
