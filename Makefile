run: build
	@./bin/cardcache

build:
	@go build -o bin/cardcache .

