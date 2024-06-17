FROM golang:1.22

WORKDIR /usr/src/app

COPY newrelic-bitmovin-analytics /usr/local/bin
COPY configs/standard-config.yml ./configs/config.yml

CMD ["newrelic-bitmovin-analytics"]
