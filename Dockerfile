FROM golang:1.13
ARG http_proxy
ARG https_proxy

ENV http_proxy=${http_proxy}
ENV https_proxy=${https_proxy}
ENV app_root=/go/src/github.com/dmittelstaedt/dpaycol

RUN mkdir -p ${app_root}

COPY  ./ ${app_root}
WORKDIR ${app_root}
RUN make build

ENTRYPOINT [ "/bin/bash" ]
