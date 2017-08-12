# Elasticsearch + packetbeat 完整测试记录

## 实验步骤

- 删除 elasticsearch 的数据目录（清除所有数据）；
- 启动 elasticsearch ；
- 通过 Dev Tools 进行查询；
- 通过 `packetbeat` 导入 mysql 数据；
- 再次通过 Dev Tools 进行查询；
- 执行 `import_dashboards` 导入相关内容；
- 将 `packetbeat-*` 选成默认的 index pattern ；
- 调整选择器范围以显示导入的 mysql 数据；
- 基于 ingest node pipeline 进行数据解析；


----------


- 删除 elasticsearch 的数据目录（清除所有数据）

```
➜  ~ brew info elasticsearch
...
Data:    /usr/local/var/elasticsearch/elasticsearch_sunfei/
Logs:    /usr/local/var/log/elasticsearch/elasticsearch_sunfei.log
Plugins: /usr/local/opt/elasticsearch/libexec/plugins/
Config:  /usr/local/etc/elasticsearch/
plugin script: /usr/local/opt/elasticsearch/libexec/bin/elasticsearch-plugin
...
➜  ~ cd /usr/local/var/elasticsearch
➜  elasticsearch rm -rf ./*
```

- 启动 elasticsearch

```
➜  ~ elasticsearch -Epath.conf=/usr/local/etc/elasticsearch
[2017-05-25T14:53:11,529][INFO ][o.e.n.Node               ] [node-1] initializing ...
[2017-05-25T14:53:11,627][INFO ][o.e.e.NodeEnvironment    ] [node-1] using [1] data paths, mounts [[/ (/dev/disk1)]], net usable_space [16.2gb], net total_space [111.8gb], spins? [unknown], types [hfs]
[2017-05-25T14:53:11,627][INFO ][o.e.e.NodeEnvironment    ] [node-1] heap size [1.9gb], compressed ordinary object pointers [true]
[2017-05-25T14:53:11,628][INFO ][o.e.n.Node               ] [node-1] node name [node-1], node ID [TLbN2l8MTBGD_TmNuqIvow]
[2017-05-25T14:53:11,629][INFO ][o.e.n.Node               ] [node-1] version[5.3.2], pid[23098], build[3068195/2017-04-24T16:15:59.481Z], OS[Mac OS X/10.11.6/x86_64], JVM[Oracle Corporation/Java HotSpot(TM) 64-Bit Server VM/1.8.0_112/25.112-b16]
[2017-05-25T14:53:13,833][INFO ][o.e.p.PluginsService     ] [node-1] loaded module [aggs-matrix-stats]
[2017-05-25T14:53:13,834][INFO ][o.e.p.PluginsService     ] [node-1] loaded module [ingest-common]
[2017-05-25T14:53:13,834][INFO ][o.e.p.PluginsService     ] [node-1] loaded module [lang-expression]
[2017-05-25T14:53:13,834][INFO ][o.e.p.PluginsService     ] [node-1] loaded module [lang-groovy]
[2017-05-25T14:53:13,834][INFO ][o.e.p.PluginsService     ] [node-1] loaded module [lang-mustache]
[2017-05-25T14:53:13,834][INFO ][o.e.p.PluginsService     ] [node-1] loaded module [lang-painless]
[2017-05-25T14:53:13,834][INFO ][o.e.p.PluginsService     ] [node-1] loaded module [percolator]
[2017-05-25T14:53:13,835][INFO ][o.e.p.PluginsService     ] [node-1] loaded module [reindex]
[2017-05-25T14:53:13,835][INFO ][o.e.p.PluginsService     ] [node-1] loaded module [transport-netty3]
[2017-05-25T14:53:13,835][INFO ][o.e.p.PluginsService     ] [node-1] loaded module [transport-netty4]
[2017-05-25T14:53:13,836][INFO ][o.e.p.PluginsService     ] [node-1] loaded plugin [x-pack]
[2017-05-25T14:53:16,552][DEBUG][o.e.a.ActionModule       ] Using REST wrapper from plugin org.elasticsearch.xpack.XPackPlugin
[2017-05-25T14:53:18,134][INFO ][o.e.n.Node               ] [node-1] initialized
[2017-05-25T14:53:18,134][INFO ][o.e.n.Node               ] [node-1] starting ...
[2017-05-25T14:53:18,584][INFO ][o.e.t.TransportService   ] [node-1] publish_address {127.0.0.1:9300}, bound_addresses {[fe80::1]:9300}, {[::1]:9300}, {127.0.0.1:9300}
[2017-05-25T14:53:21,656][INFO ][o.e.c.s.ClusterService   ] [node-1] new_master {node-1}{TLbN2l8MTBGD_TmNuqIvow}{3iKoVnM7Qj2gZ--q_Yh4Xw}{127.0.0.1}{127.0.0.1:9300}, reason: zen-disco-elected-as-master ([0] nodes joined)
[2017-05-25T14:53:21,685][INFO ][o.e.x.s.t.n.SecurityNetty4HttpServerTransport] [node-1] publish_address {127.0.0.1:9200}, bound_addresses {[fe80::1]:9200}, {[::1]:9200}, {127.0.0.1:9200}
[2017-05-25T14:53:21,693][INFO ][o.e.n.Node               ] [node-1] started
[2017-05-25T14:53:21,787][INFO ][o.e.g.GatewayService     ] [node-1] recovered [0] indices into cluster_state
[2017-05-25T14:53:22,580][INFO ][o.e.l.LicenseService     ] [node-1] license [285d1469-8bba-482d-aaf7-fd86abddc76e] mode [trial] - valid
[2017-05-25T14:53:28,291][INFO ][o.e.c.m.MetaDataCreateIndexService] [node-1] [.monitoring-es-2-2017.05.25] creating index, cause [auto(bulk api)], templates [.monitoring-es-2], shards [1]/[1], mappings [_default_, node, shards, index_stats, index_recovery, cluster_state, cluster_stats, node_stats, indices_stats]
[2017-05-25T14:53:28,411][INFO ][o.e.c.m.MetaDataCreateIndexService] [node-1] [.monitoring-data-2] creating index, cause [auto(bulk api)], templates [.monitoring-data-2], shards [1]/[1], mappings [_default_, node, logstash, kibana, cluster_info]
[2017-05-25T14:53:28,661][INFO ][o.e.c.m.MetaDataMappingService] [node-1] [.monitoring-es-2-2017.05.25/RjffXrbpTwGzqO8Y1zmFuw] update_mapping [cluster_stats]
[2017-05-25T14:53:28,730][INFO ][o.e.c.m.MetaDataMappingService] [node-1] [.monitoring-es-2-2017.05.25/RjffXrbpTwGzqO8Y1zmFuw] update_mapping [node_stats]
[2017-05-25T14:53:29,091][INFO ][o.e.c.m.MetaDataCreateIndexService] [node-1] [.kibana] creating index, cause [api], templates [], shards [1]/[1], mappings [server, config]
[2017-05-25T14:53:38,240][INFO ][o.e.c.m.MetaDataMappingService] [node-1] [.monitoring-es-2-2017.05.25/RjffXrbpTwGzqO8Y1zmFuw] update_mapping [index_stats]
[2017-05-25T14:53:38,302][INFO ][o.e.c.m.MetaDataMappingService] [node-1] [.monitoring-es-2-2017.05.25/RjffXrbpTwGzqO8Y1zmFuw] update_mapping [cluster_stats]
[2017-05-25T14:53:38,343][INFO ][o.e.c.m.MetaDataMappingService] [node-1] [.monitoring-es-2-2017.05.25/RjffXrbpTwGzqO8Y1zmFuw] update_mapping [indices_stats]
[2017-05-25T14:53:39,219][INFO ][o.e.c.m.MetaDataCreateIndexService] [node-1] [.monitoring-kibana-2-2017.05.25] creating index, cause [auto(bulk api)], templates [.monitoring-kibana-2], shards [1]/[1], mappings [_default_, kibana_stats]
[2017-05-25T14:53:51,672][INFO ][o.e.c.r.a.DiskThresholdMonitor] [node-1] low disk watermark [85%] exceeded on [TLbN2l8MTBGD_TmNuqIvow][node-1][/usr/local/var/elasticsearch/nodes/0] free: 16.2gb[14.5%], replicas will not be assigned to this node
[2017-05-25T14:53:58,899][INFO ][o.e.c.m.MetaDataMappingService] [node-1] [.kibana/7PW7HsmMRO-qKR58fBRa3g] create_mapping [index-pattern]
[2017-05-25T14:53:59,227][INFO ][o.e.c.m.MetaDataMappingService] [node-1] [.monitoring-kibana-2-2017.05.25/qEAti7cyQgq7iKG40ZDzbQ] update_mapping [kibana_stats]
[2017-05-25T14:54:21,677][INFO ][o.e.c.r.a.DiskThresholdMonitor] [node-1] low disk watermark [85%] exceeded on [TLbN2l8MTBGD_TmNuqIvow][node-1][/usr/local/var/elasticsearch/nodes/0] free: 16.2gb[14.4%], replicas will not be assigned to this node
[2017-05-25T14:54:51,686][INFO ][o.e.c.r.a.DiskThresholdMonitor] [node-1] low disk watermark [85%] exceeded on [TLbN2l8MTBGD_TmNuqIvow][node-1][/usr/local/var/elasticsearch/nodes/0] free: 16.2gb[14.4%], replicas will not be assigned to this node
[2017-05-25T14:55:21,691][INFO ][o.e.c.r.a.DiskThresholdMonitor] [node-1] low disk watermark [85%] exceeded on [TLbN2l8MTBGD_TmNuqIvow][node-1][/usr/local/var/elasticsearch/nodes/0] free: 16.2gb[14.4%], replicas will not be assigned to this node
[2017-05-25T14:55:51,699][INFO ][o.e.c.r.a.DiskThresholdMonitor] [node-1] low disk watermark [85%] exceeded on [TLbN2l8MTBGD_TmNuqIvow][node-1][/usr/local/var/elasticsearch/nodes/0] free: 16.2gb[14.5%], replicas will not be assigned to this node
[2017-05-25T14:56:21,705][INFO ][o.e.c.r.a.DiskThresholdMonitor] [node-1] low disk watermark [85%] exceeded on [TLbN2l8MTBGD_TmNuqIvow][node-1][/usr/local/var/elasticsearch/nodes/0] free: 16.2gb[14.5%], replicas will not be assigned to this node
[2017-05-25T14:56:51,713][INFO ][o.e.c.r.a.DiskThresholdMonitor] [node-1] low disk watermark [85%] exceeded on [TLbN2l8MTBGD_TmNuqIvow][node-1][/usr/local/var/elasticsearch/nodes/0] free: 16.2gb[14.5%], replicas will not be assigned to this node
[2017-05-25T14:57:21,720][INFO ][o.e.c.r.a.DiskThresholdMonitor] [node-1] low disk watermark [85%] exceeded on [TLbN2l8MTBGD_TmNuqIvow][node-1][/usr/local/var/elasticsearch/nodes/0] free: 16.2gb[14.5%], replicas will not be assigned to this node
[2017-05-25T14:57:51,725][INFO ][o.e.c.r.a.DiskThresholdMonitor] [node-1] low disk watermark [85%] exceeded on [TLbN2l8MTBGD_TmNuqIvow][node-1][/usr/local/var/elasticsearch/nodes/0] free: 16.2gb[14.4%], replicas will not be assigned to this node
[2017-05-25T14:58:21,729][INFO ][o.e.c.r.a.DiskThresholdMonitor] [node-1] low disk watermark [85%] exceeded on [TLbN2l8MTBGD_TmNuqIvow][node-1][/usr/local/var/elasticsearch/nodes/0] free: 16.2gb[14.4%], replicas will not be assigned to this node
```

