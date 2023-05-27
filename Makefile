build:
	@go build -o bin/api

run: build
	@./bin/api

up_db:
	@sudo docker run --name mongodb -p 27017:27017 mongo:latest

test:
	@go test -v ./...