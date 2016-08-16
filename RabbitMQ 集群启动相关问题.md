

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

查看 cluster 状态信息
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

## 停止服务后进行信息确认

文件内容无变化
```shell
➜  rabbitmq for i in {1..3}; do cat rabbit_$i/cluster_nodes.config; done
{['rabbit_1@sunfeideMacBook-Pro','rabbit_2@sunfeideMacBook-Pro','rabbit_3@sunfeideMacBook-Pro'],['rabbit_1@sunfeideMacBook-Pro']}.
{['rabbit_1@sunfeideMacBook-Pro','rabbit_2@sunfeideMacBook-Pro','rabbit_3@sunfeideMacBook-Pro'],['rabbit_1@sunfeideMacBook-Pro']}.
{['rabbit_1@sunfeideMacBook-Pro','rabbit_2@sunfeideMacBook-Pro','rabbit_3@sunfeideMacBook-Pro'],['rabbit_1@sunfeideMacBook-Pro']}.
➜  rabbitmq
```

文件内容的变化说明节点停止的顺序很重要！
```shell
➜  rabbitmq for i in {1..3}; do cat rabbit_$i/nodes_running_at_shutdown; done
['rabbit_1@sunfeideMacBook-Pro','rabbit_2@sunfeideMacBook-Pro','rabbit_3@sunfeideMacBook-Pro'].
['rabbit_2@sunfeideMacBook-Pro','rabbit_3@sunfeideMacBook-Pro'].
['rabbit_3@sunfeideMacBook-Pro'].
➜  rabbitmq
```


# 场景一：3 disc

