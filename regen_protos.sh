#!/bin/bash

go get -u github.com/golang/protobuf/protoc-gen-go
export PATH=$GOLANG/bin:$PATH
cd proto
rm -f *.go
protoc --go_out=. *.proto
