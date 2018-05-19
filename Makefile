PREFIX ?= $(shell pwd)
BIN_DIR ?= $(shell pwd)
DOCKER_IMAGE_NAME ?= ix2005-exporter
DOCKER_IMAGE_TAG ?= $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))

.PHONY: all build tarball promu docker

all: build

build: promu
	@promu build --prefix $(PREFIX)

tarball: promu
	@promu tarball --prefix $(PREFIX) $(BIN_DIR)

promu:
	@go get -u github.com/prometheus/promu

docker:
	@docker build -t "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" .
