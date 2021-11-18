install-go-tools:
	go mod download -x
	cat tools.go | grep _ | grep \".*\" -o | xargs -tI % go install %

format:
	gofumpt -w .

lint:
	golangci-lint run

run:
	ADDRESS="127.0.0.1:8000" go run cmd/backend/main.go

test:
	go test -v ./...

coverage_profile:
	go test -v -coverpkg=./... -coverprofile=profile.coverprofile ./...

show_coverage:
	go tool cover -html profile.coverprofile

.PHONY: install-go-tools format lint run coverage_profile show_coverage
