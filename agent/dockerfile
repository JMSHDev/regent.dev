# Docker container to automate building agent using Golang

from golang:buster

copy ./ /build_dir/

workdir /build_dir
run go mod download
run go build -ldflags="-linkmode internal" -o agent
