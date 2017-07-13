

# Redis cluster

基于 redis-3.0.3 版本；

## master 节点配置

```shell

# 为什么设置成 no 而不是 yes
daemonize no
port 7000
tcp-backlog 65535

# 若客户端空闲超过 600 秒则关闭对应连接
# 设置为 0 表示去使能
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

# 是否在上一次 bgsave 出错后不再允许写操作执行
stop-writes-on-bgsave-error yes

# 是否在生成 RDB 文件时使用 LZF 算法压缩字符串对象
rdbcompression yes

rdbchecksum yes
dbfilename dump.rdb
dir /data/redis/cluster/7000

# 确保 slave 在任何情况下都会响应 client 请求（可能回复过期数据）
slave-serve-stale-data yes

# 设置 slave 节点为只读（如何确保 slave 不暴露）
# 问题：应该在 slave 上设置才对吧？！难道是为主从切换做准备？！
slave-read-only yes

# 主从复制同步策略（disk or socket）
# no 表示同步策略为 disk
repl-diskless-sync no

# 当基于 socket 进行复制同步时，开始传输 rdb 前的等待时间
repl-diskless-sync-delay 5

# 针对复制同步，用于在低延时和高吞吐量之间进行 trade-off
# no 对应低延时
repl-disable-tcp-nodelay no

# 在 slave 断开时间内，master 用于保存 slave 数据同步信息的缓冲区大小
repl-backlog-size 32mb

# 用于 Sentinel 机制中将 slave 提升为 master 时的优先级判定
# 该值越小，优先级越高；
# 0 优先级特殊处理，标识不可提升为 master
# 问题：该优先级在 cluster 模式下，将 slave 提升为 master 时起作用么？
slave-priority 100

# 命令重命名（防止误用和攻击）
rename-command KEYS "aaa-bbb-keys"
rename-command CONFIG "aaa-bbb-config"

maxclients 50000

# 设置当前实例的内存使用上限（超过该上限时会根据 LRU 策略淘汰相应的 keys）
maxmemory 2g
maxmemory-policy allkeys-lru

# 未启动 AOF 
appendonly no
# 问题：在 master 上不启用 AOF 时下面的配置都无用，难道是为主从切换做准备？！
appendfilename "appendonly.aof"
appendfsync everysec

# 控制在 AOF 重写时是否允许 fsync 被调用
# fsync 的调用可能导致 write(2) 的阻塞
no-appendfsync-on-rewrite no

# 控制 BGREWRITEAOF 的触发条件
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb

aof-load-truncated yes
lua-time-limit 5000

# 作为集群节点运行
cluster-enabled yes

# 集群节点自己生成并维护的配置信息
cluster-config-file nodes.conf

# 判定集群节点失效的超时时间
cluster-node-timeout 15000

# 即使出现部分 slot 未被 cover ，cluster 仍能对外提供服务
cluster-require-full-coverage no

# 慢日志
slowlog-log-slower-than 10000
slowlog-max-len 1024

# latency 监控子系统
# 设置为 0 表示关闭
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

# 设置针对不同角色客户端的输出缓冲区限制，用于强制断开客户端连接
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

# 日志直接走 syslog
loglevel notice
logfile /dev/null
syslog-enabled yes
syslog-ident redis
syslog-facility local0

databases 1

# BGSAVE 触发条件
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

maxmemory 2gb
maxmemory-policy allkeys-lru

# 启用 AOF 持久化
appendonly yes
appendfilename "appendonly.aof"
appendfsync everysec

# AOF 重写配置
no-appendfsync-on-rewrite no
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb

aof-load-truncated yes

lua-time-limit 5000

# 集群配置
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

### stop-writes-on-bgsave-error

- 默认情况下，在使能了 RDB 快照功能后，Redis 会在最后一次 BGSAVE 失败的情况下（至少要存在一个 save point），拒绝接受写入操作；
- 通过这种方式可以令用户立刻发觉数据无法正常持久化到磁盘的情况，避免在将来的某个时候发生灾难；
- 如果负责 BGSAVE 的子进程再次重新开始工作，那么 Redis 将会自动允许写操作进行；
- 如果你自行构建了针对 Redis 服务和持久化功能的监控机制，你可能会希望去使能该功能，以便令 Redis 在即使遇到磁盘持久化问题时，也能像通常一样继续工作；


### slave-serve-stale-data

当 slave 与 master 连接断开，或正处于 replication 行为进行之中时，slave 会根据配置具有两种不同的行为：
- slave-serve-stale-data 设置为 `yes`（默认值）；slave 仍将继续应答来自 client 的请求，但可能回应过期数据，或者直接回应空数据集（如果是处于首次同步过程中）；
- slave-serve-stale-data 设置为 `no`；slave 将直接回复错误消息 "SYNC with master in progress" 给各种命令请求，除留 INFO 和 SLAVEOF 外；


### slave-read-only

- 配置 slave 实例是否接受 write 请求；
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


### rename-command

```shell

