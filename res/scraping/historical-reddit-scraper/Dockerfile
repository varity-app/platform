# Build go module
FROM golang:1.16.6-alpine AS builder

# Build go module
COPY scraping /src

WORKDIR /src/cmd/historical-reddit-scraper
RUN go mod download
RUN go build .

# Build prod image
FROM alpine:edge AS prod

# Copy binaries from builder
WORKDIR /app
COPY --from=builder /src/cmd/historical-reddit-scraper/historical-reddit-scraper .

ENTRYPOINT /app/historical-reddit-scraper