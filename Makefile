all: build

## Build the application
.PHONY: build 
build:
	go build -o ./bin/tunnel ./src/main.go


## Run unit tests
.PHONY: test
test: 
	go test ./src/...


## Run the app
.PHONY: run
run:
	go run ./src/main.go
