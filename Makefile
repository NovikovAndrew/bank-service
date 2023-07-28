postgres:
	docker run bank-service-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:12-alpine

createdb:
	 docker exec -it bank-service-postgres createdb --username=root --owner=root bank

dropdb:
	docker exec -it bank-service-postgres dropdb bank

migrateup:
	 migrate -path db/migration -database "postgresql://root:root@localhost:5432/bank?sslmode=disable" -verbose up

migratedown:
	 migrate -path db/migration -database "postgresql://root:root@localhost:5432/bank?sslmode=disable" -verbose down

sqlcinit:
	sqlc init

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

usermigration:
	migrate create -ext sql -dir db/migration -seq add_users

mock:
	mockgen -package mockdb -destination db/mock/store.go  bank-service/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlcinit sqlc test mock usermigration