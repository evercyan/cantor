#!/bin/bash

service=`basename $PWD`
commands=(
    "debug"
    "test"
    "build"
)

################################

debug() {
    wails serve
}

test() {
    wails build -d
    `./build/${service}`
}

build() {
    rm -rf "./build/${service}.app"
    wails build -p 
    `open ./build/${service}.app`
}

################################

command=$1
filter=`echo "${commands[@]}" | grep -w "$command"`
if [ "${command}" == "" -o "${filter}" == "" ]; then
	exit;
fi

$command