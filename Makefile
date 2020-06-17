
GOOS?=linux
GOARCH?=amd64
ENV?=dev

NAME=ingester
VERSION?=$$(git describe --abbrev=0)-$$(git rev-parse --abbrev-ref HEAD)-$$(git rev-parse --short HEAD)

REGISTRY_SERVER?=registry.videocoin.net
REGISTRY_PROJECT?=cloud

IMAGE_TAG=${REGISTRY_SERVER}/${REGISTRY_PROJECT}/${NAME}:${VERSION}
HOOKD_IMAGE_TAG=${REGISTRY_SERVER}/${REGISTRY_PROJECT}/${NAME}-hookd:${VERSION}

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
	go build -mod vendor -ldflags="-w -s -X main.Version=${VERSION}" -o bin/hookd hookd/cmd/main.go

build: build-hookd build-ingester

deps:
	GO111MODULE=on go mod vendor

lint: docker-lint

docker-lint:
	docker build -f Dockerfile.lint .

push:
	docker push ${HOOKD_IMAGE_TAG}
	docker push ${IMAGE_TAG}

tag:
	@echo ${IMAGE_TAG}
	@echo ${HOOKD_IMAGE_TAG}

release: build push

deploy:
	cd deploy && helm upgrade -i --wait --set image.tag="${VERSION}" -n console ingester ./helm
