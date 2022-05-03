DIST_DIR ?= $(shell pwd)/dist

ci: clean lint
.PHONY: ci

build:
	@echo "Building release artifacts"
	./scripts/make-releases.sh $(DIST_DIR)
.PHONY: build

clean:
	@echo "Cleaning artifacts"
	@rm -rf $(DIST_DIR)
.PHONY: clean

lint:
	@echo "Linting sources"
	@docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:v1.43.0 golangci-lint run -v
.PHONY: lint

docker-build:
	@echo "Building docker image"
	@docker build -t dyndns-netcup-go -f ./build/package/Dockerfile .
.PHONY: docker-build
