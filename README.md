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
