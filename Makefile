SHELL := /bin/bash
DATE := $(shell date +%Y%m%d-%H:%M:%S)
HASH := $(shell git rev-parse HEAD)

runserver:
	-rm ${GOPATH}/bin/mackerel-api
	go install
	GIN_MODE=debug mackerel-api

test:
	curl -v localhost:8080/api/v0/services -X POST -H "Content-Type: application/json" -d '{"name": "ExampleService", "memo": "This is an example."}'
	curl -v localhost:8080/api/v0/services/ExampleService/metadata/foobar -X PUT -H "Content-Type: application/json" -d '{"message": "this is service metadata"}'
	curl -v localhost:8080/api/v0/services/ExampleService/metadata/foobar
	curl -v localhost:8080/api/v0/services/ExampleService/metadata
	curl -v localhost:8080/api/v0/services/ExampleService/metadata/foobar -X DELETE
	curl -v localhost:8080/api/v0/services/ExampleService -X DELETE
	curl -v localhost:8080/api/v0/services
	curl -v localhost:8080/api/v0/services -X POST -H "Content-Type: application/json" -d '{"name": "ExampleService", "memo": "This is an example."}'
	curl -v localhost:8080/api/v0/services/ExampleService/roles -X POST -H "Content-Type: application/json" -d '{"name": "ExampleRole", "memo": "This is an example."}'
	curl -v localhost:8080/api/v0/services/ExampleService/roles
	curl -v localhost:8080/api/v0/services/ExampleService/roles/ExampleRole/metadata/foobar  -X PUT -H "Content-Type: application/json" -d '{"message": "this is role metadata"}'
	curl -v localhost:8080/api/v0/services/ExampleService/roles/ExampleRole/metadata/foobar
	curl -v localhost:8080/api/v0/services/ExampleService/roles/ExampleRole/metadata
	curl -v localhost:8080/api/v0/services/ExampleService/roles/ExampleRole/metadata/foobar -X DELETE
	curl -v localhost:8080/api/v0/services/ExampleService/roles/ExampleRole -X DELETE
	curl -v localhost:8080/api/v0/services/ExampleService/tsdb -X POST -H "Content-Type: application/json" -d '[{"name":"hoge", "time": 1351700030, "value": 1.234},{"name":"foobar", "time": 1351700050, "value": 5.678}]'
	curl -v "localhost:8080/api/v0/services/ExampleService/metrics?name=hoge&from=1351700000&to=1351700100"
	curl -v localhost:8080/api/v0/services/ExampleService/metric-names
	curl -v localhost:8080/api/v0/services/ExampleService -X DELETE
