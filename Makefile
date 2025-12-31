build:
	@go build -o bin/devpool .

run: build
	@./bin/devpool
