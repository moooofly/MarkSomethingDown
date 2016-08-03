

当业务和 RabbitMQ 的消息交互量大到一定程度时，RabbitMQ 的 Web 管理页面 Overview 标签中会出现如下告警信息：

![RabbitMQ management 插件告警](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/managerment%20statistics%20database%20%E8%BF%87%E8%BD%BD%E9%97%AE%E9%A2%98.png "RabbitMQ management 插件告警")

（据说最严重的情况下，积压了几十万条消息）

上述内容提供了以下几点信息：
- management 插件在内存中创建统计数据库维护了大量 web 页面展示所需的相关数据；
- management 插件基于内部 gen_server2 行为模式，通过内部 buffer 和 message queue 的组合将消息进行缓存和有序处理，随着消息的增多，势必造成内存消耗的增加，针对统计信息处理的即时性变差；
- 设置 `rates_mode` 选项参数的值为 node 可能有所改善；


本文针对 management 管理插件的原理，以及在消息量大到一定程度时的行为进行展开；

----------

# management 插件相关代码研究

在 `overview.ejs` 中，可以看到输出上述告警信息的代码

```ejs
<div class="updatable">
<% if (overview.statistics_db_event_queue > 1000) { %>
<p class="warning">
  The management statistics database currently has a queue
  of <b><%= overview.statistics_db_event_queue %></b> events to
  process. If this number keeps increasing, so will the memory used by
  the management plugin.

  <% if (overview.rates_mode != 'none') { %>
  You may find it useful to set the <code>rates_mode</code> config item
  to <code>none</code>.
  <% } %>
</p>
<% } %>
</div>
```
其中，有两个重要的 if 判定：
- 如果 `overview.statistics_db_event_queue` 中的消息量超过 `1000` 条，就会在 Web 页面上输出之前的告警信息；
- 如果 `overview.rates_mode` 的值不是 `none` ，则建议改为 `none` ；


在 `rabbit_mgmt_db.erl` 中，能够看到 management 插件获取积压消息的代码

```erlang
%% 获取 Overview 页面所需信息       
handle_call({get_overview, User, Ranges}, _From,
            State = #state{tables = Tables}) ->
    ...
    %% 将统计信息返回给前端页面
    reply([{message_stats, format_samples(Ranges, MessageStats, State)},
           {queue_totals,  format_samples(Ranges, QueueStats, State)},
           {object_totals, ObjectTotals},
           %% 获取 rabbit_mgmt_db 进程中积压消息数量
           {statistics_db_event_queue, get(last_queue_length)}], State);
...
reply(Reply, NewState) -> {reply, Reply, NewState, hibernate}.
...
```

```erlang
%% 通过该回调函数保存待处理消息数目，能够保证即使当前处于消息过载状态，也能即时获取到数值；
%% Len 的值为当前待处理消息数量
prioritise_call(_Msg, _From, Len, _State) ->
    %% 将当前 rabbit_mgmt_db 进程邮箱中消息积压的待处理消息数量保存起来
    put(last_queue_length, Len),
    %% 通过数字设定当前消息优先级，优先级越高越先得到处理，默认优先级为 0
    5.
```


----------

# esm-agent 信息采集实现方式

从源码中可以看到，esm-agent 主要通过如下 HTTP API 获取统计信息

## Overview 信息

描述整个系统的各种信息；

```shell
/api/overview
```


  field   | value
--------- | ----------
cluster_name     | cluster 名，通过 rabbitmqctl set_cluster_name 进行设置；
contexts | cluster 中包含的 web 应用 context 列表；
erlang_full_version    | erlang 版本信息 ＋ Erlang VM 信息；
erlang_version     | erlang 版本信息；
exchange_types     | 可用的所有 exchange 类型列表；
listeners     | cluster 中所有节点上的 non-HTTP 的网络 listener ；
management_version     | 当前使用的 management 插件版本；
message_stats     | 用户能够看到的、对应 message_stats 对象的所有信息；该信息显示和权限有关；
node     | 当前 management 插件实例所运行的 cluster 节点名；
object_totals     | 包含针对所有 connections, channels, exchanges, queues 和 consumers 的全局计数值；权限约束同 message_stats ；
queue_totals     | 包含针对所有 queue 中处于 messages, messages_ready 和 messages_unacknowledged 状态消息的统计数据；权限约束同 message_stats ；
rabbitmq_version     | 处理该请求的、当前节点上运行的 RabbitMQ 版本；
rates_mode     | 'none', 'basic' 或 'detailed' ；
statistics_db_event_queue     | 统计数据库中待处理的 statistics events 的数量；
statistics_db_node     | 持有 management 插件统计数据库的 cluster 节点名；

## Queue 信息

所有 queue 信息列表；

```shell
/api/queues
```



## Node 信息

RabbitMQ cluster 中包含的所有节点信息；

```shell
/api/nodes
```

  field   | value
