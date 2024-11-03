include app.env

DB_URL=$(DB_SOURCE)

server:
	go run main.go

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

sqlc:
	sqlc generate

build:
	GOOS=linux GOARCH=amd64 go build -o bin/application

.PHONY: server sqlc, build
