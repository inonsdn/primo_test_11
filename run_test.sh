#!/bin/bash

# start test postgres only
docker compose -f docker-compose-test.yml up -d

# wait for postgres to be ready
echo "Waiting for postgres..."
until docker exec postgres_test pg_isready -U test -d testdb; do
    echo "Postgres not ready, retrying..."
    sleep 1
done
echo "Postgres is ready"

# run tests locally
export DATABASE_USER=test
export DATABASE_PASSWORD=test
export DATABASE_HOST=localhost
export DATABASE_PORT=5433
export DATABASE_DBNAME=testdb
export DATABASE_SSLMODE=disable

go test ./...

# stop test postgres when done
docker compose -f docker-compose-test.yml down
