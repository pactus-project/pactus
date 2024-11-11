ARG ALPINE_VERSION=3.20
ARG GO_VERSION=1.23.2
ARG BUILDPLATFORM=linux/amd64
ARG TARGETOS
ARG TARGETARCH

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

WORKDIR /pactus

# Install dependencies
RUN apk add --no-cache git gmp-dev build-base openssl-dev

# Copy source code
COPY . .

# Build pactus-daemon, pactus-wallet, and pactus-shell
RUN GOOS={$TARGETOS} GOARCH={$TARGETARCH} CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ./build/pactus-daemon ./cmd/daemon && \
    GOOS={$TARGETOS} GOARCH={$TARGETARCH} CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ./build/pactus-wallet ./cmd/wallet && \
    GOOS={$TARGETOS} GOARCH={$TARGETARCH} CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ./build/pactus-shell ./cmd/shell

# Final stage
FROM alpine:${ALPINE_VERSION}

COPY --from=builder /pactus/build/pactus-daemon /usr/bin
COPY --from=builder /pactus/build/pactus-wallet /usr/bin
COPY --from=builder /pactus/build/pactus-shell /usr/bin

ENV WORKING_DIR="/pactus"

VOLUME $WORKING_DIR
WORKDIR $WORKING_DIR
