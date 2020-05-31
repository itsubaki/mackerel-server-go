# mackerel-server-go
Mackerel API Server in Go

[API](https://mackerel.io/api-docs/)

[API(ja)](https://mackerel.io/ja/api-docs/)


# TODO

 - API
   - [x] Hosts
   - [x] Host Metrics
   - [x] Services
   - [x] Service Metrics
   - [x] Check Monitoring
   - [x] Metadata
   - [ ] Monitors
   - [x] Alerts
   - [x] Notification Channels
   - [x] Notification Groups
   - [ ] Dashboards
   - [x] Graph Annotations
   - [ ] Users
   - [x] Invitations
   - [x] Organizations
   - [ ] Downtime
 - Persistence
   - [x] Database/Transaction
   - [ ] ORM
 - Security
   - [x] `X-Api-Key`

# Install

```
$ go get github.com/itsubaki/mackerel-server-go
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
docker run --name mysqld -e MYSQL_ROOT_PASSWORD=secret -p 3306:3306 -d mysql
6521954a39afbd1aad14729da5ee0b898ad6f5721e0e71b3f5bd8f8746fdf7af
# mysql -h127.0.0.1 -P3306 -psecret -uroot mackerel
```

```
$ make runserver
set -x
rm  ${GOPATH}/bin/mackerel-server-go
go install
GIN_MODE=debug mackerel-server-go
[GIN] 2019/05/24 - 23:52:12 | 200 |   65.133982ms |             ::1 | POST     /api/v0/hosts
[GIN] 2019/05/24 - 23:52:12 | 200 |   23.998452ms |             ::1 | GET      /api/v0/hosts/0965d1deb93
[GIN] 2019/05/24 - 23:52:14 | 200 |   55.856843ms |             ::1 | PUT      /api/v0/hosts/0965d1deb93
[GIN] 2019/05/24 - 23:53:12 | 200 |  275.695763ms |             ::1 | POST     /api/v0/tsdb
[GIN] 2019/05/24 - 23:53:26 | 200 |   60.847875ms |             ::1 | POST     /api/v0/monitoring/checks/report
```

```
$ make runclient
set -x
go get -d github.com/mackerelio/go-check-plugins
cd ${GOPATH}/src/github.com/mackerelio/go-check-plugins/check-tcp; go install
cp mackerel-agent.conf /usr/local/etc/mackerel-agent.conf
rm ~/Library/mackerel-agent/id
mackerel-agent -conf /usr/local/etc/mackerel-agent.conf -apibase=http://localhost:8080
2019/05/24 23:52:07 main.go:171: INFO <main> Starting mackerel-agent version:0.59.0, rev:, apibase:http://localhost:8080
2019/05/24 23:52:12 command.go:91: DEBUG <command> Registering new host on mackerel...
```

# Cucumber(godog)

```
$ godog features/org.feature 
Feature:
  In order to know org name
  As an API user
  I need to be able to request org

  Scenario: should get org                                                             # features/org.feature:6
    Given I set "X-Api-Key" header with "2684d06cfedbee8499f326037bb6fb7e8c22e73b16bb" # feature_test.go:68 -> *apiFeature
    When I send "GET" request to "/api/v0/org"                                         # feature_test.go:78 -> *apiFeature
    Then the response code should be 200                                               # feature_test.go:86 -> *apiFeature
    Then the response should match json:                                               # feature_test.go:94 -> *apiFeature
      """
      {
        "name": "mackerel"
      }
      """

  Scenario: forbidden                          # features/org.feature:17
    When I send "GET" request to "/api/v0/org" # feature_test.go:78 -> *apiFeature
    Then the response code should be 403       # feature_test.go:86 -> *apiFeature

2 scenarios (2 passed)
6 steps (6 passed)
1.544887852s
```
