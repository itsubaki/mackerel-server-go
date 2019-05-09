SHELL := /bin/bash
DATE := $(shell date +%Y%m%d-%H:%M:%S)
HASH := $(shell git rev-parse HEAD)

runserver:
	set -x
	-rm ${GOPATH}/bin/mackerel-api
	go install
	GIN_MODE=debug mackerel-api

runclient:
	set -x
	-rm ~/Library/mackerel-agent/id
	mackerel-agent -conf /usr/local/etc/mackerel-agent.conf -apibase=http://localhost:8080

runmysql:
	set -x
	docker pull mysql

	-docker stop mysqld
	-docker rm mysqld
	docker run --name mysqld -e MYSQL_ROOT_PASSWORD=secret -p 3307:3306 -d mysql
	# mysql -h127.0.0.1 -P3307 -psecret -uroot -e'create database mackerel;'


test:
	set -x
	curl -s localhost:8080/api/v0/org -H "X-Api-Key: secret" | jq .
	curl -s localhost:8080/api/v0/services -X POST -H "Content-Type: application/json" -d '{"name": "ExampleService", "memo": "This is an example."}' | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/metadata/foobar -X PUT -H "Content-Type: application/json" -d '{"message": "this is service metadata"}' | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/metadata/foobar | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/metadata | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/metadata/foobar -X DELETE | jq .
	curl -s localhost:8080/api/v0/services/ExampleService -X DELETE | jq .
	curl -s localhost:8080/api/v0/services | jq .
	curl -s localhost:8080/api/v0/services -X POST -H "Content-Type: application/json" -d '{"name": "ExampleService", "memo": "This is an example."}' | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/roles -X POST -H "Content-Type: application/json" -d '{"name": "ExampleRole", "memo": "This is an example."}' | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/roles | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/roles/ExampleRole/metadata/foobar -X PUT -H "Content-Type: application/json" -d '{"message": "this is role metadata"}' | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/roles/ExampleRole/metadata/foobar | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/roles/ExampleRole/metadata | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/roles/ExampleRole/metadata/foobar -X DELETE | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/roles/ExampleRole -X DELETE | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/tsdb -X POST -H "Content-Type: application/json" -d '[{"name":"hoge", "time": 1351700030, "value": 1.234},{"name":"foobar", "time": 1351700050, "value": 5.678}]' | jq .
	curl -s "localhost:8080/api/v0/services/ExampleService/metrics?name=hoge&from=1351700000&to=1351700100" | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/metric-names | jq .
	curl -s localhost:8080/api/v0/services/ExampleService -X DELETE | jq .
	$(eval HOSTID := $(shell curl -s localhost:8080/api/v0/hosts -X POST -H "Content-Type: application/json" -d '{"name":"host01", "meta": {"agent-name": "mackerel-agent/0.27.0 (Revision dfbccea)", "agent-revision": "2f531c6", "agent-version": "0.4.2"}}' | jq -r .id))
	curl -s localhost:8080/api/v0/hosts/${HOSTID} -X PUT -H "Content-Type: application/json" -d '{"name":"host01kai", "meta": {"agent-name": "mackerel-agent/0.27.0 (Revision dfbccea)", "agent-revision": "2f531c6", "agent-version": "0.4.2"}}' | jq .
	curl -s localhost:8080/api/v0/hosts/${HOSTID}/status -X POST -H "Content-Type: application/json" -d '{"status": "poweroff"}' | jq .
	curl -s localhost:8080/api/v0/hosts/${HOSTID} | jq .
	curl -s localhost:8080/api/v0/hosts/${HOSTID}/role-fullnames -X PUT -H "Content-Type: application/json" -d '{"roleFullnames": ["Hatena-Bookmark:db-master"]}' | jq .
	curl -s localhost:8080/api/v0/hosts/${HOSTID}/metadata/foobar -X PUT -H "Content-Type: application/json" -d '{"message": "this is host metadata"}' | jq .
	curl -s localhost:8080/api/v0/hosts/${HOSTID}/metadata/foobar | jq .
	curl -s localhost:8080/api/v0/hosts/${HOSTID}/metadata | jq .
	curl -s localhost:8080/api/v0/hosts/${HOSTID}/metadata/foobar -X DELETE | jq .
	curl -s localhost:8080/api/v0/tsdb -X POST -H "Content-Type: application/json" -d '[{"hostId": "${HOSTID}", "name":"hoge", "time": 1351700030, "value": 1.234},{"hostId": "${HOSTID}", "name":"foobar", "time": 1351700050, "value": 5.678}]' | jq .
	curl -s "localhost:8080/api/v0/tsdb/latest?hostId=${HOSTID}&name=foobar&name=hoge&name=piyo" | jq .
	curl -s "localhost:8080/api/v0/hosts/${HOSTID}/metrics?name=hoge&from=1351700000&to=1351700100" | jq .
	curl -s localhost:8080/api/v0/hosts/${HOSTID}/metric-names | jq .
	curl -s localhost:8080/api/v0/hosts | jq .
	curl -s localhost:8080/api/v0/hosts/${HOSTID}/retire -X POST -H "Content-Type: application/json" -d '{}' | jq .
	curl -s localhost:8080/api/v0/invitations -X POST -H "Content-Type: application/json" -d '{"email": "example@example.com","authority": "viewer"}' | jq .
	curl -s localhost:8080/api/v0/invitations
	curl -s localhost:8080/api/v0/invitations/revoke -X POST -H "Content-Type: application/json" -d '{"email": "example@example.com"}' | jq .
