#!/bin/bash -e

run:
	go run main.go --address=:8085 --traceStatus=true
test:
	go test -timeout=5s -cover -race
unlock:
	git-crypt unlock
docker-build:
	docker build -t chat .
docker-run:
	docker run --rm -p 8085:8085  chat
docker-push: