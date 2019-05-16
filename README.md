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
   - [ ] Database/Transaction
   - [ ] Amazon Timestream
 - Security
   - [ ] `X-Api-Key`

# Install

```
$ go get github.com/itsubaki/mackerel-api
$ mackerel-api
GIN_MODE=debug mackerel-api
[GIN-debug] [WARNING] Now Gin requires Go 1.6 or later and Go 1.7 will be required soon.

[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /                         --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func1 (3 handlers)
[GIN-debug] GET    /api/v0/services          --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func2 (3 handlers)
[GIN-debug] POST   /api/v0/services          --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func3 (3 handlers)
[GIN-debug] DELETE /api/v0/services/:serviceName --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func4 (3 handlers)
[GIN-debug] GET    /api/v0/services/:serviceName/metadata --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func5 (3 handlers)
[GIN-debug] GET    /api/v0/services/:serviceName/metadata/:namespace --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func6 (3 handlers)
[GIN-debug] PUT    /api/v0/services/:serviceName/metadata/:namespace --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func7 (3 handlers)
[GIN-debug] DELETE /api/v0/services/:serviceName/metadata/:namespace --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func8 (3 handlers)
[GIN-debug] GET    /api/v0/services/:serviceName/roles --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func9 (3 handlers)
[GIN-debug] POST   /api/v0/services/:serviceName/roles --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func10 (3 handlers)
[GIN-debug] DELETE /api/v0/services/:serviceName/roles/:roleName --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func11 (3 handlers)
[GIN-debug] GET    /api/v0/services/:serviceName/roles/:roleName/metadata --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func12 (3 handlers)
[GIN-debug] GET    /api/v0/services/:serviceName/roles/:roleName/metadata/:namespace --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func13 (3 handlers)
[GIN-debug] PUT    /api/v0/services/:serviceName/roles/:roleName/metadata/:namespace --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func14 (3 handlers)
[GIN-debug] DELETE /api/v0/services/:serviceName/roles/:roleName/metadata/:namespace --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func15 (3 handlers)
[GIN-debug] GET    /api/v0/services/:serviceName/metric-names --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func16 (3 handlers)
[GIN-debug] GET    /api/v0/services/:serviceName/metrics --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func17 (3 handlers)
[GIN-debug] POST   /api/v0/services/:serviceName/tsdb --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func18 (3 handlers)
[GIN-debug] GET    /api/v0/hosts             --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func19 (3 handlers)
[GIN-debug] POST   /api/v0/hosts             --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func20 (3 handlers)
[GIN-debug] GET    /api/v0/hosts/:hostId     --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func21 (3 handlers)
[GIN-debug] PUT    /api/v0/hosts/:hostId     --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func22 (3 handlers)
[GIN-debug] PUT    /api/v0/hosts/:hostId/role-fullnames --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func23 (3 handlers)
[GIN-debug] POST   /api/v0/hosts/:hostId/status --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func24 (3 handlers)
[GIN-debug] POST   /api/v0/hosts/:hostId/retire --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func25 (3 handlers)
[GIN-debug] GET    /api/v0/hosts/:hostId/metadata --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func26 (3 handlers)
[GIN-debug] GET    /api/v0/hosts/:hostId/metadata/:namespace --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func27 (3 handlers)
[GIN-debug] PUT    /api/v0/hosts/:hostId/metadata/:namespace --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func28 (3 handlers)
[GIN-debug] DELETE /api/v0/hosts/:hostId/metadata/:namespace --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func29 (3 handlers)
[GIN-debug] GET    /api/v0/hosts/:hostId/metric-names --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func30 (3 handlers)
[GIN-debug] GET    /api/v0/hosts/:hostId/metrics --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func31 (3 handlers)
[GIN-debug] GET    /api/v0/tsdb/latest       --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func32 (3 handlers)
[GIN-debug] POST   /api/v0/tsdb              --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func33 (3 handlers)
[GIN-debug] POST   /api/v0/graph-defs/create --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func34 (3 handlers)
[GIN-debug] POST   /api/v0/monitoring/checks/report --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func35 (3 handlers)
[GIN-debug] GET    /api/v0/alerts            --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func36 (3 handlers)
[GIN-debug] POST   /api/v0/alerts/:alertId/close --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func37 (3 handlers)
[GIN-debug] GET    /api/v0/invitations       --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func38 (3 handlers)
[GIN-debug] POST   /api/v0/invitations       --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func39 (3 handlers)
[GIN-debug] POST   /api/v0/invitations/revoke --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func40 (3 handlers)
[GIN-debug] GET    /api/v0/users             --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func41 (3 handlers)
[GIN-debug] DELETE /api/v0/users/:userId     --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func42 (3 handlers)
[GIN-debug] GET    /api/v0/org               --> github.com/itsubaki/mackerel-api/pkg/infrastructure.Default.func43 (3 handlers)
[GIN-debug] Listening and serving HTTP on :8080
[GIN] 2019/05/16 - 10:35:27 | 200 |    1.239153ms |             ::1 | POST     /api/v0/hosts
[GIN] 2019/05/16 - 10:35:27 | 200 |     176.587µs |             ::1 | GET      /api/v0/hosts/b9ec43060e8
[GIN] 2019/05/16 - 10:35:32 | 200 |     210.674µs |             ::1 | PUT      /api/v0/hosts/b9ec43060e8
[GIN] 2019/05/16 - 10:36:27 | 200 |     558.083µs |             ::1 | POST     /api/v0/tsdb
[GIN] 2019/05/16 - 10:36:49 | 200 |     129.602µs |             ::1 | POST     /api/v0/hosts/b9ec43060e8/status
[GIN] 2019/05/16 - 10:36:51 | 200 |     152.242µs |             ::1 | POST     /api/v0/hosts
[GIN] 2019/05/16 - 10:36:51 | 200 |     158.082µs |             ::1 | GET      /api/v0/org
[GIN] 2019/05/16 - 10:36:52 | 200 |     562.972µs |             ::1 | POST     /api/v0/services
[GIN] 2019/05/16 - 10:36:52 | 200 |     148.677µs |             ::1 | PUT      /api/v0/services/ExampleService/metadata/foobar
[GIN] 2019/05/16 - 10:36:52 | 200 |     121.124µs |             ::1 | GET      /api/v0/services/ExampleService/metadata/foobar
[GIN] 2019/05/16 - 10:36:52 | 200 |     200.926µs |             ::1 | GET      /api/v0/services/ExampleService/metadata
[GIN] 2019/05/16 - 10:36:52 | 200 |      75.446µs |             ::1 | DELETE   /api/v0/services/ExampleService/metadata/foobar
[GIN] 2019/05/16 - 10:36:52 | 200 |      61.226µs |             ::1 | DELETE   /api/v0/services/ExampleService
[GIN] 2019/05/16 - 10:36:52 | 200 |     128.698µs |             ::1 | GET      /api/v0/services
[GIN] 2019/05/16 - 10:36:52 | 200 |      537.87µs |             ::1 | POST     /api/v0/services
[GIN] 2019/05/16 - 10:36:52 | 200 |     116.294µs |             ::1 | POST     /api/v0/services/ExampleService/roles
[GIN] 2019/05/16 - 10:36:52 | 200 |      101.38µs |             ::1 | GET      /api/v0/services/ExampleService/roles
[GIN] 2019/05/16 - 10:36:52 | 200 |     101.834µs |             ::1 | PUT      /api/v0/services/ExampleService/roles/ExampleRole/metadata/foobar
[GIN] 2019/05/16 - 10:36:52 | 200 |      65.456µs |             ::1 | GET      /api/v0/services/ExampleService/roles/ExampleRole/metadata/foobar
[GIN] 2019/05/16 - 10:36:52 | 200 |     183.414µs |             ::1 | GET      /api/v0/services/ExampleService/roles/ExampleRole/metadata
[GIN] 2019/05/16 - 10:36:52 | 200 |     105.286µs |             ::1 | DELETE   /api/v0/services/ExampleService/roles/ExampleRole/metadata/foobar
[GIN] 2019/05/16 - 10:36:52 | 200 |      51.858µs |             ::1 | DELETE   /api/v0/services/ExampleService/roles/ExampleRole
[GIN] 2019/05/16 - 10:36:52 | 200 |     213.137µs |             ::1 | POST     /api/v0/services/ExampleService/tsdb
[GIN] 2019/05/16 - 10:36:52 | 200 |     133.446µs |             ::1 | GET      /api/v0/services/ExampleService/metrics?name=hoge&from=1351700000&to=1351700100
[GIN] 2019/05/16 - 10:36:52 | 200 |     195.577µs |             ::1 | GET      /api/v0/services/ExampleService/metric-names
[GIN] 2019/05/16 - 10:36:52 | 200 |       239.8µs |             ::1 | DELETE   /api/v0/services/ExampleService
[GIN] 2019/05/16 - 10:36:52 | 200 |     131.406µs |             ::1 | PUT      /api/v0/hosts/b93cf8d53fc
[GIN] 2019/05/16 - 10:36:52 | 200 |     141.595µs |             ::1 | POST     /api/v0/hosts/b93cf8d53fc/status
[GIN] 2019/05/16 - 10:36:52 | 200 |      98.827µs |             ::1 | GET      /api/v0/hosts/b93cf8d53fc
[GIN] 2019/05/16 - 10:36:52 | 200 |     100.452µs |             ::1 | PUT      /api/v0/hosts/b93cf8d53fc/role-fullnames
[GIN] 2019/05/16 - 10:36:52 | 200 |      94.374µs |             ::1 | PUT      /api/v0/hosts/b93cf8d53fc/metadata/foobar
[GIN] 2019/05/16 - 10:36:52 | 200 |      81.347µs |             ::1 | GET      /api/v0/hosts/b93cf8d53fc/metadata/foobar
[GIN] 2019/05/16 - 10:36:53 | 200 |     166.287µs |             ::1 | GET      /api/v0/hosts/b93cf8d53fc/metadata
[GIN] 2019/05/16 - 10:36:53 | 200 |     103.835µs |             ::1 | DELETE   /api/v0/hosts/b93cf8d53fc/metadata/foobar
[GIN] 2019/05/16 - 10:36:53 | 200 |     101.578µs |             ::1 | POST     /api/v0/tsdb
[GIN] 2019/05/16 - 10:36:53 | 200 |     410.567µs |             ::1 | GET      /api/v0/tsdb/latest?hostId=b93cf8d53fc&name=foobar&name=hoge&name=piyo
[GIN] 2019/05/16 - 10:36:53 | 200 |      165.86µs |             ::1 | GET      /api/v0/hosts/b93cf8d53fc/metrics?name=hoge&from=1351700000&to=1351700100
[GIN] 2019/05/16 - 10:36:53 | 200 |     128.783µs |             ::1 | GET      /api/v0/hosts/b93cf8d53fc/metric-names
[GIN] 2019/05/16 - 10:36:53 | 200 |     311.806µs |             ::1 | GET      /api/v0/hosts
[GIN] 2019/05/16 - 10:36:53 | 200 |     233.662µs |             ::1 | POST     /api/v0/hosts/b93cf8d53fc/retire
[GIN] 2019/05/16 - 10:36:53 | 200 |      171.18µs |             ::1 | POST     /api/v0/invitations
[GIN] 2019/05/16 - 10:36:53 | 200 |      83.982µs |             ::1 | GET      /api/v0/invitations
[GIN] 2019/05/16 - 10:36:53 | 200 |     146.102µs |             ::1 | POST     /api/v0/invitations/revoke
```

