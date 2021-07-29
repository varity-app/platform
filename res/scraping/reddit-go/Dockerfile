# Build module
from golang:1.16.6-buster as builder

# Copy source
WORKDIR /opt
COPY go/scraping/common/ ./common/
COPY go/scraping/protobuf/ ./protobuf/
COPY go/scraping/reddit/ ./reddit/

WORKDIR /opt/reddit
RUN go mod download
RUN go build .

# Prod module
from debian:stable-slim as prod

# Copy binary
WORKDIR /app
COPY --from=builder /opt/reddit/reddit-scraper .

ENTRYPOINT ./reddit-scraper
