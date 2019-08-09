SHELL := /bin/bash
DATE := $(shell date +%Y%m%d-%H:%M:%S)
HASH := $(shell git rev-parse HEAD)
XAPIKEY := 2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb
CONTENTTYPE := application/json

install:
	set -x
	-rm $(shell go env GOPATH)/bin/mackerel-api
	go install

runserver:
	set -x
	-rm $(shell go env GOPATH)/bin/mackerel-api
	go install
	GIN_MODE=release mackerel-api

runclient:
	set -x
	GO111MODULE=off go get -d github.com/mackerelio/go-check-plugins
	cd $(shell go env GOPATH)/src/github.com/mackerelio/go-check-plugins/check-tcp; go install
	cp mackerel-agent.conf /usr/local/etc/mackerel-agent.conf
	-rm ~/Library/mackerel-agent/id
	mackerel-agent -conf /usr/local/etc/mackerel-agent.conf -apibase=http://localhost:8080

runmysql:
	set -x
	-docker pull mysql
	-docker stop mysql
	-docker rm mysql
	docker run --name mysql -e MYSQL_ROOT_PASSWORD=secret -p 3306:3306 -d mysql
	# mysql -h127.0.0.1 -P3306 -uroot -psecret -Dmackerel

build:
	set -x
	docker build -t mackerel-api .

up: build
	set -x
	docker-compose up
	# docker exec -it ${CONTAINERID} mysql -u root -p

test:
	set -x
	go test -cover $(shell go list ./... | grep -v /vendor/ | grep -v /build/) -v

