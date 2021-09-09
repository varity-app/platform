# Build module
FROM golang:1.16.6-alpine AS builder

# Copy source
COPY scraping /src

WORKDIR /src/cmd/tiingo-scraper
RUN go mod download
RUN go build . 

# Prod module
FROM alpine:edge AS prod

# Copy binary
WORKDIR /app
COPY --from=builder /src/cmd/tiingo-scraper/tiingo-scraper .

ENTRYPOINT ./tiingo-scraper
