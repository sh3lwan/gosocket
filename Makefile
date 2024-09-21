#!make 
include .env

create:
	echo "Creating migration for " $(name)
	migrate create -ext sql -dir database/migrations/ -seq $(name)

init:
	echo "Initiating Database..."
	mysql -u ${DB_USERNAME} -p ${DB_PASSWORD} < internal/db/init.sql

up:
	echo "Migrating sql files..."
	migrate -path internal/db/migrations/ -database "mysql://${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_DATABASE}?query" -verbose up

down:
	echo "Rolling back migrations..."
	migrate -path internal/db/migrations -database "mysql://${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_DATABASE}?query" down -all
