# Build go module
FROM golang:1.16.6-alpine AS builder

# Install librdkafka
RUN apk add --no-cache librdkafka-dev build-base wget ca-certificates

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
RUN CGO_ENABLED=1 go test -tags dynamic .
RUN CGO_ENABLED=1 go build -tags dynamic .

# Build prod image
FROM alpine:edge AS prod

RUN apk add --no-cache librdkafka-dev

# Copy binaries from builder
WORKDIR /app
COPY --from=builder /app/cmd/historical-reddit-scraper/historical-reddit-scraper .

ENTRYPOINT /app/historical-reddit-scraper