- 通过 Dev Tools 进行查询

查询名为 `packetbeat-*` 的 index pattern 是否存在

```
GET /packetbeat-*

返回

{}
```

查询匹配 `packetbeat-*` 的全部 index 下的数据（默认返回 10 条）

```
GET /packetbeat-*/_search

返回

{
  "took": 1,
  "timed_out": false,
  "_shards": {
    "total": 0,
    "successful": 0,
    "failed": 0
  },
  "hits": {
    "total": 0,
    "max_score": 0,
    "hits": []
  }
}
```

查询匹配 `packetbeat-*` 的 type 为 mysql 的全部的数据（默认返回 10 条）

```
GET /packetbeat-*/mysql/_search

返回

{
  "took": 1,
  "timed_out": false,
  "_shards": {
    "total": 0,
    "successful": 0,
    "failed": 0
  },
  "hits": {
    "total": 0,
    "max_score": 0,
    "hits": []
  }
}
```

查询名为 packetbeat* 的模版的内容

```
GET /_template/packetbeat*

返回

{}
```

此时打开 Kibana Discover 界面会提示“Configure an index pattern”；


- 通过 `packetbeat` 导入 `mysql` 数据

```
➜  packetbeat git:(5.3) ✗ ./packetbeat -c ./packetbeat.yml -e -I xg-breakfast-master-1_mysql_3306.pcap -t
2017/05/25 07:20:43.717303 beat.go:285: INFO Home path: [/Users/sunfei/workspace/GIT/IdeaProjects/src/github.com/elastic/beats/packetbeat] Config path: [/Users/sunfei/workspace/GIT/IdeaProjects/src/github.com/elastic/beats/packetbeat] Data path: [/Users/sunfei/workspace/GIT/IdeaProjects/src/github.com/elastic/beats/packetbeat/data] Logs path: [/Users/sunfei/workspace/GIT/IdeaProjects/src/github.com/elastic/beats/packetbeat/logs]
2017/05/25 07:20:43.717654 beat.go:186: INFO Setup Beat: packetbeat; Version: 5.3.3
2017/05/25 07:20:43.717653 metrics.go:23: INFO Metrics logging every 30s
2017/05/25 07:20:43.719603 output.go:254: INFO Loading template enabled. Reading template file: /Users/sunfei/workspace/GIT/IdeaProjects/src/github.com/elastic/beats/packetbeat/packetbeat.template.json
2017/05/25 07:20:43.725467 output.go:265: INFO Loading template enabled for Elasticsearch 2.x. Reading template file: /Users/sunfei/workspace/GIT/IdeaProjects/src/github.com/elastic/beats/packetbeat/packetbeat.template-es2x.json
2017/05/25 07:20:43.729634 client.go:123: INFO Elasticsearch url: http://localhost:9200
2017/05/25 07:20:43.730081 outputs.go:108: INFO Activated elasticsearch as output plugin.
2017/05/25 07:20:43.730217 publish.go:295: INFO Publisher name: sunfeideMacBook-Pro.local
2017/05/25 07:20:43.730333 async.go:63: INFO Flush Interval set to: 1s
2017/05/25 07:20:43.730749 async.go:64: INFO Max Bulk Size set to: 50
2017/05/25 07:20:43.730944 procs.go:79: INFO Process matching disabled
2017/05/25 07:20:43.731627 protos.go:89: INFO registered protocol plugin: amqp
2017/05/25 07:20:43.731648 protos.go:89: INFO registered protocol plugin: dns
2017/05/25 07:20:43.731656 protos.go:89: INFO registered protocol plugin: mongodb
2017/05/25 07:20:43.731663 protos.go:89: INFO registered protocol plugin: nfs
2017/05/25 07:20:43.731670 protos.go:89: INFO registered protocol plugin: pgsql
2017/05/25 07:20:43.731676 protos.go:89: INFO registered protocol plugin: thrift
2017/05/25 07:20:43.731683 protos.go:89: INFO registered protocol plugin: cassandra
2017/05/25 07:20:43.731690 protos.go:89: INFO registered protocol plugin: http
2017/05/25 07:20:43.731697 protos.go:89: INFO registered protocol plugin: memcache
2017/05/25 07:20:43.731704 protos.go:89: INFO registered protocol plugin: mysql
2017/05/25 07:20:43.731711 protos.go:89: INFO registered protocol plugin: redis
2017/05/25 07:20:43.735139 beat.go:221: INFO packetbeat start running.
2017/05/25 07:20:43.795875 client.go:658: INFO Connected to Elasticsearch version 5.3.2
2017/05/25 07:20:43.795898 output.go:301: INFO Trying to load template for client: http://localhost:9200
2017/05/25 07:20:43.802219 output.go:308: INFO Existing template will be overwritten, as overwrite is enabled.
2017/05/25 07:20:43.907118 client.go:588: INFO Elasticsearch template with name 'packetbeat' loaded
2017/05/25 07:20:45.690553 mysql.go:208: WARN MySQL Message too short. Ignore it.
2017/05/25 07:20:45.918868 mysql.go:709: WARN Invalid response: Data too small to contain a valid length
2017/05/25 07:20:46.130855 mysql.go:208: WARN MySQL Message too short. Ignore it.
2017/05/25 07:20:52.474078 mysql.go:709: WARN Invalid response: Data too small to contain a valid length
2017/05/25 07:20:54.251135 mysql.go:208: WARN MySQL Message too short. Ignore it.
2017/05/25 07:20:56.265696 sniffer.go:384: INFO Input finish. Processed 240450 packets. Have a nice day!
2017/05/25 07:20:56.269816 util.go:48: INFO flows worker loop stopped
2017/05/25 07:20:56.270085 metrics.go:51: INFO Total non-zero values:  libbeat.es.call_count.PublishEvents=408 libbeat.es.publish.read_bytes=205711 libbeat.es.publish.write_bytes=18800825 libbeat.es.published_and_acked_events=19967 libbeat.publisher.messages_in_worker_queues=21057 libbeat.publisher.published_events=21057 mysql.unmatched_requests=1138 mysql.unmatched_responses=50014 tcp.dropped_because_of_gaps=18
2017/05/25 07:20:56.270109 metrics.go:52: INFO Uptime: 12.695584664s
2017/05/25 07:20:56.270116 beat.go:225: INFO packetbeat stopped.
➜  packetbeat git:(5.3) ✗
```

注意：导入数据后，会自动创建名为 `packetbeat-2017.05.16` 的 index（output.elasticsearch 配置段的默认配置）以及名为 mysql 的 type ，而是否自动创建索引取决于 elasticsearch 本身的配置 action.auto_create_index ，默认为 true ；


- 再次通过 Dev Tools 进行查询

确认 `packetbeat-2017.05.16` 已经存在

```
GET /packetbeat-*

返回

{
  "packetbeat-2017.05.16": {
    "aliases": {},
    "mappings": {
      "_default_": {
        "_meta": {
          "version": "5.3.3"
        },
        "_all": {
          "norms": false
        },
        "dynamic_templates": [
          {
            "strings_as_keyword": {
              "match_mapping_type": "string",
              "mapping": {
                "ignore_above": 1024,
                "type": "keyword"
              }
            }
          }
        ],
        "date_detection": false,
        "properties": {
          "@timestamp": {
            "type": "date"
          },
...
          "mysql": {
            "properties": {
              "affected_rows": {
                "type": "long"
              },
              "error_code": {
                "type": "long"
              },
              "error_message": {
                "type": "keyword",
                "ignore_above": 1024
              },
              "insert_id": {
                "type": "keyword",
                "ignore_above": 1024
              },
              "iserror": {
                "type": "boolean"
              },
              "num_fields": {
                "type": "keyword",
                "ignore_above": 1024
              },
              "num_rows": {
                "type": "keyword",
                "ignore_above": 1024
              },
              "query": {
                "type": "keyword",
                "ignore_above": 1024
              }
            }
          },
...
    },
    "settings": {
      "index": {
        "mapping": {
          "total_fields": {
            "limit": "10000"
          }
        },
        "refresh_interval": "5s",
        "number_of_shards": "5",
        "provided_name": "packetbeat-2017.05.16",
        "creation_date": "1495696843918",
        "number_of_replicas": "1",
        "uuid": "Q0hO8Y78TWqeLVFBS7lQ3w",
        "version": {
          "created": "5030299"
        }
      }
    }
  }
}
```

而名为 `packetbeat` 的 index 是不存在的

```
GET /packetbeat/_search

# 返回

{
  "error": {
    "root_cause": [
      {
        "type": "index_not_found_exception",
        "reason": "no such index",
        "resource.type": "index_or_alias",
        "resource.id": "packetbeat",
        "index_uuid": "_na_",
        "index": "packetbeat"
      }
    ],
    "type": "index_not_found_exception",
    "reason": "no such index",
    "resource.type": "index_or_alias",
    "resource.id": "packetbeat",
    "index_uuid": "_na_",
    "index": "packetbeat"
  },
  "status": 404
}
```

查询匹配 `packetbeat-*` 的全部 index 下的数据（默认返回 10 条）

