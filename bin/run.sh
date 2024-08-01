#!/bin/bash

# 配置
HOST="127.0.0.1"
PORT="6443"
MODE="release"

BASE_ROOT=$(dirname $(dirname "$(readlink -f "$0")"))
LOGPATH="$BASE_ROOT/log/go-dress.log"

start() {
    PID=$(ss -antp | grep :$PORT | awk '{print $6}' | cut -d',' -f2 | cut -d'=' -f2)
    if [ -n "$PID" ]; then
        echo "han-dress is already running"
    else
        nohup go run "$BASE_ROOT"/main.go -h "$HOST" -p "$PORT" -m "$MODE" > "$LOGPATH" 2>&1 &
        echo "han-dress started"
    fi
}

stop() {
    PID=$(ss -antp | grep :$PORT | awk '{print $6}' | cut -d',' -f2 | cut -d'=' -f2)
    if [ -n "$PID" ]; then
        kill -TERM "$PID" >/dev/null
    fi
    echo "han-dress stopped"
}

restart() {
    stop
    start
}

reload() {
    PID=$(ss -antp | grep :$PORT | awk '{print $6}' | cut -d',' -f2 | cut -d'=' -f2)
    if [ -n "$PID" ]; then
      kill -HUP "$PID"
      echo "han-dress reloaded"
    else
        start
    fi
}

case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    reload)
        reload
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|reload}"
        exit 1
esac
