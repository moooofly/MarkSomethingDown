# influxdb 开发相关

## 环境准备

> 源码安装：[这里](https://github.com/influxdata/influxdb/blob/master/CONTRIBUTING.md)


### 问题

- 执行 `gdm restore` 时，报“Error getting VCS info for golang.org/x/crypto” 错误；

gdm 的工作原理：基于 Godeps 文件，通过 `git clone https://github.com/yourname/yourreponame /go/src/github.com/yourname/yourreponame` 获取相应的代码；之后再 `cd /go/src/github.com/yourname/yourreponame` 并切换到对应 hash 的分支 `git checkout e383bb****19b6`

基于其工作原理，就可以直接在 github 上找对应[源码的 mirror](https://github.com/golang) ，然后自行下载并切换到目标 hash 值；

- `go xxx ./...` 指令

```
# To build and install the binaries
go clean ./...
go install ./...

# To run the tests
go test -v ./...

go generate ./...

# To format the codes
go fmt ./...
go vet ./...
```


- `go tool pprof` 的使用

```
# start influx with profiling
influxd -cpuprofile influxdcpu.prof -memprofile influxdmem.prof

# open up pprof to examine the profiling data.
go tool pprof ./influxd influxdcpu.prof
```

> 需要安装 `apt-get install graphviz graphviz-dev` ；


## 概念


| influxDB 中的名词 | 传统数据库中的概念 |
| -- | -- |
| `database` | 数据库 |
| `measurement` | 数据库中的表 |
| `point` | 表里面的一行数据 |

Point 由时间戳（time）、数据（field）、标签（tags）组成。


| Point 属性 | 传统数据库中的概念 |
| - | - |
| time | 每个数据的记录时间，作为数据库中的**主索引**使用（会自动生成） |
| fields | 记录的值（各种没有索引的属性） |
| tags | 各种有索引的属性 |

`tag set`：tag 在 InfluxDB 中会按照字典序排序，不管是 tag key 还是 tag value ，只要不一致就分别属于两个 tag set ，例如 `hostname=server01,device=/data` 和 `hostname=server02,device=/data` 就属于两个不同的 tag set 。

还有三个重要的名词：Series、Retention policy 和 Shard ；

`Series`：相当于是 InfluxDB 中一些**数据的集合**；在同一个 database 中，retention policy、measurement、tag sets 完全相同的数据同属于一个 series ，**同一个 series 中的数据，在物理上会按照时间顺序排列存储在一起**。

`Retention policy`：存储策略，用于设置数据保留的时间，每个数据库刚开始会自动创建一个默认的存储策略 autogen ，数据保留时间为永久，之后用户可以自己设置，例如保留最近 2 小时的数据。插入和查询数据时如果不指定存储策略，则使用默认存储策略，且默认存储策略可以修改。InfluxDB 会定期清除过期的数据。

`Shard`：在 InfluxDB 中是一个比较重要的概念，它和 Retention policy 相关联。每一个存储策略下会存在许多 shard ，**每一个 shard 存储一个指定时间段内的数据**，并且不重复，例如 7~8 点的数据落入 shard0 中，8~9 点的数据则落入 shard1 中。**每一个 shard 都对应一个底层的 tsm 存储引擎**，有独立的 cache、wal、tsm file 。

TSM 存储引擎主要由几个部分组成：

- cache
- wal
- tsm file
- compactor

`Cache`：相当于是 LSM Tree 中的 memtabl 。插入数据时，实际上是同时往 cache 与 wal 中写入数据，可以认为 cache 是 wal 文件中的数据在内存中的缓存。当 InfluxDB 启动时，会遍历所有的 wal 文件，重新构造 cache ，这样即使系统出现故障，也不会导致数据的丢失。cache 中的数据并不是无限增长的，有一个 maxSize 参数用于控制当 cache 中的数据占用多少内存后就会将数据写入 tsm 文件。如果不配置的话，默认上限为 25MB ，每当 cache 中的数据达到阀值后，会将当前的 cache 进行一次快照，之后清空当前 cache 中的内容，再创建一个新的 wal 文件用于写入，剩下的 wal 文件最后会被删除，快照中的数据会经过排序写入一个新的 tsm 文件中。

`WAL`：WAL 文件的内容与内存中的 cache 相同，其作用就是为了持久化数据，当系统崩溃后可以通过 wal 文件恢复还没有写入到 tsm 文件中的数据。

`TSM File`：单个 tsm file 大小最大为 2GB ，用于存放数据。

`Compactor`：Compactor 组件在后台持续运行，每隔 1 秒会检查一次是否有需要压缩合并的数据。

主要进行两种操作，一种是 cache 中的数据大小达到阀值后，进行快照，之后转存到一个新的 tsm 文件中。另外一种就是合并当前的 tsm 文件，将多个小的 tsm 文件合并成一个，使每一个文件尽量达到单个文件的最大大小，减少文件的数量，并且一些数据的删除操作也是在这个时候完成。

## 操作相关

### Writing Data with the HTTP API

- create a database

```
curl -i -XPOST http://localhost:8086/query --data-urlencode "q=CREATE DATABASE mydb"
```

- writing a point

```
curl -i -XPOST 'http://localhost:8086/write?db=mydb' --data-binary 'cpu_load_short,host=server01,region=us-west value=0.64 1434055562000000000'
```

> The data consist of the `measurement` **cpu_load_short**, the `tag keys` **host** and **region** with the `tag values` **server01** and **us-west**, the `field key` **value** with a `field value` of **0.64**, and the `timestamp` **1434055562000000000**. 

> The body of the POST - we call this the **`Line Protocol`** - contains the time-series data that you wish to store. They consist of a `measurement`, `tags`, `fields`, and a `timestamp`. InfluxDB requires a measurement name. Strictly speaking, tags are optional but most series include tags to differentiate data sources and to make querying both easy and efficient. Both `tag keys` and `tag values` are **strings**. `Field keys` are required and are always **strings**, and, by default, `field values` are **floats**. The `timestamp` - supplied at the end of the line in Unix time in nanoseconds since January 1, 1970 UTC - is optional. If you do not specify a timestamp InfluxDB uses the server’s local nanosecond timestamp in Unix epoch. Anything that has to do with time in InfluxDB is always **UTC**.

- Writing multiple points

> Post multiple points to multiple `series` at the same time by separating each point with a new line. 

> The following example writes three points to the database mydb. The first point belongs to the `series` with the `measurement` **cpu_load_short** and `tag set` **host=server02** and has the server’s local timestamp. The second point belongs to the `series` with the `measurement` **cpu_load_short** and `tag set` **host=server02,region=us-west** and has the specified timestamp 1422568543702900257. The third point has the same specified timestamp as the second point, but it is written to the `series` with the `measurement` **cpu_load_short** and `tag set` **direction=in,host=server01,region=us-west**. 

```
curl -i -XPOST 'http://localhost:8086/write?db=mydb' --data-binary 'cpu_load_short,host=server02 value=0.67
cpu_load_short,host=server02,region=us-west value=0.55 1422568543702900257
cpu_load_short,direction=in,host=server01,region=us-west value=2.0 1422568543702900257'
```

- Writing points from a file

> a properly-formatted file (`cpu_data.txt`):

```
[#27#root@ubuntu ~]$cat cpu_data.txt
cpu_load_short,host=server02 value=0.67
cpu_load_short,host=server02,region=us-west value=0.55 1422568543702900257
cpu_load_short,direction=in,host=server01,region=us-west value=2.0 1422568543702900257
```

基于文件写入

```
curl -i -XPOST 'http://localhost:8086/write?db=mydb' --data-binary @cpu_data.txt
```

### 其他

- Schemaless Design
- The InfluxDB API makes no attempt to be RESTful.

## Querying Data with the HTTP API

- Send a query

> To perform a query send a **GET** request to the `/query` endpoint, set the URL parameter `db` as the target database, and set the URL parameter `q` as your query. 

```
curl -G 'http://localhost:8086/query?pretty=true' --data-urlencode "db=mydb" --data-urlencode "q=SELECT \"value\" FROM \"cpu_load_short\" WHERE \"region\"='us-west'"
```

- Send multiple queries

> Send multiple queries to InfluxDB in a single API call. Simply delimit each query using a **semicolon**.

```
curl -G 'http://localhost:8086/query?pretty=true' --data-urlencode "db=mydb" --data-urlencode "q=SELECT \"value\" FROM \"cpu_load_short\" WHERE \"region\"='us-west';SELECT count(\"value\") FROM \"cpu_load_short\" WHERE \"region\"='us-west'"
```

### Timestamp Format

> Everything in InfluxDB is stored and reported in `UTC`. By default, timestamps are returned in **RFC3339 UTC** and have `nanosecond` **precision**, for example 2015-08-04T19:05:14.318570484Z. If you want timestamps in `Unix epoch` format include in your request the query string parameter `epoch` where `epoch=[h,m,s,ms,u,ns]`. 

```
curl -G 'http://localhost:8086/query' --data-urlencode "db=mydb" --data-urlencode "epoch=s" --data-urlencode "q=SELECT \"value\" FROM \"cpu_load_short\" WHERE \"region\"='us-west'"
```

### Maximum Row Limit

> The `max-row-limit` configuration option allows users to limit the maximum number of returned results to **prevent InfluxDB from running out of memory** while it aggregates the results. The `max-row-limit` configuration option is set to 0 by default. That default setting allows for an **unlimited number of rows** returned per request. 

> The maximum row limit **only applies to non-chunked queries**. Chunked queries can return an unlimited number of points.


### Chunking

> Chunking can be used to return results in streamed batches rather than as a single response by setting the query string parameter `chunked=true`. Responses will be chunked **by series** or **by every 10,000 points**, whichever occurs first. To change the maximum chunk size to a different value, set the query string parameter `chunk_size` to a different value. For example, get your results in batches of 20,000 points with:

```
curl -G 'http://localhost:8086/query' --data-urlencode "db=deluge" --data-urlencode "chunked=true" --data-urlencode "chunk_size=20000" --data-urlencode "q=SELECT * FROM liters"
```

## InfluxQL Reference

- `influx` command line interface (CLI) is a lightweight and simple way to interact with the database. The CLI communicates with InfluxDB directly by making requests to the InfluxDB HTTP API over port 8086 by default.
- `influx` takes input in the form of the Influx Query Language (a.k.a **InfluxQL**) statements. 
- Data in InfluxDB is organized by “**time series**”, which contain a measured value, like “cpu_load” or “temperature”. **Time series** have zero to many `points`, one for each discrete sample of the metric.
- `Points` consist of `time`, a `measurement`, at least one key-value `field`, and zero to many key-value `tags` containing any metadata about the value.
- Conceptually you can think of a `measurement` as an SQL table, where the **primary index** is always time. `tags` and `fields` are effectively columns in the table. tags are indexed, and fields are not.
- `Points` are written to InfluxDB using the **Line Protocol**, which follows the following format: `<measurement>[,<tag-key>=<tag-value>...] <field-key>=<field-value>[,<field2-key>=<field2-value>...] [unix-nano-timestamp]`


### Query Engine Internals

Once you understand the language itself, it’s important to know how these language constructs are implemented in the query engine. This gives you an intuitive sense for how results will be processed and how to create efficient queries. 

The **life cycle of a query** looks like this:

- InfluxQL **query string is tokenized** and then **parsed into an abstract syntax tree (`AST`)**. This is the code representation of the query itself.
- The AST is passed to the **QueryExecutor** which **directs queries to the appropriate handlers**. For example, queries related to meta data are executed by the meta service and SELECT statements are executed by the shards themselves.
- The **query engine** then **determines the shards that match the SELECT statement’s time range**. From these shards, **`iterators` are created for each field in the statement**.
- **Iterators are passed to the `emitter` which drains them and joins the resulting points**. The emitter’s job is to convert simple time/value points into the more complex result objects that are returned to the client.


## line protocol

- [Line Protocol Tutorial](https://docs.influxdata.com/influxdb/v1.3/write_protocols/line_protocol_tutorial/)
- [Line Protocol Reference](https://docs.influxdata.com/influxdb/v1.3/write_protocols/line_protocol_reference/)

## 测试

```
CREATE DATABASE mydb
SHOW DATABASES
USE mydb

# To insert a single time-series datapoint into InfluxDB
INSERT cpu,host=serverA,region=us_west value=0.64

# query for the data we just wrote
SELECT "host", "region", "value" FROM "cpu"

# store the data with two fields in the same measurement
INSERT temperature,machine=unit42,type=assembly external=25,internal=37

# return all fields and tags with a query
SELECT * FROM "temperature"

# others
SELECT * FROM /.*/ LIMIT 1
SELECT * FROM "cpu_load_short"
SELECT * FROM "cpu_load_short" WHERE "value" > 0.9
```

```
----- 获取最近更新数据，并转换为当前时间
select threads_running from mysql order by time desc limit 1;
date -d @`echo 1483441750000000000 | awk '{print substr($0,1,10)}'` +"%Y-%m-%d %H:%M:%S"

----- 检查系统是否存活
curl -sl -I localhost:8086/ping

----- 简单查询
SELECT * FROM weather ORDER BY time DESC LIMIT 3;

----- 指定时间范围，时间格式也可以为'2017-01-03 00:00:00'
SELECT usage_idle FROM cpu WHERE time >= '2017-01-03T12:40:38.708Z' AND time <= '2017-01-03T12:40:50.708Z';

----- 最近40min内的数据
SELECT * FROM mysql WHERE time >= now() - 40m;

----- 最近5分钟的秒级差值
SELECT derivative("queries", 1s) AS "queries" from "mysql" where time > now() - 5m;

----- 最近5min的秒级写入
$ influx -database '_internal' -precision 'rfc3339'
      -execute 'select derivative(pointReq, 1s) from "write" where time > now() - 5m'

----- 也可以通过日志查看
$ grep 'POST' /var/log/influxdb/influxd.log | awk '{ print $10 }' | sort | uniq -c
$ journalctl -u influxdb.service | awk '/POST/ { print $10 }' | sort | uniq -c
```

> 以下内容取自 https://github.com/xwisen/ltw/issues/30

```
1. influxdb 并发写语句
#简单curl 
curl -i -XPOST 'http://127.0.0.1:9096/write?db=test' --data-binary 'cpu_load_short,host=server11,region=us-west value=0.64'
#使用ab进行并发操作
ab -c 10 -n 100 -T "application/x-www-form-urlencoded" -p test.post 'http://127.0.0.1:9096/write?db=test'
#test.post内容
cpu_load_short,host=server12,region=us-west value=0.64

2. 通用查询
select count(*)  from cpu_load_short;
time influx -database 'dcos' -execute 'show series from container' -format 'csv'
time influx -host 127.0.0.1 -port 8086 -precision 'rfc3339' -database 'dcos' -execute 'show tag keys from container' -format 'csv'
time influx -host 127.0.0.1 -port 8086 -precision 'rfc3339' -database 'dcos' -execute 'show field keys from container' -format 'csv'
3. 排序/限制/条件
select * from cpu_load_short order by time desc limit 100;
select * from cpu_load_short order by time asc limit 100;
select * from container where time > now() - 5m and container_appid='mg-admin' order by time desc;
select container_thread_running,container_name,ipaddress from container where container_appid='scrm-web' and container_thread_running > 70 order by time;

4. 聚合查询
select max(container_thread_running),container_name,ipaddress from container where container_appid='scrm-app' order by time;
select top(container_thread_running,10),container_name,ipaddress from container where container_appid='scrm-app' and time > now() - 5m order by time;

5. grafana 中使用influxdb作为数据源
SHOW TAG KEYS FROM "container"
SHOW FIELD KEYS FROM "container"
SHOW TAG VALUES FROM "container" WITH KEY = "container_appid"	
SHOW TAG VALUES FROM "container" WITH KEY = "container_name" WHERE container_appid =~ /^$container_appid$/

6. influxdb 数据保存策略
#修改默认RP之后, 之前的数据会有'丢失',需要指定RP才能查询旧数据
#参考:https://docs.influxdata.com/influxdb/v0.9/troubleshooting/frequently_encountered_issues/#missing-data-after-creating-a-new-default-retention-policy
time influx -host 127.0.0.1 -port 8086 -precision 'rfc3339' -database 'dcos' -execute 'select * from dcos.autogen.container where time > now() - 15h and time < now() - 12h order by time desc' -format 'csv' > hours.log
#创建策略
create retention policy "rp_name" on "db_name" duration 3w replication 1 default 
#查询库数据保存策略
influx -host 127.0.0.1 -port 8086 -precision 'rfc3339' -database 'dcos' -execute 'show retention policies on dcos' -format 'csv'
#创建数据保存策略
influx -host 127.0.0.1 -port 8086 -precision 'rfc3339' -database 'dcos' -execute 'create retention policy "wz1" on "dcos" duration 3w replication 1 default' -format 'csv'
#修改数据保存策略
influx -host 127.0.0.1 -port 8086 -precision 'rfc3339' -database 'dcos' -execute 'alter retention policy "wz1" on "dcos" duration 60d' -format 'csv'
#删除数据保存策略 DROP RETENTION POLICY "rp_name" ON "db_name"
influx -host 127.0.0.1 -port 8086 -precision 'rfc3339' -database 'dcos' -execute 'drop retention policy "wz1" on "dcos"' -format 'csv'
#查询shard groups
#https://docs.influxdata.com/influxdb/v1.0/troubleshooting/frequently-asked-questions/#what-is-the-relationship-between-shard-group-durations-and-retention-policies
influx -host 127.0.0.1 -port 8086 -precision 'rfc3339' -database 'dcos' -execute 'show shard groups' -format 'csv'
influx -host 127.0.0.1 -port 8086 -precision 'rfc3339' -database 'dcos' -execute 'drop shard ${id}' -format 'csv'
7. 连续查询
#查看连续查询
influx -host 127.0.0.1 -port 8086 -precision 'rfc3339' -database 'dcos' -execute 'show continuous queries' -format 'csv'
#创建连续查询
influx -host 127.0.0.1 -port 8086 -precision 'rfc3339' -database 'dcos' -execute 'create continuous query cq_10m on dcos begin select mean(cpu) into tb10m from container group by time(10m) end' -format 'csv'
#删除连续查询
influx -host 127.0.0.1 -port 8086 -precision 'rfc3339' -database 'dcos' -execute 'drop continuous query cq_10m on dcos' -format 'csv'
8.正则表达式
# =~注意空格
time influx -host 127.0.0.1 -port 8086 -precision 'rfc3339' -database 'dcos' -execute 'select mean(container_mem_used) from dcos.wz1.container where time > now() - 60m and time < now() - 30m and container_appid =~ /^bomc-lcgl$/  group by time(1m),"container_appid","container_name" order by time desc' -format 'csv' > hours.log
```

## 参考

- [InfluxDB](https://jin-yang.github.io/post/influxdata-influxdb.html)
- [玩转时序数据库InfluxDB](http://www.ywnds.com/?p=10763)
- [linux: influxdb 常用语法记录](https://github.com/xwisen/ltw/issues/30)
- [InfluxDB安装及配置](http://www.cnblogs.com/MikeZhang/p/InfluxDBInstall20170206.html)


