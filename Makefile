build:
	go mod tidy && go build -o xkcd.exe ./cmd/xkcd

all: build