```
$ go get github.com/mackerelio/mackerel-agent/mackerel-agent
$ mackerel-agent -conf /usr/local/etc/mackerel-agent.conf -apibase=http://localhost:8080
2019/05/07 12:06:09 main.go:171: INFO <main> Starting mackerel-agent version:0.59.0, rev:, apibase:http://localhost:8080
2019/05/07 12:06:17 command.go:90: DEBUG <command> Registering new host on mackerel...
2019/05/07 12:06:17 api.go:306: DEBUG <api> POST /api/v0/hosts {"name":"host001","meta":{"agent-version":"0.59.0","agent-name":"mackerel-agent/0.59.0 (Revision )","cpu":[{"cores":"2","family":"6","model":"142","model_name":"Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz","stepping":"9","vendor_id":"GenuineIntel"},{"cores":"2","family":"6","model":"142","model_name":"Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz","stepping":"9","vendor_id":"GenuineIntel"},{"cores":"2","family":"6","model":"142","model_name":"Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz","stepping":"9","vendor_id":"GenuineIntel"},{"cores":"2","family":"6","model":"142","model_name":"Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz","stepping":"9","vendor_id":"GenuineIntel"}],"filesystem":{"/dev/disk1s1":{"kb_available":396947024,"kb_size":488245288,"kb_used":90615784,"mount":"/","percent_used":"19%"}},"kernel":{"machine":"x86_64","name":"Darwin","os":"Darwin","platform_name":"Mac OS X","platform_version":"10.13.6","release":"17.7.0","version":"Darwin Kernel Version 17.7.0: Thu Dec 20 21:47:19 PST 2018; root:xnu-4570.71.22~1/RELEASE_X86_64"},"memory":{"total":"16777216kB"}},"interfaces":[{"name":"en0","ipv4Addresses":["172.16.21.53"],"ipv6Addresses":null,"macAddress":"f0:18:98:11:45:47","defaultGateway":"172.16.16.1"}],"roleFullnames":[],"checks":[]}
2019/05/07 12:06:17 api.go:319: DEBUG <api> POST /api/v0/hosts status="200 OK"
2019/05/07 12:06:17 command.go:711: INFO <command> Start: apibase = http://localhost:8080, hostName = host001, hostID = 2f822e79-a61f-4696-8402-a7b5f4c21571
2019/05/07 12:06:17 command.go:211: DEBUG <command> wait 8 seconds before initial posting.
2019/05/07 12:06:17 command.go:605: DEBUG <command> Updating host specs...
2019/05/07 12:06:25 api.go:306: DEBUG <api> PUT /api/v0/hosts/2f822e79-a61f-4696-8402-a7b5f4c21571 {"name":"host001","meta":{"agent-version":"0.59.0","agent-name":"mackerel-agent/0.59.0 (Revision )","cpu":[{"cores":"2","family":"6","model":"142","model_name":"Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz","stepping":"9","vendor_id":"GenuineIntel"},{"cores":"2","family":"6","model":"142","model_name":"Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz","stepping":"9","vendor_id":"GenuineIntel"},{"cores":"2","family":"6","model":"142","model_name":"Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz","stepping":"9","vendor_id":"GenuineIntel"},{"cores":"2","family":"6","model":"142","model_name":"Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz","stepping":"9","vendor_id":"GenuineIntel"}],"filesystem":{"/dev/disk1s1":{"kb_available":396947020,"kb_size":488245288,"kb_used":90615788,"mount":"/","percent_used":"19%"}},"kernel":{"machine":"x86_64","name":"Darwin","os":"Darwin","platform_name":"Mac OS X","platform_version":"10.13.6","release":"17.7.0","version":"Darwin Kernel Version 17.7.0: Thu Dec 20 21:47:19 PST 2018; root:xnu-4570.71.22~1/RELEASE_X86_64"},"memory":{"total":"16777216kB"}},"interfaces":[{"name":"en0","ipv4Addresses":["172.16.21.53"],"ipv6Addresses":null,"macAddress":"f0:18:98:11:45:47","defaultGateway":"172.16.16.1"}],"roleFullnames":[],"checks":[]}
2019/05/07 12:06:25 api.go:319: DEBUG <api> PUT /api/v0/hosts/2f822e79-a61f-4696-8402-a7b5f4c21571 status="200 OK"
2019/05/07 12:06:25 command.go:618: DEBUG <command> Host specs sent.
2019/05/07 12:07:17 command.go:402: DEBUG <command> Enqueuing task to post metrics.
2019/05/07 12:07:17 command.go:306: DEBUG <command> Sleep 0 seconds before posting.
2019/05/07 12:07:17 api.go:306: DEBUG <api> POST /api/v0/tsdb [{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"interface.en2.txBytes.delta","time":1557198377,"value":0},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"interface.awdl0.rxBytes.delta","time":1557198377,"value":0},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"memory.swap_free","time":1557198377,"value":1408237568},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"cpu.idle.percentage","time":1557198377,"value":74},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"interface.en1.rxBytes.delta","time":1557198377,"value":0},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"interface.gif0.rxBytes.delta","time":1557198377,"value":0},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"memory.used","time":1557198377,"value":10787450880},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"interface.awdl0.txBytes.delta","time":1557198377,"value":0},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"interface.XHC20.rxBytes.delta","time":1557198377,"value":0},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"interface.gif0.txBytes.delta","time":1557198377,"value":0},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"interface.en0.txBytes.delta","time":1557198377,"value":875.8833333333333},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"loadavg1","time":1557198377,"value":10.71484375},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"cpu.system.percentage","time":1557198377,"value":10},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"interface.XHC20.txBytes.delta","time":1557198377,"value":0},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"interface.bridg.rxBytes.delta","time":1557198377,"value":0},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"memory.free","time":1557198377,"value":1561616384},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"memory.active","time":1557198377,"value":6444646400},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"filesystem.disk1s1.size","time":1557198377,"value":499264315392},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"interface.bridg.txBytes.delta","time":1557198377,"value":0},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"memory.inactive","time":1557198377,"value":5187747840},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"cpu.user.percentage","time":1557198377,"value":16},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"interface.utun0.txBytes.delta","time":1557198377,"value":0},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"interface.en0.rxBytes.delta","time":1557198377,"value":31699.45},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"interface.stf0.txBytes.delta","time":1557198377,"value":0},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"interface.p2p0.txBytes.delta","time":1557198377,"value":0},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"interface.en2.rxBytes.delta","time":1557198377,"value":0},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"loadavg5","time":1557198377,"value":14.10546875},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"memory.swap_total","time":1557198377,"value":2147483648},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"memory.total","time":1557198377,"value":17178972160},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"filesystem.disk1s1.used","time":1557198377,"value":92790566912},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"interface.p2p0.rxBytes.delta","time":1557198377,"value":0},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"interface.en1.txBytes.delta","time":1557198377,"value":0},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"interface.XHC0.rxBytes.delta","time":1557198377,"value":0},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"loadavg15","time":1557198377,"value":10.515625},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"memory.cached","time":1557198377,"value":4829904896},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"interface.stf0.rxBytes.delta","time":1557198377,"value":0},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"interface.XHC0.txBytes.delta","time":1557198377,"value":0},{"hostId":"2f822e79-a61f-4696-8402-a7b5f4c21571","name":"interface.utun0.rxBytes.delta","time":1557198377,"value":0}]
2019/05/07 12:07:17 api.go:319: DEBUG <api> POST /api/v0/tsdb status="200 OK"
2019/05/07 12:07:17 command.go:346: DEBUG <command> Posting metrics succeeded.
2019/05/07 12:08:00 command.go:402: DEBUG <command> Enqueuing task to post metrics.
2019/05/07 12:08:00 command.go:306: DEBUG <command> Sleep 17 seconds before posting.
```
