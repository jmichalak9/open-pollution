install-go-tools:
	go mod download -x
	cat tools.go | grep _ | grep \".*\" -o | xargs -tI % go install %

.PHONY: generate-code generate-go-code
generate-code: | generate-go-code format

generate-go-code:
	go generate ./...

generate-mocks:
	cd server/measurement && mockgen -source cache.go -destination cache_mock.go -package measurement

format:
	gofumpt -w .
	gci -w -local github.com/jmichalak9/open-pollution . 1>/dev/null

lint:
	golangci-lint run

test:
	go test -v ./...

coverage_profile:
	go test -v -coverpkg=./... -coverprofile=profile.coverprofile ./...

show_coverage:
	go tool cover -html profile.coverprofile

reset_storage:
	rm -rf ${HOME}/.local/share/pdcl/storage/*

.PHONY: install-go-tools format lint coverage_profile show_coverage reset_storage generate-mocks
