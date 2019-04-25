SHELL := /bin/bash
DATE := $(shell date +%Y%m%d-%H:%M:%S)
HASH := $(shell git rev-parse HEAD)

runserver:
	-rm ${GOPATH}/bin/mackerel-api
	go install github.com/itsubaki/mackerel-api/cmd/mackerel-api
	GIN_MODE=debug mackerel-api

test:
	curl -v localhost:8080/api/v0/services -X POST -H "Content-Type: application/json" -d '{"name": "ExampleService", "memo": "This is an example."}'
	curl -v localhost:8080/api/v0/services
	curl -v localhost:8080/api/v0/services/ExampleService/roles -X POST -H "Content-Type: application/json" -d '{"name": "ExampleRole","memo": "This is an example."}'
	curl -v localhost:8080/api/v0/services/ExampleService/roles
#	curl -v localhost:8080/api/v0/services/ExampleService/roles/ExampleRole -X DELETE
	curl -v localhost:8080/api/v0/services/ExampleService -X DELETE
