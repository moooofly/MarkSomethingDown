


Warning! This page documents an old version of InfluxDB, which is no longer actively developed.InfluxDB v0.13 is the most recent stable version of InfluxDB.


----------


# aggregation

An InfluxQL function that returns an aggregated value across a set of points. See InfluxQL Functionsfor a complete list of the available and upcoming aggregations.

Related entries: function, selector, transformation

# cluster

运行 InfluxDB node 的一组服务器；归属同一 cluster 的所有 node 共享相同的 users, databases, retention policies, 和 continuous queries ； 如何建立一个 InfluxDB cluster 详见 [Cluster Setup]() ；

Related entries: node, server

# consensus node

仅运行 consensus service 的 node 

详见 [Cluster Node Configuration]() ；

Related entries: cluster, consensus service, data node, node, hybrid node

# consensus service

会参与到 raft consensus group 中的 InfluxDB service ；一个 cluster 中必须至少存在 3 个运行 consensus service（基于 consensus 或 hybrid node）的 node ；但允许有更多的 node 运行相应的服务；在一个 cluster 中要求存在奇数运行 consensus service 的 node ；


在 cluster 功能降级（degraded）前允许失效的 consensus services 的数量为 ⌈n/2 + 1⌉ ，其中 n 为 cluster 中 consensus services 的数量；因此，偶数个 consensus services 并不会提供额外的冗余或弹性；

consensus service 用于在 cluster 中确保以下内容的一致性：node membership, databases, retention policies, users, continuous queries, shard metadata, 和 subscriptions ；

See [Cluster Node Configuration]() ；

Related entries: cluster, consensus node, data service, node, hybrid node

# continuous query (CQ)

An InfluxQL query that runs automatically and periodically within a database. Continuous queries require a function in the SELECT clause and must include a GROUP BY time() clause. See Continuous Queries.

Related entries: function

# coordinator node

能够为 cluster 分担 write 和 query 请求的 node ；

Related entries: cluster, hinted handoff, node

# data node

仅运行 data service 的 node ；

See [Cluster Node Configuration]() ；

Related entries: cluster, consensus node, data service, node, hybrid node

# data service

用于持久化时间序列数据到 node 上的 InfluxDB service ；一个 cluster 中必须至少存在一个运行 data service 的 node（可以是 data 或 hybrid node），但是允许运行更多数量的 data node ；

See [Cluster Node Configuration]() ；

Related entries: cluster, consensus node, consensus service, node, hybrid node

# database

保存 users, retention policies, continuous queries, 和时间序列数据的逻辑容器；

Related entries: continuous query, retention policy, user

# duration

The attribute of the retention policy that determines how long InfluxDB stores data. Data older than the duration are automatically dropped from the database. See Database Management for how to set duration.

Related entries: replication factor, retention policy

# field

The key-value pair in InfluxDB’s data structure that records metadata and the actual data value. Fields are required in InfluxDB’s data structure and they are not indexed - queries on field values scan all points that match the specified time range and, as a result, are not performant relative to tags.

Query tip: Compare fields to tags; tags are indexed.

Related entries: field key, field set, field value, tag

# field key

The key part of the key-value pair that makes up a field. Field keys are strings and they store metadata.

Related entries: field, field set, field value, tag key

# field set

The collection of field keys and field values on a point.

Related entries: field, field key, field value, point

# field value

The value part of the key-value pair that makes up a field. Field values are the actual data; they can be strings, floats, integers, or booleans. A field value is always associated with a timestamp.

Field values are not indexed - queries on field values scan all points that match the specified time range and, as a result, are not performant.

Query tip: Compare field values to tag values; tag values are indexed.

Related entries: field, field key, field set, tag value, timestamp

# function

InfluxQL aggregations, selectors, and transformations. See InfluxQL Functions for a complete list of InfluxQL functions.

Related entries: aggregation, selector, transformation

# hinted handoff

当数据被接收到时，若目标服务器不可用，则使用一个持久化 queue 进行数据的保存；当 write 操作的目标 node 处于短暂的 down 状态时，Coordinating nodes 将会暂时负责存储这些需要压入 queue 中的数据；

Related entries: cluster, node, server

# hybrid node

同时运行 consensus 和 data service 的 node ；

See [Cluster Node Configuration]() ；

Related entries: cluster, consensus node, consensus service, node, data node, data service

# identifier

Tokens which refer to database names, retention policy names, user names, measurement names, tag keys, and field keys. See Query Language Specification.

Related entries: database, field key, measurement, retention policy, tag key, user

# line protocol

The text based format for writing points to InfluxDB. See [Line Protocol]().

# measurement

The part of InfluxDB’s structure that describes the data stored in the associated fields. Measurements are strings.

