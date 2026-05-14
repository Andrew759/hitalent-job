#!/bin/bash
set -e

if ! command -v docker compose &> /dev/null; then
  echo "Ошибка: docker-compose не найден"
  exit 1
fi

BASE_COMPOSE="docker-compose.yml"

echo "Выберите действие $BASE_COMPOSE..."
echo "1) local go main file build"
echo "2) build (no cache)"
echo "3) build (with cache)"
echo "4) up"
echo "5) down"

read -p "Введите соответствующий номер операции над файлом (1-5): " command

case $command in
  1) COMMAND="local_go_build" ;;
  2) COMMAND="build_no_cache" ;;
  3) COMMAND="build_with_cache" ;;
  4) COMMAND="up" ;;
  5) COMMAND="down" ;;
  *) echo "Неверный номер операции"; exit 1 ;;
esac

# Переменная для фиксации контекста в корне проекта
P_DIR="--project-directory ."

if [ "$COMMAND" = "local_go_build" ]; then
  echo ">>> Локальная сборка бинарника под Linux (amd64)..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags="all=-N -l" -o main ./cmd/main.go
  echo ">>> Готово! Файл 'main' создан в корне проекта."
elif [ "$COMMAND" = "build_no_cache" ]; then
    docker compose $P_DIR -f "$BASE_COMPOSE" build --no-cache --progress=plain
elif [ "$COMMAND" = "build_with_cache" ]; then
    docker compose $P_DIR -f "$BASE_COMPOSE" build
elif [ "$COMMAND" = "up" ]; then
    docker compose $P_DIR -f "$BASE_COMPOSE" up
elif [ "$COMMAND" = "down" ]; then
    docker compose $P_DIR -f "$BASE_COMPOSE" down
fi
