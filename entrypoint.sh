#!/bin/bash

# Wait for MySQL to become available
until mysqladmin ping -h mysql_db -u root -ppassword; do
  echo "Waiting for MySQL to become available..."
  sleep 1
done

# Run database migrations
echo "Running database migrations..."
make migrate-up

# Start the Go application
echo "Starting the Go application..."
exec "$@"
