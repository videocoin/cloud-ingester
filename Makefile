.NOTPARALLEL:
.EXPORT_ALL_VARIABLES:

DOCKER_REGISTRY = us.gcr.io
APP_NAME = ingester

PROJECT_ID= videocoin-network
VERSION=$$(git describe --abbrev=0)-$$(git rev-parse --abbrev-ref HEAD)-$$(git rev-parse --short HEAD)
IMAGE_TAG=$(DOCKER_REGISTRY)/$(PROJECT_ID)/$(APP_NAME):$(VERSION)
HOOKD_IMAGE_TAG=${DOCKER_REGISTRY}/$(PROJECT_ID)/$(APP_NAME)-hookd:${VERSION}

GOOS?=linux
GOARCH?=amd64

.PHONY: deploy

default: all

all: release

version:
	@echo ${VERSION}

build-ingester:
	docker build -t ${IMAGE_TAG} -f Dockerfile.ingester .

build-hookd:
	docker build -t ${HOOKD_IMAGE_TAG} -f Dockerfile.hookd .

build-bin-hookd:
	@echo "==> Building..."
	GOOS=${GOOS} GOARCH=${GOARCH} \
	go build -ldflags="-w -s -X main.Version=${VERSION}" -o bin/hookd hookd/cmd/hookd/main.go

build: build-hookd build-ingester

deps:
	env GO111MODULE=on go mod vendor

push:
	docker push ${HOOKD_IMAGE_TAG}
	docker push ${IMAGE_TAG}

tag:
	@echo ${IMAGE_TAG}
	@echo ${HOOKD_IMAGE_TAG}

deploy:
	cd ./deploy && ./deploy.sh

release: build push
