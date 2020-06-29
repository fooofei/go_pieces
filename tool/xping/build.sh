#!/bin/bash

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"

function compile() {
    name=$1
    echo "compiling $name"
    cd $SCRIPTPATH/$name
    rm -f $name.linux.amd64
    GOOS=linux GOARCH=amd64 go build -mod=vendor -v -o $name.linux.amd64
    GOOS=windows GOARCH=amd64 go build -mod=vendor -v
}

compile httping
compile tcping
compile udping
