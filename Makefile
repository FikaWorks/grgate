export DOCKER_BUILDKIT=1

GGATE_VERSION=$(shell git rev-parse --short HEAD)
DOCKER_IMAGE=fikaworks/ggate

.PHONY: \
	all \
	build \
	build-docker \
	lint \
	test \
	test-all \
	vet

all: test-all build-binary

build-docker:
	docker build \
		--cache-from $(DOCKER_IMAGE) \
		--build-arg BUILDKIT_INLINE_CACHE=1 \
		--build-arg GGATE_VERSION=$(GGATE_VERSION) \
		-t $(DOCKER_IMAGE) .

build:
	go build -ldflags="-X 'github.com/fikaworks/ggate/cmd.Version=$(GGATE_VERSION)'" -a -o ggate .

validate: vet lint test

test:
	go test -v -parallel=4 ./...

lint:
	@go get github.com/golang/lint/golint
	go list ./... | xargs -n1 golint

vet:
	go vet ./...
