# FROM golang:1.22

# WORKDIR /usr/src/app

# COPY bin/docker/newrelic-bitmovin-analytics /usr/local/bin
# COPY configs/standard-config.yml ./configs/config.yml

# ENTRYPOINT ["newrelic-bitmovin-analytics"]

FROM dhi.io/golang:1.26.4-alpine3.24-dev AS build-stage

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 go build ./cmd/bitmovin/bitmovin.go

FROM dhi.io/golang:1.26.4-alpine3.24 AS release-stage

WORKDIR /app

COPY --from=build-stage /build/bitmovin ./

ENTRYPOINT ["bitmovin"]