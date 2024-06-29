FROM golang:1.22

WORKDIR /usr/src/app

COPY bin/docker/newrelic-bitmovin-analytics /usr/local/bin
COPY configs/standard-config.yml ./configs/config.yml

ENTRYPOINT ["newrelic-bitmovin-analytics"]
