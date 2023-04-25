start:
	go run ./cmd

test:
	go test ./... --cover

test-with-race:
	go test -race ./... --cover

build:
	docker build -t lordrahl/kvts:latest .

run: build
	docker-compose up

twr: test-with-race