migrate-create:
	echo "Creating migration for " $(name)
	migrate create -ext sql -dir database/migrations/ -seq $(name)


migrate-up:
	echo "Migrating sql files..."
	migrate -path database/migrations/ -database "mysql://root:password@tcp(mysql_db:3306)/chat?query" -verbose up


migrate-down:
	echo "Rolling back migrations..."
	migrate -path database/migrations -database "mysql://root:password@tcp(mysql_db:3306)/chat?query" -verbose down
