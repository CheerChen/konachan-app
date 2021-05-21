FROM golang:alpine as builder

RUN apk add --no-cache make git

WORKDIR /konakore-src
COPY . /konakore-src
RUN export GOPROXY=https://goproxy.io,direct && \
    go mod download && \
    make docker && \
    mv ./bin/konakore /konakore && \
    mv ./config/config.toml.sample /config.toml

FROM alpine:latest
LABEL org.opencontainers.image.source="https://github.com/CheerChen/konakore"

COPY --from=builder /konakore /
COPY --from=builder /config.toml /config/

ENTRYPOINT ["/konakore","-c","/config/config.toml"]