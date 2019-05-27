# mackerel-api
mackerel compatible monitoring service

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
   - [x] Monitors
   - [x] Alerts
   - [ ] Notification channels
   - [ ] Notification groups
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
[GIN] 2019/05/24 - 23:52:12 | 200 |   65.133982ms |             ::1 | POST     /api/v0/hosts
[GIN] 2019/05/24 - 23:52:12 | 200 |   23.998452ms |             ::1 | GET      /api/v0/hosts/0965d1deb93
[GIN] 2019/05/24 - 23:52:14 | 200 |   55.856843ms |             ::1 | PUT      /api/v0/hosts/0965d1deb93
[GIN] 2019/05/24 - 23:53:12 | 200 |  275.695763ms |             ::1 | POST     /api/v0/tsdb
[GIN] 2019/05/24 - 23:53:26 | 200 |   60.847875ms |             ::1 | POST     /api/v0/monitoring/checks/report
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

# Memo

```
mysql> explain select host_id, avg(value) from (select host_id, name, time, value, rank() over(partition by host_id order by time desc) as rnk from host_metric_values where name='memory.used') as ordtable where rnk < 10 group by host_id;
+----+-------------+--------------------+------------+------+---------------+------+---------+------+------+----------+------------------------------+
| id | select_type | table              | partitions | type | possible_keys | key  | key_len | ref  | rows | filtered | Extra                        |
+----+-------------+--------------------+------------+------+---------------+------+---------+------+------+----------+------------------------------+
|  1 | PRIMARY     | <derived2>         | NULL       | ALL  | NULL          | NULL | NULL    | NULL |  272 |    33.33 | Using where; Using temporary |
|  2 | DERIVED     | host_metric_values | NULL       | ALL  | NULL          | NULL | NULL    | NULL | 2720 |    10.00 | Using where; Using filesort  |
+----+-------------+--------------------+------------+------+---------------+------+---------+------+------+----------+------------------------------+

mysql> select host_id, avg(value) from (select host_id, name, time, value, rank() over(partition by host_id order by time desc) as rnk from host_metric_values where name='memory.used') as rnktable where rnk < 3 group by host_id;
+-------------+------------+
| host_id     | avg(value) |
+-------------+------------+
| 8b2e92d3402 | 7213408256 |
| af2dd8cee1e | 7364827136 |
| f49131deaec | 7152691200 |
+-------------+------------+
```

