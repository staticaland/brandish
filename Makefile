Here is the code for the Makefile:

.PHONY: build run

build:
	go build -o main main.go

run:
	go run main.go

