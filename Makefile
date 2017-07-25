.PHONY: install build clean

BUILD_DIR = build
BUILD_ARTIFACT = "$(BUILD_DIR)/konfigurator"

install:
	dep ensure

build:
	@echo "+++ building binary"
	go build -o $(BUILD_ARTIFACT) -ldflags "-X main.version=$(BUILDKITE_TAG)"
	chmod -R 777 $(BUILD_DIR)
	chmod +x $(BUILD_ARTIFACT)

clean:
	rm -rf vendor/ $(BUILD_DIR)
