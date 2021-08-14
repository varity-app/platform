# Build module
FROM golang:1.16.6-alpine AS builder

# Copy source
WORKDIR /opt
COPY scraping/common/ ./common/
COPY scraping/protobuf/ ./protobuf/
COPY scraping/scrapers/ ./scrapers/
COPY scraping/data/kafka/ ./data/kafka/
COPY scraping/cmd/reddit-scraper/ ./cmd/reddit-scraper/

WORKDIR /opt/cmd/reddit-scraper
RUN go mod download
RUN go build . 

# Prod module
FROM alpine:edge AS prod

# Copy binary
WORKDIR /app
COPY --from=builder /opt/cmd/reddit-scraper/reddit-scraper .

ENTRYPOINT ./reddit-scraper
