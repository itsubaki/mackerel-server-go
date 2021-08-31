# mackerel-server-go

[![Go Report Card](https://goreportcard.com/badge/github.com/itsubaki/mackerel-server-go?style=flat-square)](https://goreportcard.com/report/github.com/itsubaki/mackerel-server-go)
[![tests](https://github.com/itsubaki/mackerel-server-go/actions/workflows/tests.yml/badge.svg?branch=main)](https://github.com/itsubaki/mackerel-server-go/actions/workflows/tests.yml)
[![codecov](https://codecov.io/gh/itsubaki/mackerel-server-go/branch/main/graph/badge.svg?token=LI1C1H4D0P)](https://codecov.io/gh/itsubaki/mackerel-server-go)

- Mackerel API Server Clone written in golang

# Install

```
git clone https://github.com/itsubaki/mackerel-server-go
```

# Run

```
$ make runmysql
docker pull mysql
Using default tag: latest
latest: Pulling from library/mysql
Digest: sha256:711df5b93720801b3a727864aba18c2ae46c07f9fe33d5ce9c1f5cbc2c035101
Status: Image is up to date for mysql:latest

docker run --name mysqld -e MYSQL_ROOT_PASSWORD=secret -p 3306:3306 -d mysql
6521954a39afbd1aad14729da5ee0b898ad6f5721e0e71b3f5bd8f8746fdf7af
# mysql -h127.0.0.1 -P3306 -psecret -uroot mackerel
```

```
$ make runserver
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
mackerel-agent -conf mackerel-agent.conf -apibase=http://localhost:8080
2019/05/24 23:52:07 main.go:171: INFO <main> Starting mackerel-agent version:0.59.0, rev:, apibase:http://localhost:8080
2019/05/24 23:52:12 command.go:91: DEBUG <command> Registering new host on mackerel...
```

# Run with Docker Compose

```
$ make up
docker build -t mackerel-server-go .
[+] Building 57.3s (15/15) FINISHED
docker-compose up
Recreating mackerel-server-go_app_1 ... done
Starting mackerel-server-go_mysql_1 ... done
...
app_1    | 2021/07/23 11:38:38 main.go:54: http server listen and serve
app_1    | [GIN] 2021/07/23 - 11:40:13 | 200 |   10.606767ms |      172.18.0.1 | GET      "/api/v0/org"
app_1    | [GIN] 2021/07/23 - 11:40:13 | 200 |   24.198934ms |      172.18.0.1 | POST     "/api/v0/hosts"
app_1    | [GIN] 2021/07/23 - 11:40:13 | 200 |    8.484944ms |      172.18.0.1 | GET      "/api/v0/hosts"
app_1    | [GIN] 2021/07/23 - 11:40:13 | 200 |   18.681121ms |      172.18.0.1 | GET      "/api/v0/services"
app_1    | [GIN] 2021/07/23 - 11:40:13 | 200 |    5.317676ms |      172.18.0.1 | GET      "/api/v0/alerts"
```
