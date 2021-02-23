SHELL := /bin/bash
DATE := $(shell date +%Y%m%d-%H:%M:%S)
HASH := $(shell git rev-parse HEAD)
XAPIKEY := 2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb

runserver:
	go version
	GIN_MODE=debug SQL_MODE=debug go run main.go

runclient:
	go version
	-rm ~/Library/mackerel-agent/{id,pid}
	-go get github.com/mackerelio/mackerel-agent
	-GO111MODULE=off go get github.com/mackerelio/go-check-plugins
	cd $(shell go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-tcp; go install
	cp mackerel-agent.conf /usr/local/etc/mackerel-agent.conf
	mackerel-agent -conf /usr/local/etc/mackerel-agent.conf -apibase=http://localhost:8080

runmysql:
	-docker pull mysql
	-docker stop mysql
	-docker rm mysql
	docker run --name mysql -e MYSQL_ROOT_PASSWORD=secret -p 3306:3306 -d mysql
	docker ps
	# mysql -h127.0.0.1 -P3306 -uroot -psecret -Dmackerel

cleanup:
	docker stop $(shell docker ps -q -a)
	docker rm $(shell docker ps -q -a)
	docker rmi $(shell docker images -q)

build:
	docker build -t mackerel-server-go .

up: build
	docker-compose up
	# docker exec -it ${CONTAINERID} mysql -u root -p

test:
	go version
	go test -v -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/ | grep -v -E "mackerel-server-go$$")

godog:
	go version
	SQL_MODE=debug go test -v --godog.format=pretty
# 	go test -v --godog.format=pretty -coverprofile=coverage.out -covermode=atomic
# 	go tool cover -html=coverage.out -o coverage.html

mkr:
	go version
	go get github.com/mackerelio/mkr

	MACKEREL_APIKEY=${XAPIKEY} mkr --apibase=http://localhost:8080 org
	MACKEREL_APIKEY=${XAPIKEY} mkr --apibase=http://localhost:8080 create mkr-host
	MACKEREL_APIKEY=${XAPIKEY} mkr --apibase=http://localhost:8080 hosts
	MACKEREL_APIKEY=${XAPIKEY} mkr --apibase=http://localhost:8080 services
	MACKEREL_APIKEY=${XAPIKEY} mkr --apibase=http://localhost:8080 alerts
