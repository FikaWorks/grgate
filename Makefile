export DOCKER_BUILDKIT=1

GGATE_VERSION=$(shell git rev-parse --short HEAD)
DOCKER_IMAGE=fikaworks/ggate

.PHONY: \
	all \
	build \
	build-docker \
	lint \
	mocks \
	test \
	validate

all: \
	validate \
	build

build-docker:
	docker build \
		--cache-from $(DOCKER_IMAGE) \
		--build-arg BUILDKIT_INLINE_CACHE=1 \
		--build-arg GGATE_VERSION=$(GGATE_VERSION) \
		-t $(DOCKER_IMAGE) .

build:
	go build -ldflags="-X 'github.com/fikaworks/ggate/pkg/config.Version=$(GGATE_VERSION)'" -a -o ggate .

validate: lint test

test:
	go test -v -parallel=4 ./...

lint:
	golangci-lint run

mocks:
	mockgen -source=pkg/platforms/platforms.go -destination=pkg/platforms/mocks/platforms_mock.go
