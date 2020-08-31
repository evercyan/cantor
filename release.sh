#!/bin/bash

app=`basename $PWD`
path="./build/${app}.app"
file="./build/${app}.tar.gz"
rm -rf $path
wails build -p 
cd "./build"
tar zcvf "${app}.tar.gz" "${app}.app"