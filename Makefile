#!/bin/bash

build: binary

binary:
	@echo "building binary.."
	@go build -tags static_all .


clean:
	@echo "cleaning ..."
	@rm -f sharefood
	@rm -rf vendor
	@rm -f go.sum


install:
	@echo "Installing dependencies...."
	@rm -rf vendor
	@rm -f Gopkg.lock
	@rm -f glide.lock
	@go mod tidy && go mod download && go mod vendor

test:
	@go test $$(go list ./... | grep -v /vendor/) -cover

test-cover:
	@go test $$(go list ./... | grep -v /vendor/) -coverprofile=cover.out && go tool cover -html=cover.out ; rm -f cover.out

coverage:
	@go test -covermode=count -coverprofile=count.out fmt; rm -f count.out

start:
	@go run main.go http