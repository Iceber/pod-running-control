# Build single-platform image
# docker build -t pod-running-control:latest .
#
# Build multi-platform images
# cross-compilation prerequsites: https://docs.docker.com/build/building/multi-platform/#prerequisites
# docker build --platform linux/amd64,linux/arm64 -t pod-running-control:latest .

ARG BUILDER_IMAGE=golang:1.25
ARG BASE_IMAGE=gcr.io/distroless/static:nonroot

FROM --platform=$BUILDPLATFORM ${BUILDER_IMAGE} AS builder

WORKDIR /src
COPY . .

ARG TARGETARCH
RUN GOARCH=${TARGETARCH} go build -o running-control ./cmd

FROM ${BASE_IMAGE}

WORKDIR /
COPY --from=builder /src/running-control /running-control

ENTRYPOINT ["/running-control"]
