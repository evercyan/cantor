#/bin/bash

# mac 进程是 $4, 可视情况调整
ps aux | grep 'cantor' | grep -v 'grep' | awk '{print $4}' | xargs kill -9

path=$(cd "$(dirname "$0")"; pwd)
nohup python3 "${path}/src/main.py" &