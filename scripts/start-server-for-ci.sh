#!/bin/bash

# 为 CI 环境启动服务器的脚本
# 使用方法: ./scripts/start-server-for-ci.sh <port> <max_attempts>

set -e

PORT="${1:-8080}"
MAX_ATTEMPTS="${2:-20}"
STARTUP_DELAY="${3:-2}"

echo "Starting server on port $PORT..."
PORT=$PORT ./bin/server &
SERVER_PID=$!
echo $SERVER_PID > server.pid

echo "Server PID: $SERVER_PID"
echo "Waiting ${STARTUP_DELAY}s for server to fully start..."
sleep $STARTUP_DELAY

echo "Starting health checks (max ${MAX_ATTEMPTS} attempts)..."
for i in $(seq 1 $MAX_ATTEMPTS); do
    if curl -fsS "http://localhost:${PORT}/healthz" >/dev/null 2>&1; then
        echo "✅ Health check passed after ${STARTUP_DELAY}s + ${i}s"
        echo $SERVER_PID
        exit 0
    fi
    
    if [ $i -eq $MAX_ATTEMPTS ]; then
        echo "❌ Health check failed after ${STARTUP_DELAY}s + ${MAX_ATTEMPTS}s"
        echo "Server logs:"
        if [ -f server.pid ]; then
            PID=$(cat server.pid)
            if kill -0 "$PID" 2>/dev/null; then
                echo "Server process $PID is still running but not responding"
            else
                echo "Server process $PID has died"
            fi
        fi
        kill $SERVER_PID 2>/dev/null || true
        rm -f server.pid
        exit 1
    fi
    
    sleep 1
done