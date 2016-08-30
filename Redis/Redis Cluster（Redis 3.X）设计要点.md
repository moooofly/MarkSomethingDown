


#  Redis Cluster（Redis 3.X）设计要点


Redis 3.0.0 RC1 版本 10.9 号发布，[Release Note](https://raw.githubusercontent.com/antirez/redis/3.0/00-RELEASENOTES)；该版本支持 Redis Cluster，相信很多同学期待已久，不过这个版本只是 RC 版本，要应用到生产环境，还得等等；

Redis Cluster设计要点：

## 架构：无中心

- Redis Cluster 采用无中心结构，每个节点都保存数据和整个集群的状态；
- 每个节点都和其他所有节点连接，这些连接保持活跃；
- 使用 `gossip` 协议传播信息以及发现新节点；
- cluster 中的 node 不会为 client 代理请求，client 需要根据 node 返回的错误信息（`MOVED` 错误）重定向请求；

## 数据分布：预分桶

- 预分配 16384 个桶（slot），根据 `CRC16(key) mod 16384` 的值，决定将一个 key 放到哪个桶中；
- 每个 Redis 物理结点负责一部分桶的管理，当发生 Redis 节点的增减时，调整桶的分布即可；

例如，假设 Redis Cluster 由三个节点 A、B、C 构成，则
- Node A 包含桶的编号可以为 0 ~ 5500 ；
- Node B 包含桶的编号可以为 5500 ~ 11000 ；
- Node C 包含桶的编号可以为 11001 ~ 16384 ；

当发生 Redis 节点的增减时，只需调整桶的分布；

预分桶的方案介于“**硬 Hash**”和“**一致性 Hash**”之间，牺牲了一定的灵活性，但相比“一致性Hash“，数据的管理成本大大降低；

## 可用性：Master-Slave

- 为了保证服务的（高）可用性，Redis Cluster 采取的方案是的 Master-Slave ；
- 每个 Redis node 可以有一个或者多个 Slave ；当 Master 挂掉时，会选举一个 Slave 成为新 Master ；
- 每个 Redis  node 包含（被分配）一定量的桶；只有当这些桶对应的 Master 和 Slave 都挂掉时，这部分桶中的数据才不可用；

## 写操作

Redis Cluster 使用异步复制：
一个完整的写操作步骤：
1. client 写数据到 master 节点；
2. master 告诉 client "ok" ；
3. master 传播更新到 slave ；

存在数据丢失的风险：
- 上述写步骤 1 和 2 成功后，master crash，而此时数据还没有传播到 slave ；
-  由于分区导致同时存在两个 master ，client 向旧的 master 写入了数据；

当然，由于 Redis Cluster 存在超时及故障恢复机制，因此第 2 个风险基本上不可能发生；

## 数据迁移

- Redis Cluster 支持在线增/减节点；
- 基于桶的数据分布方式大大降低了迁移成本，只需将数据桶从一个 Redis Node 迁移到另一个 Redis Node 即可完成迁移；
- 当桶从一个 Node A 向另一个 Node B 迁移时，Node A 和 Node B 都会有这个桶，Node A 上桶的状态设置为`MIGRATING` ，Node B 上桶的状态被设置为 `IMPORTING` ；
- 当客户端请求时，所有在 Node A 上的请求都将由 A 来处理，所有不在 A 上的 key 都由 Node B 来处理。同时，Node A 上将不会创建新的 key ；

## 多 key 操作

当系统从单节点向多节点扩展时，多 key 操作总是一个非常难解决的问题；Redis Cluster 方案如下：
1. 不支持多 key 操作；
2. 如果一定要使用多 key 操作，请确保所有的 key 都在一个 slot 上，具体方法是使用“hash tag”方案

> hash tag 方案是一种数据分布的例外情况


## Reference
- [Redis cluster tutorial](http://redis.io/topics/cluster-tutorial)
- [Redis cluster Specification](http://redis.io/topics/cluster-spec)


----------


原文：[这里](http://blog.csdn.net/yfkiss/article/details/39996129)




