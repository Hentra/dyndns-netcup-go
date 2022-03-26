#!/bin/bash

mkdir -p build
rm -rf build/*

ARCHS="windows,amd64,windows.exe linux,amd64,linux linux,arm,linux-arm linux,arm64,linux-arm64 darwin,amd64,macos"

for arch in $ARCHS; do IFS=","; set -- $arch
    env GOOS=$1 GOARCH=$2 go build -o build/dyndns-netcup-go-$3 ./cmd/dyndns-netcup-go/main.go
done

for file in ./build/*; do 
    tar -czvf $file.tar.gz $file
    openssl dgst -sha256 $file > $file.sha256
    rm $file
done
