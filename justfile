build: generate
  CGO_ENABLED=0 go build -ldflags "-s -w" -v -o dist/study-server main.go
  cp -r ./docs/ ./dist/

generate:
  go generate ./...

dev:
  gin --bin local/gin-bin -i run main.go

package: clean build
  mkdir -p release/
  cd dist && zip -r ../release/study-server.zip *

clean:
  rm -rf dist/
  rm -rf release/
  rm -rf data/
  rm -rf local/gin-bin
