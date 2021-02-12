FROM golang:1.14
ARG http_proxy
ARG https_proxy

ENV http_proxy=${http_proxy} \
  https_proxy=${https_proxy} \
  app_root=/data/dpaycol

RUN mkdir -p ${app_root}

COPY  ./ ${app_root}
WORKDIR ${app_root}
RUN make build

ENTRYPOINT [ "/bin/bash" ]