```
GET /packetbeat-*/_search

返回

{
  "took": 4,
  "timed_out": false,
  "_shards": {
    "total": 5,
    "successful": 5,
    "failed": 0
  },
  "hits": {
    "total": 20017,
    "max_score": 1,
    "hits": [
      {
        "_index": "packetbeat-2017.05.16",
        "_type": "mysql",
        "_id": "AVw-eqUAbPwJc1ucUhBK",
        "_score": 1,
        "_source": {
          "@timestamp": "2017-05-16T03:14:05.277Z",
          "beat": {
            "hostname": "sunfeideMacBook-Pro.local",
            "name": "sunfeideMacBook-Pro.local",
            "version": "5.3.3"
          },
          "bytes_in": 34,
          "bytes_out": 77,
          "client_ip": "10.0.47.35",
          "client_port": 22035,
          "client_proc": "",
          "client_server": "",
          "ip": "10.0.27.244",
          "method": "SELECT",
          "mysql": {
            "affected_rows": 0,
            "error_code": 0,
            "error_message": "",
            "insert_id": 0,
            "iserror": false,
            "num_fields": 1,
            "num_rows": 1
          },
          "path": ".",
          "port": 3306,
          "proc": "",
          "query": "select @@session.tx_read_only",
          "responsetime": 0,
          "server": "",
          "status": "OK",
          "type": "mysql"
        }
      },
      {
        "_index": "packetbeat-2017.05.16",
        "_type": "mysql",
        "_id": "AVw-eqUAbPwJc1ucUhBL",
        "_score": 1,
        "_source": {
          "@timestamp": "2017-05-16T03:14:05.277Z",
          "beat": {
            "hostname": "sunfeideMacBook-Pro.local",
            "name": "sunfeideMacBook-Pro.local",
            "version": "5.3.3"
          },
          "bytes_in": 152,
          "bytes_out": 77,
          "client_ip": "10.0.47.31",
          "client_port": 57094,
          "client_proc": "",
          "client_server": "",
          "ip": "10.0.27.37",
          "method": "/*",
          "mysql": {
            "affected_rows": 0,
            "error_code": 0,
            "error_message": "",
            "insert_id": 0,
            "iserror": false,
            "num_fields": 1,
            "num_rows": 1
          },
          "path": ".",
          "port": 3306,
          "proc": "",
          "query": "/* E:rid=zaocan.ele.me^^B7AE33373BA14CE8BA18F28945C2355D|1494904445271&rpcid=1.1.2.1&appid=me.ele.breakfast.mars:E */ select @@session.tx_read_only",
          "responsetime": 0,
          "server": "",
          "status": "OK",
          "type": "mysql"
        }
      },
      {
        "_index": "packetbeat-2017.05.16",
        "_type": "mysql",
        "_id": "AVw-eqUAbPwJc1ucUhBT",
        "_score": 1,
        "_source": {
          "@timestamp": "2017-05-16T03:14:05.282Z",
          "beat": {
            "hostname": "sunfeideMacBook-Pro.local",
            "name": "sunfeideMacBook-Pro.local",
            "version": "5.3.3"
          },
          "bytes_in": 141,
          "bytes_out": 77,
          "client_ip": "10.0.47.31",
          "client_port": 60605,
          "client_proc": "",
          "client_server": "",
          "ip": "10.0.27.37",
          "method": "/*",
          "mysql": {
            "affected_rows": 0,
            "error_code": 0,
            "error_message": "",
            "insert_id": 0,
            "iserror": false,
            "num_fields": 1,
            "num_rows": 1
          },
          "path": ".",
          "port": 3306,
          "proc": "",
          "query": "/* E:rid=me.ele.breakfast.api^^1939363357903773766|1494904445280&rpcid=1.1&appid=me.ele.breakfast.api:E */ select @@session.tx_read_only",
          "responsetime": 0,
          "server": "",
          "status": "OK",
          "type": "mysql"
        }
      },
      {
        "_index": "packetbeat-2017.05.16",
        "_type": "mysql",
        "_id": "AVw-eqUAbPwJc1ucUhBU",
        "_score": 1,
        "_source": {
          "@timestamp": "2017-05-16T03:14:05.282Z",
          "beat": {
            "hostname": "sunfeideMacBook-Pro.local",
            "name": "sunfeideMacBook-Pro.local",
            "version": "5.3.3"
          },
          "bytes_in": 495,
          "bytes_out": 1175,
          "client_ip": "10.0.47.35",
          "client_port": 63136,
          "client_proc": "",
          "client_server": "",
          "ip": "10.0.27.244",
          "method": "/*",
          "mysql": {
            "affected_rows": 0,
            "error_code": 0,
            "error_message": "",
            "insert_id": 0,
            "iserror": false,
            "num_fields": 14,
            "num_rows": 1
          },
          "path": "eleme_breakfast.gu, eleme_breakfast.u",
          "port": 3306,
          "proc": "",
          "query": """
/* E:rid=arch.gateway.zaocan^^1713117257738305733|1494904445273&rpcid=1.1.2&appid=me.ele.breakfast.backend:E */ SELECT gu.user_id,
         
		u.group_id,u.group_name,u.state,u.create_time,u.last_modified,u.modified_by,u.creator_id,u.last_modified,u.modified_by,
		u.modifier_id,u.data_level,u.sms_verify,u.assign_authority
	 
        FROM t_sys_group u
        JOIN t_sys_group_user gu ON gu.group_id = u.group_id
        AND gu.user_id IN
           (  
                12304
           )
""",
          "responsetime": 0,
          "server": "",
          "status": "OK",
          "type": "mysql"
        }
      },
      {
        "_index": "packetbeat-2017.05.16",
        "_type": "mysql",
        "_id": "AVw-eqUAbPwJc1ucUhBa",
        "_score": 1,
        "_source": {
          "@timestamp": "2017-05-16T03:14:05.283Z",
          "beat": {
            "hostname": "sunfeideMacBook-Pro.local",
            "name": "sunfeideMacBook-Pro.local",
            "version": "5.3.3"
          },
          "bytes_in": 142,
          "bytes_out": 77,
          "client_ip": "10.0.47.31",
          "client_port": 58762,
          "client_proc": "",
          "client_server": "",
          "ip": "10.0.27.37",
          "method": "/*",
          "mysql": {
            "affected_rows": 0,
            "error_code": 0,
            "error_message": "",
            "insert_id": 0,
            "iserror": false,
            "num_fields": 1,
            "num_rows": 1
          },
          "path": ".",
          "port": 3306,
          "proc": "",
          "query": "/* E:rid=me.ele.breakfast.api^^-6234699590973288371|1494904445280&rpcid=1.1&appid=me.ele.breakfast.api:E */ select @@session.tx_read_only",
          "responsetime": 0,
          "server": "",
          "status": "OK",
          "type": "mysql"
        }
      },
      {
        "_index": "packetbeat-2017.05.16",
        "_type": "mysql",
        "_id": "AVw-eqUAbPwJc1ucUhBd",
        "_score": 1,
        "_source": {
          "@timestamp": "2017-05-16T03:14:05.283Z",
          "beat": {
            "hostname": "sunfeideMacBook-Pro.local",
            "name": "sunfeideMacBook-Pro.local",
            "version": "5.3.3"
          },
          "bytes_in": 34,
          "bytes_out": 77,
          "client_ip": "10.0.47.35",
          "client_port": 63138,
          "client_proc": "",
          "client_server": "",
          "ip": "10.0.27.244",
          "method": "SELECT",
          "mysql": {
            "affected_rows": 0,
            "error_code": 0,
            "error_message": "",
            "insert_id": 0,
            "iserror": false,
            "num_fields": 1,
            "num_rows": 1
          },
          "path": ".",
          "port": 3306,
          "proc": "",
          "query": "select @@session.tx_read_only",
          "responsetime": 0,
          "server": "",
          "status": "OK",
          "type": "mysql"
        }
      },
      {
        "_index": "packetbeat-2017.05.16",
        "_type": "mysql",
        "_id": "AVw-eqUAbPwJc1ucUhBp",
        "_score": 1,
        "_source": {
          "@timestamp": "2017-05-16T03:14:05.285Z",
          "beat": {
            "hostname": "sunfeideMacBook-Pro.local",
            "name": "sunfeideMacBook-Pro.local",
            "version": "5.3.3"
          },
          "bytes_in": 406,
          "bytes_out": 19,
          "client_ip": "10.0.47.31",
          "client_port": 58572,
          "client_proc": "",
          "client_server": "",
          "ip": "10.0.27.37",
          "method": "/*",
          "mysql": {
            "affected_rows": 1,
            "error_code": 0,
            "error_message": "",
            "insert_id": 281984406621,
            "iserror": false,
            "num_fields": 0,
            "num_rows": 0
          },
          "path": "",
          "port": 3306,
          "proc": "",
          "query": """
/* E:rid=me.ele.breakfast.mars^^-9169731336371294008|1494904445284&rpcid=1.1&appid=me.ele.breakfast.mars:E */ insert into t_member_cache_log (user_id, phone_no,
      device_no, grade, create_time, remark)
    values (15730072, null,
      'null', 1,
      now(), '通过缓存中获取的新老用户信息为：null,通过查询获取新老用户信息为：用户、手机、设备全部是新的')
""",
          "responsetime": 0,
          "server": "",
          "status": "OK",
          "type": "mysql"
        }
      },
      {
        "_index": "packetbeat-2017.05.16",
        "_type": "mysql",
        "_id": "AVw-eqUAbPwJc1ucUhBt",
        "_score": 1,
        "_source": {
          "@timestamp": "2017-05-16T03:14:05.286Z",
          "beat": {
            "hostname": "sunfeideMacBook-Pro.local",
            "name": "sunfeideMacBook-Pro.local",
            "version": "5.3.3"
          },
          "bytes_in": 30,
          "bytes_out": 73,
          "client_ip": "10.0.13.95",
          "client_port": 55637,
          "client_proc": "",
          "client_server": "",
          "ip": "10.0.27.244",
          "method": "SELECT",
          "mysql": {
            "affected_rows": 0,
            "error_code": 0,
            "error_message": "",
            "insert_id": 0,
            "iserror": false,
            "num_fields": 1,
            "num_rows": 1
          },
          "path": ".",
          "port": 3306,
          "proc": "",
          "query": "SELECT @@global.read_only",
          "responsetime": 0,
          "server": "",
          "status": "OK",
          "type": "mysql"
        }
      },
      {
        "_index": "packetbeat-2017.05.16",
        "_type": "mysql",
        "_id": "AVw-eqX7bPwJc1ucUhB3",
        "_score": 1,
        "_source": {
          "@timestamp": "2017-05-16T03:14:05.287Z",
          "beat": {
            "hostname": "sunfeideMacBook-Pro.local",
            "name": "sunfeideMacBook-Pro.local",
            "version": "5.3.3"
          },
          "bytes_in": 468,
          "bytes_out": 1775,
          "client_ip": "10.0.47.31",
          "client_port": 61173,
          "client_proc": "",
          "client_server": "",
          "ip": "10.0.27.37",
          "method": "/*",
          "mysql": {
            "affected_rows": 0,
            "error_code": 0,
            "error_message": "",
            "insert_id": 0,
            "iserror": false,
            "num_fields": 20,
            "num_rows": 1
          },
          "path": "eleme_breakfast.t_dish",
          "port": 3306,
          "proc": "",
          "query": """
/* E:rid=me.ele.breakfast.api^^-1226129632404483068|1494904445285&rpcid=1.1&appid=me.ele.breakfast.api:E */ select
         
    dishId, name, price, resturantId, createTime, lastModified, modifiedBy, productivity, 
    purchasePrice, status, selling_price, privilege_type, privilege_amount, new_privilege_type, 
    new_privilege_amount, new_selling_price, rating,dish_category,dish_summary,one_level_label
   
        from t_dish
     
    where dishId = 149993
""",
          "responsetime": 0,
          "server": "",
          "status": "OK",
          "type": "mysql"
        }
      },
      {
        "_index": "packetbeat-2017.05.16",
        "_type": "mysql",
        "_id": "AVw-eqX7bPwJc1ucUhB4",
        "_score": 1,
        "_source": {
          "@timestamp": "2017-05-16T03:14:05.287Z",
          "beat": {
            "hostname": "sunfeideMacBook-Pro.local",
            "name": "sunfeideMacBook-Pro.local",
            "version": "5.3.3"
          },
          "bytes_in": 207,
          "bytes_out": 63,
          "client_ip": "10.0.47.27",
          "client_port": 29601,
          "client_proc": "",
          "client_server": "",
          "ip": "10.0.27.37",
          "method": "/*",
          "mysql": {
            "affected_rows": 0,
            "error_code": 0,
            "error_message": "",
            "insert_id": 0,
            "iserror": false,
            "num_fields": 1,
            "num_rows": 1
          },
          "path": ".",
          "port": 3306,
          "proc": "",
          "query": """
/* E:rid=me.ele.breakfast.api^^8360090407405843693|1494904445283&rpcid=1.1.1&appid=me.ele.breakfast.mars:E */ SELECT COUNT(1)
    FROM t_member
    WHERE grade=2 AND
    (
      user_id=130692170 
    )
""",
          "responsetime": 0,
          "server": "",
          "status": "OK",
          "type": "mysql"
        }
      }
    ]
  }
}
```

