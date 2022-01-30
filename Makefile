COMMIT = $(shell git rev-parse HEAD)

build-sentinel-docker:
	docker build -t openpollution-sentinel:$(COMMIT) -f backend/cmd/pdcl/sentinel/Dockerfile backend
	docker tag openpollution-sentinel:$(COMMIT) openpollution-sentinel:latest

.PHONY: build-sentinel-docker
