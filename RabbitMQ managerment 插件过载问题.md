

当业务和 RabbitMQ 的消息交互量大到一定程度时，RabbitMQ 的 Web 管理页面 Overview 标签中会出现如下告警信息：

![RabbitMQ managerment 插件告警](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/managerment%20statistics%20database%20%E8%BF%87%E8%BD%BD%E9%97%AE%E9%A2%98.png "RabbitMQ managerment 插件告警")

（据说最严重的情况下，积压了几十万条消息）

上述内容提供了以下几点信息：
- managerment 插件通过一个名为 statistics 的数据库维护用于 web 页面展示的相关统计数据；
- managerment 插件使用了内部 queue 有序处理消息，随着 queue 中消息的增多，势必造成内存消耗的增加，统计信息的即时性变差，甚至可能对磁盘 I/O 造成影响（待确认）；
- 设置 `rates_mode` 选项参数的值为 node 可能有所改善；


本文针对 managerment 管理插件的原理，以及在消息量大到一定程度时的行为进行展开；


----------


在 overview.ejs 中，可以看到

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
两个重要的 if 判定：
- 如果 `overview.statistics_db_event_queue` 中的消息量超过 `1000` 条，就会在 Web 页面上输出之前的告警信息；
- 如果 `overview.rates_mode` 的值不是 `none` ，则建议改为 `none` ；

```erlang
%% 获取 Overview 页面所需信息       
handle_call({get_overview, User, Ranges}, _From,
            State = #state{tables = Tables}) ->
    ...
    %% 将统计信息返回给前端页面
    reply([{message_stats, format_samples(Ranges, MessageStats, State)},
           {queue_totals,  format_samples(Ranges, QueueStats, State)},
           {object_totals, ObjectTotals},
           %% 获取当前 rabbit_mgmt_db 进程中积压消息数量
           {statistics_db_event_queue, get(last_queue_length)}], State);
...
reply(Reply, NewState) -> {reply, Reply, NewState, hibernate}.
...
```

```erlang
%% 通过该回调函数保存待处理消息数目，能够保证即使当前处于消息过载状态，
%% 也能即时获取到数值；
%% Len 的值为当前待处理消息数量
prioritise_call(_Msg, _From, Len, _State) ->
    %% 将当前 rabbit_mgmt_db 进程邮箱中消息积压的待处理消息数量保存起来
    put(last_queue_length, Len),
    %% 通过数字设定当前消息优先级，优先级越高越先得到处理，默认优先级为 0
    5.
```

## RabbitMQ Management Plugin 配置项说明

```shell
 %% ----------------------------------------------------------------------------
 %% RabbitMQ Management Plugin
 %%
 %% See http://www.rabbitmq.com/management.html for details
 %% ----------------------------------------------------------------------------

 {rabbitmq_management,
  [%% Pre-Load schema definitions from the following JSON file
   %% {load_definitions, "/path/to/schema.json"},

   %% 将所有访问 management HTTP API 的请求记录到文件中
   %% {http_log_dir, "/path/to/access.log"},

   %% 配置 rabbitmq_management 插件的 HTTP 监听 IP 和 port
   %% Also set the listener to use SSL and provide SSL options.
   %%
   %% {listener, [{port,     12345},
   %%             {ip,       "127.0.0.1"},
   %%             {ssl,      true},
   %%             {ssl_opts, [{cacertfile, "/path/to/cacert.pem"},
   %%                         {certfile,   "/path/to/cert.pem"},
   %%                         {keyfile,    "/path/to/key.pem"}]}]},

   %% 可以设置为 'basic' 或 'detailed' 或 'none'
   %% {rates_mode, basic},

   %% 配置聚合数据被保留的时间长度；例如针对 message rates 和 queue 长度的聚合数据
   %%
   %% {sample_retention_policies,
   %%  [{global,   [{60, 5}, {3600, 60}, {86400, 1200}]},
   %%   {basic,    [{60, 5}, {3600, 60}]},
   %%   {detailed, [{10, 5}]}]}
  ]},
```

