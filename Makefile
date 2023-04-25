postgres:
	docker run --name postgres12 -p 5432:5432  -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

psql:
	docker exec -it postgres12 psql -d simple_bank

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bamk
test:
	go test -v -cover ./...
sqlc:
	sqlc generate
protoc:
	protoc --go_out=./pb --go_opt=paths=source_relative \
    --go-grpc_out=./pb --go-grpc_opt=paths=source_relative \
    proto/*.proto
dbup:
	migrate -path db/migration -database "postgresql://root:Hoang2002@localhost:5432/simple_bank?sslmode=disable" -verbose up
dbdown:
	migrate -path db/migration -database "postgresql://root:Hoang2002@localhost:5432/simple_bank?sslmode=disable" -verbose down
runTest:
	go test -v ./...
.PHONY: postgres createdb dropdb psql sqlc protoc
