

# 场景一：1 disc + 2 ram

> 前提条件：设置 rabbit_1 为 disc 节点；设置 rabbit_2 和 rabbit_3 为 ram 节点；

启动 3 节点 RabbitMQ cluster（启动顺序为 rabbit_1, rabbit_2, rabbit_3）

```shell
➜  rabbitmq cluster_up.sh
Warning: PID file not written; -detached was passed.
Warning: PID file not written; -detached was passed.
Warning: PID file not written; -detached was passed.
➜  rabbitmq
```

此时，3 个节点上看到的信息是一致的
```shell
➜  rabbitmq rabbitmqctl -n rabbit_1 cluster_status
Cluster status of node 'rabbit_1@sunfeideMacBook-Pro' ...
[{nodes,[{disc,['rabbit_1@sunfeideMacBook-Pro']},
         {ram,['rabbit_3@sunfeideMacBook-Pro',
               'rabbit_2@sunfeideMacBook-Pro']}]},
 {running_nodes,['rabbit_3@sunfeideMacBook-Pro',
                 'rabbit_2@sunfeideMacBook-Pro',
                 'rabbit_1@sunfeideMacBook-Pro']},
 {cluster_name,<<"rabbit_1@sunfeideMacBook-Pro">>},
 {partitions,[]},
 {alarms,[{'rabbit_3@sunfeideMacBook-Pro',[]},
          {'rabbit_2@sunfeideMacBook-Pro',[]},
          {'rabbit_1@sunfeideMacBook-Pro',[]}]}]
➜  rabbitmq
➜  rabbitmq
➜  rabbitmq rabbitmqctl -n rabbit_2 cluster_status
Cluster status of node 'rabbit_2@sunfeideMacBook-Pro' ...
[{nodes,[{disc,['rabbit_1@sunfeideMacBook-Pro']},
         {ram,['rabbit_3@sunfeideMacBook-Pro',
               'rabbit_2@sunfeideMacBook-Pro']}]},
 {running_nodes,['rabbit_3@sunfeideMacBook-Pro',
                 'rabbit_1@sunfeideMacBook-Pro',
                 'rabbit_2@sunfeideMacBook-Pro']},
 {cluster_name,<<"rabbit_1@sunfeideMacBook-Pro">>},
 {partitions,[]},
 {alarms,[{'rabbit_3@sunfeideMacBook-Pro',[]},
          {'rabbit_1@sunfeideMacBook-Pro',[]},
          {'rabbit_2@sunfeideMacBook-Pro',[]}]}]
➜  rabbitmq
➜  rabbitmq
➜  rabbitmq rabbitmqctl -n rabbit_3 cluster_status
Cluster status of node 'rabbit_3@sunfeideMacBook-Pro' ...
[{nodes,[{disc,['rabbit_1@sunfeideMacBook-Pro']},
         {ram,['rabbit_3@sunfeideMacBook-Pro',
               'rabbit_2@sunfeideMacBook-Pro']}]},
 {running_nodes,['rabbit_1@sunfeideMacBook-Pro',
                 'rabbit_2@sunfeideMacBook-Pro',
                 'rabbit_3@sunfeideMacBook-Pro']},
 {cluster_name,<<"rabbit_1@sunfeideMacBook-Pro">>},
 {partitions,[]},
 {alarms,[{'rabbit_1@sunfeideMacBook-Pro',[]},
          {'rabbit_2@sunfeideMacBook-Pro',[]},
          {'rabbit_3@sunfeideMacBook-Pro',[]}]}]
➜  rabbitmq
```

查看 cluster 状态文件内容
```shell
➜  rabbitmq for i in {1..3}; do cat rabbit_$i/cluster_nodes.config; done
{['rabbit_1@sunfeideMacBook-Pro','rabbit_2@sunfeideMacBook-Pro','rabbit_3@sunfeideMacBook-Pro'],['rabbit_1@sunfeideMacBook-Pro']}.
{['rabbit_1@sunfeideMacBook-Pro','rabbit_2@sunfeideMacBook-Pro','rabbit_3@sunfeideMacBook-Pro'],['rabbit_1@sunfeideMacBook-Pro']}.
{['rabbit_1@sunfeideMacBook-Pro','rabbit_2@sunfeideMacBook-Pro','rabbit_3@sunfeideMacBook-Pro'],['rabbit_1@sunfeideMacBook-Pro']}.
➜  rabbitmq
➜  rabbitmq for i in {1..3}; do cat rabbit_$i/nodes_running_at_shutdown; done
['rabbit_1@sunfeideMacBook-Pro','rabbit_2@sunfeideMacBook-Pro','rabbit_3@sunfeideMacBook-Pro'].
['rabbit_1@sunfeideMacBook-Pro','rabbit_2@sunfeideMacBook-Pro','rabbit_3@sunfeideMacBook-Pro'].
['rabbit_1@sunfeideMacBook-Pro','rabbit_2@sunfeideMacBook-Pro','rabbit_3@sunfeideMacBook-Pro'].
➜  rabbitmq
```

