# Mia Nubo
Mia Nubo → [_My cloud_](https://translate.google.com/#eo/en/Mia%20Nubo)

A cloud solution written in Go

Under heavy development. Will eventually support user permissions, metadata searching, multimedia streaming, and web-based clients.

## Compilation
1. Ensure you have [Go](https://golang.org/dl/) and [dep](https://golang.github.io/dep/docs/installation.html) installed
2. Run `dep ensure` to install dependencies
3. Run `go build main.go`

## Usage
1. Run the program `./main` (or `go run main.go`)
2. Connect to the program's webserver: `localhost:41301`

## Contributing
This project exists partially to fulfill needs that [Nextcloud](https://nextcloud.com/) doesn't _quite_ solve (e.g., consumer-grade multimedia streaming), but it's mostly just an exercise in software engineering, database design, and writing programs in Go.
