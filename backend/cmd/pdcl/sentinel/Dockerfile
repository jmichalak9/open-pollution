FROM golang:1.17-buster AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o /openpollution-sentinel cmd/pdcl/sentinel/main.go

FROM scratch
EXPOSE 9000
COPY --from=builder /openpollution-sentinel /openpollution-sentinel
ENTRYPOINT ["/openpollution-sentinel"]
