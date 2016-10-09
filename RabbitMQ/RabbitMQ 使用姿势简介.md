


# 初级用户指令


## 单机集群构建（1 disc + 2 ram）

```shell
RABBITMQ_NODE_PORT=5672 RABBITMQ_NODENAME=rabbit_1 RABBITMQ_SERVER_START_ARGS="-rabbitmq_management listener [{port,15672}]" rabbitmq-server -detached

RABBITMQ_NODE_PORT=5673 RABBITMQ_NODENAME=rabbit_2 RABBITMQ_SERVER_START_ARGS="-rabbitmq_management listener [{port,15673}]" rabbitmq-server -detached

RABBITMQ_NODE_PORT=5674 RABBITMQ_NODENAME=rabbit_3 RABBITMQ_SERVER_START_ARGS="-rabbitmq_management listener [{port,15674}]" rabbitmq-server -detached

rabbitmqctl -n rabbit_2 stop_app
rabbitmqctl -n rabbit_2 join_cluster rabbit_1@`hostname -s`
rabbitmqctl -n rabbit_2 change_cluster_node_type ram
rabbitmqctl -n rabbit_2 start_app

rabbitmqctl -n rabbit_3 stop_app
rabbitmqctl -n rabbit_3 join_cluster rabbit_1@`hostname -s`
rabbitmqctl -n rabbit_3 change_cluster_node_type ram
rabbitmqctl -n rabbit_3 start_app
```

## 创建用户+设置用户角色+设置用户权限

```shell
rabbitmqctl add_user moooofly moooofly
rabbitmqctl set_user_tags moooofly administrator
rabbitmqctl set_permissions -p / moooofly ".\*" ".\*" ".\*"
```

## 远程连接

```shell
erl -sname dbg -remsh NodeName@`hostname -s` -setcookie xxx
```

## 社区插件

