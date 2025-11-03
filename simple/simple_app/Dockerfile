FROM golang:1.23.1 AS builder

ARG APP_NAME=simple_app
ENV GOPROXY=https://goproxy.cn

COPY . /src
WORKDIR /src

RUN make generate
RUN make build

FROM debian:stable-slim

ARG APP_NAME=simple_app

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
        netbase \
        && rm -rf /var/lib/apt/lists/* \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/dist/${APP_NAME} /app/${APP_NAME}

COPY configs /app/configs
COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

WORKDIR /app

EXPOSE 8000
EXPOSE 9000

VOLUME /data/conf

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["simple_app"]
