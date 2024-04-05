FROM golang:1.22.1-alpine3.18 as builder

RUN apk add --no-cache git gmp-dev build-base g++ openssl-dev
ADD . /pactus

# Building pactus-daemon
RUN cd /pactus && \
    go build -ldflags "-s -w" -trimpath -o ./build/pactus-daemon ./cmd/daemon && \
    go build -ldflags "-s -w" -trimpath -o ./build/pactus-wallet ./cmd/wallet


## Copy binary files from builder into second container
FROM golang:1.22.1-alpine3.18

COPY --from=builder /pactus/build/pactus-daemon /usr/bin
COPY --from=builder /pactus/build/pactus-wallet /usr/bin

ENV WORKING_DIR "/pactus"

VOLUME $WORKING_DIR
WORKDIR $WORKING_DIR
