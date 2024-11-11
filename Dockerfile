ARG ALPINE_VERSION=3.20
ARG GO_VERSION=1.23.2
ARG XCPUTRANSLATE_VERSION=v0.6.0
ARG BUILDPLATFORM=linux/amd64

FROM --platform=${BUILDPLATFORM} qmcgaw/xcputranslate:${XCPUTRANSLATE_VERSION} AS xcputranslate
FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder
ARG TARGETPLATFORM

COPY --from=xcputranslate /xcputranslate /usr/local/bin/xcputranslate

# disable CGO
ENV CGO_ENABLED=0

WORKDIR /pactus

# Install dependencies
RUN apk add --no-cache git gmp-dev build-base openssl-dev

# Copy source code
COPY . .

# download go modules for build
RUN go mod download -x

# Build pactus-daemon, pactus-wallet, and pactus-shell
RUN GOARCH="$(xcputranslate translate -field arch -targetplatform ${TARGETPLATFORM})" \
    GOARM="$(xcputranslate translate -field arm -targetplatform ${TARGETPLATFORM})" \
     go build -ldflags "-s -w" -trimpath -o ./build/pactus-daemon ./cmd/daemon && \
    GOARCH="$(xcputranslate translate -field arch -targetplatform ${TARGETPLATFORM})" \
    GOARM="$(xcputranslate translate -field arm -targetplatform ${TARGETPLATFORM})" \
     go build -ldflags "-s -w" -trimpath -o ./build/pactus-wallet ./cmd/wallet && \
    GOARCH="$(xcputranslate translate -field arch -targetplatform ${TARGETPLATFORM})" \
    GOARM="$(xcputranslate translate -field arm -targetplatform ${TARGETPLATFORM})" \
     go build -ldflags "-s -w" -trimpath -o ./build/pactus-shell ./cmd/shell

# Final stage
FROM alpine:${ALPINE_VERSION}

COPY --from=builder /pactus/build/pactus-daemon /usr/bin
COPY --from=builder /pactus/build/pactus-wallet /usr/bin
COPY --from=builder /pactus/build/pactus-shell /usr/bin

ENV WORKING_DIR="/pactus"

VOLUME $WORKING_DIR
WORKDIR $WORKING_DIR
