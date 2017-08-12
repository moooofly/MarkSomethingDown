# ELK 学习笔记

## [ElasticSearch 5学习(1)——安装Elasticsearch、Kibana和X-Pack](http://www.cnblogs.com/wxw16/p/6150681.html)

要点：

- 按照文档的要求，一般情况下 kibana 的版本必须和 Elasticsearch 安装的版本一致；
- X-Pack 是一个 Elastic Stack 的扩展，将安全，警报，监视，报告和图形功能包含在一个易于安装的软件包中，也是官方推荐的。在 Elasticsearch 5.0.0 之前，您必须安装单独的 Shield，Watcher 和 Marvel 插件才能获得在 x-pack 中所有的功能；
- 若想使用 x-pack 功能，则需要执行如下操作：
    - `elasticsearch-plugin install x-pack`
    - `kibana-plugin install x-pack`
    - 在 elasticsearch.yml 中配置 action.auto_create_index 以允许 x-pack 创造以下指标：`action.auto_create_index: .security,.monitoring*,.watches,.triggered_watches,.watcher-history*`
- 启用 x-pack 后，通过 http://localhost:5601/ 登陆 kibana 时会要求输入用户名和密码，默认为 elastic 和 changeme ；
- 在刚接触 Elasticsearch 的时候，会有很多名词不能理解，或者不知道其中的关系，而其中很多内容是为不同版本的 Elasticsearch 而存在的：
    - Marvel：是 Elasticsearch 的管理和监控工具，在开发环境下免费使用，包含了 Sense ；用于在簇中从每个节点汇集数据，该插件必须每个节点都得安装；
    - Sense：交互式控制台，使用户方便的通过浏览器直接与 Elasticsearch 进行交互；
    - Hand：除了 REST 之外，另外一种查看 es 运行状态以及数据的方式；可以实现基本信息的查看，REST 请求的模拟，数据的检索等等；
- kibana 是一个与 elasticsearch 一起工作的开源的分析和可视化的平台。使用 kibana 可以查询、查看并与存储在 elasticsearch 索引的数据进行交互操作。使用 kibana 能执行高级的数据分析，并能以图表、表格和地图的形式查看数据。kibana 使得理解大容量的数据变得非常容易。它非常简单，基于浏览器的接口使我们能够快速的创建和分享显示 elasticsearch 查询结果实时变化的仪表盘。**在 Elasticsearch 5 版本之前，一般都是通过安装 Kibana，而后将 Marvel、Hand 等各种功能插件添加到 Kibana 上使用。在 Elasticsearch 5 版本之后，一般情况下只需要安装一个官方推荐的 X-pack 扩展包即可**。

## [ElasticSearch 5学习(2)——Kibana+X-Pack介绍使用（全）](http://www.cnblogs.com/wxw16/p/6156335.html)

### kibana 3/4/5 界面对比

- Kibana 3 的界面，所有的仪表盘直接放置主页

