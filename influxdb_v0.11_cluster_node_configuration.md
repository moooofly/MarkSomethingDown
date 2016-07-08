


> Warning! This page documents an old version of InfluxDB, which is no longer actively developed. 
> [InfluxDB v0.13]() is the most recent stable version of InfluxDB.


----------


下面的内容描述了 InfluxDB cluster 中 不同类型的 node ，以及如何对其进行配置；

# General node configuration

每一个 node 的 [配置文件]() 都必须制定下面各项：

`[meta]` 段中的 `bind-address`；该配置为 cluster 通信使用的地址；
`[meta]` 段中的 `http-bind-address`；该配置为 consensus service 通信使用的地址；
`[http]` 段中的 `bind-address`； 该配置为 HTTP API 访问使用的地址；

每一条配置选项都应该指明 node 的 IP 地址或主机名，和 port 信息（详见下文示例）；

> NOTE: The hostnames for each machine must be resolvable by all members of the cluster.

# Consensus node

Consensus nodes 仅用于运行 consensus service ；consensus service 用于确保 cluster 中以下内容的一致性：node membership, [databases](), [retention policies](), [users](), [continuous queries](), shard metadata, 和 [subscriptions]() ；


配置如下：

```shell
[meta]
  enabled = true #✨
  ...
  bind-address = "<IP>:8088"
  http-bind-address = "<IP>:8091"

...

[data]
  enabled = false #✨

[http]
  ...
  bind-address = "<IP>:8086"
```
  
# Data node

Data nodes 仅用于运行 data service ；data service 保存了实际的时间序列数据，并相应针对那些数据的查询请求；

配置如下：

```shell
[meta]
  enabled = false #✨
  ...
  bind-address = "<IP>:8088"
  http-bind-address = "<IP>:8091"

...

[data]
  enabled = true #✨

[http]
  ...
  bind-address = "<IP>:8086"
```
  
# Hybrid node

Hybrid nodes 同时运行了 consensus 和 data services ；

配置如下：

```shell
[meta]
  enabled = true #✨
  ...
  bind-address = "<IP>:8088"
  http-bind-address = "<IP>:8091"

...

[data]
  enabled = true #✨

[http]
  ...
  bind-address = "<IP>:8086"
```