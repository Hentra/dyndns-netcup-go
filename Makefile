DIST_DIR ?= $(shell pwd)/dist

build:
	@echo "Building release artifacts"
	./scripts/make-releases.sh $(DIST_DIR)
.PHONY: build

clean:
	@echo "Cleaning artifacts"
	@rm -rf $(DIST_DIR)
.PHONY: clean
