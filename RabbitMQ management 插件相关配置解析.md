


# management 插件相关配置项

## rabbitmq_management.app.src 中的配置项

此处为 management 插件默认启动参数设置；

```erlang
  {env, [{listener,          [{port, 15672}]},
         {http_log_dir,      none},
         {load_definitions,  none},
         {rates_mode,        basic},
         {sample_retention_policies,
          %% List of {MaxAgeInSeconds, SampleEveryNSeconds}
          [{global,   [{605, 5}, {3660, 60}, {29400, 600}, {86400, 1800}]},
           {basic,    [{605, 5}, {3600, 60}]},
           {detailed, [{10, 5}]}]}
        ]},
```


## rabbitmq.config 中的中的配置项

此为常规配置文件内容（覆盖 .app.src 文件中的配置内容）

```erlang
 %% ----------------------------------------------------------------------------
 %% RabbitMQ Management Plugin
 %%
 %% See http://www.rabbitmq.com/management.html for details
 %% ----------------------------------------------------------------------------

 {rabbitmq_management,
  [%% 可基于 JSON 文件启动时预先加载的 schema 定义信息
   %% {load_definitions, "/path/to/schema.json"},

   %% 将所有访问 management HTTP API 的请求记录到文件中
   %% {http_log_dir, "/path/to/access.log"},

   %% 配置 rabbitmq_management 插件的 HTTP 监听 IP 和 port
   %% 可以配置基于 SSL 的连接
   %%
   %% {listener, [{port,     12345},
   %%             {ip,       "127.0.0.1"},
   %%             {ssl,      true},
   %%             {ssl_opts, [{cacertfile, "/path/to/cacert.pem"},
   %%                         {certfile,   "/path/to/cert.pem"},
   %%                         {keyfile,    "/path/to/key.pem"}]}]},

   %% 可以设置为 'basic' 或 'detailed' 或 'none'
   %% {rates_mode, basic},

   %% 配置聚合数据被保留的时间长度；例如针对消息速率和 queue 长度的聚合数据
   %% {sample_retention_policies,
   %%  [{global,   [{60, 5}, {3600, 60}, {86400, 1200}]},
   %%   {basic,    [{60, 5}, {3600, 60}]},
   %%   {detailed, [{10, 5}]}]}
  ]},
```


## 和统计信息有关的其它配置项

```erlang
  {rabbit,
   ...
   %%
   %% Misc/Advanced Options
   %% =====================
   ...
   %% 设置（内部）统计信息采集粒度
   %%
   %% {collect_statistics, none},

   %% 统计信息采集时间间隔，以毫秒为单位
   %%
   %% {collect_statistics_interval, 5000},
   ...
```

## 参数解析

### load_definitions - 启动时加载预定义信息
management 插件允许你导出一个包含 broker 全部对象定义的 JSON 文件（对象包括：queues, exchanges, bindings, users, virtual hosts, permissions 和 parameters）；在一些场景中，每次启动时确保这些对象的存在是非常有必要的；

可以通过设置 `load_definitions` 变量的值为事先导出的 JSON 文件路径，来实现启动时加载；

需要注意的是，文件中定义的对象会覆盖 broker 中存在的相应对象；使用该选项不会删除已存在的其它对象；如果你启动的是一个完全重置过的 broker ，使用该选项将会阻止常规的 default user / virtual host / permissions 的创建；

### rates_mode - 消息速率的模式
management 插件默认会展示全局消息速率 ，全局消息速率针对的是所有 queue, channel, exchange 和 vhost ；这种方式称作 `basic` 消息速率模式 ；

还可以针对所有组合，例如  channel to exchange, exchange to queue 以及 queue to channel ，进行消息速率展示；这种方式称作 `detailed` 消息速率模式 ；这种方式默认是关闭的，因为当系统中存在大量这种组合时，会导致大量的 memory footprint 出现；

最后一种选择是直接关闭消息速率显示；这样就可以在 CPU-bound 的服务器上获取最佳性能；

消息速率的模式是通过 rabbitmq_management 配置段中的 `rates_mode` 配置变量进行控制的；可以设置为 `basic` (默认值), `detailed` 或 `none` ；


### collect_statistics

统计信息收集模式；主要和 management 插件有关；
可配置选项包括：
- `none` - 不发送 statistics 事件；
- `coarse` - 发送针对 per-queue / per-channel / per-connection 的统计信息；
- `fine` - 发送针对 per-queue / per-channel / per-connection / per-message 的统计信息；

该选项默认值为 none ；在不理解该参数含义的情况下，不建议修改；


### collect_statistics_interval - 统计信息采集时间间隔
默认情况下，服务器会每隔 5000ms 发送一次统计事件（包含各类统计数据）；而 management 插件所显示的消息速率值就是基于这个时间间隔计算得到的；

> 注意：此处的统计信息采集时间间隔与 web 页面上刷新页面时间间隔（默认 5s）是两回事；

你可能在两种情况下会希望增加该时间间隔：
- 为了在一段更长的时间段内进行数据采样；
- 为了降低拥有大量 queue 或 channel 的服务器的统计负载；

可以通过 collect_statistics_interval 变量进行设置，单位为毫秒；设置后需要重启 RabbitMQ ；

### http_log_dir - HTTP 请求日记记录
创建记录请求 HTTP API 时的简单访问日志；设置 http_log_dir 变量为保存相应日志的目录名，之后需要重启 RabbitMQ ；需要注意的是，只有针对 /api 的请求会被记录；默认关闭；

> 结论：个人感觉这个日志的最大用途是用来确认 web 上每个 tab 也都使用哪些 HTTP API 来获取展示数据的；

### stats_event_max_backlog - 允许事件 backlog 数目
在高负载压力下，统计事件的处理会导致内存消耗量的增加；为了缓解这种情况，可以调整 channel 和 queue 统计信息采集器的最大 backlog 消息数量；在 rabbitmq_management 配置段中的 stats_event_max_backlog 变量值对应的就是 channel 和 queue 的最大 backlog 消息数量；默认为 250 ；
> 注意：该配置参数在代码和 rabbitmq.config 文件中均未找到；


### sample_retention_policies - 采样＋保留策略
management 插件会保留一些数据采样值，例如针对消息速率和 queue 长度信息；可以通过如下配置项定制具体策略；

存在 3 种策略类型：
- global - 针对 overview 和 virtual host 页面定制策略；
- basic - 针对单独的 connections, channels, exchanges 和 queues 定制策略；
- detailed - 针对消息速率为不同的 connections, channels, exchanges 和 queues 组合定制策略；

参数 `{MaxAgeInSeconds, SampleEveryNSeconds}` 中，SampleEveryNSeconds 表示每 N 秒采样一次；MaxAgeInSeconds 表示采样数据的最长保留时间；
