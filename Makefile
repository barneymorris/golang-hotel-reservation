build:
	@go build -o bin/api

run: build
	@./bin/api

up_db:
	@sudo docker run --name mongodb -p 27017:27017 mongo:latest

seed:
	@go run scripts/seed.go

test:
	@go test -v ./...