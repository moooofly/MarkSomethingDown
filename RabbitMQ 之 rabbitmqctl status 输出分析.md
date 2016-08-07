

在使用 RabbitMQ 的过程中，经常会使用 `rabbitmqctl status` 命令查看节点状态信息，但输出的信息的具体含义是什么，以及如何判定系统是否存在隐患呢？本文试图从特定角度进行一些说明；


----------


```shell
➜  ~ rabbitmqctl status
Status of node 'rabbit_2@sunfeideMacBook-Pro' ...
[{pid,10075},
 {running_applications,
     [{rabbitmq_management,"RabbitMQ Management Console","3.6.1"},
      {rabbitmq_web_dispatch,"RabbitMQ Web Dispatcher","3.6.1"},
      {webmachine,"webmachine","1.10.3"},
      {mochiweb,"MochiMedia Web Server","2.13.0"},
      {inets,"INETS  CXC 138 49","6.3.1"},
      {compiler,"ERTS  CXC 138 10","7.0.1"},
      {rabbitmq_management_agent,"RabbitMQ Management Agent","3.6.1"},
      {rabbit,"RabbitMQ","3.6.1"},
      {os_mon,"CPO  CXC 138 46","2.4.1"},
      {mnesia,"MNESIA  CXC 138 12","4.14"},
      {syntax_tools,"Syntax tools","2.0"},
      {amqp_client,"RabbitMQ AMQP Client","3.6.1"},
      {xmerl,"XML parser","1.3.11"},
      {rabbit_common,[],"3.6.1"},
      {ssl,"Erlang/OTP SSL application","8.0"},
      {public_key,"Public key infrastructure","1.2"},
      {asn1,"The Erlang ASN1 compiler version 4.0.3","4.0.3"},
      {ranch,"Socket acceptor pool for TCP protocols.","1.2.1"},
      {crypto,"CRYPTO","3.7"},
      {sasl,"SASL  CXC 138 11","3.0"},
      {stdlib,"ERTS  CXC 138 10","3.0.1"},
      {kernel,"ERTS  CXC 138 10","5.0"}]},
 {os,{unix,darwin}},
 {erlang_version,
     "Erlang/OTP 19 [erts-8.0.2] [source] [64-bit] [smp:4:4] [async-threads:64] [hipe] [kernel-poll:true] [dtrace]\n"},
 {memory,
     [{total,56920656},
      {connection_readers,0},
      {connection_writers,0},
      {connection_channels,0},
      {connection_other,2848},
      {queue_procs,1580544},
      {queue_slave_procs,0},
      {plugins,656416},
      {other_proc,13688760},
      {mnesia,83456},
      {mgmt_db,12416},
      {msg_index,47904},
      {other_ets,1718344},
      {binary,3600952},
      {code,24306356},
      {atom,992433},
      {other_system,10230227}]},
 {alarms,[]},
 {listeners,[{clustering,25674,"::"},{amqp,5674,"0.0.0.0"}]},
 {vm_memory_high_watermark,0.6},
 {vm_memory_limit,5008664985},
 {disk_free_limit,50000000},
 {disk_free,62694096896},
 {file_descriptors,
     [{total_limit,2460},
      {total_used,3},
      {sockets_limit,2212},
      {sockets_used,0}]},
 {processes,[{limit,1048576},{used,203}]},
 {run_queue,0},
 {uptime,2039},
 {kernel,{net_ticktime,60}}]
```