查询匹配 `packetbeat-*` 的 type 为 mysql 的全部的数据（默认返回 10 条）

```
GET /packetbeat-*/mysql/_search

返回

{
  "took": 1,
  "timed_out": false,
  "_shards": {
    "total": 5,
    "successful": 5,
    "failed": 0
  },
  "hits": {
    "total": 20017,
    "max_score": 1,
    "hits": [
      {
        "_index": "packetbeat-2017.05.16",
        "_type": "mysql",
        "_id": "AVw-eqUAbPwJc1ucUhBK",
        "_score": 1,
        "_source": {
          "@timestamp": "2017-05-16T03:14:05.277Z",
          "beat": {
            "hostname": "sunfeideMacBook-Pro.local",
            "name": "sunfeideMacBook-Pro.local",
            "version": "5.3.3"
          },
          "bytes_in": 34,
          "bytes_out": 77,
          "client_ip": "10.0.47.35",
          "client_port": 22035,
          "client_proc": "",
          "client_server": "",
          "ip": "10.0.27.244",
          "method": "SELECT",
          "mysql": {
            "affected_rows": 0,
            "error_code": 0,
            "error_message": "",
            "insert_id": 0,
            "iserror": false,
            "num_fields": 1,
            "num_rows": 1
          },
          "path": ".",
          "port": 3306,
          "proc": "",
          "query": "select @@session.tx_read_only",
          "responsetime": 0,
          "server": "",
          "status": "OK",
          "type": "mysql"
        }
      },
      {
        "_index": "packetbeat-2017.05.16",
        "_type": "mysql",
        "_id": "AVw-eqUAbPwJc1ucUhBL",
        "_score": 1,
        "_source": {
          "@timestamp": "2017-05-16T03:14:05.277Z",
          "beat": {
            "hostname": "sunfeideMacBook-Pro.local",
            "name": "sunfeideMacBook-Pro.local",
            "version": "5.3.3"
          },
          "bytes_in": 152,
          "bytes_out": 77,
          "client_ip": "10.0.47.31",
          "client_port": 57094,
          "client_proc": "",
          "client_server": "",
          "ip": "10.0.27.37",
          "method": "/*",
          "mysql": {
            "affected_rows": 0,
            "error_code": 0,
            "error_message": "",
            "insert_id": 0,
            "iserror": false,
            "num_fields": 1,
            "num_rows": 1
          },
          "path": ".",
          "port": 3306,
          "proc": "",
          "query": "/* E:rid=zaocan.ele.me^^B7AE33373BA14CE8BA18F28945C2355D|1494904445271&rpcid=1.1.2.1&appid=me.ele.breakfast.mars:E */ select @@session.tx_read_only",
          "responsetime": 0,
          "server": "",
          "status": "OK",
          "type": "mysql"
        }
      },
      {
        "_index": "packetbeat-2017.05.16",
        "_type": "mysql",
        "_id": "AVw-eqUAbPwJc1ucUhBT",
        "_score": 1,
        "_source": {
          "@timestamp": "2017-05-16T03:14:05.282Z",
          "beat": {
            "hostname": "sunfeideMacBook-Pro.local",
            "name": "sunfeideMacBook-Pro.local",
            "version": "5.3.3"
          },
          "bytes_in": 141,
          "bytes_out": 77,
          "client_ip": "10.0.47.31",
          "client_port": 60605,
          "client_proc": "",
          "client_server": "",
          "ip": "10.0.27.37",
          "method": "/*",
          "mysql": {
            "affected_rows": 0,
            "error_code": 0,
            "error_message": "",
            "insert_id": 0,
            "iserror": false,
            "num_fields": 1,
            "num_rows": 1
          },
          "path": ".",
          "port": 3306,
          "proc": "",
          "query": "/* E:rid=me.ele.breakfast.api^^1939363357903773766|1494904445280&rpcid=1.1&appid=me.ele.breakfast.api:E */ select @@session.tx_read_only",
          "responsetime": 0,
          "server": "",
          "status": "OK",
          "type": "mysql"
        }
      },
      {
        "_index": "packetbeat-2017.05.16",
        "_type": "mysql",
        "_id": "AVw-eqUAbPwJc1ucUhBU",
        "_score": 1,
        "_source": {
          "@timestamp": "2017-05-16T03:14:05.282Z",
          "beat": {
            "hostname": "sunfeideMacBook-Pro.local",
            "name": "sunfeideMacBook-Pro.local",
            "version": "5.3.3"
          },
          "bytes_in": 495,
          "bytes_out": 1175,
          "client_ip": "10.0.47.35",
          "client_port": 63136,
          "client_proc": "",
          "client_server": "",
          "ip": "10.0.27.244",
          "method": "/*",
          "mysql": {
            "affected_rows": 0,
            "error_code": 0,
            "error_message": "",
            "insert_id": 0,
            "iserror": false,
            "num_fields": 14,
            "num_rows": 1
          },
          "path": "eleme_breakfast.gu, eleme_breakfast.u",
          "port": 3306,
          "proc": "",
          "query": """
/* E:rid=arch.gateway.zaocan^^1713117257738305733|1494904445273&rpcid=1.1.2&appid=me.ele.breakfast.backend:E */ SELECT gu.user_id,
         
		u.group_id,u.group_name,u.state,u.create_time,u.last_modified,u.modified_by,u.creator_id,u.last_modified,u.modified_by,
		u.modifier_id,u.data_level,u.sms_verify,u.assign_authority
	 
        FROM t_sys_group u
        JOIN t_sys_group_user gu ON gu.group_id = u.group_id
        AND gu.user_id IN
           (  
                12304
           )
""",
          "responsetime": 0,
          "server": "",
          "status": "OK",
          "type": "mysql"
        }
      },
      {
        "_index": "packetbeat-2017.05.16",
        "_type": "mysql",
        "_id": "AVw-eqUAbPwJc1ucUhBa",
        "_score": 1,
        "_source": {
          "@timestamp": "2017-05-16T03:14:05.283Z",
          "beat": {
            "hostname": "sunfeideMacBook-Pro.local",
            "name": "sunfeideMacBook-Pro.local",
            "version": "5.3.3"
          },
          "bytes_in": 142,
          "bytes_out": 77,
          "client_ip": "10.0.47.31",
          "client_port": 58762,
          "client_proc": "",
          "client_server": "",
          "ip": "10.0.27.37",
          "method": "/*",
          "mysql": {
            "affected_rows": 0,
            "error_code": 0,
            "error_message": "",
            "insert_id": 0,
            "iserror": false,
            "num_fields": 1,
            "num_rows": 1
          },
          "path": ".",
          "port": 3306,
          "proc": "",
          "query": "/* E:rid=me.ele.breakfast.api^^-6234699590973288371|1494904445280&rpcid=1.1&appid=me.ele.breakfast.api:E */ select @@session.tx_read_only",
          "responsetime": 0,
          "server": "",
          "status": "OK",
          "type": "mysql"
        }
      },
      {
        "_index": "packetbeat-2017.05.16",
        "_type": "mysql",
        "_id": "AVw-eqUAbPwJc1ucUhBd",
        "_score": 1,
        "_source": {
          "@timestamp": "2017-05-16T03:14:05.283Z",
          "beat": {
            "hostname": "sunfeideMacBook-Pro.local",
            "name": "sunfeideMacBook-Pro.local",
            "version": "5.3.3"
          },
          "bytes_in": 34,
          "bytes_out": 77,
          "client_ip": "10.0.47.35",
          "client_port": 63138,
          "client_proc": "",
          "client_server": "",
          "ip": "10.0.27.244",
          "method": "SELECT",
          "mysql": {
            "affected_rows": 0,
            "error_code": 0,
            "error_message": "",
            "insert_id": 0,
            "iserror": false,
            "num_fields": 1,
            "num_rows": 1
          },
          "path": ".",
          "port": 3306,
          "proc": "",
          "query": "select @@session.tx_read_only",
          "responsetime": 0,
          "server": "",
          "status": "OK",
          "type": "mysql"
        }
      },
      {
        "_index": "packetbeat-2017.05.16",
        "_type": "mysql",
        "_id": "AVw-eqUAbPwJc1ucUhBp",
        "_score": 1,
        "_source": {
          "@timestamp": "2017-05-16T03:14:05.285Z",
          "beat": {
            "hostname": "sunfeideMacBook-Pro.local",
            "name": "sunfeideMacBook-Pro.local",
            "version": "5.3.3"
          },
          "bytes_in": 406,
          "bytes_out": 19,
          "client_ip": "10.0.47.31",
          "client_port": 58572,
          "client_proc": "",
          "client_server": "",
          "ip": "10.0.27.37",
          "method": "/*",
          "mysql": {
            "affected_rows": 1,
            "error_code": 0,
            "error_message": "",
            "insert_id": 281984406621,
            "iserror": false,
            "num_fields": 0,
            "num_rows": 0
          },
          "path": "",
          "port": 3306,
          "proc": "",
          "query": """
/* E:rid=me.ele.breakfast.mars^^-9169731336371294008|1494904445284&rpcid=1.1&appid=me.ele.breakfast.mars:E */ insert into t_member_cache_log (user_id, phone_no,
      device_no, grade, create_time, remark)
    values (15730072, null,
      'null', 1,
      now(), '通过缓存中获取的新老用户信息为：null,通过查询获取新老用户信息为：用户、手机、设备全部是新的')
""",
          "responsetime": 0,
          "server": "",
          "status": "OK",
          "type": "mysql"
        }
      },
      {
        "_index": "packetbeat-2017.05.16",
        "_type": "mysql",
        "_id": "AVw-eqUAbPwJc1ucUhBt",
        "_score": 1,
        "_source": {
          "@timestamp": "2017-05-16T03:14:05.286Z",
          "beat": {
            "hostname": "sunfeideMacBook-Pro.local",
            "name": "sunfeideMacBook-Pro.local",
            "version": "5.3.3"
          },
          "bytes_in": 30,
          "bytes_out": 73,
          "client_ip": "10.0.13.95",
          "client_port": 55637,
          "client_proc": "",
          "client_server": "",
          "ip": "10.0.27.244",
          "method": "SELECT",
          "mysql": {
            "affected_rows": 0,
            "error_code": 0,
            "error_message": "",
            "insert_id": 0,
            "iserror": false,
            "num_fields": 1,
            "num_rows": 1
          },
          "path": ".",
          "port": 3306,
          "proc": "",
          "query": "SELECT @@global.read_only",
          "responsetime": 0,
          "server": "",
          "status": "OK",
          "type": "mysql"
        }
      },
      {
        "_index": "packetbeat-2017.05.16",
        "_type": "mysql",
        "_id": "AVw-eqX7bPwJc1ucUhB3",
        "_score": 1,
        "_source": {
          "@timestamp": "2017-05-16T03:14:05.287Z",
          "beat": {
            "hostname": "sunfeideMacBook-Pro.local",
            "name": "sunfeideMacBook-Pro.local",
            "version": "5.3.3"
          },
          "bytes_in": 468,
          "bytes_out": 1775,
          "client_ip": "10.0.47.31",
          "client_port": 61173,
          "client_proc": "",
          "client_server": "",
          "ip": "10.0.27.37",
          "method": "/*",
          "mysql": {
            "affected_rows": 0,
            "error_code": 0,
            "error_message": "",
            "insert_id": 0,
            "iserror": false,
            "num_fields": 20,
            "num_rows": 1
          },
          "path": "eleme_breakfast.t_dish",
          "port": 3306,
          "proc": "",
          "query": """
/* E:rid=me.ele.breakfast.api^^-1226129632404483068|1494904445285&rpcid=1.1&appid=me.ele.breakfast.api:E */ select
         
    dishId, name, price, resturantId, createTime, lastModified, modifiedBy, productivity, 
    purchasePrice, status, selling_price, privilege_type, privilege_amount, new_privilege_type, 
    new_privilege_amount, new_selling_price, rating,dish_category,dish_summary,one_level_label
   
        from t_dish
     
    where dishId = 149993
""",
          "responsetime": 0,
          "server": "",
          "status": "OK",
          "type": "mysql"
        }
      },
      {
        "_index": "packetbeat-2017.05.16",
        "_type": "mysql",
        "_id": "AVw-eqX7bPwJc1ucUhB4",
        "_score": 1,
        "_source": {
          "@timestamp": "2017-05-16T03:14:05.287Z",
          "beat": {
            "hostname": "sunfeideMacBook-Pro.local",
            "name": "sunfeideMacBook-Pro.local",
            "version": "5.3.3"
          },
          "bytes_in": 207,
          "bytes_out": 63,
          "client_ip": "10.0.47.27",
          "client_port": 29601,
          "client_proc": "",
          "client_server": "",
          "ip": "10.0.27.37",
          "method": "/*",
          "mysql": {
            "affected_rows": 0,
            "error_code": 0,
            "error_message": "",
            "insert_id": 0,
            "iserror": false,
            "num_fields": 1,
            "num_rows": 1
          },
          "path": ".",
          "port": 3306,
          "proc": "",
          "query": """
/* E:rid=me.ele.breakfast.api^^8360090407405843693|1494904445283&rpcid=1.1.1&appid=me.ele.breakfast.mars:E */ SELECT COUNT(1)
    FROM t_member
    WHERE grade=2 AND
    (
      user_id=130692170 
    )
""",
          "responsetime": 0,
          "server": "",
          "status": "OK",
          "type": "mysql"
        }
      }
    ]
  }
}
```

