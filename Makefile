COMMIT = $(shell git rev-parse HEAD)

build-sentinel-docker:
	docker build -t openpollution-sentinel:$(COMMIT) -f backend/cmd/pdcl/sentinel/Dockerfile backend
	docker tag openpollution-sentinel:$(COMMIT) openpollution-sentinel:latest

run-producer:
	cd backend;\
	SIGNER_ID="arek-noster-manual-hygu9uhib" PRODUCER_KEY_PATH="${HOME}/priv.pem" GRPC_HOST=35.216.150.202 GRPC_PORT=8000 go run cmd/pdcl/random-producer/main.go

build-docker:
	docker-compose build

run-docker:
	docker-compose up -d

stop-docker:
	docker-compose down

.PHONY: build-sentinel-docker run-docker stop-docker build-docker
