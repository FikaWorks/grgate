export DOCKER_BUILDKIT=1

GRGATE_COMMITSHA=$(shell git rev-parse --short HEAD)
GRGATE_VERSION=$(shell git describe --contains $(GRGATE_COMMITSHA))
DOCKER_IMAGE=fikaworks/grgate

.PHONY: \
	all \
	build \
	build-docker \
	integration \
	integration-github \
	integration-gitlab \
	lint \
	mocks \
	push-dockerhub \
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

push-dockerhub: build-docker
	docker tag $(DOCKER_IMAGE) $(DOCKER_IMAGE):$(GRGATE_VERSION)
	docker push $(DOCKER_IMAGE)
	docker push $(DOCKER_IMAGE):$(GRGATE_VERSION)

lint:
	golangci-lint run

mocks:
	mockgen -source=pkg/platforms/platforms.go -destination=pkg/platforms/mocks/platforms_mock.go

test:
	go test -v -parallel=4 ./...

validate: \
	lint \
	test

integration:
	go test -tags=integration ./...

integration-github:
	go test -tags=integrationgithub ./...

integration-gitlab:
	go test -tags=integrationgitlab ./...
