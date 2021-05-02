export DOCKER_BUILDKIT=1

GGATE_VERSION=$(shell git rev-parse --short HEAD)
DOCKER_IMAGE=fikaworks/ggate

.PHONY: \
	all \
	build \
	build-docker \
	lint \
	mocks \
	push-dockerhub \
	test \
	validate

all: \
	validate \
	build

build:
	go build -ldflags="-X 'github.com/fikaworks/ggate/pkg/config.Version=$(GGATE_VERSION)'" -a -o ggate .

build-docker:
	docker build \
		--cache-from $(DOCKER_IMAGE) \
		--build-arg BUILDKIT_INLINE_CACHE=1 \
		--build-arg GGATE_VERSION=$(GGATE_VERSION) \
		-t $(DOCKER_IMAGE) .

push-dockerhub: build-docker
	docker tag $(DOCKER_IMAGE) $(DOCKER_IMAGE):$(GGATE_VERSION)
	docker push $(DOCKER_IMAGE)
	docker push $(DOCKER_IMAGE):$(GGATE_VERSION)

lint:
	golangci-lint run

mocks:
	mockgen -source=pkg/platforms/platforms.go -destination=pkg/platforms/mocks/platforms_mock.go

test:
	go test -v -parallel=4 ./...

validate: \
	lint \
	test
