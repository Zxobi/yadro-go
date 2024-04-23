build:
	go mod tidy && go build -o xkcd ./cmd/xkcd

bench:
	go test -bench . ./benching

all: build