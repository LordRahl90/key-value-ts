start:
	go run ./cmd

test:
	go test ./... --cover

build:
	docker build -t lordrahl/kvts:latest .

run: build
	docker-compose up