对应到 `rabbit.erl` 中的代码如下
```erlang
status() ->
    S1 = [{pid,                  list_to_integer(os:getpid())},
          {running_applications, rabbit_misc:which_applications()},
          {os,                   os:type()},
          {erlang_version,       erlang:system_info(system_version)},
          {memory,               rabbit_vm:memory()},
          {alarms,               alarms()},
          {listeners,            listeners()}],
    S2 = rabbit_misc:filter_exit_map(
           fun ({Key, {M, F, A}}) -> {Key, erlang:apply(M, F, A)} end,
           [{vm_memory_high_watermark, {vm_memory_monitor,
                                        get_vm_memory_high_watermark, []}},
            {vm_memory_limit,          {vm_memory_monitor,
                                        get_memory_limit, []}},
            {disk_free_limit,          {rabbit_disk_monitor,
                                        get_disk_free_limit, []}},
            {disk_free,                {rabbit_disk_monitor,
                                        get_disk_free, []}}]),
    S3 = rabbit_misc:with_exit_handler(
           fun () -> [] end,
           fun () -> [{file_descriptors, file_handle_cache:info()}] end),
    S4 = [{processes,        [{limit, erlang:system_info(process_limit)},
                              {used, erlang:system_info(process_count)}]},
          {run_queue,        erlang:statistics(run_queue)},
          {uptime,           begin
                                 {T,_} = erlang:statistics(wall_clock),
                                 T div 1000
                             end},
          {kernel,           {net_ticktime, net_kernel:get_net_ticktime()}}],
    S1 ++ S2 ++ S3 ++ S4.
```

本文主要关注**内存使用**相关部分，因此只分析以下输出内容：

## 内存整体使用

```shell
{memory, rabbit_vm:memory()}
```

对应到 `rabbit_vm.erl` 中的代码如下

```erlang
[
	{total,               Total},         %% 总体内存分配量
	{connection_readers,  ConnsReader},   %% amqp_sup 和 ranch_conns_sup 下作为 reader 的 connection 占用的内存
	{connection_writers,  ConnsWriter},   %% amqp_sup 和 ranch_conns_sup 下作为 writer 的 connection 占用的内存
	{connection_channels, ConnsChannel},  %% amqp_sup 和 ranch_conns_sup 下 channel 占用的内存
	{connection_other,    ConnsOther},    %% amqp_sup 和 ranch_conns_sup 下其他用途 connection 占用的内存
	{queue_procs,         Qs},            %% rabbit_amqqueue_sup_sup 下 master 角色的 queue 占用的内存
	{queue_slave_procs,   QsSlave},       %% rabbit_amqqueue_sup_sup 下 slave 角色的 queue 占用的内存
	{plugins,             Plugins},       %% 启动的所有插件应用中 worker 进程占用的内存
	{other_proc,          lists:max([0, OtherProc])}, %% [1]
	{mnesia,              Mnesia},        %% mnesia 中内存表占用的内存
	{mgmt_db,             MgmtDbETS + MgmtDbProc},     %% management 插件中统计数据库   ets 表和 worker 进程占用的内存 
	{msg_index,           MsgIndexETS + MsgIndexProc}, %% 持久和临时消息索引维护 ets 表 ＋ 消息存储 worker 进程占用的内存
	{other_ets,           ETS - Mnesia - MsgIndexETS - MgmtDbETS},
	{binary,              Bin},
	{code,                Code},
	{atom,                Atom},
	{other_system,        System - ETS - Atom - Bin - Code}
].
```




作为对比，可以看一下 `erlang:memory()` 的输出；
```shell
(rabbit_1@sunfeideMacBook-Pro)7> erlang:memory().
[{total,55034248},
 {processes,16109744},
 {processes_used,16104120},
 {system,38924504},
 {atom,992433},
 {atom_used,979877},
 {binary,1146296},
 {code,24306356},
 {ets,2232600}]
(rabbit_1@sunfeideMacBook-Pro)8>
```

上述输出均为系统整体情况；




## 虚拟机内存水位设置

{vm_memory_high_watermark, {vm_memory_monitor, get_vm_memory_high_watermark, []}


## 虚拟机内存使用限制

{vm_memory_limit, {vm_memory_monitor, get_memory_limit, []}


## 磁盘使用限制

{disk_free_limit, {rabbit_disk_monitor, get_disk_free_limit, []}

## 磁盘空闲情况

{disk_free, {rabbit_disk_monitor, get_disk_free, []}


## 进程情况

{processes, [{limit, erlang:system_info(process_limit)},
                   {used, erlang:system_info(process_count)}]}

## 文件描述符使用情况

{file_descriptors, file_handle_cache:info()}



## run_queue 问题

{run_queue, erlang:statistics(run_queue)}

