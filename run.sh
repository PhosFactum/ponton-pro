#!/bin/bash

echo "🔥 Останавливаем старые процессы..."
pkill -f "go run"
pkill -f "node.*serve"
sleep 2

echo "🚀 Запускаем бэкенд (порт 8080)..."
cd /root/TechnoLotos/backend
go run cmd/api/main.go &
sleep 3

echo "🎨 Запускаем фронтенд через npx (порт 3000)..."
cd /root/TechnoLotos/frontend
npx serve -p 3000 &
sleep 2

IP=$(hostname -I | awk '{print $1}')
echo ""
echo "✅ Готово!"
echo "👉 Фронт: http://$IP:3000"
echo "👉 API:   http://$IP:8080/api"
