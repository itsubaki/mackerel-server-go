# mackerel-api
mackerel compatible monitoring service

[API](https://mackerel.io/api-docs/)

[API(ja)](https://mackerel.io/ja/api-docs/)

# TODO

 - API
   - [x] Services
   - [x] Hosts
   - [x] Host Metrics
   - [x] Service Metrics
   - [x] Check Monitoring
   - [x] Metadata
   - [x] Monitors
   - [ ] Notification channels
   - [ ] Notification groups
   - [x] Alerts
   - [ ] Dashboards
   - [x] Users
   - [x] Invitations
   - [x] Organizations
 - Persistence
   - [x] Memory
   - [x] Database/Transaction
   - [ ] Amazon Timestream
 - Security
   - [x] `X-Api-Key`

# Install

```
$ go get github.com/itsubaki/mackerel-api
$ go get github.com/mackerelio/mackerel-agent/mackerel-agent
```

# Run

```
$ make runmysql
set -x
docker pull mysql
Using default tag: latest
latest: Pulling from library/mysql
Digest: sha256:711df5b93720801b3a727864aba18c2ae46c07f9fe33d5ce9c1f5cbc2c035101
Status: Image is up to date for mysql:latest
docker stop mysqld
mysqld
docker rm mysqld
mysqld
docker run --name mysqld -e MYSQL_ROOT_PASSWORD=secret -p 3307:3306 -d mysql
6521954a39afbd1aad14729da5ee0b898ad6f5721e0e71b3f5bd8f8746fdf7af
# mysql -h127.0.0.1 -P3307 -psecret -uroot mackerel
```

```
$ make
set -x
rm  ${GOPATH}/bin/mackerel-api
go install
GIN_MODE=release mackerel-api
```

```
$ make runclient
set -x
GO111MODULE=off go get -d github.com/mackerelio/go-check-plugins
cd  ${GOPATH}/src/github.com/mackerelio/go-check-plugins/check-tcp; go install
cp mackerel-agent.conf /usr/local/etc/mackerel-agent.conf
rm ~/Library/mackerel-agent/id
mackerel-agent -conf /usr/local/etc/mackerel-agent.conf -apibase=http://localhost:8080
2019/05/24 23:52:07 main.go:171: INFO <main> Starting mackerel-agent version:0.59.0, rev:, apibase:http://localhost:8080
2019/05/24 23:52:12 command.go:91: DEBUG <command> Registering new host on mackerel...
```