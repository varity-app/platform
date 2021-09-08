# Build module
FROM golang:1.16.6-alpine AS builder

RUN apk add --no-cache build-base

# Copy source
COPY scraping /src

WORKDIR /src/cmd/proc
RUN go mod download
RUN go test . 
RUN go build . 

# Prod module
FROM alpine:edge AS prod

# Copy binary
WORKDIR /app
COPY --from=builder /src/cmd/proc/proc .

ENTRYPOINT ./proc