--------- | ----------
applications | 运行在当前节点上的所有 Erlang 应用
auth_mechanisms | 安装在当前节点上的所有 SASL 鉴权机制
cluster_links | cluster 中包含的其他节点列表；针对每个节点，都会输出到当前节点的详细 TCP 连接信息，以及针对其他数据发送的统计信息；
config_files | 当前节点所读取的配置文件列表
contexts | 当前节点上所有针对 HTTP 的 listener 列表
db_dir | 当前节点进行持久存储的位置
disk_free | Disk free space in bytes.
disk_free_alarm | 是否 disk 告警已经取消
disk_free_limit | disk 告警取消的阈值
enabled_plugins | 被显式使能且处于运行状态的插件列表
exchange_types | 当前节点上可用 exchange 类型
fd_total | 可用文件描述符数目
fd_used | 已用文件描述符数目
io_read_avg_time | 在上一次统计时间间隔内，每次 disk read 操作的平均 wall time（以毫秒为单位）
io_read_bytes | 由 persister 从 disk 上读取的总字节数
io_read_count | persister 进行 read 操作的总次数
io_reopen_count | persister 不得不在不同 queue 之间 recycle 文件句柄的总次数；在理想情况下，该值应该为 0 ；如果其值非常大，那么通过增大 RabbitMQ 可用文件句柄数目可以令性能得到改善；
io_seek_avg_time | 在上一次统计时间间隔内，每次 disk seek 操作的平均 wall time（以毫秒为单位）
io_seek_count | persister 进行 seek 操作的总次数
io_sync_avg_time | 在上一次统计时间间隔内，每次 fsync 操作的平均 wall time（以毫秒为单位）
io_sync_count | persister 进行 fsync 操作的总次数
io_write_avg_time | 在上一次统计时间间隔内，每次 disk write 操作的平均 wall time（以毫秒为单位）
io_write_bytes | 由 persister 向 disk 写入的总字节数
io_write_count | persister 进行 write 操作的总次数
log_file | main log 文件所在位置
mem_used | 内存使用，以字节为单位
mem_alarm | 是否内存告警已经取消
mem_limit | 内存告警被取消的阈值
mnesia_disk_tx_count | 请求进行 disk write 的 Mnesia 事务数（例如创建一个持久化 queue）；只有由当前节点发起的事务才被统计在内；
mnesia_ram_tx_count | 请求进行非 disk write 的 Mnesia 事务数（例如创建一个临时 queue）；只有由当前节点发起的事务才被统计在内；
msg_store_read_count | 从当前 message store 中读取的消息数量；
msg_store_write_count | 向当前 message store 中写入的消息数量；
name | 节点名字；
net_ticktime | 针对当前节点设置 net_ticktime 内核参数；
os_pid | 当前节点对应的操作系统 pid ；
partitions | 当前节点所看到的 network partitions 情况；
proc_total | 允许使用的 Erlang processe 最大数目；
proc_used | 当前已经被创建的 Erlang processe 数目；
processors | 被 Erlang 监测到可用的 CPU 核数目；
queue_index_journal_write_count | 被写入 queue index journal 中的记录数目；每条记录或者代表了被 publish 到 queue 上的消息，或者为从 queue 中 deliver 出去的消息，或者为对 queue 进行了确认的消息；
queue_index_read_count | 从 queue index 中读取的记录数目；
queue_index_write_count	 | 写入到 queue index 中的记录数目；
rates_mode | 'none', 'basic' 或 'detailed'.
run_queue | 待运行 Erlang processe 的平均数目；
running | 当前节点是否处于运行状态；很显然，如果该值为 false ，其他大部分统计信息将不会存在；
sasl_log_file | sasl 日志文件位置；
sockets_total | 可用于 socket 的文件描述符数目；
sockets_used | 已用做 socket 的文件描述符数目；
type | 当前节点的类型，'disc' 或 'ram' ；
uptime | 自 Erlang VM 启动以来过去的时间，以毫秒为单位；


----------

# management 插件使用中需要关注的点


## 插件的集群感知行为
management 插件对 cluster 是感知的；你可以在 cluster 中的某个或多个节点上启动该插件，之后通过 management 插件获取的信息将是与整个 cluster 相关的，无论你连接到 cluster 中的哪个节点；

如果你想要部署某个 cluster 节点，但不启动 management 插件的全部功能，仍然需要在每一个节点上启用 rabbitmq-management-agent 插件（这样才能通过特定节点获取到整个 cluster 的统计信息）；


## 统计数据库重启问题

统计数据库是整体保存在**内存**中的；因此其内容全部都是**临时性**的，外部访问者需要基于这个前提进行相应设计；通过重启统计数据库相关 erlang 进程，可以实现集群节点上迁移统计数据库的行为；

在 RabbitMQ 3.6.2 版本之前，统计数据库被直接保存在统计进程的内存中；
从 RabbitMQ 3.6.2 版本开始，统计数据库被保存在 ETS 表中；

