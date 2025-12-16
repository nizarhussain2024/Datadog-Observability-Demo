.PHONY: build run clean

build:
	go build -o app ./app

run:
	go run ./app

clean:
	rm -f app




