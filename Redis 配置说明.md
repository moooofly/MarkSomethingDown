

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

slave-serve-stale-data yes
slave-read-only yes
repl-diskless-sync no
repl-diskless-sync-delay 5
repl-disable-tcp-nodelay no
repl-backlog-size 32mb
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
