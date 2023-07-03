DB_URL=postgresql://devuser:admin@localhost:5432/devdb?sslmode=disable

# generate migrations, $(name) - name of the migration
generate_migrations:
	migrate create -ext sql -dir db/migrations -seq $(name)

# run up migrations, user details based on docker-compose.yml
migrate_up:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

# run down migrations, user details based on docker-compose.yml
migrate_down:
	migrate -path db/migrations -database "$(DB_URL)}" -verbose down

# generate db related go code with sqlc
sqlc:
	cmd.exe /c "docker run --rm -v ${PWD}:/src -w /src kjconroy/sqlc generate"

# generate database documentation on the dbdocs website
db_docs:
	dbdocs build docs/database.dbml

# generate .sql file with database schema
db_schema:
	dbml2sql --postgres -o docs/schema.sql docs/database.dbml

# generate mock db for testing
mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/aalug/go-bank/db/sqlc Store

# run all tests
test:
	go test -v -cover ./...

# run tests in the given path (p) and display results in the html file
test_coverage:
	go test $(p) -coverprofile=coverage.out && go tool cover -html=coverage.out

# run the HTTP and gRPC servers
runserver:
	go run main.go

# remove old files and generate new proto files. Generate swagger files
protoc:
	rm -f pb/*.go
	rm -f docs/swagger/*.swagger.json
	protoc --proto_path=protobuf --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=bank_go \
	protobuf/*.proto
	statik -src=./docs/swagger -dest=./docs

.PHONY: migrate_up, migrate_down, sqlc, test, test_coverage, runserver, mock, db_schema, db_docs, protobuf