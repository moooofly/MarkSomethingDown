

# Redis cluster

基于 redis-3.0.3 版本；

## master 节点配置

```shell
daemonize no
port 7000
tcp-backlog 65535
timeout 600
tcp-keepalive 60

# 日志直接走 syslog
loglevel notice
logfile /dev/null
syslog-enabled yes
syslog-ident redis
syslog-facility local0

# 集群节点只能使用 0 号数据库
databases 1

stop-writes-on-bgsave-error yes
rdbcompression yes
rdbchecksum yes
dbfilename dump.rdb
dir /data/redis/cluster/7000

# 确保 slave 在任何情况下都会响应 client 请求（可能回复过期数据）
slave-serve-stale-data yes

# 设置 slave 节点为只读（如何确保 slave 不暴露）
slave-read-only yes

repl-diskless-sync no
repl-diskless-sync-delay 5
repl-disable-tcp-nodelay no
repl-backlog-size 32mb

# 用于 Sentinel 机制中将 slave 提升为 master 时的优先级判定
# 该值越小，优先级越高；
# 0 优先级特殊处理，标识不可提升为 master
# 问题：该优先级在 cluster 中将 slave 提升为 master 时起作用么？
slave-priority 100

rename-command KEYS "ele-super-keys"
rename-command CONFIG "ele-super-config"
maxclients 50000
maxmemory 2g
maxmemory-policy allkeys-lru
appendonly no
appendfilename "appendonly.aof"
appendfsync everysec
no-appendfsync-on-rewrite no
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb
aof-load-truncated yes
lua-time-limit 5000
cluster-enabled yes
cluster-config-file nodes.conf
cluster-node-timeout 15000
cluster-require-full-coverage no
slowlog-log-slower-than 10000
slowlog-max-len 1024
latency-monitor-threshold 0
notify-keyspace-events ""
hash-max-ziplist-entries 512
hash-max-ziplist-value 64
list-max-ziplist-entries 512
list-max-ziplist-value 64
set-max-intset-entries 512
zset-max-ziplist-entries 128
zset-max-ziplist-value 64
hll-sparse-max-bytes 3000
activerehashing yes
client-output-buffer-limit normal 0 0 0
client-output-buffer-limit slave 2gb 0 0
client-output-buffer-limit pubsub 32mb 8mb 60
hz 10
aof-rewrite-incremental-fsync yes
```


## slave 节点配置


```shell
daemonize no
port 7101
tcp-backlog 65535
timeout 600
tcp-keepalive 60
loglevel notice
logfile /dev/null
syslog-enabled yes
syslog-ident redis
syslog-facility local0
databases 1
save 300 1
save 120 10
save 30 10000
stop-writes-on-bgsave-error yes
rdbcompression yes
rdbchecksum yes
dbfilename dump.rdb
dir /data/redis/cluster/7101
slave-serve-stale-data yes
slave-read-only yes
repl-diskless-sync no
repl-diskless-sync-delay 5
repl-disable-tcp-nodelay no
slave-priority 100
rename-command KEYS "ele-super-keys"
rename-command CONFIG "ele-super-config"
maxclients 50000
maxmemory 1gb
maxmemory-policy allkeys-lru
appendonly yes
appendfilename "appendonly.aof"
appendfsync everysec
no-appendfsync-on-rewrite no
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb
aof-load-truncated yes
lua-time-limit 5000
cluster-enabled yes
cluster-config-file nodes.conf
cluster-node-timeout 15000
cluster-require-full-coverage no
slowlog-log-slower-than 10000
slowlog-max-len 128
latency-monitor-threshold 0
notify-keyspace-events ""
hash-max-ziplist-entries 512
hash-max-ziplist-value 64
list-max-ziplist-entries 512
list-max-ziplist-value 64
set-max-intset-entries 512
zset-max-ziplist-entries 128
zset-max-ziplist-value 64
hll-sparse-max-bytes 3000
activerehashing yes
client-output-buffer-limit normal 0 0 0
client-output-buffer-limit slave 256mb 64mb 60
client-output-buffer-limit pubsub 32mb 8mb 60
hz 10
aof-rewrite-incremental-fsync yes
```

## 配置说明

### slave-serve-stale-data

当 slave 与 master 连接断开，或正处于 replication 行为进行之中时，slave 会根据配置具有两种不同的行为：
- slave-serve-stale-data 设置为 `yes`（默认值）；slave 仍将继续应答来自 client 的请求，但可能回应过期数据，或者直接回应空数据集（如果是处于首次同步过程中）；
- slave-serve-stale-data 设置为 `no`；slave 将直接回复错误消息 "SYNC with master in progress" 给各种命令请求，除留 INFO 和 SLAVEOF 外；


### slave-read-only

- 允许配置 slave 实例是否接受 write 请求；
- 向 slave 实例进行写操作在某些场景下是有意义的，如需要写入临时数据的场景（写入 slave 的数据很容易通过与 master 的 resync 清除掉）；
- 需要避免因为错误配置造成的向 slave 进行写操作的情况（主从切换时需要注意）；
- 从 Redis 2.6 开始 slave 默认配置成 read-only ；

> ⚠️ read-only slave 同样不可以随便暴露给未授信的因特网 client ；
> ⚠️ 应该将 read-only 属性当作防止 redis 实例被误用的保护层；
> ⚠️ read-only slave 同样会默认暴露全部管理命令，例如 CONFIG, DEBUG 等等；

为了进行可用命令限定，可以利用 'rename-command' 方式 shadow 所有的管理/危险命令，增加 read-only slave 的安全性；


### slave-priority

- 该整数值用于标识 slave 的优先级，可以在 INFO 命令的结果中看到；
- 该优先级配置值被 Redis Sentinel 用作将 slave 提升为 master 的参考；
- 优先级数值越低，对应的 slave 越被认为适合提升；
- 特殊优先级 0 标识对应 slave 不允许被提升为 master ；
- 默认优先级数值为 100 ；

