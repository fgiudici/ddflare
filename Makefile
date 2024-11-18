GIT_COMMIT?=$(shell git rev-parse --short HEAD)
GIT_VERSION?=$(shell git describe --tags 2>/dev/null || echo "v0.0.0-"$(GIT_COMMIT))
IMG_REG?=ghcr.io
IMG_REPO?=fgiudici/ddflare
IMG_TAG?=$(GIT_VERSION)

export ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BUILD_DIR:=bin

LDFLAGS := -w -s
LDFLAGS += -X "github.com/fgiudici/ddflare/pkg/version.Version=${GIT_VERSION}"

COVERFILE?=coverage.out

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o $(BUILD_DIR)/ddflare

.PHONY: clean
clean:
	@rm -rf $(BUILD_DIR) && rm -f $(COVERFILE)

.PHONY: docker
docker:
	DOCKER_BUILDKIT=1 docker build \
		-f Dockerfile \
		--build-arg "VERSION=${GIT_VERSION}" \
		-t ${IMG_REG}/${IMG_REPO}:${IMG_TAG}

.PHONY: unit-tests
unit-tests:
	@go test -coverprofile $(COVERFILE) ./pkg/...
