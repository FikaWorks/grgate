export DOCKER_BUILDKIT=1

GRGATE_COMMITSHA=$(shell git rev-parse --short HEAD)
GRGATE_VERSION=$(shell git describe --contains $(GRGATE_COMMITSHA))
DOCKER_IMAGE=ghcr.io/fikaworks/grgate

.PHONY: \
	all \
	build \
	build-docker \
	integration \
	integration-github \
	integration-gitlab \
	lint \
	lint-fix \
	mocks \
	push-docker \
	test \
	validate

all: \
	validate \
	build

build:
	go build \
		-ldflags="-X 'github.com/fikaworks/grgate/pkg/config.Version=$(GRGATE_VERSION)' -X 'github.com/fikaworks/grgate/pkg/config.CommitSha=$(GRGATE_COMMITSHA)'" \
		-a -o grgate .

build-docker:
	docker build \
		--cache-from $(DOCKER_IMAGE) \
		--build-arg BUILDKIT_INLINE_CACHE=1 \
		--build-arg GRGATE_COMMITSHA=$(GRGATE_COMMITSHA) \
		--build-arg GRGATE_VERSION=$(GRGATE_VERSION) \
		-t $(DOCKER_IMAGE) .

push-docker: build-docker
	docker tag $(DOCKER_IMAGE) $(DOCKER_IMAGE):$(GRGATE_VERSION)
	docker push $(DOCKER_IMAGE)
	docker push $(DOCKER_IMAGE):$(GRGATE_VERSION)

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix
	golangci-lint run --fix --build-tags unit
	golangci-lint run --fix --build-tags integration

mocks:
	go generate ./...

test:
	go test -tags=unit -v -parallel=4 ./...

validate: \
	lint \
	test

integration:
	go test -p 1 -tags=integration ./...

integration-github:
	go test -p 1 -tags=integrationgithub ./...

integration-gitlab:
	go test -p 1 -tags=integrationgitlab ./...