在 RabbitMQ 3.6.2 版本之前，重启该数据库需要执行

```erlang
rabbitmqctl eval 'exit(erlang:whereis(rabbit_mgmt_db), please_terminate).'
```

从 RabbitMQ 3.6.2 版本开始，重启该数据库需要执行

```erlang
rabbitmqctl eval 'supervisor2:terminate_child(rabbit_mgmt_sup_sup, rabbit_mgmt_sup), rabbit_mgmt_sup_sup:start_child().'
```

无论如何，上述命令必须在统计数据库所在节点上执行才能生效；

针对该问题的详细说明，可以移步另外一篇总结：《[RabbitMQ management 插件数据库重置代价问题](https://github.com/moooofly/MarkSomethingDown/blob/master/RabbitMQ%20management%20%E6%8F%92%E4%BB%B6%E6%95%B0%E6%8D%AE%E5%BA%93%E9%87%8D%E7%BD%AE%E4%BB%A3%E4%BB%B7%E9%97%AE%E9%A2%98.md)》


## 内存管理问题

management 插件中统计数据库占用的内存情况可以通过如下命令获取：

```shell
# rabbitmqctl status
...
 {memory,
     [{total,54004424},
      ...
      {mgmt_db,381184},
      ...
      {other_system,4365602}]},
...
```

或者通过 HTTP API 发送 GET 请求到 `/api/nodes/<node_name>` 进行获取；

统计信息会按照 `collect_statistics_interval` 设置的时间间隔周期性采集；也可能在某些组件被创建/声明，或者关闭/销毁时进行采集（例如打开新 connection 或 channel，或者进行 queue 声明）；

消息速率的设置不会直接对 management 插件统计数据库内存占用产生影响；

**统计数据库占用内存的总量**取决于：
- 统计信息的采集时间间隔；
- effective rates ；
- retention 策略；

行之有效的调整方案：
- 将 `collect_statistics_interval` 的值调整到 30-60s ，将会显著减少维护大量 queues/channels/connections 的系统的内存消耗；
- 调整 retention 策略以减少留存的数据量也非常有效；

channel 以及统计信息收集进程的内存使用可以通过 `stats_event_max_backlog` 参数设置最大 backlog queue 大小进行限制；如果 backlog queue 已满，则新建 channel 信息和 queue 统计信息都会被丢弃，直到 backlog queue 上尚未处理的消息被处理；

> ⚠️ `stats_event_max_backlog` 参数在配置文件和代码中均未找到；

统计信息采集时间间隔支持运行时动态调整；进行调整不会对已存在的 connections, channels 或 queues 造成影响；仅对新加入的统计实体产生影响；

运行时调整命令如下
```shell
rabbitmqctl eval 'application:set_env(rabbit, collect_statistics_interval, 60000).'
```

可以通过重启统计数据库达成强行释放所占用内存的目的（当然会丢失一部分统计数据）；



----------


刚才反馈的 publish 等曲线掉底的问题，经过 @张斌 确认，结论如下：
1.在掉底曲线的时间段内，rabbitmq 的统计信息数据库（或队列）积压了几十万条统计信息；
2.同一时间段内，业务获取 channel 超时飙高；
3.张斌重启 rabbitmq 的 management 管理插件后（等于清空积压的统计信息），统计数据库从 xg-napos-rmq-1 节点随机迁移到 xg-napos-rmq-3 节点，此时发现，整体 qps 从原来的 2000 上升到 4000 左右；此时业务获取 channel 超时时间恢复正常；


所以，建议将 rabbitmq management 插件所使用的统计数据库部署到单独一个节点上，避免对业务造成影响；应该可以立刻取得改善；之后我会深入研究下 rabbitmq management 插件的使用和调优姿势，看看内否进一步改进


----------

结论推断：
- publish 等曲线掉底是由于 esm-agent 基于 management 插件提供的 HTTP API ，从 RabbitMQ 获取统计信息失败导致的；
- 获取统计信息失败是由于 rabbit_mgmt_db 中积压了过多消息导致的；
- rabbit_mgmt_db 中积压了过多消息是由于业务针对每条 publish 消息都创建和销毁 connection 和 channel 导致的；

业务获取 channel 超时曲线飙高可能原因：
- 由于 RabbitMQ 内部忙于处理统计信息相关内容导致；
- 由于业务采用了类似短连接的访问方式 ＋ 线上 goproxy 采用了不合理的健康监测 TCP 序列导致；


----------


# 补充

## aggregate data

> In statistics, `aggregate data` are data combined from several measurements.
> 
> `Aggregate data` refers to numerical or non-numerical information that is (1) collected from multiple sources and/or on multiple measures, variables, or individuals and (2) compiled into data summaries or summary reports, typically for the purposes of public reporting or statistical analysis
> 

## Memory footprint

> Memory footprint refers to the amount of main memory that a program uses or references while running.
> In computing, the memory footprint of an executable program indicates its runtime memory requirements, while the program executes. 



