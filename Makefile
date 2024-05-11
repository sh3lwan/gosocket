#!make 
include .env

create:
	echo "Creating migration for " $(name)
	migrate create -ext sql -dir database/migrations/ -seq $(name)

up:
	echo "Migrating sql files..."
	migrate -path database/migrations/ -database "mysql://${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_DATABASE}?query" -verbose up

down:
	echo "Rolling back migrations..."
	migrate -path database/migrations -database "mysql://${DB_USERNAME}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_DATABASE}?query" down -all
