
FROM golang:1.21 AS builder

COPY . /src
WORKDIR /src

RUN GOPROXY=https://goproxy.cn make build

# 有需要就修改对应的系统
FROM ubuntu:latest

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates  \
        netbase \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y

COPY --from=builder /src/bin /app
COPY ./configs /data/conf
WORKDIR /app

EXPOSE 8000
EXPOSE 9000
# VOLUME /data/conf

CMD ["./whatshappening", "-conf", "/data/conf"]