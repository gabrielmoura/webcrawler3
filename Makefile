ALL:
	@echo "make build - build the project"

build:
	@echo "Building..."
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o wc3 ./cmd/crawler/main.go

.PHONY: build