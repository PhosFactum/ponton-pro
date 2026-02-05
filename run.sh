#!/bin/bash
set -e

echo "🔥 Останавливаем старые процессы..."

kill_port () {
	PORT=$1
	PIDS=$(lsof -ti tcp:$PORT || true)

	if [ -n "$PIDS" ]; then
	    echo "Убиваем PID на порту $PORT: $PIDS"
	    kill $PIDS
	else
	    echo "Порт $PORT свободен"
	fi
}

kill_port 8080
kill_port 3000

sleep 1


echo "🚀 Запускаем бэкенд (порт 8080)..."
cd /root/TechnoLotos/backend
go run cmd/api/main.go &

sleep 2


echo "🎨 Запускаем фронтенд через npx (порт 3000)..."
cd /root/TechnoLotos/frontend
npx serve -p 3000
sleep 2

IP=$(hostname -I | awk '{print $1}')
echo ""
echo "✅ Готово!"
echo "👉 Фронт: http://$IP:3000"
echo "👉 API:   http://$IP:8080/api"
