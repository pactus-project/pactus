FROM golang:1.23.1-alpine3.19 as builder

RUN apk add --no-cache git gmp-dev build-base g++ openssl-dev
ADD . /pactus

# Building pactus-daemon
RUN cd /pactus && \
    CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ./build/pactus-daemon ./cmd/daemon && \
    CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ./build/pactus-wallet ./cmd/wallet && \
    CGO_ENABLED=0 go build -ldflags "-s -w" -trimpath -o ./build/pactus-shell ./cmd/shell


## Copy binary files from builder into second container
FROM alpine:3.19

COPY --from=builder /pactus/build/pactus-daemon /usr/bin
COPY --from=builder /pactus/build/pactus-wallet /usr/bin
COPY --from=builder /pactus/build/pactus-shell /usr/bin

ENV WORKING_DIR "/pactus"

VOLUME $WORKING_DIR
WORKDIR $WORKING_DIR
