
CGO_ENABLED=0

.PHONY: build
build: 
	go build -a -installsuffix cgo -o bin/app -v ./cmd

.DEFAULT_GOAL := build