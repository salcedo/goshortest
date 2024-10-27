#!/bin/bash

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags=jsoniter -a -o goshortest

docker-image: build
	docker build -t goshortest .

dev:
	DATABASE_DSN='host=localhost user=goshortest password=goshortest dbname=goshortest sslmode=disable' DEFAULT_SITE="https://salcedo.dev/explore/repos" EXPIRATION='-168h' TOKEN="reallysecret" CompileDaemon -exclude-dir=.git -command=./goshortest

clean:
	go clean
	rm -f goshortest

all: build