获取地址：[这里](https://www.rabbitmq.com/community-plugins/)

## 类似文章

- [RabbitMQ基础](http://soft.dog/2015/10/23/RabbitMQ-basic/)
- [RabbitMQ 监控](http://soft.dog/2016/01/13/rabbitmq-monitoring/)
- [RabbitMQ管理](http://soft.dog/2015/11/18/RabbitMQ-management/)
- [RabbitMQ集群I](http://soft.dog/2015/10/23/RabbitMQ-cluster/)
- [RabbitMQ集群II](http://soft.dog/2015/10/26/RabbitMQ-cluster-advanced/)
- [RabbitMQ 的CLI管理工具 rabbitmqadmin](http://soft.dog/2016/04/20/RabbitMQ-cli-rabbitmqadmin/)


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
lists:sum([erlang:system_info(wordsize) * mnesia:table_info(Tab, memory) || Tab <- mnesia:system_info(tables)]).
```


## 根据 erlang 进程 pid 获取进程注册名

假设某进程 pid 为 <0,441,0>
```erlang
Pid = pid(0,441,0).
erlang:process_info(Pid,registered_name).
```

- 若进程不存在，则返回 undefined ；
- 若进程存在，但没有注册名字，则返回 [] ；
- 若进程存在，且注册了名字，则返回 {registered_name, xxx} ；


## 根据 erlang 进程注册名获取进程 pid

在当前 node 上进行获取
```erlang
erlang:whereis(Process_Register_Name).
```

- 若进程不存在，则返回 undefined ；
- 若进程存在，则返回类似 <0,xxx,0> 的 pid 值；

在 cluster 范围内进行获取
```erlang
global:whereis_name(Process_Register_Name)
```

- 若进程不存在，则返回 undefined ；
- 若进程存在，则返回进程 pid ；
    - 若目标进程位于本地节点，则返回类似 <0,xxx,0> 的 pid 值；
    - 若目标进程位于 cluster 中的其他节点，则返回类似 <xxx,xxx,0> 的 pid 值；

> 获取目标 pid 对应的节点名，使用 node(Pid) 得到，其中 Pid 可以为上述两种形式；

## 获取 RabbitMQ 当前使用的 Erlang 版本信息

```erlang
erlang:system_info(system_version).
```
或者

```erlang
list_to_binary(string:strip(erlang:system_info(system_version), both, $\n)).
```


## 如何确定进程邮箱是否存在问题

查看进程邮箱中待处理消息长度
```erlang
erlang:process_info(Pid, messages_queue_len)
```

查看进程邮箱中待处理消息内容
```erlang
erlang:process_info(Pid, messages)
```

导致邮箱爆掉的常见情况：进程一直在阻塞等待某种回应消息；在这种情况下，可以使用
```erlang
erlang:process_info(Pid, backtrace)
```
和
```erlang
erlang:process_info(Pid, current_function)
```
辅助排查目标进程到底在等什么；


## 如何抓取 Erlang 节点中所有进程的信息

```erlang
-module(fetch).
-compile(export_all).

process_infos() ->
    filelib:ensure_dir("./log/"),
    UnixTime = calendar:datetime_to_gregorian_seconds(erlang:universaltime()),
    File = "./log/processes_info_" ++ erlang:integer_to_list(UnixTime) ++ ".log",
    {ok, Fd} = file:open(File, [write, raw, binary, append]),
    Fun = fun(Pi) ->
                   Info = io_lib:format("=>~p \n\n",[Pi]),
                  case  filelib:is_file(File) of
                        true   ->
                            io:format("."),
                            file:write(Fd, Info);
                        false  ->
                            file:close(Fd),
                            {ok, NewFd} = file:open(File, [write, raw, binary, append]),
                            file:write(NewFd, Info)
                     end
                     %%timer:sleep(5)
                 end,
    [Fun(erlang:process_info(P)) || P <- erlang:processes()].
```

基于抓取到的信息可以分析整个节点内进程的运行状态，确定一些异常的进程（如 reduction 超高的进程）

以上内容取自坚强哥的[文章](http://www.cnblogs.com/me-sa/archive/2011/11/06/erlang0013.html)；



## 单个 Erlang 进程占用多少内存计算

Erlang 进程和操作系统线程和进程相比非常轻量；

在不支持 HiPE 的 non-SMP emulator 中新建 Erlang 进程会占用 309 个字的内存；若支持 SMP 和 HiPE 则会增加一些内存占用量；

计算方式如下（在 Mac 环境中）
```erlang
➜  ~ erl
Erlang/OTP 19 [erts-8.0.2] [source] [64-bit] [smp:4:4] [async-threads:10] [hipe] [kernel-poll:false] [dtrace]

Eshell V8.0.2  (abort with ^G)
1>
1> Fun = fun() -> receive after infinity -> ok end end.
#Fun<erl_eval.20.52032458>
2>
2> {_,Bytes} = process_info(spawn(Fun), memory).
{memory,2720}
3>
3> Bytes div erlang:system_info(wordsize).
340
4>
```

可以看到上面的 Erlang VM 支持 SMP 和 HiPE ，因此内存占用比 309 个字多了一点；

内存占用中的 233 个字被用于 heap 空间（其中包含了 stack 内存）；GC 会按需增加 heap 使用；


# 获取当前 Unix 时间戳

```erlang
calendar:datetime_to_gregorian_seconds(erlang:universaltime()).
```

# remsh 使用问题

It's worth to note that `remsh' should be used with extra care, especially when it comes to connecting live production systems.

Time by time, it happens (at least to me) that someone accidentally issues `q()`. instead of ^G-q. The result is remote node shutdown what is usually not what we wanted to do, although `q().` is good habit  otherwise :)

That's why I would recommend you to use extra VM options to restrict the shell and block at least `q()` and friends (well, I've never seen anybody who accidentally typed `init:stop()` yet :-) ).

In production, we do it like this: `+Bi -stdlib restricted_shell  shell_restriction_policy_mod`
