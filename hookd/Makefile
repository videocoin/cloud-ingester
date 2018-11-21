.NOTPARALLEL:
.EXPORT_ALL_VARIABLES:
.DEFAULT_GOAL := main

DOCKER_REGISTRY = us.gcr.io
CIRCLE_ARTIFACTS = ./bin
PROTOS_PATH = ./
GRPC_CPP_PLUGIN = grpc_cpp_plugin
GRPC_CPP_PLUGIN_PATH ?= `which $(GRPC_CPP_PLUGIN)`
SERVICE_NAME = hookd
GOOGLE_PROJECT = videocoin-183500


VERSION=$$(git rev-parse --short HEAD)
IMAGE_TAG=$(DOCKER_REGISTRY)/$(GOOGLE_PROJECT)/$(SERVICE_NAME):$(VERSION)
LATEST=$(DOCKER_REGISTRY)/$(GOOGLE_PROJECT)/$(SERVICE_NAME):latest

IMAGE_TARBALL_PATH=$(CIRCLE_ARTIFACTS)/$(SERVICE_NAME)-$(VERSION).tar

main: docker push

version:
	@echo $(VERSION)

image-tag:
	@echo $(IMAGE_TAG)

deps:
	 go mod vendor && go mod verify


build:
	@echo "==> Building..."
	go build -o bin/$(SERVICE_NAME) cmd/main.go

build-alpine:
	go build -o bin/$(SERVICE_NAME) --ldflags '-w -linkmode external -extldflags "-static"' cmd/main.go

test:
	@echo "==> Running tests..."
	go test -v ./...

test-coverage:
	@echo "==> Running tests..."
	go test -cover ./...

docker:
	docker build  -t ${IMAGE_TAG} -t $(LATEST) . --squash

push:
	@echo "==> Pushing $(IMAGE_TAG)..."
	gcloud docker -- push $(IMAGE_TAG)
	gcloud docker -- push $(LATEST)
	

docker-save:
	@echo "==> Saving docker image tarball..."
	gcloud auth configure-docker --quiet
	docker save -o $(IMAGE_TARBALL_PATH) $(IMAGE_TAG)