> 补充说明：
> 
> `cluster_nodes.config` 内容格式为 `{[All],[Disc]}` ，`All` 为构成 mnesia 存储的所有节点，`Disc` 为 cluster 中 disc 类型节点；
> 
> `nodes_running_at_shutdown` 文件内容的格式为 `[Running]` ，`Running` 为当前节点停止运行前，所看到的、仍处于运行状态的节点；

按照 rabbit_1, rabbit_2, rabbit_3 先后顺序停止服务；
```shell
➜  rabbitmq cluster_down.sh
Stopping and halting node 'rabbit_1@sunfeideMacBook-Pro' ...
Stopping and halting node 'rabbit_2@sunfeideMacBook-Pro' ...
Stopping and halting node 'rabbit_3@sunfeideMacBook-Pro' ...
➜  rabbitmq
```

停止服务后进行信息确认：

文件内容无变化
```shell
➜  rabbitmq for i in {1..3}; do cat rabbit_$i/cluster_nodes.config; done
{['rabbit_1@sunfeideMacBook-Pro','rabbit_2@sunfeideMacBook-Pro','rabbit_3@sunfeideMacBook-Pro'],['rabbit_1@sunfeideMacBook-Pro']}.
{['rabbit_1@sunfeideMacBook-Pro','rabbit_2@sunfeideMacBook-Pro','rabbit_3@sunfeideMacBook-Pro'],['rabbit_1@sunfeideMacBook-Pro']}.
{['rabbit_1@sunfeideMacBook-Pro','rabbit_2@sunfeideMacBook-Pro','rabbit_3@sunfeideMacBook-Pro'],['rabbit_1@sunfeideMacBook-Pro']}.
➜  rabbitmq
```

文件内容发生具有一定规律的变化说明节点停止的顺序很重要！
```shell
➜  rabbitmq for i in {1..3}; do cat rabbit_$i/nodes_running_at_shutdown; done
['rabbit_1@sunfeideMacBook-Pro','rabbit_2@sunfeideMacBook-Pro','rabbit_3@sunfeideMacBook-Pro'].
['rabbit_2@sunfeideMacBook-Pro','rabbit_3@sunfeideMacBook-Pro'].
['rabbit_3@sunfeideMacBook-Pro'].
➜  rabbitmq
```

对比实验：可以看到节点停止的顺序对该文件内容的影响；

```
➜  rabbitmq rabbitmqctl -n rabbit_2 stop
Stopping and halting node 'rabbit_2@sunfeideMacBook-Pro' ...
➜  rabbitmq
➜  rabbitmq rabbitmqctl -n rabbit_1 stop
Stopping and halting node 'rabbit_1@sunfeideMacBook-Pro' ...
➜  rabbitmq
➜  rabbitmq rabbitmqctl -n rabbit_3 stop
Stopping and halting node 'rabbit_3@sunfeideMacBook-Pro' ...
➜  rabbitmq
➜  rabbitmq
➜  rabbitmq for i in {1..3}; do cat rabbit_$i/nodes_running_at_shutdown; done
['rabbit_1@sunfeideMacBook-Pro','rabbit_3@sunfeideMacBook-Pro'].
['rabbit_1@sunfeideMacBook-Pro','rabbit_2@sunfeideMacBook-Pro','rabbit_3@sunfeideMacBook-Pro'].
['rabbit_3@sunfeideMacBook-Pro'].
```

下面基于“按照 rabbit_1, rabbit_2, rabbit_3 顺序停止“进行讨论：

## 启动 rabbit_3 节点

前台启动日志会输出如下信息
```
BOOT FAILED
===========

Error description:
   {could_not_start,rabbit,
       {{failed_to_cluster_with,
            ['rabbit_1@sunfeideMacBook-Pro','rabbit_2@sunfeideMacBook-Pro'],
            "Mnesia could not connect to any nodes."},
        {rabbit,start,[normal,[]]}}}

Log files (may contain more information):
   /Users/sunfei/workspace/WGET/rabbitmq_server-3.6.1/var/log/rabbitmq/rabbit_3.log
   /Users/sunfei/workspace/WGET/rabbitmq_server-3.6.1/var/log/rabbitmq/rabbit_3-sasl.log

{"init terminating in do_boot",{could_not_start,rabbit,{{failed_to_cluster_with,['rabbit_1@sunfeideMacBook-Pro','rabbit_2@sunfeideMacBook-Pro'],"Mnesia could not connect to any nodes."},{rabbit,start,[normal,[]]}}}}

Crash dump is being written to: erl_crash.dump...done
init terminating in do_boot ()
```

