#!/bin/bash

echo "Waiting for Postgres to be ready..."
until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do
  sleep 1
done

echo "Running migrations..."
for file in /app/migrations/*.sql; do
  echo "Applying migration: $file"
  psql -h "$DB_HOST" -U "$DB_USER" -d "$DB_NAME" -f "$file"
done

echo "Starting Go application..."
exec go run main.go
