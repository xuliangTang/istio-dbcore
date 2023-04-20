FROM golang:1.19-alpine as builder
RUN mkdir /src
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
ADD . /src
WORKDIR /src
RUN GOPROXY=https://goproxy.cn go build -o dbcore main.go && chmod +x dbcore


FROM alpine:3.12
ENV ZONEINFO=/app/zoneinfo.zip
RUN mkdir /app
WORKDIR /app

COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /app

COPY --from=builder /src/dbcore /app
COPY --from=builder /src/app.yml /app

ENTRYPOINT  ["./dbcore"]