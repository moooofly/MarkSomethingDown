


# Mac 系统

## ERROR: epmd error for host xxx: timeout (timed out)

### 问题描述

启动 RabbitMQ 服务
```shell
➜  ~ rabbitmq-server -detached
Warning: PID file not written; -detached was passed.
（卡住一段时间，大约 30s）
```

 输出如下错误信息
```shell
ERROR: epmd error for host sunfeideMacBook-Pro: timeout (timed out)
```

在卡住的过程中在 epmd 侧可以看到
```shell
➜  ~ epmd -names
epmd: up and running on port 4369 with data:
name rabbitmqprelaunch595 at port 49350
```

在 `rabbit@sunfeideMacBook-Pro.log` 中输出如下错误信息
```shell
...
=INFO REPORT==== 30-Jul-2016::15:40:01 ===
Error description:
   {could_not_start,rabbit,
       {error,
           {{shutdown,
                {failed_to_start_child,rabbit_epmd_monitor,
                    {{badmatch,noport},
                     [{rabbit_epmd_monitor,init,1,
                          [{file,"src/rabbit_epmd_monitor.erl"},{line,60}]},
                      {gen_server,init_it,6,
                          [{file,"gen_server.erl"},{line,328}]},
                      {proc_lib,init_p_do_apply,3,
                          [{file,"proc_lib.erl"},{line,247}]}]}}},
            {child,undefined,rabbit_epmd_monitor_sup,
                {rabbit_restartable_sup,start_link,
                    [rabbit_epmd_monitor_sup,
                     {rabbit_epmd_monitor,start_link,[]},
                     false]},
                transient,infinity,supervisor,
                [rabbit_restartable_sup]}}}}

Log files (may contain more information):
   /Users/sunfei/workspace/WGET/rabbitmq_server-3.6.1/var/log/rabbitmq/rabbit@sunfeideMacBook-Pro.log
   /Users/sunfei/workspace/WGET/rabbitmq_server-3.6.1/var/log/rabbitmq/rabbit@sunfeideMacBook-Pro-sasl.log
```

在 `rabbit@sunfeideMacBook-Pro-sasl.log` 中输出如下错误信息

```shell
=CRASH REPORT==== 30-Jul-2016::15:40:01 ===
  crasher:
    initial call: rabbit_epmd_monitor:init/1
    pid: <0.212.0>
    registered_name: []
    exception exit: {{badmatch,noport},
                     [{rabbit_epmd_monitor,init,1,
                          [{file,"src/rabbit_epmd_monitor.erl"},{line,60}]},
                      {gen_server,init_it,6,
                          [{file,"gen_server.erl"},{line,328}]},
                      {proc_lib,init_p_do_apply,3,
                          [{file,"proc_lib.erl"},{line,247}]}]}
      in function  gen_server:init_it/6 (gen_server.erl, line 352)
    ancestors: [rabbit_epmd_monitor_sup,rabbit_sup,<0.181.0>]
    messages: []
    links: [<0.211.0>]
    dictionary: []
    trap_exit: false
    status: running
    heap_size: 610
    stack_size: 27
    reductions: 540
  neighbours:

=SUPERVISOR REPORT==== 30-Jul-2016::15:40:01 ===
     Supervisor: {local,rabbit_epmd_monitor_sup}
     Context:    start_error
     Reason:     {{badmatch,noport},
                  [{rabbit_epmd_monitor,init,1,
                       [{file,"src/rabbit_epmd_monitor.erl"},{line,60}]},
                   {gen_server,init_it,6,[{file,"gen_server.erl"},{line,328}]},
                   {proc_lib,init_p_do_apply,3,
                       [{file,"proc_lib.erl"},{line,247}]}]}
     Offender:   [{pid,undefined},
                  {name,rabbit_epmd_monitor},
                  {mfargs,{rabbit_epmd_monitor,start_link,[]}},
                  {restart_type,transient},
                  {shutdown,4294967295},
                  {child_type,worker}]


=CRASH REPORT==== 30-Jul-2016::15:40:01 ===
  crasher:
    initial call: application_master:init/4
    pid: <0.180.0>
    registered_name: []
    exception exit: {bad_return,
                     {{rabbit,start,[normal,[]]},
                      {'EXIT',
                       {error,
                        {{shutdown,
                          {failed_to_start_child,rabbit_epmd_monitor,
                           {{badmatch,noport},
                            [{rabbit_epmd_monitor,init,1,
                              [{file,"src/rabbit_epmd_monitor.erl"},
                               {line,60}]},
                             {gen_server,init_it,6,
                              [{file,"gen_server.erl"},{line,328}]},
                             {proc_lib,init_p_do_apply,3,
                              [{file,"proc_lib.erl"},{line,247}]}]}}},
                         {child,undefined,rabbit_epmd_monitor_sup,
                          {rabbit_restartable_sup,start_link,
                           [rabbit_epmd_monitor_sup,
                            {rabbit_epmd_monitor,start_link,[]},
                            false]},
                          transient,infinity,supervisor,
                          [rabbit_restartable_sup]}}}}}}
      in function  application_master:init/4 (application_master.erl, line 134)
    ancestors: [<0.179.0>]
    messages: [{'EXIT',<0.181.0>,normal}]
    links: [<0.179.0>,<0.31.0>]
    dictionary: []
    trap_exit: true
    status: running
    heap_size: 1598
    stack_size: 27
    reductions: 98
  neighbours:
```

确认主机名配置信息
```shell
➜  ~ cat /etc/hosts
##
# Host Database
#
# localhost is used to configure the loopback interface
# when the system is booting.  Do not change this entry.
##
127.0.0.1	localhost
255.255.255.255	broadcasthost
::1             localhost
➜  ~
➜  ~ hostname
sunfeideMacBook-Pro.local
➜  ~
```

### 源码分析

在 `rabbit.erl` 中可以看到，在 boot 序列中会启动 rabbit_epmd_monitor 进程；
```erlang
-rabbit_boot_step({rabbit_epmd_monitor,
                   [{description, "epmd monitor"},
                    {mfa,         {rabbit_sup, start_restartable_child,
                                   [rabbit_epmd_monitor]}},
                    {requires,    kernel_ready},
                    {enables,     core_initialized}]}).
```

在 rabbit_epmd_monitor.erl 中，
```erlang
init([]) ->
    %% 解析 Node@Host 信息
    {Me, Host} = rabbit_nodes:parts(node()),
    %% 获取与 epmd 通信的模块名，默认 erl_epmd ，除非命令行上通过 -epmd_module 进行指定
    Mod = net_kernel:epmd_module(),
    %% 基于 Host 获取 IP 和 Port 
    %% 崩溃位置：下面函数返回 noport 导致进程崩溃，故 RabbitMQ 无法正常启动
    {port, Port, _Version} = Mod:port_please(Me, Host),
    {ok, ensure_timer(#state{mod  = Mod,
                             me   = Me,
                             host = Host,
                             port = Port})}.
```


### 解决办法

参考：[这里](http://stackoverflow.com/questions/24797947/os-x-and-rabbitmq-error-epmd-error-for-host-xxx-address-cannot-connect-to-ho)；

> try adding your hostname to your `/etc/hosts`.     
> **Sometimes Erlang distribution will get confused when your network changes**.     
> Otherwise try restarting epmd using `epmd -kill` or similar commands.    




# Ubuntu 系统



# CentOS 系统