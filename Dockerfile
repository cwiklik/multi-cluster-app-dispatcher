FROM registry.access.redhat.com/ubi8/go-toolset:1.19.10-10 AS BUILDER
ARG GO_BUILD_ARGS
WORKDIR /workdir
USER root

COPY Makefile Makefile
COPY go.mod go.mod
COPY go.sum go.sum
COPY cmd cmd
COPY pkg pkg
COPY hack hack

ENV GO_BUILD_ARGS=$GO_BUILD_ARGS
RUN echo "Go build args: $GO_BUILD_ARGS" && make mcad-controller

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

COPY --from=BUILDER /workdir/_output/bin/mcad-controller /usr/local/bin

RUN true \
    && microdnf update \
    && microdnf clean all \
    && true

WORKDIR /usr/local/bin

RUN chown -R 1000:1000 /usr/local/bin

USER 1000
