FROM golang:alpine as builder

RUN apk add --no-cache make git

WORKDIR /kanachan-src
COPY . /kanachan-src
RUN export GOPROXY=https://goproxy.io,direct && \
    go mod download && \
    make docker && \
    mv ./bin/konachan-app /kanachan-app && \
    mv ./config/config.toml.sample /config.toml

FROM alpine:latest
LABEL org.opencontainers.image.source="https://github.com/CheerChen/konachan-app"

COPY --from=builder /kanachan-app /
COPY --from=builder /config.toml /config/

ENTRYPOINT ["/kanachan-app","-c","/config/config.toml"]