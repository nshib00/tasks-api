#!/bin/bash

if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
else
    echo "Файл .env не найден!"
    exit 1
fi
DB_HOST=${POSTGRES_HOST}
DB_USER=${POSTGRES_USER}
DB_PASSWORD=${POSTGRES_PASSWORD}
DB_NAME=${POSTGRES_DB}
DB_PORT=${POSTGRES_PORT:-5432}

docker run --rm -v $(pwd)/migrations:/migrations --network host \
   migrate/migrate \
  -path=/migrations/ \
  -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" \
  "$@"