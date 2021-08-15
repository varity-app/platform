# Build module
FROM golang:1.16.6-alpine AS builder

# Copy source
WORKDIR /opt
COPY scraping/ ./

WORKDIR /opt/cmd/tiingo-scraper
RUN go mod download
RUN go build . 

# Prod module
FROM alpine:edge AS prod

# Copy binary
WORKDIR /app
COPY --from=builder /opt/cmd/tiingo-scraper/tiingo .

ENTRYPOINT ./tiingo
