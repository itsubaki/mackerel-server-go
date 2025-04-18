SHELL := /bin/bash
DATE := $(shell date +%Y%m%d-%H:%M:%S)
HASH := $(shell git rev-parse HEAD)
XAPIKEY := 2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb

update:
	go get -u
	go mod tidy

runserver:
	go version
	GIN_MODE=debug SQL_MODE=debug RUN_FIXTURE=true go run main.go

runclient:
	go version
	-rm ~/Library/mackerel-agent/{id,pid}
	-go install github.com/mackerelio/mackerel-agent@latest
	git clone https://github.com/mackerelio/go-check-plugins
	cd go-check-plugins/check-tcp; go install
	mackerel-agent -conf mackerel-agent.conf -apibase=http://localhost:8080

runmysql:
	-docker pull mysql
	-docker stop mysql
	-docker rm mysql
	docker run --name mysql -e MYSQL_ROOT_PASSWORD=secret -p 3306:3306 -d mysql
	docker ps
	# mysql -h127.0.0.1 -P3306 -uroot -psecret -Dmackerel

testpkg:
	go version
	go test -v -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/ | grep -v -E "mackerel-server-go$$") -coverprofile=coverage-pkg.out -covermode=atomic
	go tool cover -html=coverage-pkg.out -o coverage-pkg.html

test:
	go version
	SQL_MODE=debug go test -v -coverprofile=coverage.out -covermode=atomic -coverpkg ./...
	go tool cover -html=coverage.out -o coverage.html

mkr:
	go version
	go install github.com/mackerelio/mkr@latest

	MACKEREL_APIKEY=${XAPIKEY} mkr --apibase=http://localhost:8080 org
	MACKEREL_APIKEY=${XAPIKEY} mkr --apibase=http://localhost:8080 create mkr-host
	MACKEREL_APIKEY=${XAPIKEY} mkr --apibase=http://localhost:8080 hosts
	MACKEREL_APIKEY=${XAPIKEY} mkr --apibase=http://localhost:8080 services
	MACKEREL_APIKEY=${XAPIKEY} mkr --apibase=http://localhost:8080 alerts

build:
	docker build -t mackerel-server-go .

up: build
	docker compose up
	# docker exec -it ${CONTAINERID} mysql -u root -p

down:
	docker compose down

cleanup:
	docker stop $(shell docker ps -q -a)
	docker rm   $(shell docker ps -q -a)
	docker rmi  $(shell docker images -q)
	docker volume rm $(docker volume ls -qf dangling=true)

package:
	docker build -t ghcr.io/itsubaki/mackerel-server-go .
	docker push     ghcr.io/itsubaki/mackerel-server-go

dbdoc:
	tbls diff
	tbls doc --force
	tbls coverage
