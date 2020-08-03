#/bin/bash

ps aux | grep 'cantor/main.py' | grep -v 'grep' | awk '{print $2}' | xargs kill -9

path=$(cd "$(dirname "$0")"; pwd)
nohup python3 "${path}/main.py" >/dev/null 2>&1 &

ps aux | grep 'cantor/main.py'