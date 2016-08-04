


# 初级用户指令

## 创建用户

```shell
rabbitmqctl add_user moooofly moooofly
```

## 设置用户角色

```shell
rabbitmqctl set_user_tags moooofly administrator
```

## 设置用户权限

```shell
rabbitmqctl set_permissions -p / moooofly ".\*" ".\*" ".\*"
```

## 单机集群构建

```shell
RABBITMQ_NODE_PORT=5672 RABBITMQ_NODENAME=rabbit_1 RABBITMQ_SERVER_START_ARGS="-rabbitmq_management listener [{port,15672}]" rabbitmq-server -detached

RABBITMQ_NODE_PORT=5673 RABBITMQ_NODENAME=rabbit_2 RABBITMQ_SERVER_START_ARGS="-rabbitmq_management listener [{port,15673}]" rabbitmq-server -detached

RABBITMQ_NODE_PORT=5674 RABBITMQ_NODENAME=rabbit_3 RABBITMQ_SERVER_START_ARGS="-rabbitmq_management listener [{port,15674}]" rabbitmq-server -detached

rabbitmqctl -n rabbit_2 stop_app
rabbitmqctl -n rabbit_2 join_cluster rabbit_1@`hostname -s`
rabbitmqctl -n rabbit_2 start_app

rabbitmqctl -n rabbit_3 stop_app
rabbitmqctl -n rabbit_3 join_cluster rabbit_1@`hostname -s`
rabbitmqctl -n rabbit_3 start_app
```


----------

# 高级用户指令

## ets - 获取属于指定进程的表名字

```erlang
ProcName = msg_store_persistent.
[ets:info(T,name) || T <- ets:all(), O <- [ets:info(T, owner)], O =:= erlang:whereis(ProcName) ].
```

## ets - 获取属于指定进程的表的内存占用情况

根据 erlang 进程注册名统计 owner 为对应目标进程的 ets 表占用的内存（字节为单位）
```erlang
Procs = [msg_store_persistent, msg_store_transient].
Owners = [whereis(N) || N <- Procs].
lists:sum([erlang:system_info(wordsize) * ets:info(T, memory) || T <- ets:all(), O <- [ets:info(T, owner)], lists:member(O, Owners)]).
```

## mnesia - 获取内存表的内存占用情况

```erlang
lists:sum([ erlang:system_info(wordsize) * mnesia:table_info(Tab, memory) || Tab <- mnesia:system_info(tables)]).
```


