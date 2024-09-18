#!/bin/bash
set -x
mkdir -p bin

go build -o bin/chip8 ./cmd/chip8/
go build -o bin/dis ./cmd/dis/
GOOS=js GOARCH=wasm go build -o static/chip8.wasm ./cmd/wasm/

# strip minor version
GOVERSION=$(go env GOVERSION | awk -F. '{print $1"."$2}')
wget -nc -O static/wasm_exec.js \
    "https://raw.githubusercontent.com/golang/go/release-branch.${GOVERSION}/misc/wasm/wasm_exec.js"
