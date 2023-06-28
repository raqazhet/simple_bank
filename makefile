mock:
	mockgen -destination storage/mock/store.go bank/storage Store
migrate:
	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
	mv migrate.linux-amd64 $GOPATH/bin/migrate
postgres:
	docker run --name test1 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD='secret' -d postgres
createdb:
	docker exec -it test1 createdb --username=root --owner=root test
dropdb:
	docker exec -it test1 dropdb test
migrationup:
	migrate -path ./migrations -database 'postgres://root:secret@localhost:5432/test?sslmode=disable' up
migrationdown:
	migrate -path ./migrations -database 'postgres://root:secret@localhost:5432/test?sslmode=disable' down
server:
	go run main.go
.PHONY: postgres createdb dropdb migrationup migrationdown server migrate mock proto evans

test:
	go test -v -cover ./...

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
        --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
        --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
        proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl