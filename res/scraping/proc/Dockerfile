# Build module
from golang:1.16.6-buster as builder

# Copy source
WORKDIR /opt
COPY go/scraping/common/ ./common/
COPY go/scraping/protobuf/ ./protobuf/
COPY go/scraping/proc/ ./proc/

WORKDIR /opt/proc
RUN go mod download
RUN go build .

# Prod module
from debian:stable-slim as prod

# Copy binary
WORKDIR /app
COPY --from=builder /opt/proc/proc .

ENTRYPOINT ./proc
