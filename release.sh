#!/bin/bash

set -e

version=$1
if [ "${version}" == "" ];
then
    echo "sh release.sh v0.0.1"
    exit
fi

rm -rf "./build/cantor.app"
wails build -p 
cd "./build"
tar zcvf "cantor-${version}.tar.gz" "cantor.app"