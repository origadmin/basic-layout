FROM golang:1.23.1 AS builder

ENV GOPROXY=https://goproxy.cn

COPY . /src
WORKDIR /src

RUN make generate
RUN make build

FROM debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
        netbase \
        && rm -rf /var/lib/apt/lists/* \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/dist/helloworld /app/helloworld
COPY --from=builder /src/dist/secondworld /app/secondworld
COPY --from=builder /src/dist/gateway /app/gateway

COPY resources/configs /app/resources/configs

COPY docker-entrypoint.sh /usr/local/bin/docker-entrypoint.sh
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

WORKDIR /app

EXPOSE 8000
EXPOSE 9000

VOLUME /data/conf

ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["gateway"]