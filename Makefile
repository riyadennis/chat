#!/bin/bash -e

run:
	go run main.go --address=:8080 --traceStatus=true
test:
	go test -timeout=5s -cover -race
unlock:
	git-crypt unlock