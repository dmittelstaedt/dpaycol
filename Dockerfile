FROM golang:1.11
ARG http_proxy
ARG https_proxy

ENV http_proxy=${http_proxy}
ENV https_proxy=${https_proxy}

WORKDIR /go/src/app

COPY .git .
COPY dpaycol.go .

RUN go get ./...
RUN VERSION=$(git tag --list | tail -1 | cut -c 2-) && \
GIT_COMMIT=$(git rev-parse --short HEAD) && \
BUILD_DATE=$(date +"%Y-%m-%d %T") && \
go build -ldflags "-X main.versionNumber=$VERSION -X main.gitCommit=$GIT_COMMIT -X 'main.buildDate=$BUILD_DATE'" dpaycol.go

ENTRYPOINT [ "/bin/bash" ]
