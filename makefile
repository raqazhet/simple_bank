postgres:
	docker run --name test1 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD='secret' -d postgres
createdb:
	docker exec -it test1 createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it test1 dropdb simple_bank
migrationup:
	migrate -path ./migrations -database 'postgres://root:secret@localhost:5432/test?sslmode=disable' up
migrationdown:
	migrate -path ./migrations -database 'postgres://root:secret@localhost:5432/test?sslmode=disable' down

.PHONY: postgres createdb dropdb migrationup migrationdown

test:
	go test -v -cover ./...