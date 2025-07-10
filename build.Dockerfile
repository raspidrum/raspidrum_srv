FROM golang:1.24.4-bookworm AS builder

WORKDIR /src

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=arm64

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

CMD ["go", "build", "-o", "./build/raspidrum",  "./cmd/server/main.go"]