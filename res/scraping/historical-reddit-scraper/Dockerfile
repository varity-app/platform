# Build go module
FROM golang:1.16.6-alpine AS builder

# Build go module
WORKDIR /app
COPY scraping/common/ ./common/
COPY scraping/protobuf/ ./protobuf/
COPY scraping/transforms/ ./transforms/
COPY scraping/scrapers/ ./scrapers/
COPY scraping/data/kafka/ ./data/kafka/
COPY scraping/cmd/historical-reddit-scraper/ ./cmd/historical-reddit-scraper/

WORKDIR /app/cmd/historical-reddit-scraper
RUN go mod download
RUN go build .

# Build prod image
FROM alpine:edge AS prod

# Copy binaries from builder
WORKDIR /app
COPY --from=builder /app/cmd/historical-reddit-scraper/historical-reddit-scraper .

ENTRYPOINT /app/historical-reddit-scraper