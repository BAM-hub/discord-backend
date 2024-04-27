include .env

steps := 1


create_migration:
	migrate create -ext=sql -dir=internal/database/migrations -seq init

migrate_up:
	migrate -path=internal/database/migrations -database "mysql://${DB_USER}:${DB_PASS}@tcp(127.0.0.1:3306)/discord" -verbose up

migrate_down:
	migrate -path=internal/database/migrations -database "mysql://${DB_USER}:${DB_PASS}@tcp(127.0.0.1:3306)/discord" -verbose down

migrate_down_step:
	migrate -path=internal/database/migrations -database "mysql://${DB_USER}:${DB_PASS}@tcp(127.0.0.1:3306)/discord" -verbose down $(steps)

.PHONY: create_migration migrate_up migrate_down