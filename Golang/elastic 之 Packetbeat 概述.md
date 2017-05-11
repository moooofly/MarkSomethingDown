# elastic 之 Packetbeat 概述

## [elastic/beats](https://github.com/elastic/beats)

特点：

- 轻量（轻微安装痕迹，占用少量系统资源，无运行时依赖）；
- 基于 Go 语言；
- 捕获全类型的操作数据（operational data），比如 logs ， metrics ，或网络抓包数据；
- 数据发往 Elasticsearch ，或者直接交互，或者通过 Logstash ；数据可视化基于 Kibana ；

## [libbeat](https://github.com/elastic/beats/tree/master/libbeat)

libbeat 为用于创建各种 beats 的 Go 框架；

基于该框架，目前官方支持的 Beats 如下：

Beat  | Description
--- | ---
[Filebeat](https://github.com/elastic/beats/tree/master/filebeat) | Tails and ships log files
[Metricbeat](https://github.com/elastic/beats/tree/master/metricbeat) | Fetches sets of metrics from the operating system and services
[Packetbeat](https://github.com/elastic/beats/tree/master/packetbeat) | Monitors the network and applications by sniffing packets
[Winlogbeat](https://github.com/elastic/beats/tree/master/winlogbeat) | Fetches and ships Windows Event logs


----------


## [Packetbeat 参考手册](https://www.elastic.co/guide/en/beats/packetbeat/current/index.html)

### 概览

实时网络包分析 ＋ 基于 Elasticsearch 实现 APM ；

> Packetbeat is a real-time network packet analyzer that you can use with Elasticsearch to provide an application monitoring and performance analytics system.

为 Beats 平台提供了针对服务器间通信的“视野”；

> Packetbeat completes the Beats platform by providing visibility between the servers of your network.

针对应用层协议进行解析，关联请求和应答，针对每个 transaction 记录关心的 field ；

> Packetbeat works by capturing the network traffic between your application servers, decoding the application layer protocols (HTTP, MySQL, Redis, and so on), correlating the requests with the responses, and recording the interesting fields for each transaction.

帮助发现和排查后端服务器的 bug 和性能问题；

> Packetbeat can help you easily notice issues with your back-end application, such as bugs or performance problems, and it makes troubleshooting them - and therefore fixing them - much faster.

即时解析应用层协议，以 transaction 为单位关联相应的消息；

> Packetbeat sniffs the traffic between your servers, parses the application-level protocols on the fly, and correlates the messages into transactions. Currently, Packetbeat supports the following protocols:

- ICMP (v4 and v6)
- DNS
- HTTP
- **AMQP 0.9.1**
- Cassandra
- **Mysql**
- PostgreSQL
- **Redis**
- Thrift-RPC
- **MongoDB**
- Memcache

可以基于 Elasticsearch 或 Redis 或 Logstash 进行存储和分析；

> Packetbeat can insert the correlated transactions directly into Elasticsearch or into a central queue created with Redis and Logstash.

部署和使用方案（合并 or 单独）；

> Packetbeat can run on the same servers as your application processes or on its own servers. When running on dedicated servers, Packetbeat can get the traffic from the **switch’s mirror ports** or from **tapping devices**. In such a deployment, there is zero overhead on the monitored application. See [Setting Traffic Capturing Options](https://www.elastic.co/guide/en/beats/packetbeat/current/capturing-options.html) for details.

使用 JSON 文档进行保存；

> After decoding the Layer 7 messages, Packetbeat correlates the requests with the responses in what we call transactions. For each transaction, Packetbeat inserts a **JSON** document into Elasticsearch. See the [Exported Fields](https://www.elastic.co/guide/en/beats/packetbeat/current/exported-fields.html) section for details about which fields are indexed.

Packetbeat 和 Logstash 可以使用相同的 Elasticsearch 和 Kibana 实例；

> The same Elasticsearch and Kibana instances that are used for analysing the network traffic gathered by Packetbeat can be used for analysing the log files gathered by Logstash. This way, you can have network traffic and log analysis in the same system.


这里需要知道的是：使用 Packetbeat 前，需要先搞定 **Elastic Stack** 的安装，即

- **Elasticsearch** 进行数据的存储和索引查询；
- **Kibana** 提供 UI 供查询和展示；
- **Logstash** 用于插入数据到 Elasticsearch 中（可选）；

> **Elastic Stack** 的安装详见 [Getting Started with Beats and the Elastic Stack](https://github.com/moooofly/MarkSomethingDown/blob/master/Golang/Elastic%20Stack%20%E5%AE%89%E8%A3%85.md) ；

在完成 **Elastic Stack** 安装后，则可进行 Packetbeat 的安装、配置和运行，详见 [Getting Started With Packetbeat](https://www.elastic.co/guide/en/beats/packetbeat/5.0/packetbeat-getting-started.html) ；
