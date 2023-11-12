postgres:
	docker run --name bank-service-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:12-alpine

postgresnetwork:
	docker run --name bank-service-postgres --netwotk bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:12-alpine

createdb:
	 docker exec -it bank-service-postgres createdb --username=root --owner=root bank

dropdb:
	docker exec -it bank-service-postgres dropdb bank

migrateup:
	 migrate -path db/migration -database "postgresql://root:root@localhost:5432/bank?sslmode=disable" -verbose up

migrateup1:
	 migrate -path db/migration -database "postgresql://root:root@localhost:5432/bank?sslmode=disable" -verbose up 1

migratedown:
	 migrate -path db/migration -database "postgresql://root:root@localhost:5432/bank?sslmode=disable" -verbose down

migratedown1:
	 migrate -path db/migration -database "postgresql://root:root@localhost:5432/bank?sslmode=disable" -verbose down 1

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

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
        --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
        proto/*.proto

evants:
	evans -r --host localhost --port 9092  repl

buildimage:
	docker build -t bank:latest .

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlcinit sqlc test mock usermigration proto evants