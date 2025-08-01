postgres:
	sudo docker run --name postgres-17 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine3.22

createdb:
	sudo docker exec -it postgres-17 createdb --username=root --owner=root simple_bank

dropdb:
	sudo docker exec -it postgres-17 dropdb simple_bank

migrateup:
	migrate --path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate --path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate --path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate --path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

test-out:
	go test -v -cover -coverprofile=coverage.out ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/shevgn/simplebank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc test test-out server mock
