#!/bin/bash

app=`basename $PWD`
path="./build/${app}.app"
file="./build/${app}.tar.gz"
rm -rf $path
wails build -p 
tar zcvf $file $path