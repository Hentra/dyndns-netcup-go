#!/bin/bash

BIN_DIR="$1"

[ -n "$BIN_DIR" ] || {
    echo "No directory for the binary files specified."
    exit 1
}

mkdir -p "$BIN_DIR"
rm -rf "$BIN_DIR/*"

ARCHS="windows,amd64,windows.exe linux,amd64,linux linux,arm,linux-arm linux,arm64,linux-arm64 darwin,amd64,macos"

for arch in $ARCHS; do IFS=","; set -- $arch
    env GOOS=$1 GOARCH=$2 go build -o "$BIN_DIR/dyndns-netcup-go-$3" ./cmd/dyndns-netcup-go/main.go
done

for file in "$BIN_DIR"/*; do 
    tar -czvf $file.tar.gz $file
    openssl dgst -sha256 $file > $file.sha256
    rm $file
done
