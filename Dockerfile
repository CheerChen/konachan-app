FROM golang:alpine as builder

RUN apk add --no-cache make git

WORKDIR /kanachan-src
COPY . /kanachan-src
RUN go mod download && \
    make docker && \
    mv ./bin/konachan-app /kanachan-app && \
    mv ./config/config.toml.sample /config.toml

FROM alpine:latest
LABEL org.opencontainers.image.source="https://github.com/CheerChen/konachan-app"

COPY --from=builder /kanachan-app /
COPY --from=builder /config.toml /

ENTRYPOINT ["/kanachan-app","-c","config.toml"]