### 启动时加载预定义信息
management 插件允许你导出一个包含 broker 全部对象定义的 JSON 文件（对象包括：queues, exchanges, bindings, users, virtual hosts, permissions 和 parameters）；在一些场景中，每次启动时确保这些对象的存在是非常有必要的；

可以通过设置 `load_definitions` 变量的值为事先导出的 JSON 文件路径，来实现启动时加载；

需要注意的是，文件中定义的对象会覆盖 broker 中存在的相应对象；使用该选项不会删除已存在的其它对象；如果你启动的是一个完全重置过的 broker ，使用该选项将会阻止常规的 default user / virtual host / permissions 的创建；

### 消息速率
management 插件默认会展示全局消息速率 ，全局消息速率针对的是每个 queue, channel, exchange 和 vhost ；这种方式称作 `basic` 消息速率 ；

还可以针对所有组合，例如  channel to exchange, exchange to queue 以及 queue to channel ，进行消息速率展示；这种方式称作 `detailed` 消息速率 ；这种方式默认是关闭的，因为当系统中存在大量这种组合时，会导致大量的 memory footprint 出现；

最后一种选择是直接关闭消息速率显示；这样就可以在 CPU-bound 的服务器上获取最佳性能；

消息速率的模式是通过 rabbitmq_management 配置段中的 rates_mode 配置变量进行控制的；可以设置为 basic (默认值), detailed 或 none ；

### Statistics interval
默认情况下，服务器会每隔 5000ms 发送一次统计事件（包含各类统计数据）；而 management 插件所显示的消息速率值就是基于这个时间间隔计算得到的；

你可能在两种情况下会希望增加该时间间隔：
- 为了在一段更长的时间段内进行数据采样；
- 为了降低拥有大量 queue 或 channel 的服务器的统计负载；

可以通过 collect_statistics_interval 变量进行设置，单位为毫秒；设置后需要重启 RabbitMQ ；

### HTTP request logging
To create simple access logs of requests to the HTTP API, set the value of the http_log_dir variable in the rabbitmq_management application to the name of a directory in which logs can be created and restart RabbitMQ. Note that only requests to the API at /api are logged, not requests to the static files which make up the browser-based GUI.

### Events backlog
在高负载压力下，统计事件的处理会导致内存消耗量的增加；为了缓解这种情况，可以调整 channel 和 queue 统计信息采集器的最大 backlog 消息数量；在 rabbitmq_management 配置段中的 stats_event_max_backlog 变量值对应的就是 channel 和 queue 的最大 backlog 消息数量；默认为 250 ；





> In statistics, `aggregate data` are data combined from several measurements.
> 
> `Aggregate data` refers to numerical or non-numerical information that is (1) collected from multiple sources and/or on multiple measures, variables, or individuals and (2) compiled into data summaries or summary reports, typically for the purposes of public reporting or statistical analysis
> 


> Memory footprint refers to the amount of main memory that a program uses or references while running.
> In computing, the memory footprint of an executable program indicates its runtime memory requirements, while the program executes. 


刚才反馈的 publish 等曲线掉底的问题，经过 @张斌 确认，结论如下：
1.在掉底曲线的时间段内，rabbitmq 的统计信息数据库（或队列）积压了几十万条统计信息；
2.同一时间段内，业务获取 channel 超时飙高；
3.张斌重启 rabbitmq 的 managerment 管理插件后（等于清空积压的统计信息），统计数据库从 xg-napos-rmq-1 节点随机迁移到 xg-napos-rmq-3 节点，此时发现，整体 qps 从原来的 2000 上升到 4000 左右；此时业务获取 channel 超时时间恢复正常；


所以，建议将 rabbitmq managerment 插件所使用的统计数据库部署到单独一个节点上，避免对业务造成影响；应该可以立刻取的改善；


之后我会深入研究下 rabbitmq managerment 插件的使用和调优姿势，看看内否进一步改进


$sudo rabbitmqctl eval 'application:stop(rabbitmq_management), application:start(rabbitmq_management).'
or
$sudo rabbitmqctl eval 'exit(erlang:whereis(rabbit_mgmt_db), please_terminate).'
1 Comment
帮忙看下这两条命令的区别在哪