COMMIT = $(shell git rev-parse HEAD)

build-sentinel-docker:
	docker build -t openpollution-sentinel:$(COMMIT) -f backend/cmd/pdcl/sentinel/Dockerfile backend
	docker tag openpollution-sentinel:$(COMMIT) openpollution-sentinel:latest

build-docker:
	docker-compose build

run-docker:
	docker-compose up -d

stop-docker:
	docker-compose down

.PHONY: build-sentinel-docker run-docker stop-docker build-docker
