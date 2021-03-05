#!/bin/bash

mkdir -p build
rm -rf build/*

env GOOS=windows GOARCH=amd64 go build -o build/dyndns-netcup-go-windows.exe
env GOOS=linux GOARCH=amd64 go build -o build/dyndns-netcup-go-linux
env GOOS=linux GOARCH=arm go build -o build/dyndns-netcup-go-linux-arm
env GOOS=linux GOARCH=arm64 go build -o build/dyndns-netcup-go-linux-arm64
env GOOS=darwin GOARCH=amd64 go build -o build/dyndns-netcup-go-macos

for file in ./build/*; do 
    tar -czvf $file.tar.gz $file
    openssl dgst -sha256 $file > $file.sha256
    rm $file
done