```
mysql> explain select host_id, name, time, value, rank() over(partition by host_id order by time desc) as rnk from host_metric_values where name='memory.used';
+----+-------------+--------------------+------------+------+---------------+------+---------+------+------+----------+-----------------------------+
| id | select_type | table              | partitions | type | possible_keys | key  | key_len | ref  | rows | filtered | Extra                       |
+----+-------------+--------------------+------------+------+---------------+------+---------+------+------+----------+-----------------------------+
|  1 | SIMPLE      | host_metric_values | NULL       | ALL  | NULL          | NULL | NULL    | NULL | 2856 |    10.00 | Using where; Using filesort |
+----+-------------+--------------------+------------+------+---------------+------+---------+------+------+----------+-----------------------------+
1 row in set, 2 warnings (0.01 sec)


mysql> select host_id, name, time, value, rank() over(partition by host_id order by time desc) as rnk from host_metric_values where name='memory.used';
+-------------+-------------+------------+------------+-----+
| host_id     | name        | time       | value      | rnk |
+-------------+-------------+------------+------------+-----+
| 8b2e92d3402 | memory.used | 1558786260 | 7212445696 |   1 |
| 8b2e92d3402 | memory.used | 1558786200 | 7214370816 |   2 |
| 8b2e92d3402 | memory.used | 1558786140 | 7210688512 |   3 |
| 8b2e92d3402 | memory.used | 1558786108 | 7227314176 |   4 |
| af2dd8cee1e | memory.used | 1558822500 | 7365353472 |   1 |
| af2dd8cee1e | memory.used | 1558822440 | 7364300800 |   2 |
| af2dd8cee1e | memory.used | 1558822380 | 7362813952 |   3 |
| af2dd8cee1e | memory.used | 1558822320 | 7362846720 |   4 |
| af2dd8cee1e | memory.used | 1558822260 | 7366426624 |   5 |
| af2dd8cee1e | memory.used | 1558822200 | 7362777088 |   6 |
| af2dd8cee1e | memory.used | 1558822140 | 7363547136 |   7 |
| af2dd8cee1e | memory.used | 1558822080 | 7363190784 |   8 |
| af2dd8cee1e | memory.used | 1558822020 | 7363055616 |   9 |
| af2dd8cee1e | memory.used | 1558821960 | 7362027520 |  10 |
| af2dd8cee1e | memory.used | 1558821900 | 7358636032 |  11 |
| af2dd8cee1e | memory.used | 1558821840 | 7367450624 |  12 |
| af2dd8cee1e | memory.used | 1558817580 | 7355977728 |  13 |
| af2dd8cee1e | memory.used | 1558817520 | 7356424192 |  14 |
| af2dd8cee1e | memory.used | 1558817460 | 7357734912 |  15 |
| af2dd8cee1e | memory.used | 1558817400 | 7356313600 |  16 |
| af2dd8cee1e | memory.used | 1558817340 | 7355748352 |  17 |
| af2dd8cee1e | memory.used | 1558817280 | 7355252736 |  18 |
| af2dd8cee1e | memory.used | 1558817220 | 7354626048 |  19 |
| af2dd8cee1e | memory.used | 1558817160 | 7354888192 |  20 |
| af2dd8cee1e | memory.used | 1558810680 | 7288905728 |  21 |
| af2dd8cee1e | memory.used | 1558810620 | 7288299520 |  22 |
| af2dd8cee1e | memory.used | 1558810560 | 7287832576 |  23 |
| af2dd8cee1e | memory.used | 1558810500 | 7287054336 |  24 |
| af2dd8cee1e | memory.used | 1558810440 | 7298199552 |  25 |
| af2dd8cee1e | memory.used | 1558810380 | 7292084224 |  26 |
| af2dd8cee1e | memory.used | 1558804860 | 7421829120 |  27 |
| af2dd8cee1e | memory.used | 1558804800 | 7418277888 |  28 |
| af2dd8cee1e | memory.used | 1558804740 | 7419596800 |  29 |
| af2dd8cee1e | memory.used | 1558803420 | 7295332352 |  30 |
| af2dd8cee1e | memory.used | 1558803360 | 7296831488 |  31 |
| af2dd8cee1e | memory.used | 1558802460 | 6947278848 |  32 |
| af2dd8cee1e | memory.used | 1558798260 | 7124975616 |  33 |
| af2dd8cee1e | memory.used | 1558797840 | 7103434752 |  34 |
| af2dd8cee1e | memory.used | 1558797652 | 6997856256 |  35 |
| af2dd8cee1e | memory.used | 1558786500 | 7264423936 |  36 |
| f49131deaec | memory.used | 1558845120 | 7210176512 |   1 |
| f49131deaec | memory.used | 1558845060 | 7122694144 |   2 |
| f49131deaec | memory.used | 1558845000 | 7118041088 |   3 |
| f49131deaec | memory.used | 1558844940 | 7187341312 |   4 |
| f49131deaec | memory.used | 1558844880 | 7024566272 |   5 |
| f49131deaec | memory.used | 1558844820 | 7170715648 |   6 |
| f49131deaec | memory.used | 1558844760 | 7153893376 |   7 |
| f49131deaec | memory.used | 1558844700 | 7142801408 |   8 |
| f49131deaec | memory.used | 1558844640 | 7428108288 |   9 |
| f49131deaec | memory.used | 1558844580 | 7422115840 |  10 |
| f49131deaec | memory.used | 1558844520 | 7417364480 |  11 |
| f49131deaec | memory.used | 1558844460 | 7368708096 |  12 |
| f49131deaec | memory.used | 1558844040 | 7450918912 |  13 |
| f49131deaec | memory.used | 1558843980 | 7244079104 |  14 |
| f49131deaec | memory.used | 1558843920 | 7305256960 |  15 |
| f49131deaec | memory.used | 1558843860 | 7276478464 |  16 |
| f49131deaec | memory.used | 1558843800 | 7102087168 |  17 |
| f49131deaec | memory.used | 1558843740 | 7255044096 |  18 |
| f49131deaec | memory.used | 1558843680 | 7130230784 |  19 |
| f49131deaec | memory.used | 1558843620 | 7409401856 |  20 |
| f49131deaec | memory.used | 1558843560 | 7341084672 |  21 |
| f49131deaec | memory.used | 1558843500 | 7363649536 |  22 |
| f49131deaec | memory.used | 1558843440 | 7386136576 |  23 |
| f49131deaec | memory.used | 1558843380 | 7332016128 |  24 |
| f49131deaec | memory.used | 1558843320 | 7489540096 |  25 |
| f49131deaec | memory.used | 1558843260 | 7227633664 |  26 |
| f49131deaec | memory.used | 1558843240 | 7440502784 |  27 |
| f49131deaec | memory.used | 1558841280 | 7389683712 |  28 |
| f49131deaec | memory.used | 1558841220 | 7332749312 |  29 |
| f49131deaec | memory.used | 1558841160 | 7418433536 |  30 |
| f49131deaec | memory.used | 1558841100 | 7416475648 |  31 |
| f49131deaec | memory.used | 1558841040 | 7394238464 |  32 |
| f49131deaec | memory.used | 1558840980 | 7314087936 |  33 |
| f49131deaec | memory.used | 1558840920 | 7366483968 |  34 |
| f49131deaec | memory.used | 1558840860 | 7587901440 |  35 |
| f49131deaec | memory.used | 1558840800 | 7494623232 |  36 |
| f49131deaec | memory.used | 1558840740 | 7484448768 |  37 |
| f49131deaec | memory.used | 1558840680 | 7471005696 |  38 |
| f49131deaec | memory.used | 1558840620 | 7483002880 |  39 |
| f49131deaec | memory.used | 1558840560 | 7459889152 |  40 |
| f49131deaec | memory.used | 1558840500 | 7456268288 |  41 |
| f49131deaec | memory.used | 1558840440 | 7418056704 |  42 |
| f49131deaec | memory.used | 1558840433 | 7434199040 |  43 |
+-------------+-------------+------------+------------+-----+
83 rows in set (0.01 sec)
```