在 `rabbit_3.log` 日志中输出如下信息
```shell
=INFO REPORT==== 16-Aug-2016::18:14:57 ===
Error description:
   {could_not_start,rabbit,
       {{failed_to_cluster_with,
            ['rabbit_1@sunfeideMacBook-Pro','rabbit_2@sunfeideMacBook-Pro'],
            "Mnesia could not connect to any nodes."},
        {rabbit,start,[normal,[]]}}}

Log files (may contain more information):
   /Users/sunfei/workspace/WGET/rabbitmq_server-3.6.1/var/log/rabbitmq/rabbit_3.log
   /Users/sunfei/workspace/WGET/rabbitmq_server-3.6.1/var/log/rabbitmq/rabbit_3-sasl.log
```

在 `rabbit_3-sasl.log` 日志中输出如下信息
```shell
=CRASH REPORT==== 16-Aug-2016::18:14:57 ===
  crasher:
    initial call: application_master:init/4
    pid: <0.122.0>
    registered_name: []
    exception exit: {{failed_to_cluster_with,
                         ['rabbit_1@sunfeideMacBook-Pro',
                          'rabbit_2@sunfeideMacBook-Pro'],
                         "Mnesia could not connect to any nodes."},
                     {rabbit,start,[normal,[]]}}
      in function  application_master:init/4 (application_master.erl, line 134)
    ancestors: [<0.121.0>]
    messages: [{'EXIT',<0.123.0>,normal}]
    links: [<0.121.0>,<0.31.0>]
    dictionary: []
    trap_exit: true
    status: running
    heap_size: 1598
    stack_size: 27
    reductions: 98
  neighbours:
```

从上述信息中看到，RabbitMQ 在启动的过程中触发异常；错误原因也很清楚：**在 rabbit_3 节点启动时，无法和 cluster 中的 rabbit_1 和 rabbit_2 建立连接**；

上述错误很好理解：此时 cluster 中所有节点都被停止了，因此 rabbit_3 启动后无法与其它节点建立联系；


抛出异常的代码在 `rabbit_mnesia.erl` 中
```erlang
change_extra_db_nodes(ClusterNodes0, CheckOtherNodes) ->
    ClusterNodes = nodes_excl_me(ClusterNodes0),
    %% ClusterNodes => a list of nodes that Mnesia is to try to connect to
    case {mnesia:change_config(extra_db_nodes, ClusterNodes), ClusterNodes} of
        {{ok, []}, [_|_]} when CheckOtherNodes ->
            throw({error, {failed_to_cluster_with, ClusterNodes,
                           "Mnesia could not connect to any nodes."}});
        {{ok, Nodes}, _} -> %% Nodes => those nodes in ClusterNodes that Mnesia is connected to
            Nodes
    end.
```

## 启动 rabbit_2 节点

结果和上面完全相同，除了将日志中的 rabbit_2 换成 rabbit_3 ；

## 启动 rabbit_1 节点

启动成功，没有什么特别的输出信息；

此时 `nodes_running_at_shutdown` 中的内容发生了变化
```shell
➜  rabbitmq for i in {1..3}; do cat rabbit_$i/nodes_running_at_shutdown; done
['rabbit_1@sunfeideMacBook-Pro'].
['rabbit_2@sunfeideMacBook-Pro','rabbit_3@sunfeideMacBook-Pro'].
['rabbit_3@sunfeideMacBook-Pro'].
➜  rabbitmq
```

上述内容变化说明：节点停止时会根据所看到 cluster 状态进行记录，节点启动后会根据所看到的 cluster 状态进行更新；


### 启动 rabbit_3 节点

在 `rabbit_1.log` 中会看到

```shell
=INFO REPORT==== 16-Aug-2016::18:43:16 ===
node 'rabbit_3@sunfeideMacBook-Pro' up

=INFO REPORT==== 16-Aug-2016::18:43:18 ===
rabbit on node 'rabbit_3@sunfeideMacBook-Pro' up
```
> rabbit_1 发现 rabbit_3 的启动

在 `rabbit_3.log` 中会看到

```shell
...
=INFO REPORT==== 16-Aug-2016::18:43:18 ===
cluster contains disc nodes again
```
> rabbit_3 在 cluster 中有 disc 节点时能够正常启动；


### 停止 rabbit_1 节点

此时 rabbit_3.log 中会输出

```shell
=INFO REPORT==== 16-Aug-2016::18:54:11 ===
Statistics database started.

=INFO REPORT==== 16-Aug-2016::18:54:11 ===
rabbit on node 'rabbit_1@sunfeideMacBook-Pro' down

=INFO REPORT==== 16-Aug-2016::18:54:11 ===
Keep rabbit_1@sunfeideMacBook-Pro listeners: the node is already back

=INFO REPORT==== 16-Aug-2016::18:54:11 ===
only running disc node went down

=INFO REPORT==== 16-Aug-2016::18:54:12 ===
node 'rabbit_1@sunfeideMacBook-Pro' down: connection_closed
```
> 发现集群中惟一的 disc 节点停止了


# 场景一：3 disc

