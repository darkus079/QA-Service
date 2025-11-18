#!/bin/sh

echo "Waiting for PostgreSQL to be ready..."
while ! nc -z ${DB_HOST:-localhost} ${DB_PORT:-5432}; do
  sleep 1
done

echo "PostgreSQL is ready. Checking migration status..."

DB_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=$DB_SSLMODE"

echo "Checking current migration status..."
goose -dir ./migrations postgres "$DB_URL" status

echo "Running migrations..."
if goose -dir ./migrations postgres "$DB_URL" up; then
    echo "Migrations completed successfully"
    goose -dir ./migrations postgres "$DB_URL" status
else
    echo "Migration failed. Checking if tables already exist..."

    if goose -dir ./migrations postgres "$DB_URL" status 2>/dev/null | grep -q "Applied"; then
        echo "Some migrations already applied. Checking database state..."

        if goose -dir ./migrations postgres "$DB_URL" version 2>/dev/null; then
            echo "Migration version check successful"
            exit 0
        else
            echo "Migration state is inconsistent. Manual intervention may be required."
            exit 1
        fi
    else
        echo "Migration failed completely"
        exit 1
    fi
fi