```
mysql> explain select host_id, name, time, value, rank() over(order by time desc) as rnk from host_metric_values where name='memory.used' and host_id='f49131deaec' limit 3;
+----+-------------+--------------------+------------+------+---------------+---------+---------+-------------+------+----------+----------------+
| id | select_type | table              | partitions | type | possible_keys | key     | key_len | ref         | rows | filtered | Extra          |
+----+-------------+--------------------+------------+------+---------------+---------+---------+-------------+------+----------+----------------+
|  1 | SIMPLE      | host_metric_values | NULL       | ref  | PRIMARY       | PRIMARY | 580     | const,const |   54 |   100.00 | Using filesort |
+----+-------------+--------------------+------------+------+---------------+---------+---------+-------------+------+----------+----------------+
1 row in set, 2 warnings (0.00 sec)

mysql> select host_id, name, time, value, rank() over(order by time desc) as rnk from host_metric_values where name='memory.used' and host_id='f49131deaec' limit 3;
+-------------+-------------+------------+------------+-----+
| host_id     | name        | time       | value      | rnk |
+-------------+-------------+------------+------------+-----+
| f49131deaec | memory.used | 1558845780 | 7224537088 |   1 |
| f49131deaec | memory.used | 1558845720 | 7271268352 |   2 |
| f49131deaec | memory.used | 1558845660 | 7191089152 |   3 |
+-------------+-------------+------------+------------+-----+
3 rows in set (0.00 sec)

mysql> explain select host_id, avg(value) from(select host_id, name, time, value, rank() over(order by time desc) as rnk from host_metric_values where name='memory.used' and host_id='f49131deaec' limit 3) as rnktable group by host_id;
+----+-------------+--------------------+------------+------+---------------+---------+---------+-------------+------+----------+-----------------+
| id | select_type | table              | partitions | type | possible_keys | key     | key_len | ref         | rows | filtered | Extra           |
+----+-------------+--------------------+------------+------+---------------+---------+---------+-------------+------+----------+-----------------+
|  1 | PRIMARY     | <derived2>         | NULL       | ALL  | NULL          | NULL    | NULL    | NULL        |    3 |   100.00 | Using temporary |
|  2 | DERIVED     | host_metric_values | NULL       | ref  | PRIMARY       | PRIMARY | 580     | const,const |   56 |   100.00 | Using filesort  |
+----+-------------+--------------------+------------+------+---------------+---------+---------+-------------+------+----------+-----------------+
2 rows in set, 2 warnings (0.00 sec)

mysql> select host_id, avg(value) from(select host_id, name, time, value, rank() over(order by time desc) as rnk from host_metric_values where name='memory.used' and host_id='f49131deaec' limit 3) as rnktable group by host_id;
+-------------+-------------------+
| host_id     | avg(value)        |
+-------------+-------------------+
| f49131deaec | 7231591765.333333 |
+-------------+-------------------+
1 row in set (0.00 sec)
```