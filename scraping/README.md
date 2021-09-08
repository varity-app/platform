# Scraping Go Module

Most custom scraping infrastructure is contained within this Golang module.

## Project Structure

### `api/`
Protobuf, JSON, and other API resource definitions.

Compile and lint protobufs with the `buf` CLI.  See more instructions on the [buf website](https://docs.buf.build/introduction).

### `cmd/`
Application entrypoints and executables.

### `internal/`
Private application and library code.  This is code that shouldn't be seen by other files.