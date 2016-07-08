

> Warning! This page documents an old version of InfluxDB, which is no longer actively developed.InfluxDB v0.13 is the most recent stable version of InfluxDB.


----------


> NOTE: InfluxDB 0.11 is the last open source version that includes clustering. For more information, please see Paul Dix’s blog post on InfluxDB Clustering, High-Availability, and Monetization. Please note that the 0.11 version of clustering is still considered experimental, and there are still quite a few rough edges.

This guide briefly introduces the InfluxDB cluster model and provides step-by-step instructions for setting up a cluster.

# InfluxDB cluster model

InfluxDB supports arbitrarily sized clusters and any replication factor from 1 to the number of nodes in the cluster.

There are three types of nodes in an InfluxDB cluster: consensus nodes, data nodes, and hybrid nodes. A cluster must have an odd number of nodes running the consensus service to form a Raft consensus group and remain in a healthy state.

Hardware requirements vary for the different node types. See Hardware Sizing for cluster hardware requirements.

# Cluster setup

The following steps configure and start up an InfluxDB cluster with three hybrid nodes. If you’re interested in having any of the different node types, see Cluster Node Configuration for their configuration details. Note that your first three nodes must be either hybrid nodes or consensus nodes.

We assume that you are running some version of Linux, and, while it is possible to build a cluster on a single server, it is not recommended.

1  在 3 台机器上安装 InfluxDB ；但暂时不要在任何一台机器上启动 daemon 程序

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

3  将所有 node 都配置成可以与其它 node 进行通信

针对上述 3 个 node ，在 `/etc/default/influxdb` 中设置 `INFLUXD_OPTS` 的值如下：

```shell
INFLUXD_OPTS="-join <IP1>:8091,<IP2>:8091,<IP3>:8091"
```

where IP1 is the first node’s IP address or hostname, IP2 is the second nodes’s IP address orhostname, and IP3 is the third node’s IP address or hostname.

如果文件 `/etc/default/influxdb` 不存在，则需手动创建；

4  在每一个 node 上都启动 InfluxDB ：

```shell
sudo service influxdb start
```

5  校验 cluster 的健康性

Issue a SHOW SERVERS query to each node in your cluster using the influx CLI. The output should show that your cluster is made up of three hybrid nodes (hybrid nodes appear as both data_nodesand meta_nodes in the SHOW SERVERS query results):

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

Note: The SHOW SERVERS query groups results into data_nodes and meta_nodes. The termmeta_nodes is outdated and refers to a node that runs the consensus service.

至此，已经成功创建了 3 node cluster ！

If you believe that you did the above steps correctly, but are still experiencing problems, try restarting each node in your cluster.

## Adding nodes to your cluster

Once your initial cluster is healthy and running appropriately, you can start adding nodes to the cluster. Additional nodes can be consensus nodes, data nodes, or hybrid nodes. See Cluster Node Configuration for how to configure the different node types.

Adding a node to your cluster follows the same procedure that we outlined above. Note that in step 4, when you point your new node to the cluster, you must set INFLUXD_OPTS to every node in the cluster, including itself.

## Removing nodes from your cluster

Please see the reference documentation on DROP SERVER.