查询名为 packetbeat* 的模版的内容

```
GET /_template/packetbeat*

返回

{
  "packetbeat": {
    "order": 0,
    "template": "packetbeat-*",
    "settings": {
      "index": {
        "mapping": {
          "total_fields": {
            "limit": "10000"
          }
        },
        "refresh_interval": "5s"
      }
    },
    "mappings": {
      "_default_": {
        "_meta": {
          "version": "5.3.3"
        },
        "dynamic_templates": [
          {
            "strings_as_keyword": {
              "mapping": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "match_mapping_type": "string"
            }
          }
        ],
        "_all": {
          "norms": false
        },
        "date_detection": false,
        "properties": {
          "real_ip": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "notes": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "release": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "source": {
            "properties": {
              "outer_ipv6_location": {
                "type": "geo_point"
              },
              "port": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "stats": {
                "properties": {
                  "net_bytes_total": {
                    "type": "long"
                  },
                  "net_packets_total": {
                    "type": "long"
                  }
                }
              },
              "ipv6": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "ip": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "outer_ip": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "outer_ipv6": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "ip_location": {
                "type": "geo_point"
              },
              "ipv6_location": {
                "type": "geo_point"
              },
              "mac": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "outer_ip_location": {
                "type": "geo_point"
              }
            }
          },
          "dest": {
            "properties": {
              "outer_ipv6_location": {
                "type": "geo_point"
              },
              "port": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "stats": {
                "properties": {
                  "net_bytes_total": {
                    "type": "long"
                  },
                  "net_packets_total": {
                    "type": "long"
                  }
                }
              },
              "ipv6": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "ip": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "outer_ip": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "outer_ipv6": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "ip_location": {
                "type": "geo_point"
              },
              "ipv6_location": {
                "type": "geo_point"
              },
              "mac": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "outer_ip_location": {
                "type": "geo_point"
              }
            }
          },
          "icmp": {
            "properties": {
              "request": {
                "properties": {
                  "code": {
                    "type": "long"
                  },
                  "message": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "type": {
                    "type": "long"
                  }
                }
              },
              "response": {
                "properties": {
                  "code": {
                    "type": "long"
                  },
                  "message": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "type": {
                    "type": "long"
                  }
                }
              },
              "version": {
                "ignore_above": 1024,
                "type": "keyword"
              }
            }
          },
          "type": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "redis": {
            "properties": {
              "return_value": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "error": {
                "ignore_above": 1024,
                "type": "keyword"
              }
            }
          },
          "path": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "domloadtime": {
            "type": "long"
          },
          "flow_id": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "beat": {
            "properties": {
              "hostname": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "name": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "version": {
                "ignore_above": 1024,
                "type": "keyword"
              }
            }
          },
          "client_ip": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "mysql": {
            "properties": {
              "error_message": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "insert_id": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "query": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "num_fields": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "num_rows": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "error_code": {
                "type": "long"
              },
              "affected_rows": {
                "type": "long"
              },
              "iserror": {
                "type": "boolean"
              }
            }
          },
          "memcache": {
            "properties": {
              "protocol_type": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "request": {
                "properties": {
                  "count_values": {
                    "type": "long"
                  },
                  "opaque": {
                    "type": "long"
                  },
                  "sleep_us": {
                    "type": "long"
                  },
                  "noreply": {
                    "type": "boolean"
                  },
                  "initial": {
                    "type": "long"
                  },
                  "keys": {
                    "properties": {}
                  },
                  "line": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "values": {
                    "properties": {}
                  },
                  "delta": {
                    "type": "long"
                  },
                  "flags": {
                    "type": "long"
                  },
                  "cas_unique": {
                    "type": "long"
                  },
                  "automove": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "opcode": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "type": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "command": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "raw_args": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "exptime": {
                    "type": "long"
                  },
                  "bytes": {
                    "type": "long"
                  },
                  "dest_class": {
                    "type": "long"
                  },
                  "source_class": {
                    "type": "long"
                  },
                  "vbucket": {
                    "type": "long"
                  },
                  "opcode_value": {
                    "type": "long"
                  },
                  "quiet": {
                    "type": "boolean"
                  },
                  "verbosity": {
                    "type": "long"
                  }
                }
              },
              "response": {
                "properties": {
                  "count_values": {
                    "type": "long"
                  },
                  "error_msg": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "opaque": {
                    "type": "long"
                  },
                  "status_code": {
                    "type": "long"
                  },
                  "keys": {
                    "properties": {}
                  },
                  "values": {
                    "properties": {}
                  },
                  "flags": {
                    "type": "long"
                  },
                  "cas_unique": {
                    "type": "long"
                  },
                  "opcode": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "type": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "version": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "command": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "stats": {
                    "properties": {}
                  },
                  "bytes": {
                    "type": "long"
                  },
                  "opcode_value": {
                    "type": "long"
                  },
                  "value": {
                    "type": "long"
                  },
                  "status": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  }
                }
              }
            }
          },
          "client_geoip": {
            "properties": {
              "location": {
                "type": "geo_point"
              }
            }
          },
          "loadtime": {
            "type": "long"
          },
          "method": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "resource": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "ip": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "query": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "dns": {
            "properties": {
              "op_code": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "additionals": {
                "properties": {
                  "data": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "name": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "type": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "class": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "ttl": {
                    "type": "long"
                  }
                }
              },
              "opt": {
                "properties": {
                  "ext_rcode": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "udp_size": {
                    "type": "long"
                  },
                  "do": {
                    "type": "boolean"
                  },
                  "version": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  }
                }
              },
              "response_code": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "question": {
                "properties": {
                  "etld_plus_one": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "name": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "type": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "class": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  }
                }
              },
              "answers_count": {
                "type": "long"
              },
              "authorities_count": {
                "type": "long"
              },
              "answers": {
                "properties": {
                  "data": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "name": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "type": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "class": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "ttl": {
                    "type": "long"
                  }
                }
              },
              "flags": {
                "properties": {
                  "authoritative": {
                    "type": "boolean"
                  },
                  "truncated_response": {
                    "type": "boolean"
                  },
                  "checking_disabled": {
                    "type": "boolean"
                  },
                  "recursion_available": {
                    "type": "boolean"
                  },
                  "recursion_desired": {
                    "type": "boolean"
                  },
                  "authentic_data": {
                    "type": "boolean"
                  }
                }
              },
              "additionals_count": {
                "type": "long"
              },
              "id": {
                "type": "long"
              },
              "authorities": {
                "properties": {
                  "name": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "type": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "class": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  }
                }
              }
            }
          },
          "params": {
            "norms": false,
            "type": "text"
          },
          "pgsql": {
            "properties": {
              "error_message": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "error_severity": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "query": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "num_fields": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "num_rows": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "error_code": {
                "type": "long"
              },
              "iserror": {
                "type": "boolean"
              }
            }
          },
          "tags": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "start_time": {
            "type": "date"
          },
          "bytes_out": {
            "type": "long"
          },
          "cassandra": {
            "properties": {
              "request": {
                "properties": {
                  "headers": {
                    "properties": {
                      "op": {
                        "ignore_above": 1024,
                        "type": "keyword"
                      },
                      "stream": {
                        "ignore_above": 1024,
                        "type": "keyword"
                      },
                      "flags": {
                        "ignore_above": 1024,
                        "type": "keyword"
                      },
                      "length": {
                        "type": "long"
                      },
                      "version": {
                        "type": "long"
                      }
                    }
                  },
                  "query": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  }
                }
              },
              "response": {
                "properties": {
                  "result": {
                    "properties": {
                      "keyspace": {
                        "ignore_above": 1024,
                        "type": "keyword"
                      },
                      "schema_change": {
                        "properties": {
                          "args": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          },
                          "keyspace": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          },
                          "change": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          },
                          "name": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          },
                          "table": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          },
                          "object": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          },
                          "target": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          }
                        }
                      },
                      "prepared": {
                        "properties": {
                          "prepared_id": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          },
                          "req_meta": {
                            "properties": {
                              "keyspace": {
                                "ignore_above": 1024,
                                "type": "keyword"
                              },
                              "pkey_columns": {
                                "type": "long"
                              },
                              "paging_state": {
                                "ignore_above": 1024,
                                "type": "keyword"
                              },
                              "flags": {
                                "ignore_above": 1024,
                                "type": "keyword"
                              },
                              "col_count": {
                                "type": "long"
                              },
                              "table": {
                                "ignore_above": 1024,
                                "type": "keyword"
                              }
                            }
                          },
                          "resp_meta": {
                            "properties": {
                              "keyspace": {
                                "ignore_above": 1024,
                                "type": "keyword"
                              },
                              "pkey_columns": {
                                "type": "long"
                              },
                              "paging_state": {
                                "ignore_above": 1024,
                                "type": "keyword"
                              },
                              "flags": {
                                "ignore_above": 1024,
                                "type": "keyword"
                              },
                              "col_count": {
                                "type": "long"
                              },
                              "table": {
                                "ignore_above": 1024,
                                "type": "keyword"
                              }
                            }
                          }
                        }
                      },
                      "rows": {
                        "properties": {
                          "meta": {
                            "properties": {
                              "keyspace": {
                                "ignore_above": 1024,
                                "type": "keyword"
                              },
                              "pkey_columns": {
                                "type": "long"
                              },
                              "paging_state": {
                                "ignore_above": 1024,
                                "type": "keyword"
                              },
                              "flags": {
                                "ignore_above": 1024,
                                "type": "keyword"
                              },
                              "col_count": {
                                "type": "long"
                              },
                              "table": {
                                "ignore_above": 1024,
                                "type": "keyword"
                              }
                            }
                          },
                          "num_rows": {
                            "type": "long"
                          }
                        }
                      },
                      "type": {
                        "ignore_above": 1024,
                        "type": "keyword"
                      }
                    }
                  },
                  "headers": {
                    "properties": {
                      "op": {
                        "ignore_above": 1024,
                        "type": "keyword"
                      },
                      "stream": {
                        "ignore_above": 1024,
                        "type": "keyword"
                      },
                      "flags": {
                        "ignore_above": 1024,
                        "type": "keyword"
                      },
                      "length": {
                        "type": "long"
                      },
                      "version": {
                        "type": "long"
                      }
                    }
                  },
                  "warnings": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "error": {
                    "properties": {
                      "msg": {
                        "ignore_above": 1024,
                        "type": "keyword"
                      },
                      "code": {
                        "type": "long"
                      },
                      "details": {
                        "properties": {
                          "alive": {
                            "type": "long"
                          },
                          "stmt_id": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          },
                          "received": {
                            "type": "long"
                          },
                          "write_type": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          },
                          "num_failures": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          },
                          "required": {
                            "type": "long"
                          },
                          "read_consistency": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          },
                          "keyspace": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          },
                          "function": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          },
                          "arg_types": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          },
                          "data_present": {
                            "type": "boolean"
                          },
                          "table": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          },
                          "blockfor": {
                            "type": "long"
                          }
                        }
                      },
                      "type": {
                        "ignore_above": 1024,
                        "type": "keyword"
                      }
                    }
                  },
                  "event": {
                    "properties": {
                      "schema_change": {
                        "properties": {
                          "args": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          },
                          "keyspace": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          },
                          "change": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          },
                          "name": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          },
                          "table": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          },
                          "object": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          },
                          "target": {
                            "ignore_above": 1024,
                            "type": "keyword"
                          }
                        }
                      },
                      "port": {
                        "type": "long"
                      },
                      "change": {
                        "ignore_above": 1024,
                        "type": "keyword"
                      },
                      "host": {
                        "ignore_above": 1024,
                        "type": "keyword"
                      },
                      "type": {
                        "ignore_above": 1024,
                        "type": "keyword"
                      }
                    }
                  },
                  "authentication": {
                    "properties": {
                      "class": {
                        "ignore_above": 1024,
                        "type": "keyword"
                      }
                    }
                  },
                  "supported": {
                    "properties": {}
                  }
                }
              }
            }
          },
          "port": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "meta": {
            "properties": {
              "cloud": {
                "properties": {
                  "machine_type": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "availability_zone": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "instance_id": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "project_id": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "provider": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "region": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  }
                }
              }
            }
          },
          "final": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "http": {
            "properties": {
              "request": {
                "properties": {
                  "headers": {
                    "properties": {}
                  },
                  "body": {
                    "norms": false,
                    "type": "text"
                  },
                  "params": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  }
                }
              },
              "response": {
                "properties": {
                  "headers": {
                    "properties": {}
                  },
                  "code": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "phrase": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "body": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  }
                }
              }
            }
          },
          "nfs": {
            "properties": {
              "tag": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "opcode": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "version": {
                "type": "long"
              },
              "minor_version": {
                "type": "long"
              },
              "status": {
                "ignore_above": 1024,
                "type": "keyword"
              }
            }
          },
          "fields": {
            "properties": {}
          },
          "status": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "request": {
            "norms": false,
            "type": "text"
          },
          "server": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "bytes_in": {
            "type": "long"
          },
          "client_service": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "icmp_id": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "client_port": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "outer_vlan": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "vlan": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "client_location": {
            "type": "geo_point"
          },
          "last_time": {
            "type": "date"
          },
          "thrift": {
            "properties": {
              "return_value": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "service": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "params": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "exceptions": {
                "ignore_above": 1024,
                "type": "keyword"
              }
            }
          },
          "dnstime": {
            "type": "long"
          },
          "responsetime": {
            "type": "long"
          },
          "mongodb": {
            "properties": {
              "fullCollectionName": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "numberReturned": {
                "type": "long"
              },
              "query": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "numberToSkip": {
                "type": "long"
              },
              "update": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "selector": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "error": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "returnFieldsSelector": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "startingFrom": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "cursorId": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "numberToReturn": {
                "type": "long"
              }
            }
          },
          "direction": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "cpu_time": {
            "type": "long"
          },
          "proc": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "rpc": {
            "properties": {
              "cred": {
                "properties": {
                  "gids": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "uid": {
                    "type": "long"
                  },
                  "gid": {
                    "type": "long"
                  },
                  "machinename": {
                    "ignore_above": 1024,
                    "type": "keyword"
                  },
                  "stamp": {
                    "type": "long"
                  }
                }
              },
              "xid": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "reply_size": {
                "type": "long"
              },
              "auth_flavor": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "call_size": {
                "type": "long"
              },
              "time": {
                "type": "long"
              },
              "time_str": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "status": {
                "ignore_above": 1024,
                "type": "keyword"
              }
            }
          },
          "amqp": {
            "properties": {
              "content-encoding": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "correlation-id": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "method-id": {
                "type": "long"
              },
              "no-wait": {
                "type": "boolean"
              },
              "reply-code": {
                "type": "long"
              },
              "consumer-tag": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "type": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "mandatory": {
                "type": "boolean"
              },
              "consumer-count": {
                "type": "long"
              },
              "durable": {
                "type": "boolean"
              },
              "class-id": {
                "type": "long"
              },
              "delivery-tag": {
                "type": "long"
              },
              "exclusive": {
                "type": "boolean"
              },
              "message-id": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "content-type": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "no-ack": {
                "type": "boolean"
              },
              "no-local": {
                "type": "boolean"
              },
              "reply-to": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "timestamp": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "headers": {
                "properties": {}
              },
              "message-count": {
                "type": "long"
              },
              "app-id": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "user-id": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "multiple": {
                "type": "boolean"
              },
              "if-unused": {
                "type": "boolean"
              },
              "priority": {
                "type": "long"
              },
              "passive": {
                "type": "boolean"
              },
              "redelivered": {
                "type": "boolean"
              },
              "delivery-mode": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "reply-text": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "auto-delete": {
                "type": "boolean"
              },
              "immediate": {
                "type": "boolean"
              },
              "arguments": {
                "properties": {}
              },
              "exchange": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "expiration": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "routing-key": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "exchange-type": {
                "ignore_above": 1024,
                "type": "keyword"
              },
              "if-empty": {
                "type": "boolean"
              },
              "queue": {
                "ignore_above": 1024,
                "type": "keyword"
              }
            }
          },
          "connecttime": {
            "type": "long"
          },
          "transport": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "client_proc": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "client_server": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "@timestamp": {
            "type": "date"
          },
          "connection_id": {
            "ignore_above": 1024,
            "type": "keyword"
          },
          "response": {
            "norms": false,
            "type": "text"
          },
          "service": {
            "ignore_above": 1024,
            "type": "keyword"
          }
        }
      }
    },
    "aliases": {}
  }
}
```

