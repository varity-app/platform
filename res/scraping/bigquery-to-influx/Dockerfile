# Build module
FROM golang:1.16.6-alpine AS builder

RUN apk add --no-cache build-base

# Copy source
COPY scraping /src

WORKDIR /src/cmd/bigquery2influx
RUN go mod download
RUN go build . 

# Prod module
FROM alpine:edge AS prod

# Copy binary
WORKDIR /app
COPY --from=builder /src/cmd/bigquery2influx/bigquery2influx .

ENTRYPOINT ./bigquery2influx
