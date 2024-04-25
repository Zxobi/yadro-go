cli:
	go mod tidy && go build -o xkcd ./cmd/xkcd

server:
	go mod tidy && go build -o xkcd-server ./cmd/server

bench:
	go test -bench . ./benching

all: cli