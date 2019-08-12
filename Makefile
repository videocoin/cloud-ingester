.NOTPARALLEL:
.EXPORT_ALL_VARIABLES:
.DEFAULT_GOAL := docker

DOCKER_REGISTRY = us.gcr.io
CIRCLE_ARTIFACTS = ./bin
SERVICE_NAME = ingester

PROJECT_ID= videocoin-network
VERSION=$$(git describe --abbrev=0)-$$(git rev-parse --short HEAD)
IMAGE_TAG=$(DOCKER_REGISTRY)/$(PROJECT_ID)/$(SERVICE_NAME):$(VERSION)

PACKAGE_FILENAME=ingester-${VERSION}.deb
DEB_REPO_PASSWORD?=
REPO_ADDR?=http://aptly.videocoin.io

GOOS?=linux
GOARCH?=amd64

default: all

main: build-docker-image push-docker-image

all: build-docker-image push-docker-image

version:
	@echo ${VERSION}

docker:
	docker build -t ${IMAGE_TAG} .

local-docker:
	docker build -t ingester -f Dockerfile.local .

push:
	docker push ${IMAGE_TAG}

tag:
	@echo ${IMAGE_TAG}

deploy:
	cd ./deploy && ./deploy.sh

build-deb:
	docker build \
		--build-arg VERSION=${VERSION} \
		-t lp-stream-ingester-deb:${VERSION} \
		-f Dockerfile.deb . && \
	docker run \
		--rm -it \
		-v `pwd`/build/pkg:/build/pkg \
		lp-stream-ingester-deb:${VERSION} \
		cp /pkg/${PACKAGE_FILENAME} /build/pkg
