ALL:
	@echo "make build - build the project"
	@echo "make build-front - build the front"

build:
	@echo "Building..."
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o wc3 ./cmd/crawler/main.go

build-front:
	@echo "Building front..."
	cd infra/mgr/front && pnpm install && pnpm run build && cd -

.PHONY: build build-front