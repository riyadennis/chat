#!/bin/bash -e

run:
	go run main.go --address=:8080 --traceStatus=true
deps:
	glide install
test:
	go test -timeout=5s -cover -race $$(glide novendor)