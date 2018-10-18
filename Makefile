APP_NAME?=stream-ingester
VERSION?=$$(git describe --abbrev=0)-$$(git rev-parse --short HEAD)
GOOGLE_PROJECT?=liveplanet-cloud-staging
DOCKER_REGISTRY?=us.gcr.io
IMAGE_TAG=${DOCKER_REGISTRY}/${GOOGLE_PROJECT}/${APP_NAME}:${VERSION}
PACKAGE_FILENAME=liveplanet-cloud-stream-ingester-${VERSION}.deb
DEB_REPO_PASSWORD?=
REPO_ADDR?=http://aptly.liveplanetstage.net

GOOS?=linux
GOARCH?=amd64

.PHONY: deploy

default: all

all: build-docker-image push-docker-image

version:
	@echo ${VERSION}

build-docker-image:
	docker build -t ${IMAGE_TAG} .

push-docker-image:
	gcloud docker -- push ${IMAGE_TAG}

docker-image-tag:
	@echo ${IMAGE_TAG}

deploy:
	cd ./deploy && ./deploy.sh

build-hookd:
	@echo "==> Building..."
	GOOS=${GOOS} GOARCH=${GOARCH} \
	go build \
		-ldflags="-w -s -X main.Version=${VERSION}" \
		-o build/release/stream-ingester/bin/stream-ingester-hookd \
		gitlab.videocoin.io/ingester/hookd

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