FROM alpine:3.15 as builder

RUN apk add --no-cache git go gmp-dev build-base g++ openssl-dev
ADD . /zarb-go

# Building herumi && zarb-daemon
RUN cd /zarb-go && \
    make herumi && \
    export CGO_LDFLAGS="-L$(pwd)/.herumi/bls/lib -lbls384_256 -lm -g -O2" && \
    go env && \
    go build -ldflags "-s -w" -o ./build/zarb-daemon ./cmd/daemon


## Copy binary files from builder into second container
FROM alpine:3.15

COPY --from=builder /zarb-go/build/zarb-daemon /usr/bin

ENV WORKING_DIR "/zarb"

VOLUME $WORKING_DIR
WORKDIR $WORKING_DIR

ENTRYPOINT ["zarb-daemon"]