最后，查询匹配 packetbeat-* 的模版，不存在，查询匹配 packetbeat* 的模版，存在（因为默认模版名为 packetbeat）；

```
GET /_template/packetbeat-*

返回

{}
```

此时打开 Kibana Discover 界面仍会提示“Configure an index pattern”；

- 执行 import_dashboards 导入相关内容

```
➜  packetbeat git:(5.3) ✗ cd ../libbeat/dashboards
➜  dashboards git:(5.3) ✗
➜  dashboards git:(5.3) ✗ ll
total 25056
-rw-r--r--  1 sunfei  staff    91B  5 24 17:40 Makefile
drwxr-xr-x  7 sunfei  staff   238B  5 24 17:44 dashboards
-rwxr-xr-x  1 sunfei  staff    12M  5 12 17:47 import_dashboards
-rw-r--r--  1 sunfei  staff   7.3K  5 24 17:40 import_dashboards.go
➜  dashboards git:(5.3) ✗
➜  dashboards git:(5.3) ✗ ./import_dashboards -user elastic -pass changeme -url https://artifacts.elastic.co/downloads/beats/beats-dashboards/beats-dashboards-5.3.2.zip
Create temporary directory /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008
Downloading https://artifacts.elastic.co/downloads/beats/beats-dashboards/beats-dashboards-5.3.2.zip
Unzip archive /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008
Importing Kibana from /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat
Import directory /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/index-pattern
Import index to /.kibana/index-pattern/filebeat-* from /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/index-pattern/filebeat.json

Import directory /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/dashboard
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/dashboard/Filebeat-Apache2-Dashboard.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/visualization/Apache2-access-unique-IPs-map.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/search/Apache2-access-logs.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/visualization/Apache2-response-codes-of-top-URLs.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/search/Apache2-access-logs.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/visualization/Apache2-browsers.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/search/Apache2-access-logs.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/visualization/Apache2-operating-systems.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/search/Apache2-access-logs.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/visualization/Apache2-error-logs-over-time.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/search/Apache2-errors-log.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/visualization/Apache2-response-codes-over-time.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/search/Apache2-access-logs.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/search/Apache2-errors-log.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/dashboard/Filebeat-MySQL-Dashboard.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/visualization/MySQL-slowest-queries.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/search/Filebeat-MySQL-Slow-log.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/visualization/MySQL-Slow-queries-over-time.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/search/Filebeat-MySQL-Slow-log.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/visualization/MySQL-error-logs.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/search/Filebeat-MySQL-error-log.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/search/Filebeat-MySQL-error-log.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/visualization/MySQL-Error-logs-levels.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/search/Filebeat-MySQL-error-log.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/visualization/MySQL-Slow-logs-by-count.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/search/Filebeat-MySQL-Slow-log.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/dashboard/Filebeat-Nginx-Dashboard.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/visualization/Errors-over-time.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/visualization/Nginx-Access-Browsers.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/visualization/Nginx-Access-OSes.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/visualization/New-Visualization.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/search/Filebeat-Nginx-module.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/visualization/Nginx-Access-Response-codes-by-top-URLs.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/visualization/Sent-sizes.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/visualization/Nginx-Access-Map.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/search/Filebeat-Nginx-module.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/dashboard/Filebeat-syslog-dashboard.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/visualization/Syslog-events-by-hostname.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/search/Syslog-system-logs.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/visualization/Syslog-hostnames-and-processes.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/search/Syslog-system-logs.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/filebeat/search/Syslog-system-logs.json
Importing Kibana from /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/heartbeat
Import directory /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/heartbeat/index-pattern
Import index to /.kibana/index-pattern/heartbeat-* from /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/heartbeat/index-pattern/heartbeat.json

Import directory /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/heartbeat/dashboard
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/heartbeat/dashboard/f3e771c0-eb19-11e6-be20-559646f8b9ba.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/heartbeat/visualization/c65ef340-eb19-11e6-be20-559646f8b9ba.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/heartbeat/search/c49bd160-eb17-11e6-be20-559646f8b9ba.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/heartbeat/visualization/920e8140-eb1a-11e6-be20-559646f8b9ba.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/heartbeat/search/c49bd160-eb17-11e6-be20-559646f8b9ba.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/heartbeat/visualization/1738dbc0-eb1d-11e6-be20-559646f8b9ba.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/heartbeat/search/c49bd160-eb17-11e6-be20-559646f8b9ba.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/heartbeat/visualization/091c3a90-eb1e-11e6-be20-559646f8b9ba.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/heartbeat/search/c49bd160-eb17-11e6-be20-559646f8b9ba.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/heartbeat/visualization/0f4c0560-eb20-11e6-9f11-159ff202874a.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/heartbeat/search/c49bd160-eb17-11e6-be20-559646f8b9ba.json
Importing Kibana from /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat
Import directory /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/index-pattern
Import index to /.kibana/index-pattern/metricbeat-* from /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/index-pattern/metricbeat.json

Import directory /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/dashboard
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/dashboard/CPU-slash-Memory-per-container.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Container-CPU-usage.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/System-Navigation.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Container-Memory-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Container-Block-IO.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/dashboard/Metricbeat-Apache-HTTPD-server-status.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Apache-HTTPD-CPU.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Apache-HTTPD.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Apache-HTTPD-Hostname-list.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Apache-HTTPD.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Apache-HTTPD-Load1-slash-5-slash-15.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Apache-HTTPD.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Apache-HTTPD-Scoreboard.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Apache-HTTPD.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Apache-HTTPD-Total-accesses-and-kbytes.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Apache-HTTPD.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Apache-HTTPD-Uptime.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Apache-HTTPD.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Apache-HTTPD-Workers.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Apache-HTTPD.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/dashboard/Metricbeat-Docker.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Docker-containers.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Metricbeat-Docker.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Docker-Number-of-Containers.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Metricbeat-Docker.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Docker-containers-per-host.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Metricbeat-Docker.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Docker-images-and-names.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Metricbeat-Docker.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Docker-CPU-usage.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Docker-memory-usage.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Docker-Network-IO.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/dashboard/Metricbeat-MongoDB.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/MongoDB-hosts.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/MongoDB-search.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/MongoDB-Engine-ampersand-Version.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/MongoDB-search.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/MongoDB-operation-counters.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/MongoDB-search.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/MongoDB-Concurrent-transactions-Read.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/MongoDB-search.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/MongoDB-Concurrent-transactions-Write.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/MongoDB-search.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/MongoDB-memory-stats.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/MongoDB-search.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/MongoDB-asserts.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/MongoDB-search.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/MongoDB-WiredTiger-Cache.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/MongoDB-search.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/dashboard/Metricbeat-Redis.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Redis-Clients-Metrics.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Metricbeat-Redis.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Redis-Connected-clients.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Metricbeat-Redis.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Redis-hosts.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Metricbeat-Redis.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Redis-Server-Versions.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Metricbeat-Redis.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Redis-server-mode.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Metricbeat-Redis.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Redis-multiplexing-API.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Metricbeat-Redis.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Redis-Keyspaces.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Metricbeat-Redis.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/dashboard/Metricbeat-cpu.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/System-Navigation.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/CPU-usage-over-time.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Cpu-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/System-load.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Load-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/System-Load-over-time.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Load-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Top-hosts-by-CPU-usage.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Cpu-Load-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/CPU-Usage.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Cpu-stats.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/dashboard/Metricbeat-filesystem-per-Host.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Top-disks-by-memory-usage.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Filesystem-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Disk-utilization-over-time.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Filesystem-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/System-Navigation.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Disk-space-distribution.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Filesystem-stats.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/dashboard/Metricbeat-filesystem.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/System-Navigation.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Top-hosts-by-disk-size.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Fsstats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Disk-space-overview.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Fsstats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Free-disk-space-over-days.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Fsstats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Total-files-over-days.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Fsstats.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/dashboard/Metricbeat-memory.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/System-Navigation.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Top-hosts-by-memory-usage.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Memory-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Memory-usage-over-time.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Memory-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Swap-usage-over-time.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Memory-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Total-Memory.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Memory-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Available-Memory.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Memory-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Memory-usage.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Memory-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Swap-usage.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Memory-stats.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/dashboard/Metricbeat-network.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/In-vs-Out-Network-Bytes.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Network-data.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Top-10-interfaces.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Network-data.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Network-Packetloss.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Network-data.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Packet-loss-on-interfaces.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Network-data.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/System-Navigation.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Network-Bytes.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Network-data.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/dashboard/Metricbeat-overview.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Servers-overview.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/System-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/System-Navigation.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/dashboard/Metricbeat-processes.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/System-Navigation.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Number-of-processes.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Process-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Process-state-by-host.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Process-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Number-of-processes-by-host.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Process-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/CPU-usage-per-process.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Process-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Memory-usage-per-process.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Process-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Top-processes-by-memory-usage.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Process-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Top-processes-by-CPU-usage.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Process-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Number-of-processes-over-time.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Process-stats.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/dashboard/Metricbeat-system-overview.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Network-Bytes.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Network-data.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Network-Packetloss.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Network-data.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/System-Navigation.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Total-Memory.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Memory-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/Available-Memory.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Memory-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/System-overview-by-host.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/System-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/System-load.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Load-stats.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/visualization/CPU-Usage.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/metricbeat/search/Cpu-stats.json
Importing Kibana from /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat
Import directory /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/index-pattern
Import index to /.kibana/index-pattern/packetbeat-* from /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/index-pattern/packetbeat.json

Import directory /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/dashboard
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/dashboard/Packetbeat-Cassandra.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Cassandra-ResponseKeyspace.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Cassandra-ResponseType.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Cassandra-ResponseTime.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Cassandra-RequestCount.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Cassandra-Ops.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Cassandra-RequestCountStackByType.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Cassandra-ResponseCountStackByType.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Cassandra-RequestCountByType.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Cassandra-ResponseCountByType.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Navigation.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/Cassandra-QueryView.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/dashboard/Packetbeat-Dashboard.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Web-transactions.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/Web-transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/DB-transactions.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/DB-transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Cache-transactions.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/Cache-transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/RPC-transactions.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/RPC-transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Response-times-percentiles.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/Packetbeat-Search.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Errors-count-over-time.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/Transactions-errors.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Errors-vs-successful-transactions.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/Packetbeat-Search.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Latency-histogram.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/Packetbeat-Search.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Client-locations.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Response-times-repartition.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/Packetbeat-Search.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Navigation.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/dashboard/Packetbeat-Flows.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Navigation.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Connections-over-time.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/Packetbeat-Flows-Search.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Top-hosts-creating-traffic.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/Packetbeat-Flows-Search.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Top-hosts-receiving-traffic.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/Packetbeat-Flows-Search.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Network-traffic-between-your-hosts.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/Packetbeat-Flows-Search.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/dashboard/Packetbeat-HTTP.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Web-transactions.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/Web-transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/HTTP-error-codes.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/HTTP-error-codes-evolution.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Navigation.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Total-number-of-HTTP-transactions.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/Web-transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/HTTP-codes-for-the-top-queries.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/Web-transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Top-10-HTTP-requests.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/Web-transactions.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/dashboard/Packetbeat-MongoDB-performance.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Navigation.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/MongoDB-errors.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/MongoDB-errors.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/MongoDB-commands.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/MongoDB-transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/MongoDB-errors-per-collection.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/MongoDB-errors.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/MongoDB-in-slash-out-throughput.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/MongoDB-transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/MongoDB-response-times-by-collection.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/MongoDB-transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Top-slowest-MongoDB-queries.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/MongoDB-transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Number-of-MongoDB-transactions-with-writeConcern-w-equal-0.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/MongoDB-transactions-with-write-concern-0.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/dashboard/Packetbeat-MySQL-performance.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/MySQL-Errors.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/MySQL-errors.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/MySQL-Methods.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/MySQL-Transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Navigation.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/MySQL-throughput.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/MySQL-Transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Most-frequent-MySQL-queries.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/MySQL-Transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Slowest-MySQL-queries.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/MySQL-Transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Mysql-response-times-percentiles.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/MySQL-Transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/MySQL-Reads-vs-Writes.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/MySQL-Transactions.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/dashboard/Packetbeat-NFS.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/NFS-clients-pie-chart.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/nfs.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/NFS-operations-area-chart.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/nfs.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/NFS-top-group-pie-chart.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/nfs.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/NFS-top-users-pie-chart.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/nfs.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/NFS-response-times.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/nfs.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/NFS-errors.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/NFS-errors-search.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/NFS-operation-table.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/nfs.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/NFS-bytes-in-slash-out.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/nfs.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Navigation.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/dashboard/Packetbeat-PgSQL-performance.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Navigation.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/PgSQL-Errors.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/PgSQL-errors.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/PgSQL-Methods.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/PgSQL-transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/PgSQL-response-times-percentiles.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/PgSQL-transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/PgSQL-throughput.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/PgSQL-transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/PgSQL-Reads-vs-Writes.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/PgSQL-transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Most-frequent-PgSQL-queries.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/PgSQL-transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Slowest-PgSQL-queries.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/PgSQL-transactions.json
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/dashboard/Packetbeat-Thrift-performance.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Navigation.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Thrift-requests-per-minute.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/Thrift-transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Thrift-RPC-Errors.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/Thrift-errors.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Slowest-Thrift-RPC-methods.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/Thrift-transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Thrift-response-times-percentiles.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/Thrift-transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Top-Thrift-RPC-methods.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/Thrift-transactions.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/visualization/Top-Thrift-RPC-calls-with-errors.json
Import search /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/packetbeat/search/Thrift-errors.json
Importing Kibana from /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/winlogbeat
Import directory /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/winlogbeat/index-pattern
Import index to /.kibana/index-pattern/winlogbeat-* from /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/winlogbeat/index-pattern/winlogbeat.json

Import directory /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/winlogbeat/dashboard
Import dashboard /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/winlogbeat/dashboard/Winlogbeat-Dashboard.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/winlogbeat/visualization/Number-of-Events-Over-Time-By-Event-Log.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/winlogbeat/visualization/Number-of-Events.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/winlogbeat/visualization/Top-Event-IDs.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/winlogbeat/visualization/Event-Levels.json
Import visualization /var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/tmp414985008/beats-dashboards-5.3.2/winlogbeat/visualization/Sources.json
➜  dashboards git:(5.3) ✗
```

