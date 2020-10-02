
FROM alpine:3.12

ENV WORKING_DIR /zarb


RUN apk add --no-cache bash libstdc++ gmp 


ENV GOPATH /go
ENV PATH "$PATH:/go/bin"

RUN mkdir -p /go/src/github.com/zarbchain/zarb-go
ADD . /go/src/github.com/zarbchain/zarb-go

RUN apk add --no-cache curl git go gmp-dev build-base g++ openssl-dev && \
    mkdir -p /go/src/github.com/zarbchain/zarb-go/ && \
    cd /go/src/github.com/zarbchain/zarb-go/ && \
    make tools build && \
    cp ./build/zarb /usr/bin && \
    apk del curl git go gmp-dev build-base g++ openssl-dev && \
    rm -rf /go /root/.cache
    
EXPOSE 35470
EXPOSE 37620

VOLUME $WORKING_DIR
WORKDIR $WORKING_DIR
ENTRYPOINT ["zarb"]
