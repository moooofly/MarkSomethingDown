

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

因为我主要关注**内存使用**相关信息，因此只分析一下几个输出内容：

## 内存整体使用

 {memory, rabbit_vm:memory()}


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