此时打开 Kibana Discover 界面，已经有 index pattern 可供选择了（默认会导入 `packetbeat-*/filebeat-*/heartbeat-*/metricbeat-*/winlogbeat-*`）；

- 将 `packetbeat-*` 选成默认的 index pattern ；
- 调整选择器范围以显示导入的 mysql 数据；
- 基于 ingest node pipeline 进行数据解析；

先通过 _simulate API 进行模拟

```
GET /_ingest/pipeline/_simulate
{
  "pipeline": {
    "processors": [
      {
        "grok": {
          "field": "query",
          "patterns": [
            "(%{DATA}) (E:rid=%{DATA:rid})&(rpcid=%{DATA:rpcid})&(appid=%{DATA:appid}):E (%{DATA}) (%{DATA:sql})$", "%{DATA:sql}$"
          ]
        }
      },
      {
      "remove": {
        "field": "bytes_in"
      }
    },
    {
      "remove": {
        "field": "bytes_out"
      }
    },
    {
      "remove": {
        "field": "beat"
      }
    },
    {
      "remove": {
        "field": "mysql"
      }
    },
    {
      "remove": {
        "field": "method"
      }
    },
    {
      "remove": {
        "field": "path"
      }
    },
    {
      "remove": {
        "field": "proc"
      }
    },
    {
      "remove": {
        "field": "client_proc"
      }
    },
    {
      "remove": {
        "field": "server"
      }
    },
    {
      "remove": {
        "field": "client_server"
      }
    }
    ]
  }
  ,
  "docs": [
    {
        "_index": "packetbeat-2017.05.16",
        "_type": "mysql",
        "_id": "AVw-eqUAbPwJc1ucUhBL",
        "_score": 1,
        "_source": {
          "@timestamp": "2017-05-16T03:14:05.277Z",
          "beat": {
            "hostname": "sunfeideMacBook-Pro.local",
            "name": "sunfeideMacBook-Pro.local",
            "version": "5.3.3"
          },
          "bytes_in": 152,
          "bytes_out": 77,
          "client_ip": "10.0.47.31",
          "client_port": 57094,
          "client_proc": "",
          "client_server": "",
          "ip": "10.0.27.37",
          "method": "/*",
          "mysql": {
            "affected_rows": 0,
            "error_code": 0,
            "error_message": "",
            "insert_id": 0,
            "iserror": false,
            "num_fields": 1,
            "num_rows": 1
          },
          "path": ".",
          "port": 3306,
          "proc": "",
          "query": "/* E:rid=zaocan.ele.me^^B7AE33373BA14CE8BA18F28945C2355D|1494904445271&rpcid=1.1.2.1&appid=me.ele.breakfast.mars:E */ select @@session.tx_read_only",
          "responsetime": 0,
          "server": "",
          "status": "OK",
          "type": "mysql"
        }
    },
    {
        "_index": "packetbeat-2017.05.16",
        "_type": "mysql",
        "_id": "AVw-eqUAbPwJc1ucUhBL",
        "_score": 1,
        "_source": {
          "@timestamp": "2017-05-16T03:14:05.277Z",
          "beat": {
            "hostname": "sunfeideMacBook-Pro.local",
            "name": "sunfeideMacBook-Pro.local",
            "version": "5.3.3"
          },
          "bytes_in": 152,
          "bytes_out": 77,
          "client_ip": "10.0.47.31",
          "client_port": 57094,
          "client_proc": "",
          "client_server": "",
          "ip": "10.0.27.37",
          "method": "/*",
          "mysql": {
            "affected_rows": 0,
            "error_code": 0,
            "error_message": "",
            "insert_id": 0,
            "iserror": false,
            "num_fields": 1,
            "num_rows": 1
          },
          "path": ".",
          "port": 3306,
          "proc": "",
          "query": "select @@session.tx_read_only",
          "responsetime": 0,
          "server": "",
          "status": "OK",
          "type": "mysql"
        }
    }
  ]
}

返回

{
  "docs": [
    {
      "doc": {
        "_id": "AVw-eqUAbPwJc1ucUhBL",
        "_type": "mysql",
        "_index": "packetbeat-2017.05.16",
        "_source": {
          "type": "mysql",
          "rid": "zaocan.ele.me^^B7AE33373BA14CE8BA18F28945C2355D|1494904445271",
          "sql": "select @@session.tx_read_only",
          "client_port": 57094,
          "responsetime": 0,
          "client_ip": "10.0.47.31",
          "ip": "10.0.27.37",
          "query": "/* E:rid=zaocan.ele.me^^B7AE33373BA14CE8BA18F28945C2355D|1494904445271&rpcid=1.1.2.1&appid=me.ele.breakfast.mars:E */ select @@session.tx_read_only",
          "@timestamp": "2017-05-16T03:14:05.277Z",
          "port": 3306,
          "appid": "me.ele.breakfast.mars",
          "rpcid": "1.1.2.1",
          "status": "OK"
        },
        "_ingest": {
          "timestamp": "2017-06-06T08:01:24.656Z"
        }
      }
    },
    {
      "doc": {
        "_id": "AVw-eqUAbPwJc1ucUhBL",
        "_type": "mysql",
        "_index": "packetbeat-2017.05.16",
        "_source": {
          "ip": "10.0.27.37",
          "query": "select @@session.tx_read_only",
          "type": "mysql",
          "sql": "select @@session.tx_read_only",
          "client_port": 57094,
          "@timestamp": "2017-05-16T03:14:05.277Z",
          "port": 3306,
          "responsetime": 0,
          "client_ip": "10.0.47.31",
          "status": "OK"
        },
        "_ingest": {
          "timestamp": "2017-06-06T08:01:24.656Z"
        }
      }
    }
  ]
}
```

