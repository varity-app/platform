# Install kafka
FROM golang:1.16.6-buster AS base

WORKDIR /root
RUN git clone https://github.com/edenhill/librdkafka.git
WORKDIR /root/librdkafka
RUN ./configure --prefix=/opt/librdkafka
RUN make -j4
RUN make install

# Build module
FROM base AS builder

# Copy source
WORKDIR /opt
COPY go/scraping/common/ ./common/
COPY go/scraping/protobuf/ ./protobuf/
COPY go/scraping/reddit/ ./reddit/

WORKDIR /opt/reddit
RUN go mod download
RUN go build .

# Prod module
FROM base AS prod

# Copy binary
WORKDIR /app
COPY --from=builder /opt/reddit/reddit-scraper .

ENTRYPOINT ./reddit-scraper
