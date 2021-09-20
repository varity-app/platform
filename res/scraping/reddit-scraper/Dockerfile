# Build module
FROM golang:1.16.6-alpine AS builder

# Copy source
COPY scraping /src

WORKDIR /src/cmd/reddit-scraper
RUN go mod download
RUN go build . 

# Prod module
FROM alpine:edge AS prod

# Copy binary
WORKDIR /app
COPY --from=builder /src/cmd/reddit-scraper/reddit-scraper .

ENTRYPOINT ./reddit-scraper
