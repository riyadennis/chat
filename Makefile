#!/bin/bash -e

run:
	go run main.go --address=:8080 --traceStatus=true
test:
	go test -timeout=5s -cover -race
unlock:
	git-crypt unlock
docker-build:
	docker build -t riyadennis/chat .
docker-run:
	docker run --rm -p 8080:8080  chat
docker-push:
	docker push riyadennis/chat:latest