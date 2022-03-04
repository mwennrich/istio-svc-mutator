GO111MODULE := on
DOCKER_TAG := $(or ${GIT_TAG_NAME}, latest)

all: istio-svc-mutator

.PHONY: istio-svc-mutator
istio-svc-mutator:
	go build istio-svc-mutator.go
	strip istio-svc-mutator

.PHONY: dockerimages
dockerimages: istio-svc-mutator
	docker build -t mwennrich/istio-svc-mutator:${DOCKER_TAG} .

.PHONY: dockerpush
dockerpush: istio-svc-mutator
	docker push mwennrich/istio-svc-mutator:${DOCKER_TAG}
