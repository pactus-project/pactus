FROM alpine:3.15 as builder

RUN apk add --no-cache git go gmp-dev build-base g++ openssl-dev
ADD . /zarb-go
RUN cd /zarb-go && make release

## Copy binary files from builder into second container
FROM alpine:3.15

COPY --from=builder /zarb-go/zarbd /usr/bin

ENV WORKING_DIR "/zarb"

VOLUME $WORKING_DIR
WORKDIR $WORKING_DIR

ENTRYPOINT ["zarb"]
