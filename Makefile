# Variables
BINARY_NAME=wallet-generator
RECONSTRUCT_DIR=reconstructor
RECONSTRUCT_BINARY_NAME=reconstruct

# Build targets
.PHONY: all set-env download build build-reconstruct clean

all: start set-env download clean build build-reconstruct end

start:
	@echo "\n-----STARTING BUILD-----\n"

end:
	@echo "\n-----BUILD COMPLETE-----\n"

set-env:
	@echo "Setting environment variables"
	@go env -w GOFLAGS=-trimpath

download:
	@echo "Downloading dependencies"
	@go mod download

build:
	@echo "Building wallet-generator binary"
	@go build -o build/$(BINARY_NAME) wallet-generator.go

build-reconstruct:
	@echo "Building reconstruct binary"
	@go build -o build/$(RECONSTRUCT_BINARY_NAME) $(RECONSTRUCT_DIR)/reconstruct.go

clean:
	@echo "Cleaning up and creating required directories"
	@rm -rf build shares
	@mkdir -p build
	@mkdir -p shares