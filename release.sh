#!/bin/bash

service=`basename $PWD`
version=$1
if [ "${version}" == "" ]; 
then
    echo "invalid version";
	exit;
fi

rm -rf "./build/${service}.app"
rm -rf "./dist/${service}.app"
wails build -p 
mv "./build/${service}.app" "./dist/${service}.app"

git add . && git commit -m $version && git push origin master
git tag -a $version -m $version
git push origin $version