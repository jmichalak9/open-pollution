install-go-tools:
	go mod download -x
	cat tools.go | grep _ | grep \".*\" -o | xargs -tI % go install %

format:
	gofumpt -w .

lint:
	golangci-lint run