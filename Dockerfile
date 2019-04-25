FROM golang:1.11
ARG http_proxy
ARG https_proxy

ENV http_proxy=${http_proxy}
ENV https_proxy=${https_proxy}

WORKDIR /go/src/app

COPY dpaycol.go .

RUN go get ./...
RUN go build dpaycol.go

ENTRYPOINT [ "/bin/bash" ]
