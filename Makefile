.NOTPARALLEL:
.EXPORT_ALL_VARIABLES:
.DEFAULT_GOAL := main

APP_NAME?=ingester
VERSION?=$$(git rev-parse --short HEAD)
PROJECT_ID?=
DOCKER_REGISTRY?=us.gcr.io
IMAGE_TAG=${DOCKER_REGISTRY}/${PROJECT_ID}/${APP_NAME}:${VERSION}
LATEST=${DOCKER_REGISTRY}/${PROJECT_ID}/${APP_NAME}:latest
PACKAGE_FILENAME=ingester-${VERSION}.deb
DEB_REPO_PASSWORD?=
REPO_ADDR?=http://aptly.videocoin.io

GOOS?=linux
GOARCH?=amd64

.PHONY: deploy

default: all

main: build-docker-image push-docker-image

all: build-docker-image push-docker-image

version:
	@echo ${VERSION}

build-docker-image:
	docker build -t ${IMAGE_TAG} -t ${LATEST} .

local-docker:
	docker build -t ingester -f Dockerfile.local .

push-docker-image:
	gcloud docker -- push ${IMAGE_TAG}
	gcloud docker -- push ${LATEST}

docker-image-tag:
	@echo ${IMAGE_TAG}

deploy:
	cd ./deploy && ./deploy.sh

build-hookd:
	@echo "==> Building..."
	cd hookd/cmd
	GOOS=${GOOS} GOARCH=${GOARCH} \
	go build \
		-ldflags="-w -s -X main.Version=${VERSION}" \
		-o build/release/stream-ingester/bin/stream-ingester-hookd \
		github.com/videocoin/ingester/hookd

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

publish-deb:
	PACKAGE_FILE=./build/pkg/${PACKAGE_FILENAME} \
	PACKAGE_NAME=${PACKAGE_FILENAME} \
	REPO_USER=circleci \
	REPO_PWD=${DEB_REPO_PASSWORD} \
	REPO_ADDR=${REPO_ADDR} \
	./repo_publish.sh