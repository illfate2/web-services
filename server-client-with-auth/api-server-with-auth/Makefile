bin/museum:
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/museum ./cmd/

.PHONY: dep
dep:
	go mod tidy

.PHONY: run-migration
run-migration:
	migrate -path migration -database ${DB_ADDR} -verbose up

.PHONY: graphql-gen
graphql-gen:
	gqlgen generate -config .gqlgen.yml
