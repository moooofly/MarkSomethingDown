

> Warning! This page documents an old version of InfluxDB, which is no longer actively developed. 
> [InfluxDB v0.13]() is the most recent stable version of InfluxDB.


----------


> 注意: InfluxDB 0.11 是包含 cluster 功能的最后一个开源版本（此处的 cluster 是指 raft）；更多关于此问题的信息请阅读 Paul Dix 的博客文章[InfluxDB Clustering, High-Availability, and Monetization](). 请注意，0.11 版本的 cluster 仍被认为是实验性的，仍存在很多犄角旮旯的地方需要处理。

本文简要介绍了 InfluxDB 的 cluster 模型，并提供了逐步建立 cluster 的步骤。

# InfluxDB cluster model

InfluxDB 支持创建任意（节点）数量的 cluster ，并且允许 [replication factor]() 的值被设置为 1 至 cluster 中的 node 数目；

在 InfluxDB cluster 中存在 3 种类型的 node ：[consensus nodes](), [data nodes]() 和 [hybrid nodes]() 。一个 cluster 中要求必须存在奇数个 node 运行 [consensus service]() 以便构成 Raft consensus group ，且保持在一个健康的状态中。

Hardware requirements vary for the different node types. See Hardware Sizing for cluster [hardware requirements]().

# Cluster setup

The following steps configure and start up an InfluxDB cluster with three [hybrid nodes](). If you’re interested in having any of the different node types, see [Cluster Node Configuration]() for their configuration details. Note that your first three nodes must be either hybrid nodes or consensus nodes.

We assume that you are running some version of Linux, and, while it is possible to build a cluster on a single server, it is not recommended.

1  在 3 台机器上[安装]() InfluxDB ；但暂时不要在任何一台机器上启动 daemon 程序

2  配置 3 个 node

其中 `IP` 对应的是 node 的 IP 地址或主机名，每一个 node 均需要一个 `/etc/influxdb/influxdb.conf` 文件，其中包含如下设置信息：

```shell
[meta]
  enabled = true
  ...
  bind-address = "<IP>:8088"
  ...
  http-bind-address = "<IP>:8091"

...

[data]
  enabled = true

...

[http]
  ...
  bind-address = "<IP>:8086"
```

同时设置 `[meta] enabled = true` 和 `[data] enabled = true` 将使得当前 node 成为一个 hybrid node ；
`[meta] bind-address` 用作当前 node  在 cluster 中的通信地址；
`[meta] http-bind-address` 用作 meta node 通信地址；
`[http] bind-address` 用作 HTTP API 通信地址；

> 注意：针对每台机器配置的主机名必须能够被 cluster 中的所有 node 所解析；

3  将所有 node 都配置成可以与其它 node 进行通信;

针对上述 3 个 node ，在 `/etc/default/influxdb` 中设置 `INFLUXD_OPTS` 的值如下：

```shell
INFLUXD_OPTS="-join <IP1>:8091,<IP2>:8091,<IP3>:8091"
```

其中 `IP1` 为第一个 node 的 IP 地址或主机名，`IP2` 为第二个 node 的 IP 地址或主机名，`IP3` 为第二个 node 的 IP 地址或主机名；

如果文件 `/etc/default/influxdb` 不存在，则需手动创建；

4  在每一个 node 上都启动 InfluxDB ：

```shell
sudo service influxdb start
```

5  校验 cluster 的健康性

可以使用 `influx` [CLI]() 工具向 cluster 中的每一个 node 发送 `SHOW SERVERS` 查询指令；对应的输出信息显示出你的 cluster 由 3 个 hybrid node 构成（hybrid node 在 `SHOW SERVERS` 的查询结果中既作为 `data_nodes` 又作为 `meta_nodes` 显示）：

```shell
> SHOW SERVERS
name: data_nodes
----------------
id	 http_addr		  tcp_addr
1	  <IP1>:8086	  <IP1>:8088
2	  <IP2>:8086	  <IP2>:8088
3	  <IP3>:8086	  <IP3>:8088


name: meta_nodes
----------------
id	 http_addr		  tcp_addr
1	  <IP1>:8091	  <IP1>:8088
2	  <IP2>:8091	  <IP2>:8088
3	  <IP3>:8091	  <IP3>:8088
```

> Note: The `SHOW SERVERS` query groups results into `data_nodes` and `meta_nodes`. The term `meta_nodes` is outdated and refers to a node that runs the consensus service.

至此，你已经成功创建了 3 node cluster ！

If you believe that you did the above steps correctly, but are still experiencing problems, try restarting each node in your cluster.

## Adding nodes to your cluster

Once your initial cluster is healthy and running appropriately, you can start adding nodes to the cluster. Additional nodes can be consensus nodes, data nodes, or hybrid nodes. See [Cluster Node Configuration]() for how to configure the different node types.

Adding a node to your cluster follows the same procedure that we outlined above. Note that in step 4, when you point your new node to the cluster, you must set `INFLUXD_OPTS` to every node in the cluster, including itself.

## Removing nodes from your cluster

Please see the [reference documentation]() on `DROP SERVER.`

