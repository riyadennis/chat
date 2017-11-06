#!/bin/bash -e

run_server:
	go run main.go --address=:8081
deps:
	glide install
test:
	go test -timeout=5s -cover -race $$(glide novendor)