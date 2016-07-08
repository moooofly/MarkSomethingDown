

```shell
### Welcome to the InfluxDB configuration file.

# Once every 24 hours InfluxDB will report anonymous data to m.influxdb.com
# The data includes raft id (random 8 bytes), os, arch, version, and metadata.
# We don't track ip addresses of servers reporting. This is only used
# to track the number of instances running and the versions, which
# is very helpful for us.
# Change this option to true to disable reporting.
reporting-disabled = false

# we'll try to get the hostname automatically, but if it the os returns something
# that isn't resolvable by other servers in the cluster, use this option to
# manually set the hostname
# hostname = "localhost"

###
### [meta]
###
### 针对 Raft consensus group 的参数进行设置；其中保存了和 InfluxDB cluster
### 相关的 metadata
###

[meta]
  # 控制当前 node 是否应该运行 metaservice 并加入到 Raft 组中
  enabled = true

  # metadata/raft database 被保存的位置
  dir = "/var/lib/influxdb/meta"

  # metaservice 默认绑定的 tcp 通信地址
  bind-address = ":8088"

  # metaservice 供外部访问的 HTTP API 地址
  http-bind-address = ":8091"
  https-enabled = false
  https-certificate = ""

  retention-autocreate = true
  
  # 用于 store 的默认选举超时
  election-timeout = "1s"
  # 用于 store 的默认心跳超时
  heartbeat-timeout = "1s"

  # 用于 store 的默认 leader 租期
  leader-lease-timeout = "500ms"

  # 用于 store 的默认 commit 超时
  commit-timeout = "50ms"

  # 针对 meta service 是否打印跟踪日志消息
  cluster-tracing = false

  # 在必要时，是否自动提升一个普通 node 成为 raft node
  raft-promotion-enabled = true

  # 针对 meta service 是否打印日志消息
  logging-enabled = true
  pprof-enabled = false

  # 默认租期时长
  lease-duration = "1m0s"

###
### [data]
###
### 针对 InfluxDB 的 shard 数据所在位置进行控制；
### 针对 shard 数据如何从 WAL 中 flush 到磁盘进行控制；
### "dir" 可能需要按照系统的实际情况，变更到一个合适的位置；
### WAL 属于高级设置范畴，其默认值对于大多数系统来说都能正常工作；
###

[data]
  # 决定当前 node 是否在 cluster 中持有保存时间序列数据的 shard ；
  enabled = true

  dir = "/var/lib/influxdb/data"

  # 下面这些 WAL 设置针对的是 storage engine >= 0.9.3
  wal-dir = "/var/lib/influxdb/wal"
  wal-logging-enabled = true
  data-logging-enabled = true

  # 是否在 query 执行之前进行日志记录；
  # 该设置对于排查问题非常有用，但是可能会将包含在 query 中的敏感数据也记录进去
  # query-log-enabled = true

  # 针对 TSM 引擎的设置

  # 在 shard 所使用的 cache 开始拒绝 write 之前允许使用的最大内存量
  # cache-max-memory-size = 524288000

  # TSM 引擎开始将 cache 中的内容进行 snapshot ，再写入 TSM 文件，之后释放内存的临界值
  # cache-snapshot-memory-size = 26214400

  # 在当前 shard 没有接收到任何 write 或 delete 请求的情况下，
  # 若达到该参数指定的时间长度，则引擎会将 cache 进行 snapshot 并写入
  # 一个新的 TSM 文件中；
  # cache-snapshot-write-cold-duration = "1h"

  # 在一个压缩 cycle 运行前，TSM 文件需要存在的最少数量
  # compact-min-file-count = 3

  # 在当前 shard 没有接收到任何 write 或 delete 请求的情况下，
  # 引擎开始压缩当前 shard 中全部 TSM 文件的时间长度临界值；
  # compact-full-write-cold-duration = "24h"

  # 在一个 TSM 文件的编码块中，允许存在的 point 的最大数目；
  # 更大的数值可能可以产生更好的压缩效果，但是可能导致 query 时的性能损失
  # max-points-per-block = 1000

###
### [hinted-handoff]
###
### 控制 hinted handoff 特性；
### 该特性允许当前 node 在 cluster 中的某个 node 短时间 down 掉的情况下，
### 临时针对 queued data 进行存储； 
###

[hinted-handoff]
  enabled = true
  dir = "/var/lib/influxdb/hh"
  max-size = 1073741824
  max-age = "168h"
  retry-rate-limit = 0

  # Hinted handoff 特性会以每秒一次的速率对已 down 掉 node 进行写重试；
  # 若重试过程中发生了错误，将按照指数退避掉方式进行重试速度调整，直到时间间隔达到
  # retry-max-interval 设定的值；
  # 一旦针对所有 node 的写操作都成功完成，时间间隔将被重置回 retry-interval 设置的值；
  retry-interval = "1s"
  retry-max-interval = "1m"

  # 针对 data 是否应该被 purge 所运行的 check 的时间间隔；
  # data 会基于两种原因从 hinted-handoff 队列中被 purge 出来
  # 1) data 比 max age 还要 older
  # 2) 目标 node 被 cluster 给 drop 掉了
  # data 在 max-age 到达前绝对不会被 drop ，但针对被 drop 的 node 来说不成立；
  purge-interval = "1h"

###
### [cluster]
###
### 控制 non-Raft cluster 行为，主要是关于数据如何进行 shard 的设置
###

[cluster]
  shard-writer-timeout = "5s"  # remote shard 必须应答写请求的超时时间
  write-timeout = "10s"        # 写请求必须在 cluster 中完成的超时时间

###
### [retention]
###
### Controls the enforcement of retention policies for evicting old data.
### 针对老旧数据进行淘汰的 retention 策略控制
###

[retention]
  enabled = true
  check-interval = "30m"

###
### [shard-precreation]
###
### Controls the precreation of shards, so they are available before data arrives.
### 针对 shard 的预创建进行控制，因此在数据到达前 shard 就已经被创建好了

### Only shards that, after creation, will have both a start- and end-time in ### the future, will ever be created.
### 在创建后，只有 shard 将会同时具有 start-time 和（未来的） end-time 属性；

### Shards are never precreated that would be wholly or partially in the past.
### shard 不可能在全部或者部分属于过去的情况下被预创建；

[shard-precreation]
  enabled = true
  check-interval = "10m"
  advance-period = "30m"

###
### 控制系统的自我监控，信息统计和诊断功能的使用
###
### 用于保存监控数据的内部数据库会在不存在的情况下自动创建； 
### 在上述被创建数据库内的 target retention 被称作 "monitor" ；
### 同样会创建一个 7 天的 retention period 和值为 1 的复制因子（如果不存在的话）；
### 在所有情况下，retention policy 都会作为数据库的默认配置而存在

[monitor]
  store-enabled = true           # 是否在内部进行统计信息记录
  store-database = "_internal"   # 记录统计信息的数据库
  store-interval = "10s"         # 统计信息记录的时间间隔

###
### [admin]
###
### 控制内置的，基于 web 的 admin 接口是否可用；
### 如果针对 admin 接口使能了 HTTPS ，则必须同时在 [http] 段中使能 HTTPS ；
###

[admin]
  enabled = true
  bind-address = ":8083"
  https-enabled = false
  https-certificate = "/etc/ssl/influxdb.pem"

###
### [http]
###
### 针对 HTTP 实体进行配置；
### 此为向 InfluxDB 写入或从其中读取数据的主要手段
###

[http]
  enabled = true
  bind-address = ":8086"
  auth-enabled = false
  log-enabled = true
  write-tracing = false
  pprof-enabled = false
  https-enabled = false
  https-certificate = "/etc/ssl/influxdb.pem"

###
### [[graphite]]
###
### 控制针对 Graphite 数据的 listener 使用
###

[[graphite]]
  enabled = false
  # database = "graphite"
  # bind-address = ":2003"
  # protocol = "tcp"
  # consistency-level = "one"

  # These next lines control how batching works. You should have this enabled
  # otherwise you could get dropped metrics or poor performance. Batching
  # will buffer points in memory if you have many coming in.

  # batch-size = 5000 # will flush if this many points get buffered
  # batch-pending = 10 # number of batches that may be pending in memory
  # batch-timeout = "1s" # will flush at least this often even if we haven't hit buffer limit
  # udp-read-buffer = 0 # UDP Read buffer size, 0 means OS default. UDP listener will fail if set above OS max.

  ### This string joins multiple matching 'measurement' values providing more control over the final measurement name.
  # separator = "."

  ### Default tags that will be added to all metrics.  These can be overridden at the template level
  ### or by tags extracted from metric
  # tags = ["region=us-east", "zone=1c"]

  ### Each template line requires a template pattern.  It can have an optional
  ### filter before the template and separated by spaces.  It can also have optional extra
  ### tags following the template.  Multiple tags should be separated by commas and no spaces
  ### similar to the line protocol format.  There can be only one default template.
  # templates = [
  #   "*.app env.service.resource.measurement",
  #   # Default template
  #   "server.*",
  # ]

###
### [collectd]
###
### 控制针对 collectd 数据的 listener 使用
###

[collectd]
  enabled = false
  # bind-address = ""
  # database = ""
  # typesdb = ""

  # These next lines control how batching works. You should have this enabled
  # otherwise you could get dropped metrics or poor performance. Batching
  # will buffer points in memory if you have many coming in.

  # batch-size = 1000 # will flush if this many points get buffered
  # batch-pending = 5 # number of batches that may be pending in memory
  # batch-timeout = "1s" # will flush at least this often even if we haven't hit buffer limit
  # read-buffer = 0 # UDP Read buffer size, 0 means OS default. UDP listener will fail if set above OS max.

###
### [opentsdb]
###
### 控制针对 OpenTSDB 数据的 listener 使用
###

[opentsdb]
  enabled = false
  # bind-address = ":4242"
  # database = "opentsdb"
  # retention-policy = ""
  # consistency-level = "one"
  # tls-enabled = false
  # certificate= ""
  # log-point-errors = true # Log an error for every malformed point.

  # These next lines control how batching works. You should have this enabled
  # otherwise you could get dropped metrics or poor performance. Only points
  # metrics received over the telnet protocol undergo batching.

  # batch-size = 1000 # will flush if this many points get buffered
  # batch-pending = 5 # number of batches that may be pending in memory
  # batch-timeout = "1s" # will flush at least this often even if we haven't hit buffer limit

###
### [[udp]]
###
### 控制针对 InfluxDB line protocol data via UDP 的 listener 使用
###

[[udp]]
  enabled = false
  # bind-address = ""
  # database = "udp"
  # retention-policy = ""

  # These next lines control how batching works. You should have this enabled
  # otherwise you could get dropped metrics or poor performance. Batching
  # will buffer points in memory if you have many coming in.

  # batch-size = 1000 # will flush if this many points get buffered
  # batch-pending = 5 # number of batches that may be pending in memory
  # batch-timeout = "1s" # 在未达到缓冲区限制的情况下，也要执行 flush 操作的时间间隔
  # read-buffer = 0 # UDP 读缓冲区大小，0 表示由 OS 决定；若设置的值超过了 OS 允许设置的最大值，UDP listener 将会失效；

  # 设置预估的 UDP 负载大小；
  # 更小的值可能会产生更好的性能，默认值为 UDP 允许的最大大小，即 65536
  # udp-payload-size = 65536

###
### [continuous_queries]
###
### 控制 continuous queries 在 InfluxDB 中如何运行
###

[continuous_queries]
  log-enabled = true
  enabled = true
  # run-interval = "1s" # 在 continuous queries 需要运行的情况下，需要对其进行 check 的 时间间隔
```