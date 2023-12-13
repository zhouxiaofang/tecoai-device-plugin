IMAGE_VERSION = latest
REGISTRY = docker.io/joyme
IMAGE = ${REGISTRY}/tecoai-device-plugin:${IMAGE_VERSION}

.PHONY: build deploy

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/tecoai cmd/server/app.go

buildImage:
	docker build -t ${IMAGE} .

kindLoad:
	kind load docker-image ${IMAGE}

pushImage:
	docker push ${IMAGE}

deploy:
	helm install tecoai deploy/helm/tecoai

upgrade:
	helm upgrade tecoai deploy/helm/tecoai

dry-run:
	helm install tecoai deploy/helm/tecoai --dry-run
