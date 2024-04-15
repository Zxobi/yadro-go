build:
	go mod tidy && go build -o xkcd ./cmd/xkcd

all: build