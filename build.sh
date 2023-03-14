#!/bin/bash

for OS_ARCH in darwin:amd64 linux:amd64 linux:arm linux:arm64 windows:amd64
do
    arrOS_ARCH=(${OS_ARCH//:/ })
    GOOS=${arrOS_ARCH[0]}
    GOARCH=${arrOS_ARCH[1]}
    docker run --rm -it --env GOOS=${GOOS} --env GOARCH=${GOARCH} --workdir /usr/src/build -v $(pwd):/usr/src/build golang:1.13.15 go build -o bin/cloudflare-dynamic-dns-${GOOS}-${GOARCH}
done