在 elasticsearch 中创建名为 query_split 的 pipeline ；

```
PUT /_ingest/pipeline/query_split
{
  "processors": [
    {
      "grok": {
        "field": "query",
        "patterns": [
          "(%{DATA}) (E:rid=%{DATA:rid})&(rpcid=%{DATA:rpcid})&(appid=%{DATA:appid}):E (%{DATA}) (%{DATA:sql})$", "%{DATA:sql}$"
        ]
      }
    },
      {
      "remove": {
        "field": "bytes_in"
        }
      },
      {
        "remove": {
          "field": "bytes_out"
        }
      },
      {
        "remove": {
          "field": "beat"
        }
      },
      {
        "remove": {
          "field": "mysql"
        }
      },
      {
        "remove": {
          "field": "method"
        }
      },
      {
        "remove": {
          "field": "path"
        }
      },
      {
        "remove": {
          "field": "proc"
        }
      },
      {
        "remove": {
          "field": "client_proc"
        }
      },
      {
        "remove": {
          "field": "server"
        }
      },
      {
        "remove": {
          "field": "client_server"
        }
      }
  ]
}
```

之后需要修改 `packetbeat.yml` 配置文件，指定使用上面创建的 pipeline ；

若想仅查到“干净”数据（即后面 grok 过的），则需要在重新导入 mysql 数据前，删除已有数据内容（即 document 内容）；删除后，重新导入 mysql 数据，之后重新执行上面的步骤，就可以看到干净数据了；

> 针对 Elasticsearch 的 CURL 操作详见《[ElasticSearch cURL API](https://github.com/moooofly/MarkSomethingDown/blob/master/Golang/ElasticSearch%20cURL%20API.md)》


----------


## 常用测试 API

```
GET /_search
GET /packetbeat
GET /packetbeat-*
GET /packetbeat/_search
GET /packetbeat-*/_search
GET /packetbeat-*/mysql/_search

GET /_template
GET /_template/packetbeat
GET /_template/packetbeat*
GET /_template/packetbeat-*
```


----------

## Grok 示例

- Case 1

```
2015-03-27 16:04:31.701 warning zeus.eos.dispatcher[25859]: [zeus.eos 0.1.2 3DC67F4505FB4E0D8045FD1C5A4491B8] ## 订单金额过低 => make_order
```

对应

```
(?<datetime>\d{4}-\d{2}-\d{2} %{TIME}) (?<level>%{LOGLEVEL}) (?<logger>\S+)\[(?<thread>\S+)\]: \[(?<app_id>\S+) (?<rpc_request_id>.*)\] (?<ext_meta>.*)## (?<msg>.+)$
```

- Case 2

```
/* E:rid=zaocan.ele.me^^B7AE33373BA14CE8BA18F28945C2355D|1494904445271&rpcid=1.1.2.1&appid=me.ele.breakfast.mars:E */ select @@session.tx_read_only
```

对应

```
(E:rid=%{DATA:rid})&(rpcid=%{DATA:rpcid})&(appid=%{DATA:appid}):E
```

> 关于 Grok 的说明详见《[Elastic 之 Grok](https://github.com/moooofly/MarkSomethingDown/blob/master/Golang/Elastic%20%E4%B9%8B%20Grok.md)》

