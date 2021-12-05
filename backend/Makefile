install-go-tools:
	go mod download -x
	cat tools.go | grep _ | grep \".*\" -o | xargs -tI % go install %

format:
	gofumpt -w .
	gci -w -local github.com/jmichalak9/open-pollution . 1>/dev/null

lint:
	golangci-lint run

run:
	ADDRESS="0.0.0.0:9000" PDCL_HOST="localhost" PDCL_PORT="8000" go run cmd/backend/main.go

test:
	go test -v ./...

coverage_profile:
	go test -v -coverpkg=./... -coverprofile=profile.coverprofile ./...

show_coverage:
	go tool cover -html profile.coverprofile

.PHONY: install-go-tools format lint run coverage_profile show_coverage
