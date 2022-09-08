build: generate
  CGO_ENABLED=0 go build -ldflags "-s -w" -v -o dist/study-server main.go
  cp -r ./docs/ ./dist/

generate:
  go generate ./...

dev:
  gin --bin local/gin-bin -i run main.go

package suffix="":
  mkdir -p release/
  cd dist && zip -r ../release/study-server{{suffix}}.zip *

build-for-target os arch: generate
  CGO_ENABLED=0 GOOS={{os}} GOARCH={{arch}} go build -ldflags "-s -w" -v -o dist/study-server_{{os}}_{{arch}}{{ if os == "windows" { ".exe" } else { "" } }} main.go
  cp -r ./docs/ ./dist/

build-for-all-targets:
  just build-for-target linux amd64
  just build-for-target linux 386
  just build-for-target linux arm64
  just build-for-target linux arm
  just build-for-target windows amd64
  just build-for-target windows 386
  just build-for-target darwin amd64
  just build-for-target darwin arm64

package-gh-release: clean build-for-all-targets
  just package -`git rev-parse --short HEAD`

clean:
  rm -rf dist/
  rm -rf release/
  rm -rf data/
  rm -rf local/gin-bin
