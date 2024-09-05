all: auth_service

auth_service: clean
	go build -o $@ cmd/sso/main.go

migrate:
	go run ./cmd/migrator --storage-path=./storage/sso.db --migrations-path=./migrations

clean:
	rm -rf auth_service