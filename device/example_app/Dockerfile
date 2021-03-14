# Docker container to automate building ExampleApp using Golang

from golang:buster

copy ./ /build_dir/

workdir /build_dir
run go build -ldflags="-linkmode internal" -o ExampleApp -i ExampleApp.go
