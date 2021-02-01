bin/museum:
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/museum ./cmd/

.PHONY: dep
dep:
	go mod tidy

.PHONY: run-migration
run-migration:
	migrate -source migrations -database ${DB_ADDR} up -verbose
