server:
	go mod tidy && go build -o xkcd-server ./cmd/server

bench:
	go test -bench . ./benching

migrate:
	migrate -database "sqlite3://./database.db" -path "./migrations" up

all: server