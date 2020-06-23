#/bin/bash

path=$(cd "$(dirname "$0")"; pwd)
echo $path
cd $path
git add . && git commit -m 'cantor auto deploy' && git push origin master