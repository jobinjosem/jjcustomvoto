include ./common.mk

.PHONY: web integration-tests push

all: build integration-tests

build-base-docker-image:
	docker build . -f Dockerfile-base -t "jobinjosem/emojivoto-svc-base:$(IMAGE_TAG)"

web:
	$(MAKE) -C emojivoto-web

build: web

multi-arch:
	$(MAKE) -C emojivoto-web build-multi-arch

push-%:
	docker push jobinjosem/emojivoto-$*:$(IMAGE_TAG)

push: push-svc-base push-web