curl:
	set -x

	curl -s localhost:8080/api/v0/org -H "X-Api-Key: ${XAPIKEY}" | jq .

	$(eval HOSTID := $(shell curl -s localhost:8080/api/v0/hosts -X POST -H "X-Api-Key: ${XAPIKEY}" -H "Content-Type: application/json" -d '{"name":"host01", "meta": {"agent-name": "mackerel-agent/0.27.0 (Revision dfbccea)", "agent-revision": "2f531c6", "agent-version": "0.4.2"}}' | jq -r .id))
	curl -s localhost:8080/api/v0/hosts/${HOSTID} -X PUT -H "X-Api-Key: ${XAPIKEY}" -H "Content-Type: application/json" -d '{"name":"host01kai", "meta": {"agent-name": "mackerel-agent/0.27.0 (Revision dfbccea)", "agent-revision": "2f531c6", "agent-version": "0.4.2"}}' | jq .
	curl -s localhost:8080/api/v0/hosts/${HOSTID}/status -X POST -H "X-Api-Key: ${XAPIKEY}" -H "Content-Type: application/json" -d '{"status": "poweroff"}' | jq .
	curl -s localhost:8080/api/v0/hosts/${HOSTID} -H "X-Api-Key: ${XAPIKEY}" | jq .
	curl -s localhost:8080/api/v0/hosts/${HOSTID}/role-fullnames -X PUT -H "X-Api-Key: ${XAPIKEY}" -H "Content-Type: application/json" -d '{"roleFullnames": ["Hatena-Bookmark:db-master"]}' | jq .
	curl -s localhost:8080/api/v0/hosts/${HOSTID}/metadata/foobar -X PUT -H "X-Api-Key: ${XAPIKEY}" -H "Content-Type: application/json" -d '{"message": "this is host metadata"}' | jq .
	curl -s localhost:8080/api/v0/hosts/${HOSTID}/metadata/foobar -H "X-Api-Key: ${XAPIKEY}" | jq .
	curl -s localhost:8080/api/v0/hosts/${HOSTID}/metadata -H "X-Api-Key: ${XAPIKEY}" | jq .
	curl -s localhost:8080/api/v0/hosts/${HOSTID}/metadata/foobar -X DELETE -H "X-Api-Key: ${XAPIKEY}" | jq .
	curl -s localhost:8080/api/v0/tsdb -X POST -H "X-Api-Key: ${XAPIKEY}" -H "Content-Type: application/json" -d '[{"hostId": "${HOSTID}", "name":"hoge", "time": 1351700030, "value": 1.234},{"hostId": "${HOSTID}", "name":"foobar", "time": 1351700050, "value": 5.678}]' | jq .
	curl -s "localhost:8080/api/v0/tsdb/latest?hostId=${HOSTID}&name=foobar&name=hoge&name=piyo" -H "X-Api-Key: ${XAPIKEY}" | jq .
	curl -s "localhost:8080/api/v0/hosts/${HOSTID}/metrics?name=hoge&from=1351700000&to=1351700100" -H "X-Api-Key: ${XAPIKEY}" | jq .
	curl -s localhost:8080/api/v0/hosts/${HOSTID}/metric-names -H "X-Api-Key: ${XAPIKEY}" | jq .
	curl -s localhost:8080/api/v0/hosts -H "X-Api-Key: ${XAPIKEY}" | jq .
	curl -s localhost:8080/api/v0/hosts/${HOSTID}/retire -X POST -H "X-Api-Key: ${XAPIKEY}" -H "Content-Type: application/json" -d '{}' | jq .

	curl -s localhost:8080/api/v0/graph-annotations -H "X-Api-Key: ${XAPIKEY}" | jq .
	$(eval ANNOTATIONID := $(shell curl -s localhost:8080/api/v0/graph-annotations -X POST -H "X-Api-Key: ${XAPIKEY}" -d '{ "title": "deploy application", "description": "link: https://example.com/", "from": 1484000000, "to": 1484000030, "service": "ExampleService", "roles": [ "ExampleRole1", "ExampleRole2" ] }'| jq -r .id))
	curl -s localhost:8080/api/v0/graph-annotations/${ANNOTATIONID} -X PUT -H "X-Api-Key: ${XAPIKEY}" -d '{ "title": "deploy application", "description": "link: https://example.com/example", "from": 1484000000, "to": 1484000030, "service": "ExampleService", "roles": [ "ExampleRole1", "ExampleRole2" ] }'| jq .
	curl -s localhost:8080/api/v0/graph-annotations/${ANNOTATIONID} -X DELETE -H "X-Api-Key: ${XAPIKEY}" | jq .

	curl -s localhost:8080/api/v0/notification-groups -H "X-Api-Key: ${XAPIKEY}" | jq .
	$(eval GROUPID := $(shell 	curl -s localhost:8080/api/v0/notification-groups -X POST -H "X-Api-Key: ${XAPIKEY}" -d '{ "name": "Example notification group", "notificationLevel": "all", "childNotificationGroupIds": [], "childChannelIds": [ "2vh7AZ21abc" ], "monitors": [ { "id": "2qtozU21abc", "skipDefault": false } ], "services": [ { "name": "Example-Service-1" }, { "name": "Example-Service-2" } ] }' | jq -r .id) )
	curl -s localhost:8080/api/v0/notification-groups/${GROUPID} -X PUT -H "X-Api-Key: ${XAPIKEY}" -d '{ "name": "Example notification group", "notificationLevel": "critical", "childNotificationGroupIds": [], "childChannelIds": [ "2vh7AZ21abc" ], "monitors": [ { "id": "2qtozU21abc", "skipDefault": false } ], "services": [ { "name": "Example-Service-1" }, { "name": "Example-Service-2" } ] }' | jq -r .id
	curl -s localhost:8080/api/v0/notification-groups/${GROUPID} -X DELETE -H "X-Api-Key: ${XAPIKEY}"

	curl -s localhost:8080/api/v0/channels -H "X-Api-Key: ${XAPIKEY}" | jq .
	$(eval CHANNELID := $(shell curl -s localhost:8080/api/v0/channels -X POST -H "X-Api-Key: ${XAPIKEY}" -d '{ "type": "email", "name": "My Channel", "emails": ["myaddress@example.com"], "userIds": ["userId"], "events": ["alert"] }' | jq -r .id ))
	curl -s localhost:8080/api/v0/channels/${CHANNELID} -X DELETE -H "X-Api-Key: ${XAPIKEY}" | jq .

	curl -s localhost:8080/api/v0/alerts -H "X-Api-Key: ${XAPIKEY}" | jq .
	$(eval ALERTID := $(shell curl -s localhost:8080/api/v0/alerts -H "X-Api-Key: ${XAPIKEY}" | jq '.alerts[0].id' -r ))
	curl -s localhost:8080/api/v0/alerts/${ALERTID}/close -X POST -H "X-Api-Key: ${XAPIKEY}" -d '{ "reason": "manual" }' | jq .

	$(eval MONITORID := $(shell curl -s localhost:8080/api/v0/monitors -X POST -H "X-Api-Key: ${XAPIKEY}" -H "Content-Type: application/json" -d '{ "type": "host", "name": "disk.aa-00.writes.delta", "memo": "This monitor is for Hatena Blog.", "duration": 3, "metric": "disk.aa-00.writes.delta", "operator": ">", "warning": 20000.0, "critical": 400000.0, "maxCheckAttempts": 3, "notificationInterval": 60, "scopes": [ "Hatena-Blog" ], "excludeScopes": [ "Hatena-Bookmark: db-master" ] }' | jq -r .id))
	curl -s localhost:8080/api/v0/monitors/${MONITORID} -H "X-Api-Key: ${XAPIKEY}" | jq .
	curl -s localhost:8080/api/v0/monitors/${MONITORID} -X DELETE -H "X-Api-Key: ${XAPIKEY}" | jq .

	curl -s localhost:8080/api/v0/invitations -X POST -H "X-Api-Key: ${XAPIKEY}" -H "Content-Type: application/json" -d '{"email": "example@example.com","authority": "viewer"}' | jq .
	curl -s localhost:8080/api/v0/invitations -H "X-Api-Key: ${XAPIKEY}" | jq .
	curl -s localhost:8080/api/v0/invitations/revoke -X POST -H "X-Api-Key: ${XAPIKEY}" -H "Content-Type: application/json" -d '{"email": "example@example.com"}' | jq .

	curl -s localhost:8080/api/v0/services -X POST -H "X-Api-Key: ${XAPIKEY}" -H "Content-Type: application/json" -d '{"name": "ExampleService", "memo": "This is an example."}' | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/metadata/foobar -X PUT -H "X-Api-Key: ${XAPIKEY}" -H "Content-Type: application/json" -d '{"message": "this is service metadata"}' | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/metadata/foobar -H "X-Api-Key: ${XAPIKEY}" | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/metadata -H "X-Api-Key: ${XAPIKEY}" | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/metadata/foobar -X DELETE -H "X-Api-Key: ${XAPIKEY}" | jq .
	curl -s localhost:8080/api/v0/services/ExampleService -X DELETE -H "X-Api-Key: ${XAPIKEY}" | jq .
	curl -s localhost:8080/api/v0/services -H "X-Api-Key: ${XAPIKEY}" | jq .
	curl -s localhost:8080/api/v0/services -X POST -H "X-Api-Key: ${XAPIKEY}" -H "Content-Type: application/json" -d '{"name": "ExampleService", "memo": "This is an example."}' | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/roles -X POST -H "X-Api-Key: ${XAPIKEY}" -H "Content-Type: application/json" -d '{"name": "ExampleRole", "memo": "This is an example."}' | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/roles -H "X-Api-Key: ${XAPIKEY}" | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/roles/ExampleRole/metadata/foobar -X PUT -H "X-Api-Key: ${XAPIKEY}" -H "Content-Type: application/json" -d '{"message": "this is role metadata"}' | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/roles/ExampleRole/metadata/foobar -H "X-Api-Key: ${XAPIKEY}" | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/roles/ExampleRole/metadata -H "X-Api-Key: ${XAPIKEY}" | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/roles/ExampleRole/metadata/foobar -X DELETE -H "X-Api-Key: ${XAPIKEY}" | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/roles/ExampleRole -X DELETE -H "X-Api-Key: ${XAPIKEY}" | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/tsdb -X POST -H "X-Api-Key: ${XAPIKEY}" -H "Content-Type: application/json" -d '[{"name":"hoge", "time": 1351700030, "value": 1.234},{"name":"foobar", "time": 1351700050, "value": 5.678}]' | jq .
	curl -s "localhost:8080/api/v0/services/ExampleService/metrics?name=hoge&from=1351700000&to=1351700100" -H "X-Api-Key: ${XAPIKEY}" | jq .
	curl -s localhost:8080/api/v0/services/ExampleService/metric-names -H "X-Api-Key: ${XAPIKEY}" | jq .
	curl -s localhost:8080/api/v0/services/ExampleService -X DELETE -H "X-Api-Key: ${XAPIKEY}" | jq .

	curl -s -v localhost:8080/api/v0/monitors -X POST -H "X-Api-Key: ${XAPIKEY}" -H "Content-Type: application/json" -d '{ "type": "host", "name": "loadavg1", "duration": 3, "metric": "loadavg1", "operator": ">", "warning": 3.0, "critical": 6.0 }' | jq .
	curl -s -v localhost:8080/api/v0/monitoring/checks/host-metric -H "X-Api-Key: ${XAPIKEY}" -H "Content-Type: application/json" } | jq .
#	curl -s -v localhost:8080/api/v0/alerts/0d76964bd7b/close -X POST -H "X-Api-Key: ${XAPIKEY}" -H "Content-Type: application/json" } -d '{ "reason": "manual" }' | jq .