# It is possible to change the name of dangerous commands in a shared
# environment. For instance the CONFIG command may be renamed into something
# hard to guess so that it will still be available for internal-use tools
# but not available for general clients.
#
# Example:
#
# rename-command CONFIG b840fc02d524045429941cc15f59e41cb7be6c52
#
# It is also possible to completely kill a command by renaming it into
# an empty string:
#
# rename-command CONFIG ""
#
# Please note that changing the name of commands that are logged into the
# AOF file or transmitted to slaves may cause problems.
```

### repl-backlog-size

- 该配置项用于设置复制（同步）缓冲区大小；
- backlog 表示在 slave 断开时间内，master 用于保存 slave 数据同步信息的缓冲区大小；
- 当发生 slave 重连时，全量重同步可能不是必须的，因为部分重同步可能就足够了；此时只需传输连接断开时 slave 缺失的那部分数据变更；
- 设置的复制缓冲区越大，允许 slave 在断开后，通过部分重同步进行恢复的时间窗口就越长；
- 至少有一个 slave 连接上 master 时，才会分配  backlog 对应的空间；

### repl-ping-slave-period

Slave 以预定义的时间间隔发送 PING 到主服务器； 默认值为 10 秒；

### maxmemory

- 限定可用的最大内存量；
- 当内存上限被达到时，Redis 将会根据 eviction 策略进行 key 的移除；
- 如果 Redis 根据策略无法移除 keys ，或者策略被设置成 'noeviction'，那么 Redis 将会对需要消耗内存的命令（例如 SET,LPUSH 等）回复错误应答，而对只读命令（例如 GET）回复成功应答；
- 将 Redis 用作 LRU cache 时，该选项非常有用；或者也可用于为单个实例设置内存使用上限（使用 'noeviction' 策略）；

```shell
# WARNING: If you have slaves attached to an instance with maxmemory on,
# the size of the output buffers needed to feed the slaves are subtracted
# from the used memory count, so that network problems / resyncs will
# not trigger a loop where keys are evicted, and in turn the output
# buffer of slaves is full with DELs of keys evicted triggering the deletion
# of more keys, and so forth until the database is completely emptied.
```

- 简单来说，如果当前节点是配置了 slave 的，则建议设置一个更低的 maxmemory 值，以便系统中留有一些空余的 RAM 用于 slave 的 output buffer（在策略配置成 'noeviction' 时，则没有保留的必要）

### maxmemory-policy

MAXMEMORY POLICY: 在达到 maxmemory 设定的值时，决定了 Redis 移除内容的选择方式；

存在五种可以配置的策略：
- volatile-lru -> remove the key with an expire set using an LRU algorithm
- allkeys-lru -> remove any key according to the LRU algorithm
- volatile-random -> remove a random key with an expire set
- allkeys-random -> remove a random key, any key
- volatile-ttl -> remove the key with the nearest expire time (minor TTL)
- noeviction -> don't expire at all, just return an error on write operations（默认值）

⚠️ 在上述任意一种策略配置下，当不存在符合淘汰条件的键时，Redis 将对写操作返回错误；



### cluster-config-file

- 每一个 cluster 节点都有一个集群配置文件；
- 该文件不应该手动进行编辑；
- 该文件由 Redis 节点自行创建和更新；
- 每一个 Redis Cluster 节点都要求一个单独的集群配置文件；

请确保运行在同一个系统中的不同实例不会因为集群配置文件名字相同而相互覆盖；


### cluster-node-timeout

- 定义判定集群节点处于失效状态（`PFAIL`）的前，允许的最长不可达时间；以毫秒为单位；
- 大部分内部时间限定值都是该值的倍数；


### cluster-slave-validity-factor

从属于失效 master 的 slave ，在发现自身数据过于老旧的情况下，将不会进行 failover 处理；

无法通过一种简单的方式准确测量出 slave 中“数据的年龄“，因此实际中会执行如下两种检测：
- 如果存在多个 slave 能够进行 failover ，那么它们之间会通过信息交换的方式，确定具有最佳复制偏移位置的 slave（即持有更多来自 master 的数据）；全部 slave 将会根据偏移量进行排名，并在 failover 启动前增加一个正比于排名顺序的延时值；
- 每个 slave 都需要计算自身与 master 上次交互的时间点；该时间点可能为最后一次 ping 发生的时刻，或者接收到其他集群消息到时刻（如果 master 仍旧处于 "connected" 状态），或者自从和 master 断开连接后，到目前为止流逝的时间（如果用于进行复制的链路当前是 down 状态）；如果最后一次交互发生的时间过于久远，对应的 slave 将不再进行 failover 行为；

上述第 2 点可以由用户进行调节；特别是，如果一个 slave 自上次与 master 交互后，已流逝的时间超过了如下公式对应的数值时，将会不再执行 failover 操作：

    (node-timeout * slave-validity-factor) + repl-ping-slave-period

例如，如果 node-timeout 设置为 30 秒，并且 slave-validity-factor 设置为 10 ，并假定 repl-ping-slave-period 采用默认值 10 秒，那么当 slave 与 master 无法通信的时间超过 310 秒时，将不会执行 failover 操作；

设置更大的 slave-validity-factor 值等价于允许 slaves 使用更加老旧的数据通过 failover 方式成为 master ，而设置更小的值，则可能造成 cluster 无法在规定时间内成功选出合适的 slave 进行提升；

从最大可用性角度考虑，将 slave-validity-factor 设置为 0 也是一种可能情况，意味着对应的 slaves 将总是会尝试通过 failover 接管 master ，而不管其最后一次与 master 的交互时间；（然而，即使设置为 0 ，也还是会添加正比与偏移量排名的延时值）

只有设置成 0 值，才能确保当全部分区被治愈后 cluster 总是能继续工作；



### cluster-require-full-coverage

- 默认情况下，Redis Cluster 中的节点会停止接受查询请求，如果检测到存在至少一个 hash slot 未被 cover（即没有节点负责该 slot）；
- 在这种情况下，如果 cluster 中的节点出现部分 down 掉的情况，会造成一定范围内的 hash slots 未被 cover 的情况，此时会导致 cluster 中的所有节点都不可用，即 cluster 失效；
- 只要所有 slot 能够被重新 cover 到，cluster 就会自动变回可用状态；
- 然而，在某些情况下，你可能希望在出现问题时，cluster 的某个子集仍能工作，即针对特定 key 空间范围（cover 的部分）的请求继续提供服务，此时可以将 cluster-require-full-coverage 设置为 no ；


### no-appendfsync-on-rewrite

当 AOF 的 fsync 策略设置为 always 或者 everysec 时，并且存在后台进程正进行大量磁盘 I/O 时（对应 BGSAVE 或 BGREWRITEAOF），在某些 Linux 配置下，Redis 可能会因为 fsync() 调用的原因发生长时间阻塞；
需要注意的是，当前没有针对此问题的解决办法，因为即使在不同的线程中执行 fsync ，也同样会阻塞 write(2) 这个同步调用；

为了缓解这个问题，可以通过当前配置选项，在存在执行中的 BGSAVE 或 BGREWRITEAOF 时，阻止 fsync() 在主进程中被调用；

在配置了该功能后，就意味着尽管存在另一个子进程正在保存信息，Redis 的持久性仍等价于 "appendfsync none"；从实际的角度来看，这意味着在最差场景下，可能会丢失多达 30 秒的日志信息（基于默认的 Linux 设置）；

如果你需要解决延迟问题，就设置该选项为 "yes" ，否则设置为 "no" （从持久性的角度最安全的选择）；


### client-output-buffer-limit

定义客户端输出缓冲区（buffer）限制值，用于强制断开由于某种原因无法快速从服务器上读取数据的客户端连接（一种常见的原因为：Pub/Sub 客户端消费消息的速度无法跟上生产者产生消息的速度）；

可以针对三种不同类型（class）的客户端设置不同的限制方式：
- **normal** -> 包括 MONITOR 在内的所有普通客户端；
- **slave**  -> 作为 slave 的客户端；
- **pubsub** -> 订阅到至少一个 pubsub channel 或 pattern 的客户端；

配置指令如下：

    client-output-buffer-limit <class> <hard limit> <soft limit> <soft seconds>

当 hard limit 被达到时，客户端会被立即断开连接；    
当 soft limit 被达到时，并保持在此达到状态特定时间长度后（连续），客户端被断开连接；    

举例来说，如果 hard limit 为 32 MB ，而 soft limit 为 16 MB 和 10 秒持续时间，则一旦客户端输出缓冲区大小达到 32 MB ，则会被立即断开连接；另外，若客户端输出缓冲区达到 16 MB ，并且持续达到此限制长达 10 秒，则会被断开连接；

默认情况下，普通客户端不受限制，因为在未做要求的情况下，其不会进行数据接收（基于 push 方式），但是即使进行发起了请求，也只有异步客户端才可能创造出数据请求快于其读取速度的场景；

作为对比，对于 pubsub 和 slave 客户端是存在默认限制的，因为订阅者和 slaves 都是以一种 push 方式接收数据的；

hard limit 或 soft limit 都可以通过设置成 0 以去使能；



