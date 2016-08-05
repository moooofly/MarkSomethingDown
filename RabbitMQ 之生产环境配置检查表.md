


# Production Checklist

## Introduction

像 RabbitMQ 类似的数据服务通常都有很多调优参数；一些配置参数对于开发环境来说非常有意义，但是并不适用于生产环境；没有任何一种配置能够适用于每一种使用场景；因此，在真正进入生产环境使用前，需要仔细确认你的配置内容；本文档就是针对此目的；

## Virtual Hosts, Users, Permissions

### Virtual Hosts
In a single-tenant environment, for example, when your RabbitMQ cluster is dedicated to power a single system in production, using default virtual host (/) is perfectly fine.

In multi-tenant environments, use a separate vhost for each tenant/environment, e.g. project1_development, project1_production, project2_development, project2_production, and so on.

### Users
For production environments, delete the default user (guest). Default user only can connect from localhost by default, because it has well-known credentials. Instead of enabling remote connections, consider using a separate user with administrative permissions and a generated password.

It is recommended to use a separate user per application. For example, if you have a mobile app, a Web app, and a data aggregation system, you'd have 3 separate users. This makes a number of things easier:

Correlating client connections with applications
Using fine-grained permissions
Credentials roll-over (e.g. periodically or in case of a breach)
In case there are many instances of the same application, there's a trade-off between better security (having a set of credentials per instance) and convenience of provisioning (sharing a set of credentials between some or all instances). For IoT applications that involve many clients performing the same or similar function and having fixed IP addresses, it may make sense to authenticate using x509 certificates or source IP addresse ranges.

## Resource Limits

RabbitMQ 采用资源驱动的方式进行告警，进而在 consumer 跟不上 publisher 时对后者进行限制；在真正应用到生产环境前，评估资源使用配置非常重要；

### Memory

默认情况下，RabbitMQ 将使用最多 40% 的可用 RAM ；对于希望运行 RabbitMQ 的节点来说，增大资源限制值是很有必要的；然而，需要注意操作系统和文件系统中的 cache 也需要占用 RAM 资源；若不加考虑，则可能会由于操作系统的 swapping 行为导致严重的吞吐量下降，甚至可能导致RabbitMQ 进程被操作系统直接终止；

下面给出一些基本准则以确定何种 RAM 限制值是被推荐的：

- 只要要有 128 MB
- 75% of the configured RAM limit when the limit is up to 4 GB of RAM
- 80% of the configured RAM limit when the limit is between 4 and 8 GB of RAM
- 85% of the configured RAM limit when the limit is between 8 and 16 GB of RAM
- 90% of the configured RAM limit when the limit is above 16 GB of RAM
- 高于 0.9 的配置可能会导致问题，不建议使用；

### Free Disk Space

Some free disk space should be available to avoid disk space alarms. By default RabbitMQ requires 50 MiB of free disk space at all times. This improves developer experience on many popular Linux distributions which may place the /var directory on a small partition. However, this is not a value recommended for production environments, since they may have significantly higher RAM limits. Below are some basic guidelines for determining how much free disk space is recommended:

- At least 2 GB
- 50% of the configured RAM limit when the limit is between 1 and 8 GB of RAM
- 40% of the configured RAM limit when the limit is between 8 and 32 GB of RAM
- 30% of the configured RAM limit when the limit is above 32 GB of RAM

The rabbit.disk_free_limit configuration setting can be set to {mem_relative, N} to make it calculated as a percentage of the RAM limit. For example, use {mem_relative, 0.5} for 50%, {mem_relative, 0.25} for 25%, and so on.

### Open File Handles Limit

Operating systems limit maximum number of concurrently open file handles, which includes network sockets. Make sure that you have limits set high enough to allow for expected number of concurrent connections and queues.

Make sure your environment allows for at least 50K open file descriptors for effective RabbitMQ user, including in development environments.

As a rule of thumb, multiple the 95th percentile number of concurrent connections by 2 and add total number of queues to calculate recommended open file handle limit. Values as high as 500K are not inadequate and won't consume a lot of hardware resources, and therefore are recommended for production setups. See Networking guide for more information.

## Security Considerations

### Users and Permissions
See the section on vhosts, users, and credentials above.

### Erlang Cookie
On Linux and BSD systems, it is necessary to restrict Erlang cookie access only to the users that will run RabbitMQ and tools such as rabbitmqctl.

### TLS
We recommend using TLS connections when possible, at least to encrypt traffic. Peer verification (authentication) is also recommended. Development and QA environments can use self-signed TLS certificates. Self-signed certificates can be appropriate in production environments when RabbitMQ and all applications run on a trusted network or isolated using technologies such as VMware NSX.

While RabbitMQ tries to offer a secure TLS configuration by default (e.g. SSLv3 is disabled), we recommend evaluating what TLS versions and cipher suites are enabled. Please see our TLS guide for more information.

## Networking Configuration

Production environments may require network configuration tuning. Please refer to the Networking Guide for details.

### Automatic Connection Recovery
Some client libraries, for example, Java, .NET, and Ruby ones, support automatic connection recovery after network failures. If the client used provides this feature, it is recommended to use it instead of developing your own recovery mechanism.

## Clustering Considerations

### Cluster Size
When determining cluster size, it is important to take several factors into consideration:

Expected throughput
Expected replication (number of mirrors)
Data locality
Since clients can connect to any node, RabbitMQ may need to perform inter-cluster routing of messages and internal operations. Try making consumers and producers connect to the same node, if possible: this will reduce inter-node traffic. Equally helpful would be making consumers connect to the node that currently hosts queue master (can be inferred using HTTP API). When data locality is taken into consideration, total cluster throughput can reach non-trivial volumes.

For most environments, mirroring to more than half of cluster nodes is sufficient. It is recommended to use clusters with an odd number of nodes (3, 5, and so on).

### Partition Handling Strategy
It is important to pick a partition handling strategy before going into production. When in doubt, use the autoheal strategy.