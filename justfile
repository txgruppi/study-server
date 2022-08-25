build: generate
  go build -ldflags "-s -w" -v -o dist/tasks-server main.go

generate:
  go generate ./...

dev:
  gin --bin local/gin-bin -i run main.go
