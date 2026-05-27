VERSION := $(shell git describe --tags --always 2>/dev/null || echo "unknown")
BINARY := smtp-checker-$(VERSION)-amd64
BUILD_DIR := build

.PHONY: build-PC
build-PC: | $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY) ./src/

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

.PHONY: release
release: | $(BUILD_DIR)
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY) ./src/

.PHONY: install
install: release
	./install.sh $(BUILD_DIR)/$(BINARY)

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)/.[^.]* $(BUILD_DIR)/*