![](http://images2015.cnblogs.com/blog/763363/201612/763363-20161210173802335-1435351087.png)


- Kibana 4 的界面，将 3 原来的主体分成三个部分，分别对应发现页、可视化、仪表盘

![](http://images2015.cnblogs.com/blog/763363/201612/763363-20161210173809819-268441203.png)

- Kibana 5 的界面，除了界面风格的变化，最主要是功能栏上添加了 Timeline、Management 和 Dev Tools 选项

![](http://images2015.cnblogs.com/blog/763363/201612/763363-20161210173819679-1909820200.png)


### Kibana 5 标签页详解

#### Discover

> You can interactively explore your data from the Discover page. You have access to every document in every index that matches the selected index pattern. You can submit search queries, filter the search results, and view document data. You can also see the number of documents that match the search query and get field value statistics. If a time field is configured for the selected index pattern, the distribution of documents over time is displayed in a histogram at the top of the page.

从 Discover 页可以交互地探索 ES 的数据：

- 可以访问与所选索引模式（index pattern）相匹配的每一个索引（index）中的每一个文档（document）；
- 可以**提交搜索查询**、**筛选搜索结果**和**查看文档数据**；
- 还可以看到匹配搜索查询的文档的数量，以及获取字段值统计信息。如果一个时间字段被配置为所选择的索引模式，则文档的分布随着时间的推移显示在页面顶部的直方图中。

![](http://images2015.cnblogs.com/blog/763363/201612/763363-20161210174039319-2128363229.png)

#### Visualize

> Visualize enables you to create visualizations of the data in your Elasticsearch indices. You can then build dashboards that display related visualizations. Kibana visualizations are based on Elasticsearch queries. By using a series of Elasticsearch aggregations to extract and process your data, you can create charts that show you the trends, spikes, and dips you need to know about. You can create visualizations from a search saved from Discover or start with a new search query.

Visualize 允许你针对 Elasticsearch 的索引数据创建**可视化的指标数据**。之后，你久可以建立仪表板显示相应的可视化内容。Kibana 的可视化是基于 Elasticsearch 查询的，通过使用一系列 Elasticsearch 聚合功能进行数据的提取和处理，您久可以创建图表来显示关于**趋势**，**峰值**和**骤降**信息。您可以基于在 Discover 页中保存的一次“search”创建可视化，或者基于一个全新 search 查询创建可视化。

![](http://images2015.cnblogs.com/blog/763363/201612/763363-20161210174332616-1801083291.png)

#### Dashboard

> A Kibana dashboard displays a collection of saved visualizations. You can arrange and resize the visualizations as needed and save dashboards so they be reloaded and shared.

Kibana dashboard 用于显示一组保存好的可视化实体。你可以根据需要布置和调整可视化，以便于再次加载和用于共享。

![](http://images2015.cnblogs.com/blog/763363/201612/763363-20161210173846694-1170257472.png)

#### Monitoring

默认情况下，Kibana 是没有提供该选项。Monitoring 功能是由 X-Pack 扩展提供的。

> The X-Pack monitoring components enable you to easily monitor Elasticsearch through Kibana. You can view cluster health and performance in real time as well as analyze past cluster, index, and node metrics. In addition, you can monitor the performance of Kibana itself.When you install X-Pack on your cluster, a monitoring agent runs on each node to collect and index metrics from Elasticsearch. With X-Pack installed in Kibana, you can then view the monitoring data through a set of specialized dashboards.

该 X-pack 监控组件使您可以通过 Kibana 轻松地监控 ElasticSearch 。您可以实时查看集群的健康和性能，以及分析过去的集群、索引和节点度量。此外，您可以监视 Kibana 本身性能。当你安装 X-pack 在群集上，监控代理运行在每个节点上收集和指数指标从 Elasticsearch 。安装在 X-pack 在 Kibana 上，您可以查看通过一套专门的仪表板监控数据。

![](http://images2015.cnblogs.com/blog/763363/201612/763363-20161210173853976-845961527.png)

#### Graph

> The X-Pack graph capabilities enable you to discover how items in an Elasticsearch index are related. You can explore the connections between indexed terms and see which connections are the most meaningful. This can be useful in a variety of applications, from fraud detection to recommendation engines.For example, graph exploration could help you uncover website vulnerabilities that hackers are targeting so you can harden your website. Or, you might provide graph-based personalized recommendations to your e-commerce customers. X-Pack provides a simple, yet powerful graph exploration API, and an interactive graph visualization tool for Kibana. Both work with out of the box with existing Elasticsearch indices—you don’t need to store any additional data to use the X-Pack graph features.

X-Pack 图的能力使你能够发现 Elasticsearch  索引中的所有 item 是如何相关联的。你可以探索 indexed terms 之间的关联，以便确认哪些关联是最有意义的。从欺诈检测到推荐引擎，对各种应用中这都是有用的；例如，图的探索可以帮助你发现网站上黑客的目标的漏洞，所以你可以硬化你的网站。或者，您可以为您的电子商务客户提供基于图表的个性化推荐。X-pack 提供了一套简单但功能强大的图形开发 API ，以及一个能够与 Kibana 交互式的图形可视化工具。这两者均工作于 Elasticsearch 中已存在的索引数据之上，因此无需保存任何额外信息就可以使用 x-pack 功能。

![](http://images2015.cnblogs.com/blog/763363/201612/763363-20161210173904679-826322201.png)

#### Timelion

> Timelion is a time series data visualizer that enables you to combine totally independent data sources within a single visualization. It’s driven by a simple expression language you use to retrieve time series data, perform calculations to tease out the answers to complex questions, and visualize the results.

Timelion 是一个时序数据可视化工具，允许你在一个单独的可视化展示中将完全独立的数据源绑定在一起；其通过一种简单的表达式语言进行驱动，可用来检索时间序列数据，通过计算梳理出复杂问题的答案，并将结果进行可视化。

这个功能由一系列的功能函数组成，同样的查询结果，也可以通过Dashboard 显示查看。

![](http://images2015.cnblogs.com/blog/763363/201612/763363-20161210173913632-907340971.png)

#### Management

> The Management application is where you perform your runtime configuration of Kibana, including both the initial setup and ongoing configuration of index patterns, advanced settings that tweak the behaviors of Kibana itself, and the various "objects" that you can save throughout Kibana such as searches, visualizations, and dashboards. This section is pluginable, so in addition to the out of the box capabitilies, packs such as X-Pack can add additional management capabilities to Kibana.

管理应用是你进行 kibana 运行时配置的地方，包括初始设置和针对 index patterns 持续存在的（ongoing）配置，能够改变 Kibana 自身行为的高级设置，以及能够通过 Kibana 进行保存的各种“对象”（诸如 searches, visualizations 和 dashboards），你可以查看保存在整个Kibana的内容如发现页，可视化和仪表板。这部分功能是 pluginable 的，因此除了上述提及的功能外，使能诸如 X-pack 这类扩展可以为 Kibana 增加额外的管理能力。

> You can use X-Pack Security to control what Elasticsearch data users can access through Kibana. When you install X-Pack, Kibana users have to log in. They need to have the kibana_user role as well as access to the indices they will be working with in Kibana. If a user loads a Kibana dashboard that accesses data in an index that they are not authorized to view, they get an error that indicates the index does not exist. X-Pack Security does not currently provide a way to control which users can load which dashboards.

你可以基于 X-pack 安全机制控制哪些 Elasticsearch 数据用户能够通过 Kibana 进行访问。当你安装了 X-pack 之后，Kibana 用户就必须进行登录操作。其必须具有 kibana_user 角色，以及工作在 Kibana 中所需的 indices 访问的权限。如果用户加载了某个 Kibana dashboard ，后者访问了某个索引中的数据，但对该索引的访问未被授权查看，则会得到一个表示索引不存在的错误，表明指数不存在。X-pack安全功能当前不能提供针对哪些用户可以加载哪些 dashboards 的控制。

![](http://images2015.cnblogs.com/blog/763363/201612/763363-20161210173922663-216149537.png)

#### Dev Tools

原先的交互式控制台 Sense ，使用户能够方便的通过浏览器直接与 Elasticsearch 进行交互。从 Kibana 5 开始改名并直接内建在 Kibana 中，就是 Dev Tools 选项。

注意：如果是 Kibana 5 以上，不能通过以下命令安装 Sense。(踩过的坑)

```
./bin/kibana plugin --install elastic/sense
./bin/kibana-plugin install elastic/sense instead
```

![](http://images2015.cnblogs.com/blog/763363/201612/763363-20161210220924335-776040200.png)


## [ElasticSearch 5学习(3)——单台服务器部署多个节点](http://www.cnblogs.com/wxw16/p/6160186.html)

单机部署多节点存在一些问题：

- elasticsearch 只能通过 `-Epath.conf=xxx` 指定配置文件所在目录，而不能直接指定配置文件；
- 即使拆分成不同目录，仍存在如下 lock 相关问题（未解决）；

```
➜  elasticsearch2 elasticsearch -Epath.conf=/usr/local/etc/elasticsearch2
[2017-05-15T11:32:54,370][INFO ][o.e.n.Node               ] [node-2] initializing ...
[2017-05-15T11:32:54,421][WARN ][o.e.b.ElasticsearchUncaughtExceptionHandler] [node-2] uncaught exception in thread [main]
org.elasticsearch.bootstrap.StartupException: java.lang.IllegalStateException: failed to obtain node locks, tried [[/usr/local/var/elasticsearch/elasticsearch_sunfei]] with lock id [0]; maybe these locations are not writable or multiple nodes were started without increasing [node.max_local_storage_nodes] (was [1])?
	at org.elasticsearch.bootstrap.Elasticsearch.init(Elasticsearch.java:127) ~[elasticsearch-5.3.2.jar:5.3.2]
	at org.elasticsearch.bootstrap.Elasticsearch.execute(Elasticsearch.java:114) ~[elasticsearch-5.3.2.jar:5.3.2]
	at org.elasticsearch.cli.EnvironmentAwareCommand.execute(EnvironmentAwareCommand.java:58) ~[elasticsearch-5.3.2.jar:5.3.2]
	at org.elasticsearch.cli.Command.mainWithoutErrorHandling(Command.java:122) ~[elasticsearch-5.3.2.jar:5.3.2]
	at org.elasticsearch.cli.Command.main(Command.java:88) ~[elasticsearch-5.3.2.jar:5.3.2]
	at org.elasticsearch.bootstrap.Elasticsearch.main(Elasticsearch.java:91) ~[elasticsearch-5.3.2.jar:5.3.2]
	at org.elasticsearch.bootstrap.Elasticsearch.main(Elasticsearch.java:84) ~[elasticsearch-5.3.2.jar:5.3.2]
Caused by: java.lang.IllegalStateException: failed to obtain node locks, tried [[/usr/local/var/elasticsearch/elasticsearch_sunfei]] with lock id [0]; maybe these locations are not writable or multiple nodes were started without increasing [node.max_local_storage_nodes] (was [1])?
	at org.elasticsearch.env.NodeEnvironment.<init>(NodeEnvironment.java:260) ~[elasticsearch-5.3.2.jar:5.3.2]
	at org.elasticsearch.node.Node.<init>(Node.java:262) ~[elasticsearch-5.3.2.jar:5.3.2]
	at org.elasticsearch.node.Node.<init>(Node.java:242) ~[elasticsearch-5.3.2.jar:5.3.2]
	at org.elasticsearch.bootstrap.Bootstrap$6.<init>(Bootstrap.java:242) ~[elasticsearch-5.3.2.jar:5.3.2]
	at org.elasticsearch.bootstrap.Bootstrap.setup(Bootstrap.java:242) ~[elasticsearch-5.3.2.jar:5.3.2]
	at org.elasticsearch.bootstrap.Bootstrap.init(Bootstrap.java:360) ~[elasticsearch-5.3.2.jar:5.3.2]
	at org.elasticsearch.bootstrap.Elasticsearch.init(Elasticsearch.java:123) ~[elasticsearch-5.3.2.jar:5.3.2]
	... 6 more
➜  elasticsearch2
...
...
➜  ~ lsof | grep node.lock
java      35247 sunfei   58w     REG                1,4           0 6006718 /usr/local/var/elasticsearch/nodes/0/node.lock
➜  ~
```

- 启动 x-pack 之后，访问 elasticsearch 均需要提供用户名和密码；

```
➜  ~ curl -i -XGET 'http://localhost:9200/_cluster/health?pretty'\;
HTTP/1.1 401 Unauthorized
WWW-Authenticate: Basic realm="security" charset="UTF-8"
content-type: application/json; charset=UTF-8
content-length: 546

{
  "error" : {
    "root_cause" : [
      {
        "type" : "security_exception",
        "reason" : "missing authentication token for REST request [/_cluster/health?pretty;]",
        "header" : {
          "WWW-Authenticate" : "Basic realm=\"security\" charset=\"UTF-8\""
        }
      }
    ],
    "type" : "security_exception",
    "reason" : "missing authentication token for REST request [/_cluster/health?pretty;]",
    "header" : {
      "WWW-Authenticate" : "Basic realm=\"security\" charset=\"UTF-8\""
    }
  },
  "status" : 401
}
➜  ~
➜  ~
➜  ~ curl -u elastic:changeme -i -XGET 'http://localhost:9200/_cluster/health?pretty'\;
HTTP/1.1 200 OK
content-type: application/json; charset=UTF-8
content-length: 473

{
  "cluster_name" : "elasticsearch_sunfei",
  "status" : "yellow",
  "timed_out" : false,
  "number_of_nodes" : 1,
  "number_of_data_nodes" : 1,
  "active_primary_shards" : 4,
  "active_shards" : 4,
  "relocating_shards" : 0,
  "initializing_shards" : 0,
  "unassigned_shards" : 4,
  "delayed_unassigned_shards" : 0,
  "number_of_pending_tasks" : 0,
  "number_of_in_flight_fetch" : 0,
  "task_max_waiting_in_queue_millis" : 0,
  "active_shards_percent_as_number" : 50.0
}
➜  ~
```

- 踩过的坑（看原文）


## [ElasticSearch 5学习(4)——简单搜索笔记](http://www.cnblogs.com/wxw16/p/6171016.html)

- 空搜索
    - `total` ：匹配的总数
    - `hits` 数组：给出匹配度最高的前 10 个数据
    - `_score` ：相关性得分（relevance score），衡量文档与查询的匹配程度
    - `took` ：整个搜索请求花费的毫秒数
    - `_shards` ：参与查询的分片数（total 字段），有多少是成功的（successful 字段），有多少的是失败的（failed 字段）
    - `timed_out` ：查询超时与否（可通过 `?timeout=10ms` 进行控制）

> 注意：即使发生了 timeout ，也不会停止执行查询；超时发生后，得到返回值仅仅是告诉你当前顺利返回结果的节点，然后关闭连接。在后台，其他分片可能依旧执行查询，尽管结果已经被发送。使用超时是因为对于你的业务需求来说非常重要，而不是因为你想中断执行长时间运行的查询。

- 多索引和多类别
- 分页
- 简易搜索
    - 结构化查询语句（DSL）
    - `+name:john +tweet:mary` 中，'+' 前缀表示语句匹配条件必须被满足；'-' 前缀表示条件必须不被满足；
- `_all` 字段
    - 当你索引一个文档，Elasticsearch 会把所有字符串字段值连接起来放在一个大字符串中，并被索引为一个特殊的字段 `_all` ；
    - 若没有指定字段，查询字符串搜索（即 `q=xxx`）使用 `_all` 字段搜索；
- 更复杂的语句
    - 编码前：`+name:(mary john) +date:>2014-09-10 +(aggregations geo)`
    - 编码后：`?q=%2Bname%3A(mary+john)+%2Bdate%3A%3E2014-09-10+%2B(aggregations+geo)`
- 可能的隐患
    - 简单查询字符串搜索惊人的强大，允许我们简洁明快的表示复杂的查询，在命令行下一次性查询或者开发模式下非常有用；
    - 但简洁带来了隐晦和调试困难，查询字符串中一个细小的语法错误就会导致返回错误而不是结果；
    - 查询字符串搜索允许任意用户在索引中任何一个字段上运行潜在的慢查询语句，可能暴露私有信息甚至使你的集群瘫痪；


## [ElasticSearch 5学习(5)——第一个例子（很实用）](http://www.cnblogs.com/wxw16/p/6185378.html)

在 Elasticsearch 中**存储数据（文档）的行为**称为**索引**，但是在索引文档之前，我们需要决定在哪里存储它。

在 Elasticsearch 中，**文档**属于某个**类型**，这些**类型**位于**索引**中。

Elasticsearch 与传统关系数据库的粗略对比
```
Relational DB  ⇒ Databases ⇒ Tables ⇒ Rows      ⇒ Columns
Elasticsearch  ⇒ Indices   ⇒ Types  ⇒ Documents ⇒ Fields
```

Elasticsearch 集群可以包含多个索引（数据库），这些索引又包含多个类型（表）。这些类型包含多个文档（行），每个文档都有多个字段（列）。

在 Elasticsearch 的上下文中，索引被重载了几个含义：

- **索引（名词）**：正如前面所解释的那样，索引就像传统的关系数据库中的数据库一样，它是存储相关文档的地方。index 的复数形式是 indices 或 indexes 。
- **索引（动词）**：索引一个文档是将一个文档存储在索引（名词）中，以便它可以检索和查询。它很像插入关键词 SQL 。此外，如果文档已经存在，新的文档将取代旧的。
- **倒排索引**：关系数据库中增加一个索引，如 B-树索引，对特定列为了提高数据检索的速度。Elasticsearch 和 Lucene 提供相同目的的索引称为倒排索引。默认情况下，文档中的每个字段索引（有一个倒排索引）这样的搜索。一个没有倒排索引字段不可搜索。

### 构造数据

测试数据添加

```
# 第一次执行时，返回 "201 Created"
➜  ~ curl -u elastic:changeme -i -XPUT 'http://localhost:9200/megacorp/employee/1' -d '{
quote>     "first_name" : "John",
quote>     "last_name" :  "Smith",
quote>     "age" :        25,
quote>     "about" :      "I love to go rock climbing",
quote>     "interests": [ "sports", "music" ]
quote> }
quote> '
HTTP/1.1 201 Created
Location: /megacorp/employee/1
Warning: 299 Elasticsearch-5.3.2-3068195 "Content type detection for rest requests is deprecated. Specify the content type using the [Content-Type] header." "Wed, 17 May 2017 07:30:06 GMT"
content-type: application/json; charset=UTF-8
content-length: 145

{"_index":"megacorp","_type":"employee","_id":"1","_version":1,"result":"created","_shards":{"total":2,"successful":1,"failed":0},"created":true}%            ➜  ~
➜  ~
# 第二次执行时，返回 "200 OK"
➜  ~ curl -u elastic:changeme -i -XPUT 'http://localhost:9200/megacorp/employee/1' -d '{
    "first_name" : "John",
    "last_name" :  "Smith",
    "age" :        25,
    "about" :      "I love to go rock climbing",
    "interests": [ "sports", "music" ]
}
'
HTTP/1.1 200 OK
Warning: 299 Elasticsearch-5.3.2-3068195 "Content type detection for rest requests is deprecated. Specify the content type using the [Content-Type] header." "Wed, 17 May 2017 07:30:18 GMT"
content-type: application/json; charset=UTF-8
content-length: 146

{"_index":"megacorp","_type":"employee","_id":"1","_version":2,"result":"updated","_shards":{"total":2,"successful":1,"failed":0},"created":false}%           ➜  ~
➜  ~
# 第三次执行时，指定 "Content-Type: application/json" 解决报错问题
➜  ~ curl -u elastic:changeme -i -H "Content-Type: application/json" -XPUT 'http://localhost:9200/megacorp/employee/1' -d '{
    "first_name" : "John",
    "last_name" :  "Smith",
    "age" :        25,
    "about" :      "I love to go rock climbing",
    "interests": [ "sports", "music" ]
}
'
HTTP/1.1 200 OK
content-type: application/json; charset=UTF-8
content-length: 146

{"_index":"megacorp","_type":"employee","_id":"1","_version":3,"result":"updated","_shards":{"total":2,"successful":1,"failed":0},"created":false}%           ➜  ~
```

> 注意：
>
> - 如果上述命令执行失败了，可能的原因是 elasticsearch 默认配置中不允许自动创建索引，所以我们可以先简单在 `elasticsearch.yml` 配置文件添加 `action.auto_create_index: true` ，允许自动创建索引。
> - 没有必要先创建一个索引或指定每个字段所包含的数据类型（即无需执行任何管理任务），可以直接索引一个目标文档。

再创建两个测试数据

```
➜  ~ curl -u elastic:changeme -i -H "Content-Type: application/json" -XPUT 'http://localhost:9200/megacorp/employee/2' -d '{
    "first_name" :  "Jane",
    "last_name" :   "Smith",
    "age" :         32,
    "about" :       "I like to collect rock albums",
    "interests":  [ "music" ]
}
'
HTTP/1.1 201 Created
Location: /megacorp/employee/2
content-type: application/json; charset=UTF-8
content-length: 145

{"_index":"megacorp","_type":"employee","_id":"2","_version":1,"result":"created","_shards":{"total":2,"successful":1,"failed":0},"created":true}%
➜  ~
➜  ~ curl -u elastic:changeme -i -H "Content-Type: application/json" -XPUT 'http://localhost:9200/megacorp/employee/3' -d '{
    "first_name" :  "Douglas",
    "last_name" :   "Fir",
    "age" :         35,
    "about":        "I like to build cabinets",
    "interests":  [ "forestry" ]
}
'
HTTP/1.1 201 Created
Location: /megacorp/employee/3
content-type: application/json; charset=UTF-8
content-length: 145

{"_index":"megacorp","_type":"employee","_id":"3","_version":1,"result":"created","_shards":{"total":2,"successful":1,"failed":0},"created":true}%
➜  ~
```

### 检索文档

略

### DSL 查询

查询字符串搜索（query string）对于从命令行进行搜索非常方便，但它有其局限性。Elasticsearch 提供了一种丰富，灵活的查询语言，称为**查询 DSL** ，它允许我们构建更复杂，更健壮的查询。

```
GET /megacorp/employee/_search
{
    "query" : {
        "match" : {
            "last_name" : "Smith"
        }
    }
}
```

等价于

```
curl -u elastic:changeme -i -XGET 'http://localhost:9200/megacorp/employee/_search?q=last_name:Smith'
```

> 注意：
> 
> - 上面这个内容可以直接贴到 Dev tools 的 console 上执行；
> - 上面使用了匹配查询 `match` ；


```
GET /megacorp/employee/_search
{
    "query" : {
        "bool" : {
            "must" : {
                "match" : {
                    "last_name" : "smith" 
                }
            },
            "filter" : {
                "range" : {
                    "age" : { "gt" : 30 } 
                }
            }
        }
    }
}
```

> 添加了一个过滤器 filter ，并执行了 range 范围搜索；

### 全文搜索（Full-Text Search）

```
GET /megacorp/employee/_search
{
    "query" : {
        "match" : {
            "about" : "rock climbing"
        }
    }
}
```

得到

```
{
  "took": 7,
  "timed_out": false,
  "_shards": {
    "total": 5,
    "successful": 5,
    "failed": 0
  },
  "hits": {
    "total": 2,
    "max_score": 0.53484553,
    "hits": [
      {
        "_index": "megacorp",
        "_type": "employee",
        "_id": "1",
        "_score": 0.53484553,
        "_source": {
          "first_name": "John",
          "last_name": "Smith",
          "age": 25,
          "about": "I love to go rock climbing",
          "interests": [
            "sports",
            "music"
          ]
        }
      },
      {
        "_index": "megacorp",
        "_type": "employee",
        "_id": "2",
        "_score": 0.26742277,
        "_source": {
          "first_name": "Jane",
          "last_name": "Smith",
          "age": 32,
          "about": "I like to collect rock albums",
          "interests": [
            "music"
          ]
        }
      }
    ]
  }
}
```

默认情况下，Elasticsearch 按匹配结果的相关性分值（即每个文档与查询匹配程度）对匹配结果进行排序。可以看出，这里并不是完全匹配目标字段；

### 精确字段搜索

当你想要匹配字词或短语的确切序列时，需要使用 match_phrase ；

```
GET /megacorp/employee/_search
{
    "query" : {
        "match_phrase" : {
            "about" : "rock climbing"
        }
    }
}
```

### 高亮搜索结果

许多应用程序喜欢从每个搜索结果突出显示文本片段，以便用户可以看到文档与查询匹配的原因。添加一个新的 highlight 参数；

```
GET /megacorp/employee/_search
{
    "query" : {
        "match_phrase" : {
            "about" : "rock climbing"
        }
    },
    "highlight": {
        "fields" : {
            "about" : {}
        }
    }
}
```

执行后得到

```
{
  "took": 73,
  "timed_out": false,
  "_shards": {
    "total": 5,
    "successful": 5,
    "failed": 0
  },
  "hits": {
    "total": 1,
    "max_score": 0.53484553,
    "hits": [
      {
        "_index": "megacorp",
        "_type": "employee",
        "_id": "1",
        "_score": 0.53484553,
        "_source": {
          "first_name": "John",
          "last_name": "Smith",
          "age": 25,
          "about": "I love to go rock climbing",
          "interests": [
            "sports",
            "music"
          ]
        },
        "highlight": {
          "about": [
            "I love to go <em>rock</em> <em>climbing</em>"
          ]
        }
      }
    ]
  }
}
```

返回结果中包含了一个称为**突出显示**的内容；

### 分析

Elasticsearch 具有称为**聚合**的功能，允许您对数据生成复杂的分析。它类似于 SQL 中的 GROUP BY ，但功能更强大。

```
GET /megacorp/employee/_search
{
  "aggs": {
    "all_interests": {
      "terms": { "field": "interests" }
    }
  }
}
```

执行后得到

```
{
  "error": {
    "root_cause": [
      {
        "type": "illegal_argument_exception",
        "reason": "Fielddata is disabled on text fields by default. Set fielddata=true on [interests] in order to load fielddata in memory by uninverting the inverted index. Note that this can however use significant memory."
      }
    ],
    "type": "search_phase_execution_exception",
    "reason": "all shards failed",
    "phase": "query",
    "grouped": true,
    "failed_shards": [
      {
        "shard": 0,
        "index": "megacorp",
        "node": "Cyn9wLoOQWyrkIpBhGY7Iw",
        "reason": {
          "type": "illegal_argument_exception",
          "reason": "Fielddata is disabled on text fields by default. Set fielddata=true on [interests] in order to load fielddata in memory by uninverting the inverted index. Note that this can however use significant memory."
        }
      }
    ],
    "caused_by": {
      "type": "illegal_argument_exception",
      "reason": "Fielddata is disabled on text fields by default. Set fielddata=true on [interests] in order to load fielddata in memory by uninverting the inverted index. Note that this can however use significant memory."
    }
  },
  "status": 400
}
```

上述聚合操作再 5.0 版本之前是能够成功返回的；

5.0 之后[返回上述错误的原因](https://www.elastic.co/guide/en/elasticsearch/reference/5.0/fielddata.html#_fielddata_is_disabled_on_literal_text_literal_fields_by_default)在于：Fielddata 会消耗大量的堆空间（内存），尤其当加载高基（cardinality）text 域（field）时。一旦 fielddata 加载到堆中，其在该段（segment）的整个生存期内都存在。此外，加载 fielddata 也是一个昂贵的过程，可能导致用户经历延迟命中问题。这就是 fielddata 被默认禁用的原因。如果尝试基于脚本对 `text` 域中的值进行排序（sort），聚合（aggregate）或访问（access），就会看到上述异常。

> 注意：聚合结果不是预先计算的，而是从与当前查询匹配的文档即时生成的；
>
> ES最大的使用场景是对实时数据的查询，而这种数据可以是秒级上万数据的插入，这对数据库来说是比较困难的。所以对于存储来说，可以这样看待，首先数据库一般会存储一些基本信息，就拿上面的例子来说，员工的个人信息一般来说是存Mysql之类的关系数据库（例子为了介绍所以使用比较简单的数据），但是比方说每个员工一天的工作流水等实时性的数据，就可以直接存储到ES，并提供实时的查询。因为数据量大，不利于查询等原因，这些数据一般也不需要存到数据库，所以也没有什么数据需要ES和数据库同步（不排除特殊情况）。现在工作中也有碰到一些需要存储的，不过主要就是用于备份，或者真的ES出现问题，临时数据库顶替作为降级方法处理，将实时数据里面用户比较关心的数据备份到数据库。


## [ElasticSearch 5学习(6)——分布式集群学习分享1](http://www.cnblogs.com/wxw16/p/6188044.html)

略

## [ElasticSearch 5学习(7)——分布式集群学习分享2](http://www.cnblogs.com/wxw16/p/6188560.html)

略

## [ElasticSearch 5学习(8)——分布式文档存储（wait_for_active_shards新参数分析）](http://www.cnblogs.com/wxw16/p/6192549.html)

略

## [ElasticSearch 5学习(9)——映射和分析（string类型废弃）](http://www.cnblogs.com/wxw16/p/6195284.html)

在 ElasticSearch 中，存入文档（document）的内容类似于传统数据库中的每个字段一样，都会有一个指定的属性，为了能够把日期字段处理成日期，把数字字段处理成数字，把字符串字段处理成字符串值，Elasticsearch 需要知道每个字段里面都包含了什么类型。这些**类型**和**字段**的信息存储（包含）在**映射（mapping）**中。

当你索引一个之前没有的字段时，Elasticsearch 将使用**动态映射**猜测字段类型，这类型来自于 JSON 的基本数据类型；


| JSON type | Field type |
| -- | -- |
| Boolean: `true` of `false` | `"boolean"` |
| whole number: `123` | `"long"` |
| Floating point: `123.456` | `"double"` |
| String, Valid date: `"2014-09-15"` | `"date"` |
| String: `"foo bar"` | `"string"` |

### 查看映射

使用 `_mapping` 后缀来查看 Elasticsearch 中的映射

```
GET /megacorp/_mapping/employee
```

执行后得到

```
{
  "megacorp": {
    "mappings": {
      "employee": {
        "properties": {
          "about": {
            "type": "text",
            "fields": {
              "keyword": {
                "type": "keyword",
                "ignore_above": 256
              }
            }
          },
          "age": {
            "type": "long"
          },
          "first_name": {
            "type": "text",
            "fields": {
              "keyword": {
                "type": "keyword",
                "ignore_above": 256
              }
            }
          },
          "interests": {
            "type": "text",
            "fields": {
              "keyword": {
                "type": "keyword",
                "ignore_above": 256
              }
            }
          },
          "last_name": {
            "type": "text",
            "fields": {
              "keyword": {
                "type": "keyword",
                "ignore_above": 256
              }
            }
          }
        }
      }
    }
  }
}
```

上面输出的**映射**内容是 Elasticsearch 创建索引时动态生成的；

在 Elasticsearch 中的每种字段（field）都有不同的索引处理方式；Elasticsearch 对每一种核心数据类型（string, number, boolean 及 date）以不同的方式进行索引。

Elasticsearch 中的数据可以大致分为两种类型：

- **确切值**
- **全文文本**

为了方便在全文文本字段中进行这些类型的查询，Elasticsearch 首先对文本进行分析（analyze），然后使用结果建立一个**倒排索引**。

### 倒排索引

Elasticsearch 使用一种叫做**倒排索引（inverted index）**的结构来做快速的全文搜索。倒排索引由在文档中出现的唯一的单词列表，以及对于每个单词在文档中的位置组成。

为了创建倒排索引，我们首先切分每个文档为单独的单词（称作 `terms` 或者 `tokens`）；在搜索时，若存在多个文档都匹配，则通过加入简单的相似度算法（similarity algorithm）进行排序；

倒排索引中会遇到的问题：

- 大小写问题（Quick 和 quick）
- 同根词问题（fox 和 foxes）
- 同义词问题（jumped 和 leap）

注意：我们只可以找到确实存在于索引中的词，所以**索引文本**和**查询字符串**都要标准化为相同的形式，而这个标记化和标准化的过程叫做`分词(analysis)`；

### 分析和分析器

**分析（analysis）**是这样一个过程：

- 首先，标记化一个文本块为适用于倒排索引单独的词（term）；
- 然后，标准化这些词为标准形式，提高它们的“可搜索性”或“查全率”。

这个工作是`分析器（analyzer）`完成的。一个分析器包含三个功能：

- **字符过滤器（character filter）**：在标记化前，字符串经过字符过滤器进行处理（去除和转换）；
- **分词器（tokenizer）**：字符串被分词器标记化成独立的词；
- **标记过滤（token filters）**：每个独立的词都会通过所有标记过滤，用以**修改**、**去掉**和**增加**词；

> Elasticsearch 提供很多开箱即用的字符过滤器，分词器和标记过滤器。

Elasticsearch 还附带了一些预装的分析器，你可以直接使用它们：

- 标准分析器：对于文本分析，它对于任何语言都是最佳选择；
- 简单分析器：将非单个字母的文本切分，然后把每个词转为小写；
- 空格分析器：依据空格切分文本，不转换小写；
- 语言分析器：特定语言分析器适用于很多语言，能够考虑到特定语言的特性；
- 测试分析器：用于查看文本是如何被分析的；

### 当分析器被使用

在索引中有 12 个 tweets ，只有一个包含日期 2014-09-15 ，但是我们看看下面查询中的结果 total 值。

```
GET /_search?q=2014              # 12 个结果
GET /_search?q=2014-09-15        # 还是 12 个结果 !
GET /_search?q=date:2014-09-15   # 1  一个结果
GET /_search?q=date:2014         # 0  个结果 !
```

问题：

- 为什么全日期的查询返回所有的 tweets ，而针对 `date` 字段进行年度查询却什么都不返回？
- 为什么我们的结果因查询 `_all` 字段（默认所有字段中进行查询）或 `date` 字段而变得不同？

现在我们来看为什么会产生这样的结果：

- `date` 字段包含一个**确切值**：单独的一个词 "2014-09-15" 。
- `_all` 字段是一个**全文字段**，所以分析过程会将日期转为三个词："2014"、"09" 和 "15" 。
- 当我们在 `_all` 字段查询 2014 时，一定会匹配到 12 条推文，因为这些推文都包含词 2014 。
- 当我们在 `_all` 字段中查询 2014-09-15 时，首先分析查询字符串，产生匹配任一词 2014、09 或 15 的查询语句，它依旧匹配 12 个推文，因为它们都包含词 2014 。
- 当我们在 `date` 字段中查询 2014-09-15 时，它查询一个确切的日期，然后只找到一条推文。
- 当我们在 `date` 字段中查询 2014 时，没有找到文档，因为没有文档包含那个确切的日期。

### 指定分析器

当 Elasticsearch 在你的文档中探测到一个新的字符串字段（field）时，它将自动设置它为**全文 string（弃用）字段**并用**标准分析器**分析。

你不可能总是想要这样做，也许你想使用一个更适合这个数据的语言分析器。或者你只想把字符串字段（field）当作一个普通的字段，不做任何分析，只存储确切值，就像字符串类型的用户 ID 或者内部状态字段或者标签。

为了达到这种效果，我们必须通过自定义**映射（mapping）**人工设置。

映射中最重要的字段参数是 `type` ，除了 string（弃用）类型的字段，你可能很少需要映射其他的 `type` ，因为一般情况下，Elasticsearch 自动帮我们映射的类型都能满足我们需求；对于 string（弃用）字段，两个最重要的映射参数是 `index` 和 `analyzer` 。`index` 参数控制字符串以何种方式被索引（`analyzed`/`not_analyzed`/`no`）。对于 `analyzed` 类型的字符串字段，使用 `analyzer` 参数来指定哪一种分析器将在搜索和索引的时候使用，默认 Elasticsearch 使用 standard 分析器，但是你可以通过指定一个内建的分析器来更改它，例如 whitespace、simple 或 english 等。

> 注意：其他简单类型（long、double、date 等等）也接受 index 参数，但相应的值只能是 no 和 not_analyzed ，它们的值不能被分析。

### 更新映射

你可以在第一次创建索引的时候指定映射的类型。此外，你也可以晚些时候为新类型添加映射（或者为已有的类型更新映射）。

重要：你可以向已有映射中增加字段（field），但你不能修改它。如果一个字段在映射中已经存在，这可能意味着那个字段的数据已经被索引。如果你改变了字段映射，那已经被索引的数据将错误并且不能被正确的搜索到。


## [ElasticSearch 5学习(10)——结构化查询（包括新特性）](http://www.cnblogs.com/wxw16/p/6204644.html)

之前我们所有的查询都属于**命令行查询**，这种方式不利于完成复杂的查询，而且一般在项目开发中也不会使用命令行查询方式，只有在调试测试时才使用；因此，如果想要用好搜索功能，则必须使用**请求体查询（request body search）**API 。之所以这么称呼，是因为大多数参数都是以 JSON 格式指定的，而非查询字符串（query string）。请求体查询并不仅仅用来处理查询，而且还可以高亮返回结果中的片段，并且能够给出帮助用户找寻最好结果相关数据的建议。

### 空查询

空查询将会返回索引中所有的文档；

```
GET /_search
{}
```

可以使用 `from` 及 `size` 参数进行分页：

```
GET /_search
{
  "from": 30,
  "size": 10
}
```

或者

```
POST /_search
{
  "from": 30,
  "size": 10
}
```

具体说明详见《[HTTP 之 GET 在 body 中携带参数问题](https://github.com/moooofly/MarkSomethingDown/blob/master/nonsense/HTTP%20%E4%B9%8B%20GET%20%E5%9C%A8%20body%20%E4%B8%AD%E6%90%BA%E5%B8%A6%E5%8F%82%E6%95%B0%E9%97%AE%E9%A2%98.md)》；

### 结构化查询 Query DSL

**结构化查询**是一种**灵活**的，**多表现形式**的查询语言。Elasticsearch 在一个简单的 JSON 结构中用结构化查询来展现 Lucene 绝大多数能力。

你应当在产品中采用这种方式进行查询，会使得你的查询更加灵活，精准，易于阅读并且易于 debug 。

使用结构化查询，你需要传递 `query` 参数：

```
GET /_search
{
    "query": 发查询体放置于此即可
}
```

**空查询**在功能上等同于使用 `match_all` 查询子句的结构化查询，后者正如其名字一样，会匹配所有的文档：

```
GET /_search
{
    "query": {
        "match_all": {}  # 查询体
    }
}
```

#### 查询子句

一个**查询子句**一般使用这种结构：

```
# 整个属于查询体
{
    QUERY_NAME（查询命令）: {
        ARGUMENT: VALUE,
        ARGUMENT: VALUE,...
    }
}
```

或指向一个指定的字段：

```
#整个属于查询体
{
    QUERY_NAME（查询命令）: {
        FIELD_NAME（匹配字段）: {
            ARGUMENT: VALUE,
            ARGUMENT: VALUE,...
        }
    }
}
```

例如，你可以使用 `match` **查询子句**用来找寻在 tweet 字段中找寻包含 elasticsearch 的成员：

```
GET /_search
{
    "query": {
        "match": {
            "tweet": "elasticsearch"
        }
    }
}
```

#### 多子句合并

查询子句就像是搭积木一样，可以合并简单的子句为一个复杂的查询语句，比如：

- 叶子子句（`leaf`）：用以将查询字符串与一个字段或多字段进行比较，比如 `match` 子句；
- 复合子句（`compound`）：用以合并其他的子句。例如，`bool` 子句允许你合并其他的合法子句，`must`，`must_not` 或者 `should` ，如果可能的话；

```
{
    "bool": {
        "must":     { "match": { "tweet": "elasticsearch" }},
        "must_not": { "match": { "name":  "mary" }},
        "should":   { "match": { "tweet": "full text" }}
    }
}
```

复合子句能合并任意其他查询子句，包括其他的复合子句。这就意味着复合子句可以相互嵌套，从而实现非常复杂的逻辑。

以下实例查询的是：邮件正文中含有 “business opportunity” 字样的**星标**邮件，或**收件箱**中正文中含有 “business opportunity” 字样的**非垃圾**邮件：

```
#整个属于查询体
{
    "bool": {
        "must": { "match":      { "email": "business opportunity" }},
        "should": [
             { "match":         { "starred": true }},
             { "bool": {
                   "must":      { "folder": "inbox" },
                   "must_not":  { "spam": true }
             }}
        ],
        "minimum_should_match": 1
    }
}
```

### 查询与过滤

Elasticsearch 使用的 DSL 具有一组称为**查询**的组件，它们可以混合并以无穷组合进行匹配。这一组组件可以在两个上下文中使用：**过滤上下文**和**查询上下文**。

当用于**过滤上下文**时，该查询被称为“**非评分**”或“**过滤**”查询。也就是说，查询只询问问题：“此文档是否匹配？”。答案总是一个简单的二进制 `yes|no` 。

- `created` 的日期范围是否在 2013 到 2014 ?
- `status` 字段中是否包含单词 "published" ?
- `lat_lon` 字段中的地理位置与目标点相距是否不超过 10km ?

当在**查询上下文**中使用时，查询变为“**评分**”查询。类似于其非评分兄弟，这确定文档是否匹配以及文档匹配的程度。

查询的典型用法：

- 查找与 `full text search` 这个词语最佳匹配的文档；
- 查找包含单词 `run` ，但是也包含 `runs`, `running`, `jog` 或 `sprint` 的文档；
- 同时包含着 `quick`, `brown` 和 `fox` --- 单词间离得越近，该文档的相关性越高；
- 标识着 `lucene`, `search` 或 `java` --- 标识词越多，该文档的相关性越高；

评分查询计算每个文档与查询的相关程度，并为其分配相关性 `_score` ，稍后用于按相关性对匹配文档进行排序。这种相关性的概念非常适合于全文搜索，其中很少有完全“正确”的答案。

> 新特性：历史上，**查询**和**过滤器**是 Elasticsearch 中的单独组件。从 Elasticsearch 2.0 开始，过滤器在技术上被消除，并且所有查询都获得了成为非评分的能力。

然而，为了清楚和简单，将使用 `term` “过滤器”来表示在**非评分过滤上下文**中使用的查询。可以将 `term` “过滤器”，“过滤查询”和“非评分查询”视为相同。

类似地，如果单独使用 `term` “查询”而不使用限定符，指的是“评分查询”。

#### 性能差异

**过滤查询**是对集合包含/排除的简单检查，这使得计算非常快。当您的过滤查询中至少有一个是“稀疏”（匹配文档较少）时，可以利用各种优化，并且可以将经常使用的非评分查询缓存在内存中以便更快地访问。

相比之下，**评分查询**不仅必须找到匹配的文档，而且还要计算每个文档的相关程度，这通常使得他们比他们的非评分对手更重。此外，查询结果不可缓存。

由于倒排索引，只匹配几个文档的**简单评分查询**可能与跨越数百万个文档的**过滤器**一样好或更好。然而**一般来说，过滤器将胜过评分查询**。

**过滤的目的**是减少必须由评分查询检查的文档的数量。

#### 使用原则

作为一般规则，对全文搜索或任何会影响相关性分数的条件使用查询子句，并对其他所有条件使用过滤器。

### 最重要的查询过滤语句

#### match_all 查询

`match_all` 查询只匹配所有文档。如果未指定任何查询，则是使用的默认查询：

```
{“match_all”：{}}
```

此查询经常与过滤器结合使用，例如，用于检索收件箱文件夹中的所有电子邮件。所有文件被认为是同等相关的，所以他们都获得 1 的中性分数。

#### match 查询

`match` 查询是一个标准查询，不管你需要**全文本查询**还是**精确查询**基本上都要用到它。

如果你使用 `match` 查询一个全文本字段，它会在真正查询之前用分析器先 `match` 一下查询字符：

```
{
    "match": {
        "tweet": "About Search"
    }
}
```

如果在 `match` 下指定了一个确切值，在遇到数字，日期，布尔值或者 `not_analyzed` 的字符串时，它将为你搜索你给定的值：

```
{ "match": { "age":    26           }}
{ "match": { "date":   "2014-09-01" }}
{ "match": { "public": true         }}
{ "match": { "tag":    "full_text"  }}
```

> 提示：做精确匹配搜索时，你最好用过滤语句，因为过滤语句可以缓存数据。

不像我们之前介绍的字符查询，`match` 查询不可以用类似 "`+usid:2 +tweet:search`" 这样的语句。它只能针对指定的某个确切字段、某个确切的值进行搜索，而你要做的就是为它指定正确的字段名以避免语法错误。

#### multi_match 查询

`multi_match` 查询允许你在 `match` 查询的基础上，同时搜索多个字段：

```
{
    "multi_match": {
        "query":    "full text search",
        "fields":   [ "title", "body" ]
    }
}
```

#### range 过滤

`range` 过滤允许我们按照指定范围查找一批数据：

```
{
    "range": {
        "age": {
            "gte":  20,
            "lt":   30
        }
    }
}
```

范围操作符包含：

- gt ：大于
- gte：大于等于
- lt ：小于
- lte：小于等于

#### term 查询

`term` 用于按照**精确值**进行搜索，无论是数字，日期，布尔值，还是未分析的精确值字符串字段：

```
{ "term": { "age":    26           }}
{ "term": { "date":   "2014-09-01" }}
{ "term": { "public": true         }}
{ "term": { "tag":    "full_text"  }}
```

`term` 不对输入文本执行分析，因此它将精确查找提供的值。

#### terms 查询

`terms` 查询与 `term` 查询相同，但允许您指定多个值进行匹配。如果字段包含任何指定的值，则文档匹配：

```
{ "terms": { "tag": [ "search", "full_text", "nosql" ] }}
```

与 `term` 查询类似，不对输入文本执行分析。它正在寻找精确匹配（包括大小写，重音，空格等）。

#### exists 和 missing 查询

`exists` 和 `missing` 查询用于查找指定字段是否具有一个或多个值（`exists`），或者是否没有任何值（`missing`）的文档。

它在本质上类似于 SQL 中的 `IS_NULL`（缺失）和 `NOT IS_NULL`（存在）：

```
{
    "exists":   {
        "field":    "title"
    }
}
```

这些查询经常用于仅在存在字段时应用条件，以及在缺少条件时应用不同的条件。

### 查询与过滤条件的合并

现实世界的搜索请求从来不简单，他们使用各种输入文本搜索多个字段，并根据条件数组进行过滤。要构建复杂的搜索，您需要一种将多个查询组合到一个搜索请求中的方法。

要做到这一点，你可以使用 `bool` 查询。此查询在用户定义的布尔组合中将多个查询组合在一起。此查询接受以下参数：

#### bool 过滤

`bool` 过滤可以用来合并多个过滤条件的查询结果，包含一下操作符：

- `must` ：多个查询条件的完全匹配，相当于 `and` 。
- `must_not` ：多个查询条件的相反匹配，相当于 `not` 。
- `should` ：至少有一个查询条件匹配，相当于 `or` 。

这些参数可以分别继承一个过滤条件或者一个过滤条件的数组：

```
{
    "bool": {
        "must":     { "term": { "folder": "inbox" }},
        "must_not": { "term": { "tag":    "spam"  }},
        "should": [
                    { "term": { "starred": true   }},
                    { "term": { "unread":  true   }}
        ]
    }
}
```

因为这是我们见过的第一个包含其他查询的查询，所以我们需要谈论**分数是如何组合的**。**每个子查询子句将单独计算文档的相关性分数**。一旦计算了这些分数，`bool` 查询将会把分数合并在一起，并返回表示布尔运算的总分数的单个分数。

以下查询将会找到 title 字段中包含 "how to make millions"，并且 tag 字段没有被标为 "spam" 的内容。如果有标识为 "starred" 或者发布日期为 2014 年之前，那么这些匹配的文档将比同类网站等级高：

```
{
    "bool": {
        "must":     { "match": { "title": "how to make millions" }},
        "must_not": { "match": { "tag":   "spam" }},
        "should": [
            { "match": { "tag": "starred" }},
            { "range": { "date": { "gte": "2014-01-01" }}}
        ]
    }
}
```

> 提示：如果 `bool` 查询下没有 `must` 子句，那至少应该有一个 `should` 子句。但是如果有 `must` 子句，那么没有 `should` 子句都可以进行查询。

#### 添加过滤查询

如果我们不希望文档的日期影响评分，我们可以重新排列前面的示例，例如使用**过滤子句**实现：

```
{
    "bool": {
        "must":     { "match": { "title": "how to make millions" }},
        "must_not": { "match": { "tag":   "spam" }},
        "should": [
            { "match": { "tag": "starred" }}
        ],
        "filter": {
          "range": { "date": { "gte": "2014-01-01" }} 
        }
    }
}
```

范围查询已从 `should` 子句中移出，并进入过滤器子句。

通过将范围查询**移动**到**过滤子句**中，我们将其**转换**为**非评分查询**，它将不再为文档的相关性排名贡献分数。并且因为它现在是一个非评分查询，它可以使用可用于过滤器的各种优化，这应该提高性能。

任何查询都可以以这种方式使用。只需将查询移动到 `bool` 查询的过滤器子句中，它就会自动转换为非评分过滤器。

如果你需要过滤许多不同的标准，**`bool` 查询本身可以用作非评分查询**。只需将它放在过滤器子句中，并继续构建布尔逻辑：

```
{
    "bool": {
        "must":     { "match": { "title": "how to make millions" }},
        "must_not": { "match": { "tag":   "spam" }},
        "should": [
            { "match": { "tag": "starred" }}
        ],
        "filter": {
          "bool": { 
              "must": [
                  { "range": { "date": { "gte": "2014-01-01" }}},
                  { "range": { "price": { "lte": 29.99 }}}
              ],
              "must_not": [
                  { "term": { "category": "ebooks" }}
              ]
          }
        }
    }
}
```

通过在 `filter` 子句中嵌入 `bool` 查询，我们可以为我们的过滤条件添加布尔逻辑。

#### constant_score 查询

尽管不像 `bool` 查询那样经常使用，但是 `constant_score` 查询在你的工具箱中仍然有用。查询对所有匹配的文档应用**静态、常数得分**。它主要用于当你想执行一个过滤器，没有别的（例如没有评分查询）。你可以使用它，而不是一个只有过滤器子句的 `bool` 。性能将是相同的，但它以帮助查询简单/清晰。

```
{
    "constant_score":   {
        "filter": {
            "term": { "category": "ebooks" } 
        }
    }
}
```

### 验证查询

查询语句可以变得非常复杂，特别是与不同的**分析器**和**字段映射**相结合后，就会有些难度。

**`validate API` 可以验证一条查询语句是否合法**。

```
GET /gb/tweet/_validate/query
{
   "query": {
      "tweet" : {
         "match" : "really powerful"
      }
   }
}
```

以上请求的返回值告诉我们这条语句是非法的：

```
{
  "valid" :         false,
  "_shards" : {
    "total" :       1,
    "successful" :  1,
    "failed" :      0
  }
}
```

#### 理解错误信息

要找出为什么它无效，请将 `explain` 参数添加到查询字符串：

```
GET /gb/tweet/_validate/query?explain 
{
   "query": {
      "tweet" : {
         "match" : "really powerful"
      }
   }
}
```

显然，我们已经将查询（match）类型与字段名称（tweet）混淆：

```
{
  "valid" :     false,
  "_shards" :   { ... },
  "explanations" : [ {
    "index" :   "gb",
    "valid" :   false,
    "error" :   "org.elasticsearch.index.query.QueryParsingException:
                 [gb] No query registered for [tweet]"
  } ]
}
```

#### 理解查询语句

`explain` 参数的使用具有返回（获取）与查询相关的可读描述的附加优点，这对理解 Elasticsearch 如何解释查询非常有用：

```
GET /gb/tweet/_validate/query?explain
{
   "query": {
      "tweet" : {
         "match" : "really powerful"
      }
   }
}
```

为每个我们查询的索引返回一个解释，因为每个索引可以有不同的映射和分析器：

```
{
  "valid" :         true,
  "_shards" :       { ... },
  "explanations" : [ {
    "index" :       "us",
    "valid" :       true,
    "explanation" : "tweet:really tweet:powerful"
  }, {
    "index" :       "gb",
    "valid" :       true,
    "explanation" : "tweet:realli tweet:power"
  } ]
}
```

从解释中，您可以看到查询字符串的 `match` 查询 "really powerful" 已被重写为对 tweet 字段的两个单项查询，每个 `term` 一个。

此外，对于我们的 us 索引，这两个 `term` 是 really 和 powerful 的，而对于 gb 索引，`term` 是 realli 和 power 。原因在于我们已经将 gb 索引中 tweet 字段的分析器改成了 english 分析器。
