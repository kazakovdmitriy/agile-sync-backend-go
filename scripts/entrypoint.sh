#!/bin/bash

set -e

# Проверяем, что DATABASE_URL задан
if [ -z "$DATABASE_URL" ]; then
  echo "Ошибка: переменная окружения DATABASE_URL не задана"
  exit 1
fi

# Извлекаем хост и порт из DATABASE_URL для проверки доступности
# Пример URL: postgres://user:pass@host:port/dbname?...
# Регулярное выражение для извлечения host:port
if [[ $DATABASE_URL =~ postgresql?://[^:]+:[^@]+@([^:/]+)(:([0-9]+))? ]]; then
  DB_HOST="${BASH_REMATCH[1]}"
  DB_PORT="${BASH_REMATCH[3]:-5432}"  # порт по умолчанию — 5432
else
  echo "Ошибка: не удалось распарсить DATABASE_URL"
  exit 1
fi

echo "Ожидание готовности базы данных ($DB_HOST:$DB_PORT)..."
while ! nc -z "$DB_HOST" "$DB_PORT" 2>/dev/null; do
  sleep 1
done
echo "База данных доступна"

# Применяем миграции (они используют DATABASE_URL из окружения)
echo "Применение миграций..."
./migrate up

# Запускаем сервер
echo "Запуск сервера..."
exec ./server