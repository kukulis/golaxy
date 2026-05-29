# web game 

Test multiple files

    go test ./...

## run

    go run cmd/server/main.go

Browser:

    http://localhost:8080

## swagger

    http://localhost:8080/swagger/index.html

# pitfalls

The Go json.Unmarshal works case-insensitive, while JS decoding json is case-sensitive.

