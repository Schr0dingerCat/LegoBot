FROM golang:alpine as builder

RUN apk --no-cache --no-progress add make git

WORKDIR /go/lego

ENV GO11MODULE=on \
    GOPROXY=https://goproxy.cn,direct

COPY . .
RUN make build

FROM alpine:latest
RUN apk update \
    && apk add --no-cache ca-certificates tzdata \
    && update-ca-certificates

COPY --from=builder /go/lego/dist/legobot /usr/bin/legobot

ENTRYPOINT [ "/usr/bin/legobot" ]