Related entries: field, series

# metastore

保存系统状态相关的内部信息的地方；包括：user information, database 和 shard metadata, 以及哪些 retention policies 被使能了；

Related entries: database, retention policy, user

# node

一个独立的 influxd 进程；

Related entries: cluster, server

# point

The part of InfluxDB’s data structure that consists of a single collection of fields in a series. Each point is uniquely identified by its series and timestamp.

You cannot store more than one point with the same timestamp in the same series. Instead, when you write a new point to the same series with the same timestamp as an existing point in that series, the field set becomes the union of the old field set and the new field set, where any ties go to the new field set. For an example, see Frequently Encountered Issues.

Related entries: field set, series, timestamp

# query

An operation that retrieves data from InfluxDB. See Data Exploration, Schema Exploration, Database Management.

# replication factor

The attribute of the retention policy that determines how many copies of the data are stored in the cluster. InfluxDB replicates data across N data nodes, where N is the replication factor.

To maintain data availability for queries, the replication factor should be less than or equal to the number of data nodes in the cluster:

Data are fully available when the replication factor is greater than the number of unavailable data nodes.
Data may be unavailable when the replication factor is less than the number of unavailable data nodes.
Note that there are no query performance benefits from replication. Replication is for ensuring data availability when a data node or nodes are unavailable. See Database Management for how to set the replication factor.

Related entries: cluster, duration, node, retention policy

# retention policy (RP)

The part of InfluxDB’s data structure that describes for how long InfluxDB keeps data (duration) and how many copies of those data are stored in the cluster (replication factor). RPs are unique per database and along with the measurement and tag set define a series.

When you create a database, InfluxDB automatically creates a retention policy called default with an infinite duration and a replication factor set to the number of nodes in the cluster. See Database Management for retention policy management.

Related entries: duration, measurement, replication factor, series, tag set

# schema

How the data are organized in InfluxDB. The fundamentals of the InfluxDB schema are databases, retention policies, series, measurements, tag keys, tag values, and field keys. See Schema Design for more information.

Related entries: database, field key, measurement, retention policy, series, tag key, tag value

# selector

An InfluxQL function that returns a single point from the range of specified points. See InfluxQL Functions for a complete list of the available and upcoming selectors.

Related entries: aggregation, function, transformation

# series

The collection of data in InfluxDB’s data structure that share a measurement, tag set, and retention policy.

Note: The field set is not part of the series identification!

Related entries: field set, measurement, retention policy, tag set

# series cardinality

The count of all combinations of measurements and tags within a given data set. For example, take measurement mem_available with tags host and total_mem. If there are 35 different hosts and 15 different total_mem values then series cardinality for that measurement is 35 * 15 = 525. To calculate series cardinality for a database add the series cardinalities for the individual measurements together.

Related entries: tag set, measurement, tag key

# server

A machine, virtual or physical, that is running InfluxDB. There should only be one InfluxDB process per server.

Related entries: cluster, node

# tag

The key-value pair in InfluxDB’s data structure that records metadata. Tags are an optional part of InfluxDB’s data structure but they are useful for storing commonly-queried metadata; tags are indexed so queries on tags are performant. Query tip: Compare tags to fields; fields are not indexed.

Related entries: field, tag key, tag set, tag value

# tag key

The key part of the key-value pair that makes up a tag. Tag keys are strings and they store metadata. Tag keys are indexed so queries on tag keys are performant.

Query tip: Compare tag keys to field keys; field keys are not indexed.

Related entries: field key, tag, tag set, tag value

# tag set

The collection of tag keys and tag values on a point.

Related entries: point, series, tag, tag key, tag value

# tag value

The value part of the key-value pair that makes up a tag. Tag values are strings and they store metadata. Tag values are indexed so queries on tag values are performant.

Related entries: tag, tag key, tag set

# timestamp

The date and time associated with a point. All time in InfluxDB is UTC.

For how to specify time when writing data, see Write Syntax. For how to specify time when querying data, see Data Exploration.

Related entries: point

# transformation

An InfluxQL function that returns a value or a set of values calculated from specified points, but does not return an aggregated value across those points. See InfluxQL Functions for a complete list of the available and upcoming aggregations.

Related entries: aggregation, function, selector

# user

在 InfluxDB 中存在两种类型的用户：

Admin user 拥有针对所有数据库的 READ 和 WRITE 权限；拥有 administrative 类型 queries 的全部权限；拥有用户管理命令的权限；

Non-admin user 拥有针对特定数据库的 READ, WRITE, 或 ALL (both READ and WRITE) 权限；
当使能量 authentication 机制，InfluxDB 只能执行带有合法用户名和密码的 HTTP 请求；See [Authentication and Authorization]() ；

