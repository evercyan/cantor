#/bin/bash

ps aux | grep 'cantor/src/main.py' | grep -v 'grep' | awk '{print $2}' | xargs kill -9

path=$(cd "$(dirname "$0")"; pwd)
nohup python3 "${path}/src/main.py" >/dev/null 2>&1 &