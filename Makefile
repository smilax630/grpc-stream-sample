VERSION := dev-latest

image-build:
	docker build -t asia.gcr.io/projectID/streamer:${VERSION} .

image-push:
	docker push asia.gcr.io/projectID/streamer:${VERSION}

docker: image-build image-push
	@echo finish

.PHONY: image-build image-push docker
