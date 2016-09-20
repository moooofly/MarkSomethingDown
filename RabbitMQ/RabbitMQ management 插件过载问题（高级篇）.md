


本文主要针对“RabbitMQ 中 rabbit_mgmt_db 进程状态信息获取“进行研究，并源码分析了 `sys:get_state/1,2` 的使用；

----------

# 开发环境

- 测试环境：MacBook Pro OS X EI Captitan 版本 10.11.6 (15G31)
- RabbitMQ 版本：rabbitmq_server-3.6.1
- Erang 版本：Erlang/OTP 19

启动 3 节点构成的 RabbitMQ cluster ，节点名分别为
- rabbit_1 (disk node)
- rabbit_2 (ram node)
- rabbit_3 (ram node)


----------


# 常用信息获取

## 查看 rabbit_mgmt_db 进程 pid

```erlang
(rabbit_2@sunfeideMacBook-Pro)2> whereis(rabbit_mgmt_db).
<0.395.0>
(rabbit_2@sunfeideMacBook-Pro)3>
```

## 查看 rabbit_mgmt_db 的进程信息

```erlang
(rabbit_2@sunfeideMacBook-Pro)3> erlang:process_info(whereis(rabbit_mgmt_db)).
[{registered_name,rabbit_mgmt_db},
 {current_function,{gen_server2,process_next_msg,1}},
 {initial_call,{proc_lib,init_p,5}},
 {status,waiting},
 {message_queue_len,0},
 {messages,[]},
 {links,[<0.390.0>]},
 {dictionary,[{'$initial_call',{rabbit_mgmt_db,init,1}},
              {'$ancestors',[<0.390.0>,rabbit_mgmt_sup,
                             rabbit_mgmt_sup_sup,<0.360.0>]},
              {last_queue_length,0}]},
 {trap_exit,false},
 {error_handler,error_handler},
 {priority,high},
 {group_leader,<0.359.0>},
 {total_heap_size,13544},
 {heap_size,6772},
 {stack_size,7},
 {reductions,381020191},
 {garbage_collection,[{max_heap_size,#{error_logger => true,
                                       kill => true,
                                       size => 0}},
                      {min_bin_vheap_size,46422},
                      {min_heap_size,233},
                      {fullsweep_after,65535},
                      {minor_gcs,5}]},
 {suspending,[]}]
(rabbit_2@sunfeideMacBook-Pro)4>
```

有价值的内容：
- {**current_function**,{gen_server2,process_next_msg,1}} 表明当前指定的函数；
- {**message_queue_len**,0} 和 {messages,[]} 提供了进程邮箱状态信息；
- {**dictionary**, [...]} 提供当前进程的进程字典中保存的内容；
- {**last_queue_length**,0} 保存在进程字典中的、当前优先级队列中留存的消息数量；
- {**priority**,high} 表明了进程运行优先级；
- {**reductions**,381020191} 表明了进程在系统中运行所耗费的时间度量；


## 查看与 rabbit_mgmt_db 相关进程的信息

```erlang
(rabbit_2@sunfeideMacBook-Pro)4> i().
Pid                   Initial Call                          Heap     Reds Msgs
Registered            Current Function                     Stack
...
<0.360.0>             application_master:start_it/4          233    29999    0
                      application_master:loop_it/4             5
...
<0.387.0>             supervisor2:init/1                     233      757    0
rabbit_mgmt_sup_sup   gen_server:loop/6                        9
<0.389.0>             supervisor2:init/1                     233     2040    0
rabbit_mgmt_sup       gen_server:loop/6                        9
<0.390.0>             supervisor2:init/1                     233     2158    0
                      gen_server:loop/6                        9
...
<0.395.0>             rabbit_mgmt_db:init/1                  175 38154667    0
rabbit_mgmt_db        erlang:hibernate/3                       0
```

## 通过 `sys:get_state/1` 获取 rabbit_mgmt_db 进程信息

```erlang
(rabbit_2@sunfeideMacBook-Pro)1> sys:get_state(rabbit_mgmt_db).
{gs2_state,<0.390.0>,rabbit_mgmt_db,
           {state,[{channel_stats,462941},
                   {connection_stats,458844},
                   {consumers_by_channel,471135},
                   {consumers_by_queue,467038},
                   {node_node_stats,479329},
                   {node_stats,475232},
                   {queue_stats,454747}],
                  483426,487523,491620,#Ref<0.0.1.148153>,
                  {{queue_stats,{resource,<<"/">>,queue,<<"q5">>}},
                   messages_unacknowledged},
                  [{exchange,#Fun<rabbit_exchange.lookup.1>},
                   {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                  10000,#Ref<0.0.1.1044>,detailed},
           rabbit_mgmt_db,hibernate,
           {backoff,5,1000,10000,{23942,-18676,-7813}},
           {queue,[],[],0},
           [],
           {#Fun<gen_server2.9.133231575>,
            #Fun<gen_server2.10.133231575>,
            #Fun<gen_server2.8.133231575>}}
(rabbit_2@sunfeideMacBook-Pro)2>
```

有价值的内容：
- {**queue**,[],[],0} 给出了 gen_server2 自行实现的、优先级队列中尚未处理的内容；


## 通过 `sys:get_status/1` 获取 rabbit_mgmt_db 进程信息

```erlang
(rabbit_2@sunfeideMacBook-Pro)5> sys:get_status(rabbit_mgmt_db).
{status,<0.794.0>,
        {module,gen_server2},
        [[{'$initial_call',{rabbit_mgmt_db,init,1}},
          {'$ancestors',[<0.401.0>,rabbit_mgmt_sup,
                         rabbit_mgmt_sup_sup,<0.379.0>]},
          {last_queue_length,0}],
         running,<0.401.0>,[],
         [{header,"Status for generic server rabbit_mgmt_db"},
          {data,[{"Status",running},
                 {"Parent",<0.401.0>},
                 {"Logged events",[]},
                 {"Queued messages",{0,[]}}]},
          {data,[{"State",
                  {state,[{channel_stats,458774},
                          {connection_stats,454677},
                          {consumers_by_channel,467037},
                          {consumers_by_queue,462940},
                          {node_node_stats,475231},
                          {node_stats,471134},
                          {queue_stats,450580}],
                         479328,483425,487522,#Ref<0.0.1.3805>,
                         {{node_stats,'rabbit_3@sunfeideMacBook-Pro'},disk_free},
                         [{exchange,#Fun<rabbit_exchange.lookup.1>},
                          {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                         10000,#Ref<0.0.1.1448>,detailed}}]}]]}
(rabbit_2@sunfeideMacBook-Pro)6>
```

有价值的内容：
- {**last_queue_length**,0} 保存在进程字典中的、当前优先级队列中留存的消息数量；
- {"**Queued messages**",{0,[]}} 给出了 gen_server2 自行实现的、优先级队列中尚未处理的内容；

----------

# 源码分析

## sys:get_state/1

当在 erlang console 中执行 `sys:get_state(rabbit_mgmt_db).` 命令时，对应的命令执行过程如下：

在 `sys.erl` 模块中实现了高层次的系统状态获取接口；

```erlang
get_state(Name) ->
    %% 向注册名为 Name 的进程发送 get_state 系统消息
    case send_system_msg(Name, get_state) of
        {error, Reason} -> error(Reason);
        State -> State
    end.
...
send_system_msg(Name, Request) ->
    %% 向 generic server 发送 system 消息
    case catch gen:call(Name, system, Request) of
        {ok,Res} -> Res;
        {'EXIT', Reason} -> exit({Reason, mfa(Name, Request)})
    end.
```

在 `gen.erl` 模块中实现了基于 **Pid** 或**进程注册名**进行**本地**或**远端**的 call 调用；

```erlang
%% 对应上面的请求，此处有
%% Process => rabbit_mgmt_db
%% Label => system
%% Request => get_state
call(Process, Label, Request) -> 
    call(Process, Label, Request, ?default_timeout).

%% Local or remote by pid
%% 基于进程 Pid 进行本地或远端调用
%% 例如 gen:call(global:whereis_name(rabbit_mgmt_db), system, get_state, 5000).
call(Pid, Label, Request, Timeout) 
  when is_pid(Pid), Timeout =:= infinity;
       is_pid(Pid), is_integer(Timeout), Timeout >= 0 ->
    do_call(Pid, Label, Request, Timeout);
%% Local by name
%% 基于进程注册名进行本地调用
%% 例如 gen:call(rabbit_mgmt_db, system, get_state, 5000).
call(Name, Label, Request, Timeout) 
  when is_atom(Name), Timeout =:= infinity;
       is_atom(Name), is_integer(Timeout), Timeout >= 0 ->
    %% 本地（当前节点）查找
    %% 注意：这里说明基于进程注册名进行 gen:call/3,4 调用时，仅针对本地进程才有效
    case whereis(Name) of
        Pid when is_pid(Pid) ->
            do_call(Pid, Label, Request, Timeout);
        undefined ->
            exit(noproc)
    end;
%% Global by name
%% 基于进程注册名进行全局调用
%% 例如 gen:call({global, rabbit_mgmt_db}, system, get_state, 5000).
call(Process, Label, Request, Timeout)
  when ((tuple_size(Process) == 2 andalso element(1, Process) == global)
    orelse
      (tuple_size(Process) == 3 andalso element(1, Process) == via))
       andalso
       (Timeout =:= infinity orelse (is_integer(Timeout) andalso Timeout >= 0)) ->
    %% 全局查找
    case where(Process) of
        Pid when is_pid(Pid) ->
            Node = node(Pid),
                try do_call(Pid, Label, Request, Timeout)
                catch
                    exit:{nodedown, Node} ->
                        %% A nodedown not yet detected by global,
                        %% pretend that it was.
                        exit(noproc)
            end;
        undefined ->
            exit(noproc)
    end;
```

最终进行请求发送，应答接收的地方；

```erlang
%% Process => Pid | Process_Registered_Name
%% Label => system
%% Request => get_state
do_call(Process, Label, Request, Timeout) ->
    try erlang:monitor(process, Process) of
        Mref ->
            %% If the monitor/2 call failed to set up a connection to a
            %% remote node, we don't want the '!' operator to attempt
            %% to set up the connection again. (If the monitor/2 call
            %% failed due to an expired timeout, '!' too would probably
            %% have to wait for the timeout to expire.) Therefore,
            %% use erlang:send/3 with the 'noconnect' option so that it
            %% will fail immediately if there is no connection to the
            %% remote node.

            %% 发送 system 消息到目标进程
            %% 对应我们的调用，推出
            %% Process => rabbit_mgmt_db 进程 pid
            %% Label => system
            catch erlang:send(Process, {Label, {self(), Mref}, Request}, [noconnect]),
            receive
                {Mref, Reply} ->  %% 收到应答，返回给最开始的 sys:get_state/1,2 调用
                    erlang:demonitor(Mref, [flush]),
                    {ok, Reply};
                {'DOWN', Mref, _, _, noconnection} ->
                    Node = get_node(Process),
                    exit({nodedown, Node});
                {'DOWN', Mref, _, _, Reason} ->
                    exit(Reason)
            after Timeout ->
                    erlang:demonitor(Mref, [flush]),
                    exit(timeout)
            end
    catch
        error:_ ->
            %% Node (C/Java?) is not supporting the monitor.
            %% The other possible case -- this node is not distributed
            %% -- should have been handled earlier.
            %% Do the best possible with monitor_node/2.
            %% This code may hang indefinitely if the Process 
            %% does not exist. It is only used for featureweak remote nodes.
            Node = get_node(Process),
            monitor_node(Node, true),
            receive
                {nodedown, Node} -> 
                    monitor_node(Node, false),
                    exit({nodedown, Node})
                after 0 -> 
                    Tag = make_ref(),
                    Process ! {Label, {self(), Tag}, Request},
                    wait_resp(Node, Tag, Timeout)
            end
    end.
```

通过 `erlang:send/3` 发送请求后，rabbit_mgmt_db 进程的邮箱中会收到对应的 *{system, {self(), Mref}, get_state}* 消息；最终 `gen_server2.erl` 中的如下代码进行处理

```erlang
in({system, _From, _Req} = Input, GS2State) ->
    in(Input, infinity, GS2State);
...
%% 根据 Priority 放入不同的优先级队列
in(Input, Priority, GS2State = #gs2_state { queue = Queue }) ->
    GS2State # gs2_state { queue = priority_queue:in(Input, Priority, Queue) }.
...
%% 处理优先级队列中的每一条消息
%% From => {self(), Mref}
%% Req => get_state
process_msg({system, From, Req},
            GS2State = #gs2_state { parent = Parent, debug  = Debug }) ->
    %% gen_server puts Hib on the end as the 7th arg, but that version
    %% of the fun seems not to be documented so leaving out for now.
    sys:handle_system_msg(Req, From, Parent, ?MODULE, Debug, GS2State);
```

可以看到，针对 system 消息到处理，最终调用的是 `sys:handle_system_msg/6` 函数；

在 `sys.erl` 中

```erlang
%% Msg => get_state
%% From => {self(), Mref}
%% Module => gen_server2
%% Misc => GS2State
handle_system_msg(Msg, From, Parent, Module, Debug, Misc) ->
    handle_system_msg(running, Msg, From, Parent, Module, Debug, Misc, false).

handle_system_msg(Msg, From, Parent, Mod, Debug, Misc, Hib) ->
   handle_system_msg(running, Msg, From, Parent, Mod, Debug, Misc, Hib).

%% SysState => running | suspended
%% 此处 SysState 为 running
handle_system_msg(SysState, Msg, From, Parent, Mod, Debug, Misc, Hib) ->
    case do_cmd(SysState, Msg, Parent, Mod, Debug, Misc) of
        {suspended, Reply, NDebug, NMisc} -> %% 进程被挂起
            _ = gen:reply(From, Reply),
            suspend_loop(suspended, Parent, Mod, NDebug, NMisc, Hib);
        {running, Reply, NDebug, NMisc} -> %% 进程继续执行
            _ = gen:reply(From, Reply),    %% 针对 gen_server2 模块，此处 Reply 对应 GS2State
            Mod:system_continue(Parent, NDebug, NMisc)
    end.
...
%% 获取 Mod 模块中的 state 信息
%% SysState => running
%% Mod => gen_server2
%% Misc => GS2State
do_cmd(SysState, get_state, _Parent, Mod, Debug, Misc) ->
    {SysState, do_get_state(Mod, Misc), Debug, Misc};
...
%% Mod => gen_server2
%% Misc => GS2State
do_get_state(Mod, Misc) ->
    %% gen_server2 中没有定义 system_get_state 函数
    case erlang:function_exported(Mod, system_get_state, 1) of
        true ->
            try
                {ok, State} = Mod:system_get_state(Misc),
                State
            catch
                Cl:Exc ->
                    {error, {callback_failed,{Mod,system_get_state},{Cl,Exc}}}
            end;
        false ->
            Misc
    end.
```

最后看下，应答消息回复实现，位于 `gen.erl` 模块中

```erlang
%%
%% Send a reply to the client.
%%
%% 对应上述过程，可以推出
%% {To, Tag} => {self(), Mref}
%% Reply => GS2State
reply({To, Tag}, Reply) ->
    Msg = {Tag, Reply},
    %% 恢复应答给调用进程
    try To ! Msg catch _:_ -> Msg end.
```

之后，消息 *{Mref, GS2State}* 被发回给发出 system 消息的进程，即之前的

```erlang
    ...
    %% Label => system
    %% Request => get_state
    catch erlang:send(Process, {Label, {self(), Mref}, Request}, [noconnect]),
    receive
        {Mref, Reply} ->  %% 收到应答，返回给调用 sys:get_state/1,2 的进程
            erlang:demonitor(Mref, [flush]),
            {ok, Reply};
    ...
```

至此，我们已经清楚的知道，通过 `sys:get_state/1` 获取 rabbit_mgmt_db 进程（即 gen_server2 进程）的 state 信息全部过程；


## 详解 gen_server2 的 state 信息

对比 gen_server2.erl 的 state 定义，可以更好的理解获取信息的含义；

```erlang
%% State record
%%
-record(gs2_state, {
                        parent,         %% 当前 gen_server2 进程的父进程
                        name,           %% 当前 gen_server2 进程的注册名字，如 rabbit_mgmt_db
                        state,          %% 当前 gen_server2 进程保存的外部 state 记录，如 rabbit_mgmt_db 中的 state 值
                        mod,            %% 当前 gen_server2 进程对应的模块，如 rabbit_mgmt_db
                        time,           %% 
                        timeout_state,  %% backoff 策略配置
                        queue,          %% 优先级队列
                        debug,          %% 与 sys 模块中 Debug 结构相关的元组列表
                        prioritisers    %% {PCall, PCast, PInfo}
                    }).
```

这里再次给出实际调用结果作为对比

```erlang
(rabbit_2@sunfeideMacBook-Pro)1> sys:get_state(rabbit_mgmt_db).
{gs2_state,<0.390.0>,rabbit_mgmt_db,
           {state,[{channel_stats,462941},
                   {connection_stats,458844},
                   {consumers_by_channel,471135},
                   {consumers_by_queue,467038},
                   {node_node_stats,479329},
                   {node_stats,475232},
                   {queue_stats,454747}],
                  483426,487523,491620,#Ref<0.0.1.148153>,
                  {{queue_stats,{resource,<<"/">>,queue,<<"q5">>}},
                   messages_unacknowledged},
                  [{exchange,#Fun<rabbit_exchange.lookup.1>},
                   {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                  10000,#Ref<0.0.1.1044>,detailed},
           rabbit_mgmt_db,hibernate,
           {backoff,5,1000,10000,{23942,-18676,-7813}},
           {queue,[],[],0},
           [],
           {#Fun<gen_server2.9.133231575>,
            #Fun<gen_server2.10.133231575>,
            #Fun<gen_server2.8.133231575>}}
(rabbit_2@sunfeideMacBook-Pro)2>
```


## sys:get_status/1

由于相关代码几乎与 `sys:get_state/1` 相同，故只给出差异部分的代码；

在 `sys.erl` 中

```erlang
%% 获取 Mod 模块的状态信息
do_cmd(SysState, get_status, Parent, Mod, Debug, Misc) ->
    Res = get_status(SysState, Parent, Mod, Debug, Misc),
    {SysState, Res, Debug, Misc};
...
%% 针对 rabbit_mgmt_db.erl 模块的情况
%% Mod => gen_server2
%% Misc => GS2State
get_status(SysState, Parent, Mod, Debug, Misc) ->
    PDict = get(),
    FmtMisc =
        %% gen_server2 中定义了 format_status 函数
        case erlang:function_exported(Mod, format_status, 2) of
            true ->
                FmtArgs = [PDict, SysState, Parent, Debug, Misc],
                Mod:format_status(normal, FmtArgs);
            _ ->
                Misc
        end,
    {status, self(), {module, Mod}, [PDict, SysState, Parent, Debug, FmtMisc]}.
```

在 `gen_server2.erl` 中

```erlang
format_status(Opt, StatusData) ->
    [PDict, SysState, Parent, Debug,
     #gs2_state{name = Name, state = State, mod = Mod, queue = Queue}] =
        StatusData,
    NameTag = if is_pid(Name) ->
                      pid_to_list(Name);
                 is_atom(Name) ->
                      Name
              end,
    Header = lists:concat(["Status for generic server ", NameTag]),
    Log = sys:get_debug(log, Debug, []),
    %% 此处的 Mod 为基于 gen_server2 行为模式运行的进程名，如 rabbit_mgmt_db ，下同
    %%
    %% 由于 rabbit_mgmt_db 没有导出 format_status 函数，因此这里输出的是
    %% rabbit_mgmt_db 模块中的 state 记录内容
    Specfic = callback(Mod, format_status, [Opt, [PDict, State]],
                       fun () -> [{data, [{"State", State}]}] end),
    %% 输出优先级队列相关内容
    %% rabbit_mgmt_db 中导出了 format_message_queue 函数
    Messages = callback(Mod, format_message_queue, [Opt, Queue],
                        fun () -> priority_queue:to_list(Queue) end),
    [{header, Header},
     {data, [{"Status", SysState},
             {"Parent", Parent},
             {"Logged events", Log},
             {"Queued messages", Messages}]} |
     Specfic].

%% 调用指定 Mod 中的 FunName/Args
callback(Mod, FunName, Args, DefaultThunk) ->
    case erlang:function_exported(Mod, FunName, length(Args)) of
		%% 函数调用实际发生的地方
        true  -> case catch apply(Mod, FunName, Args) of
                     {'EXIT', _} -> DefaultThunk();
                     Success     -> Success
                 end;
        false -> DefaultThunk()
    end.
```

在 `rabbit_mgmt_db.erl` 中

```erlang
format_message_queue(Opt, MQ) -> rabbit_misc:format_message_queue(Opt, MQ).
```

在 `rabbit_misc.erl` 中

```erlang
format_message_queue(_Opt, MQ) ->
	%% 获取优先级队列中消息的总数
    Len = priority_queue:len(MQ),
    {Len,
     case Len > 100 of
         %% 100 条之内，直接完整输出消息内容
         false -> priority_queue:to_list(MQ);
         %% 超过 100 条，则以统计数据形式输出消息内容
         true  -> {summary,
                   orddict:to_list(
                     lists:foldl(
                       fun ({P, V}, Counts) ->
                               orddict:update_counter(
                                 {P, format_message_queue_entry(V)}, 1, Counts)
                       end, orddict:new(), priority_queue:to_list(MQ)))}
     end}.
```

## status 信息

需要知道的是，status 信息不是以 record 的形式定义的，而是在 sys:get_status/5 接口中直接指定的

```erlang
get_status(SysState, Parent, Mod, Debug, Misc) ->
	...
    {status, self(), {module, Mod}, [PDict, SysState, Parent, Debug, FmtMisc]}
```

这里再次给出实际调用结果作为对比

```erlang
(rabbit_1@sunfeideMacBook-Pro)8> sys:get_status(whereis(rabbit_mgmt_db)).
{status,<0.434.0>,
        {module,gen_server2},
        [[{'$initial_call',{rabbit_mgmt_db,init,1}},
          {'$ancestors',[<0.429.0>,rabbit_mgmt_sup,
                         rabbit_mgmt_sup_sup,<0.407.0>]},
          {last_queue_length,0}],
         running,<0.429.0>,[],
         [{header,"Status for generic server rabbit_mgmt_db"},
          {data,[{"Status",running},
                 {"Parent",<0.429.0>},
                 {"Logged events",[]},
                 {"Queued messages",{0,[]}}]},
          {data,[{"State",
                  {state,[{channel_stats,598109},
                          {connection_stats,594010},
                          {consumers_by_channel,606303},
                          {consumers_by_queue,602206},
                          {node_node_stats,614497},
                          {node_stats,610400},
                          {queue_stats,589915}],
                         618594,622691,626788,#Ref<0.0.1.63852>,
                         {{queue_stats,{resource,<<...>>,...}},
                          messages_unacknowledged},
                         [{exchange,#Fun<rabbit_exchange.lookup.1>},
                          {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                         10000,#Ref<0.0.2.1293>,detailed}}]}]]}
```


----------


# 信息获取方法

由于本文的目标是获取 rabbit_mgmt_db 进程的状态信息，因此可以采用两种方式进行操作：
- 先确定 rabbit_mgmt_db 进程（即统计数据库的位置）运行于哪个 RabbitMQ 节点上，再调用访问本地节点的接口进行信息获取；
- 通过全局可用接口从 cluster 中任意节点上进行信息获取；

使用的命令来自 `sys.erl` 模块（可以在 erlang console 上直接使用），可以使用如下两种全局调用形式：

```erlang
sys:get_state(global:whereis_name(rabbit_mgmt_db)).   %% 基于 pid
sys:get_state({global, rabbit_mgmt_db}).              %% 基于进程注册名
```

对应的底层接口为（可以在 erlang console 上直接使用）：

```erlang
gen:call(global:whereis_name(rabbit_mgmt_db),system,get_state).
gen:call({global,rabbit_mgmt_db},system,get_state).
```

对应到 RabbitMQ 上可用的 shell 命令如下

```shell
sudo rabbitmqctl eval "sys:get_state(global:whereis_name(rabbit_mgmt_db))."
sudo rabbitmqctl eval "sys:get_state({global, rabbit_mgmt_db})."
```

> ⚠️ 由于 `sys:get_state/1,2` 在 OTP R16 的各个版本中无法使用，建议改为使用 `sys:get_status/1,2`;


# 实验

## 向 rabbit_mgmt_db 进程压入大量消息

### 初始状态


```erlang
Erlang/OTP 19 [erts-8.0.2] [source] [64-bit] [smp:4:4] [async-threads:10] [hipe] [kernel-poll:false] [dtrace]

Eshell V8.0.2  (abort with ^G)
(rabbit_2@sunfeideMacBook-Pro)1>
(rabbit_2@sunfeideMacBook-Pro)1>
(rabbit_2@sunfeideMacBook-Pro)1> whereis(rabbit_mgmt_db).
<0.394.0>
(rabbit_2@sunfeideMacBook-Pro)2>
(rabbit_2@sunfeideMacBook-Pro)2> erlang:process_info(whereis(rabbit_mgmt_db)).
[{registered_name,rabbit_mgmt_db},
 {current_function,{erlang,hibernate,3}},
 {initial_call,{proc_lib,init_p,5}},
 {status,waiting},
 {message_queue_len,0},            %% 观察点
 {messages,[]},                    %% 观察点
 {links,[<0.389.0>]},
 {dictionary,[{'$initial_call',{rabbit_mgmt_db,init,1}},
              {'$ancestors',[<0.389.0>,rabbit_mgmt_sup,
                             rabbit_mgmt_sup_sup,<0.367.0>]},
              {last_queue_length,0}]},      %% 观察点
 {trap_exit,false},
 {error_handler,error_handler},
 {priority,high},
 {group_leader,<0.366.0>},
 {total_heap_size,172},
 {heap_size,172},
 {stack_size,0},
 {reductions,2340944},                     %% 观察点
 {garbage_collection,[{max_heap_size,#{error_logger => true,
                                       kill => true,
                                       size => 0}},
                      {min_bin_vheap_size,46422},
                      {min_heap_size,233},
                      {fullsweep_after,65535},
                      {minor_gcs,50}]},
 {suspending,[]}]
(rabbit_2@sunfeideMacBook-Pro)3>
(rabbit_2@sunfeideMacBook-Pro)3>
(rabbit_2@sunfeideMacBook-Pro)3> sys:get_state({global, rabbit_mgmt_db}).
{gs2_state,<0.389.0>,rabbit_mgmt_db,
           {state,[{channel_stats,471133},
                   {connection_stats,467036},
                   {consumers_by_channel,479327},
                   {consumers_by_queue,475230},
                   {node_node_stats,487521},
                   {node_stats,483424},
                   {queue_stats,462939}],
                  491618,495715,499812,#Ref<0.0.2.6308>,
                  {{queue_stats,{resource,<<"/">>,queue,<<"q9">>}},messages},
                  [{exchange,#Fun<rabbit_exchange.lookup.1>},
                   {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                  10000,#Ref<0.0.3.14>,detailed},
           rabbit_mgmt_db,hibernate,
           {backoff,11,1000,10000,{26480,-6868,-4853}},
           {queue,[],[],0},                %% 观察点
           [],
           {#Fun<gen_server2.9.133231575>,
            #Fun<gen_server2.10.133231575>,
            #Fun<gen_server2.8.133231575>}}
(rabbit_2@sunfeideMacBook-Pro)4>
(rabbit_2@sunfeideMacBook-Pro)4>
(rabbit_2@sunfeideMacBook-Pro)4> sys:get_status({global, rabbit_mgmt_db}).
{status,<0.394.0>,
        {module,gen_server2},
        [[{'$initial_call',{rabbit_mgmt_db,init,1}},
          {'$ancestors',[<0.389.0>,rabbit_mgmt_sup,
                         rabbit_mgmt_sup_sup,<0.367.0>]},
          {last_queue_length,0}],            %% 观察点
         running,<0.389.0>,[],
         [{header,"Status for generic server rabbit_mgmt_db"},
          {data,[{"Status",running},
                 {"Parent",<0.389.0>},
                 {"Logged events",[]},
                 {"Queued messages",{0,[]}}]},           %% 观察点
          {data,[{"State",
                  {state,[{channel_stats,471133},
                          {connection_stats,467036},
                          {consumers_by_channel,479327},
                          {consumers_by_queue,475230},
                          {node_node_stats,487521},
                          {node_stats,483424},
                          {queue_stats,462939}],
                         491618,495715,499812,#Ref<0.0.2.6495>,
                         {{node_stats,'rabbit_2@sunfeideMacBook-Pro'},io_read_bytes},
                         [{exchange,#Fun<rabbit_exchange.lookup.1>},
                          {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                         10000,#Ref<0.0.3.14>,detailed}}]}]]}
(rabbit_2@sunfeideMacBook-Pro)5>
```

放开 statistics 开关，便于观察 messages 的 in 和 out 情况；

```erlang
(rabbit_2@sunfeideMacBook-Pro)5> sys:statistics({global, rabbit_mgmt_db},true).
ok
(rabbit_2@sunfeideMacBook-Pro)6>
(rabbit_2@sunfeideMacBook-Pro)6>
(rabbit_2@sunfeideMacBook-Pro)6> sys:get_status({global, rabbit_mgmt_db}).
{status,<0.394.0>,
        {module,gen_server2},
        [[{'$initial_call',{rabbit_mgmt_db,init,1}},
          {'$ancestors',[<0.389.0>,rabbit_mgmt_sup,
                         rabbit_mgmt_sup_sup,<0.367.0>]},
          {last_queue_length,0}],
         running,<0.389.0>,
         [{statistics,{{{2016,9,19},{16,25,17}},
                       {reductions,9513561},
                       20,0}}],     %% 两个值分别对应 messages_in 和 messages_out
         [{header,"Status for generic server rabbit_mgmt_db"},
          {data,[{"Status",running},
                 {"Parent",<0.389.0>},
                 {"Logged events",[]},
                 {"Queued messages",{0,[]}}]},
          {data,[{"State",
                  {state,[{channel_stats,471133},
                          {connection_stats,467036},
                          {consumers_by_channel,479327},
                          {consumers_by_queue,475230},
                          {node_node_stats,487521},
                          {node_stats,483424},
                          {queue_stats,462939}],
                         491618,495715,499812,#Ref<0.0.1.6794>,
                         {{node_stats,'rabbit_2@sunfeideMacBook-Pro'},io_read_bytes},
                         [{exchange,#Fun<rabbit_exchange.lookup.1>},
                          {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                         10000,#Ref<0.0.3.14>,detailed}}]}]]}
(rabbit_2@sunfeideMacBook-Pro)7>
```

RabbitMQ Web 控制台状态信息

![初始状态 web](https://github.com/moooofly/ImageCache/blob/master/Pictures/%E5%88%9D%E5%A7%8B%E7%8A%B6%E6%80%81%20web.png "初始状态 web")

通过 entop 查看的状态信息（基于 reduction 数值进行排序）

![初始状态 entop](https://github.com/moooofly/ImageCache/blob/master/Pictures/%E5%88%9D%E5%A7%8B%E7%8A%B6%E6%80%81%20entop.png "初始状态 entop")



### 压入 1w 条消息时的状态

在另一个 erlang console 中，向 rabbit_mgmt_db 进程压入 10,000 条消息；

```erlang
(rabbit_2@sunfeideMacBook-Pro)1> P=whereis(rabbit_mgmt_db).
<0.394.0>
(rabbit_2@sunfeideMacBook-Pro)2>
(rabbit_2@sunfeideMacBook-Pro)2>
(rabbit_2@sunfeideMacBook-Pro)2> [P!{test,N}||N<-lists:seq(1,10000)].
[{test,1},
 {test,2},
 {test,3},
 {test,4},
 {test,5},
 {test,6},
 {test,7},
 {test,8},
 {test,9},
 {test,10},
 {test,11},
 {test,12},
 {test,13},
 {test,14},
 {test,15},
 {test,16},
 {test,17},
 {test,18},
 {test,19},
 {test,20},
 {test,21},
 {test,22},
 {test,23},
 {test,24},
 {test,25},
 {test,26},
 {test,27},
 {test,...},
 {...}|...]
(rabbit_2@sunfeideMacBook-Pro)3>
```

可以看到 messages_in 对应的数值增加了 10,000（由于在默认配置下，RabbitMQ 同样会定时上报统计信息，所以该值在正常情况下也会缓慢增长）；

```erlang
(rabbit_2@sunfeideMacBook-Pro)13> sys:get_status({global, rabbit_mgmt_db}).
{status,<0.394.0>,
        {module,gen_server2},
        [[{'$initial_call',{rabbit_mgmt_db,init,1}},
          {'$ancestors',[<0.389.0>,rabbit_mgmt_sup,
                         rabbit_mgmt_sup_sup,<0.367.0>]},
          {last_queue_length,0}],
         running,<0.389.0>,
         [{statistics,{{{2016,9,19},{16,25,17}},
                       {reductions,9513561},
                       11114,0}}],        %% 已压入 1w+ 消息
         [{header,"Status for generic server rabbit_mgmt_db"},
          {data,[{"Status",running},
                 {"Parent",<0.389.0>},
                 {"Logged events",[]},
                 {"Queued messages",{0,[]}}]},
          {data,[{"State",
                  {state,[{channel_stats,471133},
                          {connection_stats,467036},
                          {consumers_by_channel,479327},
                          {consumers_by_queue,475230},
                          {node_node_stats,487521},
                          {node_stats,483424},
                          {queue_stats,462939}],
                         491618,495715,499812,#Ref<0.0.2.14474>,
                         {{queue_stats,{resource,<<...>>,...}},
                          messages_unacknowledged},
                         [{exchange,#Fun<rabbit_exchange.lookup.1>},
                          {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                         10000,#Ref<0.0.3.14>,detailed}}]}]]}
(rabbit_2@sunfeideMacBook-Pro)14>
```

![压入 10000 消息后 web](https://github.com/moooofly/ImageCache/blob/master/Pictures/%E5%8E%8B%E5%85%A5%2010000%20%E6%B6%88%E6%81%AF%E5%90%8E%20web.png "压入 10000 消息后 web")

![压入 10000 消息后 entop](https://github.com/moooofly/ImageCache/blob/master/Pictures/%E5%8E%8B%E5%85%A5%2010000%20%E6%B6%88%E6%81%AF%E5%90%8E%20entop.png "压入 10000 消息后 entop")


### 压入 100w 条消息时的状态

在另一个 erlang console 中，再向 rabbit_mgmt_db 进程压入 1,000,000 条消息；

```erlang
(rabbit_2@sunfeideMacBook-Pro)15> sys:get_status({global, rabbit_mgmt_db}).
{status,<0.394.0>,
        {module,gen_server2},
        [[{'$initial_call',{rabbit_mgmt_db,init,1}},
          {'$ancestors',[<0.389.0>,rabbit_mgmt_sup,
                         rabbit_mgmt_sup_sup,<0.367.0>]},
          {last_queue_length,18}],          %% 调用该命令的瞬间，进程邮箱中尚仍在的消息数量
         running,<0.389.0>,
         [{statistics,{{{2016,9,19},{16,25,17}},
                       {reductions,9513561},
                       427079,0}}],             %% 已压入 40w+ 消息
         [{header,"Status for generic server rabbit_mgmt_db"},
          {data,[{"Status",running},
                 {"Parent",<0.389.0>},
                 {"Logged events",[]},
                 {"Queued messages",
                  {2250,{summary,[{{0,{test,'_'}},2250}]}}}]},  %% 优先级队列中积压了 2000+ 消息，此时会触发 management 插件在 web 页面上的黄色告警（>1000）
          {data,[{"State",
                  {state,[{channel_stats,471133},
                          {connection_stats,467036},
                          {consumers_by_channel,479327},
                          {consumers_by_queue,475230},
                          {node_node_stats,487521},
                          {node_stats,483424},
                          {queue_stats,462939}],
                         491618,495715,499812,#Ref<0.0.1.18409>,
                         {{node_node_stats,{'rabbit_3@sunfeideMacBook-Pro','rabbit_1@sunfeideMacBook-Pro'}},
                          send_bytes},
                         [{exchange,#Fun<rabbit_exchange.lookup.1>},
                          {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                         10000,#Ref<0.0.3.14>,detailed}}]}]]}
(rabbit_2@sunfeideMacBook-Pro)16>
(rabbit_2@sunfeideMacBook-Pro)16>
(rabbit_2@sunfeideMacBook-Pro)16> sys:get_status({global, rabbit_mgmt_db}).
{status,<0.394.0>,
        {module,gen_server2},
        [[{'$initial_call',{rabbit_mgmt_db,init,1}},
          {'$ancestors',[<0.389.0>,rabbit_mgmt_sup,
                         rabbit_mgmt_sup_sup,<0.367.0>]},
          {last_queue_length,18}],
         running,<0.389.0>,
         [{statistics,{{{2016,9,19},{16,25,17}},
                       {reductions,9513561},
                       856133,0}}],            %% 已压入 80w+ 消息
         [{header,"Status for generic server rabbit_mgmt_db"},
          {data,[{"Status",running},
                 {"Parent",<0.389.0>},
                 {"Logged events",[]},
                 {"Queued messages",{0,[]}}]},
          {data,[{"State",
                  {state,[{channel_stats,471133},
                          {connection_stats,467036},
                          {consumers_by_channel,479327},
                          {consumers_by_queue,475230},
                          {node_node_stats,487521},
                          {node_stats,483424},
                          {queue_stats,462939}],
                         491618,495715,499812,#Ref<0.0.1.18409>,
                         {{node_node_stats,{'rabbit_3@sunfeideMacBook-Pro','rabbit_1@sunfeideMacBook-Pro'}},
                          send_bytes},
                         [{exchange,#Fun<rabbit_exchange.lookup.1>},
                          {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                         10000,#Ref<0.0.3.14>,detailed}}]}]]}
(rabbit_2@sunfeideMacBook-Pro)17>
(rabbit_2@sunfeideMacBook-Pro)17>
(rabbit_2@sunfeideMacBook-Pro)17> sys:get_status({global, rabbit_mgmt_db}).
{status,<0.394.0>,
        {module,gen_server2},
        [[{'$initial_call',{rabbit_mgmt_db,init,1}},
          {'$ancestors',[<0.389.0>,rabbit_mgmt_sup,
                         rabbit_mgmt_sup_sup,<0.367.0>]},
          {last_queue_length,18}],
         running,<0.389.0>,
         [{statistics,{{{2016,9,19},{16,25,17}},
                       {reductions,9513561},
                       1012948,0}}],              %% 已压入 100w+ 消息
         [{header,"Status for generic server rabbit_mgmt_db"},
          {data,[{"Status",running},
                 {"Parent",<0.389.0>},
                 {"Logged events",[]},
                 {"Queued messages",{0,[]}}]},
          {data,[{"State",
                  {state,[{channel_stats,471133},
                          {connection_stats,467036},
                          {consumers_by_channel,479327},
                          {consumers_by_queue,475230},
                          {node_node_stats,487521},
                          {node_stats,483424},
                          {queue_stats,462939}],
                         491618,495715,499812,#Ref<0.0.1.18409>,
                         {{node_node_stats,{'rabbit_3@sunfeideMacBook-Pro','rabbit_1@sunfeideMacBook-Pro'}},
                          send_bytes},
                         [{exchange,#Fun<rabbit_exchange.lookup.1>},
                          {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                         10000,#Ref<0.0.3.14>,detailed}}]}]]}
(rabbit_2@sunfeideMacBook-Pro)18>
```

![压入 100w 消息后 web](https://github.com/moooofly/ImageCache/blob/master/Pictures/%E5%8E%8B%E5%85%A5%20100w%20%E6%B6%88%E6%81%AF%E5%90%8E%20web.png "压入 100w 消息后 web")

![压入 100w 消息后 entop](https://github.com/moooofly/ImageCache/blob/master/Pictures/%E5%8E%8B%E5%85%A5%20100w%20%E6%B6%88%E6%81%AF%E5%90%8E%20entop.png "压入 100w 消息后 entop")


### 压入 1000w 条消息时的状态

再向 rabbit_mgmt_db 进程压入 10,000,000 条消息；


```erlang
(rabbit_2@sunfeideMacBook-Pro)23> sys:get_status({global, rabbit_mgmt_db}).
{status,<0.394.0>,
        {module,gen_server2},
        [[{'$initial_call',{rabbit_mgmt_db,init,1}},
          {'$ancestors',[<0.389.0>,rabbit_mgmt_sup,
                         rabbit_mgmt_sup_sup,<0.367.0>]},
          {last_queue_length,0}],
         running,<0.389.0>,
         [{statistics,{{{2016,9,19},{16,25,17}},
                       {reductions,9513561},
                       1646573,0}}],         %% 已压入 160w+ 消息
         [{header,"Status for generic server rabbit_mgmt_db"},
          {data,[{"Status",running},
                 {"Parent",<0.389.0>},
                 {"Logged events",[]},
                 {"Queued messages",{0,[]}}]},
          {data,[{"State",
                  {state,[{channel_stats,471133},
                          {connection_stats,467036},
                          {consumers_by_channel,479327},
                          {consumers_by_queue,475230},
                          {node_node_stats,487521},
                          {node_stats,483424},
                          {queue_stats,462939}],
                         491618,495715,499812,#Ref<0.0.3.19520>,
                         {{queue_stats,{resource,<<...>>,...}},
                          messages_unacknowledged},
                         [{exchange,#Fun<rabbit_exchange.lookup.1>},
                          {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                         10000,#Ref<0.0.3.14>,detailed}}]}]]}
(rabbit_2@sunfeideMacBook-Pro)24>
...
(rabbit_2@sunfeideMacBook-Pro)25>
(rabbit_2@sunfeideMacBook-Pro)25> sys:get_status({global, rabbit_mgmt_db}).
{status,<0.394.0>,
        {module,gen_server2},
        [[{'$initial_call',{rabbit_mgmt_db,init,1}},
          {'$ancestors',[<0.389.0>,rabbit_mgmt_sup,
                         rabbit_mgmt_sup_sup,<0.367.0>]},
          {last_queue_length,10}],
         running,<0.389.0>,
         [{statistics,{{{2016,9,19},{16,25,17}},
                       {reductions,9513561},
                       2741179,0}}],        %% 已压入 270w+ 消息
         [{header,"Status for generic server rabbit_mgmt_db"},
          {data,[{"Status",running},
                 {"Parent",<0.389.0>},
                 {"Logged events",[]},
                 {"Queued messages",        %% 在该命令执行瞬间，看到了优先级队列中的内容
                  {21,
                   [{0,{test,1724065}},
                    {0,{test,1724066}},
                    {0,{test,1724067}},
                    {0,{test,1724068}},
                    {0,{test,...}},
                    {0,{...}},
                    {0,...},
                    {...}|...]}}]},
          {data,[{"State",
                  {state,[{channel_stats,471133},
                          {connection_stats,467036},
                          {consumers_by_channel,479327},
                          {consumers_by_queue,475230},
                          {node_node_stats,487521},
                          {node_stats,483424},
                          {queue_stats,462939}],
                         491618,495715,499812,#Ref<0.0.3.19655>,
                         {{node_node_stats,{'rabbit_3@sunfeideMacBook-Pro','rabbit_1@sunfeideMacBook-Pro'}},
                          send_bytes},
                         [{exchange,#Fun<rabbit_exchange.lookup.1>},
                          {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                         10000,#Ref<0.0.3.14>,detailed}}]}]]}
(rabbit_2@sunfeideMacBook-Pro)26>

(rabbit_2@sunfeideMacBook-Pro)27>
(rabbit_2@sunfeideMacBook-Pro)27> sys:get_status({global, rabbit_mgmt_db}).
{status,<0.394.0>,
        {module,gen_server2},
        [[{'$initial_call',{rabbit_mgmt_db,init,1}},
          {'$ancestors',[<0.389.0>,rabbit_mgmt_sup,
                         rabbit_mgmt_sup_sup,<0.367.0>]},
          {last_queue_length,10}],
         running,<0.389.0>,
         [{statistics,{{{2016,9,19},{16,25,17}},
                       {reductions,9513561},
                       3993400,0}}],          %% 已压入 390w+ 消息
         [{header,"Status for generic server rabbit_mgmt_db"},
          {data,[{"Status",running},
                 {"Parent",<0.389.0>},
                 {"Logged events",[]},
                 {"Queued messages",
                  {44,
                   [{0,{test,2976273}},
                    {0,{test,2976274}},
                    {0,{test,2976275}},
                    {0,{test,2976276}},
                    {0,{test,...}},
                    {0,{...}},
                    {0,...},
                    {...}|...]}}]},
          {data,[{"State",
                  {state,[{channel_stats,471133},
                          {connection_stats,467036},
                          {consumers_by_channel,479327},
                          {consumers_by_queue,475230},
                          {node_node_stats,487521},
                          {node_stats,483424},
                          {queue_stats,462939}],
                         491618,495715,499812,#Ref<0.0.1.34513>,
                         {{node_stats,'rabbit_2@sunfeideMacBook-Pro'},
                          io_sync_avg_time},
                         [{exchange,#Fun<rabbit_exchange.lookup.1>},
                          {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                         10000,#Ref<0.0.3.14>,detailed}}]}]]}
(rabbit_2@sunfeideMacBook-Pro)28>
(rabbit_2@sunfeideMacBook-Pro)28>

(rabbit_2@sunfeideMacBook-Pro)31>
(rabbit_2@sunfeideMacBook-Pro)31> sys:get_status({global, rabbit_mgmt_db}).
{status,<0.394.0>,
        {module,gen_server2},
        [[{'$initial_call',{rabbit_mgmt_db,init,1}},
          {'$ancestors',[<0.389.0>,rabbit_mgmt_sup,
                         rabbit_mgmt_sup_sup,<0.367.0>]},
          {last_queue_length,3}],
         running,<0.389.0>,
         [{statistics,{{{2016,9,19},{16,25,17}},
                       {reductions,9513561},
                       9178618,0}}],         %% 已压入 910w+ 消息
         [{header,"Status for generic server rabbit_mgmt_db"},
          {data,[{"Status",running},
                 {"Parent",<0.389.0>},
                 {"Logged events",[]},
                 {"Queued messages",
                  {11,
                   [{0,{test,8161436}},
                    {0,{test,8161437}},
                    {0,{test,8161438}},
                    {0,{test,8161439}},
                    {0,{test,...}},
                    {0,{...}},
                    {0,...},
                    {...}|...]}}]},
          {data,[{"State",
                  {state,[{channel_stats,471133},
                          {connection_stats,467036},
                          {consumers_by_channel,479327},
                          {consumers_by_queue,475230},
                          {node_node_stats,487521},
                          {node_stats,483424},
                          {queue_stats,462939}],
                         491618,495715,499812,#Ref<0.0.2.31881>,
                         {{node_node_stats,{'rabbit_3@sunfeideMacBook-Pro','rabbit_1@sunfeideMacBook-Pro'}},
                          send_bytes},
                         [{exchange,#Fun<rabbit_exchange.lookup.1>},
                          {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                         10000,#Ref<0.0.3.14>,detailed}}]}]]}
(rabbit_2@sunfeideMacBook-Pro)32>
(rabbit_2@sunfeideMacBook-Pro)32>

(rabbit_2@sunfeideMacBook-Pro)34>
(rabbit_2@sunfeideMacBook-Pro)34>
(rabbit_2@sunfeideMacBook-Pro)34> sys:get_status({global, rabbit_mgmt_db}).
{status,<0.394.0>,
        {module,gen_server2},
        [[{'$initial_call',{rabbit_mgmt_db,init,1}},
          {'$ancestors',[<0.389.0>,rabbit_mgmt_sup,
                         rabbit_mgmt_sup_sup,<0.367.0>]},
          {last_queue_length,1}],
         running,<0.389.0>,
         [{statistics,{{{2016,9,19},{16,25,17}},
                       {reductions,9513561},
                       11017212,0}}],       %% 已压入 1100w+ 消息
         [{header,"Status for generic server rabbit_mgmt_db"},
          {data,[{"Status",running},
                 {"Parent",<0.389.0>},
                 {"Logged events",[]},
                 {"Queued messages",{0,[]}}]},
          {data,[{"State",
                  {state,[{channel_stats,471133},
                          {connection_stats,467036},
                          {consumers_by_channel,479327},
                          {consumers_by_queue,475230},
                          {node_node_stats,487521},
                          {node_stats,483424},
                          {queue_stats,462939}],
                         491618,495715,499812,#Ref<0.0.1.34888>,
                         {{node_stats,'rabbit_2@sunfeideMacBook-Pro'},
                          io_sync_avg_time},
                         [{exchange,#Fun<rabbit_exchange.lookup.1>},
                          {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                         10000,#Ref<0.0.3.14>,detailed}}]}]]}
(rabbit_2@sunfeideMacBook-Pro)35>
```

![压入 1000w 消息后 web](https://github.com/moooofly/ImageCache/blob/master/Pictures/%E5%8E%8B%E5%85%A5%201000w%20%E6%B6%88%E6%81%AF%E5%90%8E%20web.png "压入 1000w 消息后 web")

![压入 1000w 消息后 entop](https://github.com/moooofly/ImageCache/blob/master/Pictures/%E5%8E%8B%E5%85%A5%201000w%20%E6%B6%88%E6%81%AF%E5%90%8E%20entop.png "压入 1000w 消息后 entop")

### 压入 10000w 条消息时的状态

再向 rabbit_mgmt_db 进程压入 100,000,000 条消息；


```erlang
(rabbit_2@sunfeideMacBook-Pro)63> sys:get_status({global, rabbit_mgmt_db}).
{status,<0.394.0>,
        {module,gen_server2},
        [[{'$initial_call',{rabbit_mgmt_db,init,1}},
          {'$ancestors',[<0.389.0>,rabbit_mgmt_sup,
                         rabbit_mgmt_sup_sup,<0.367.0>]},
          {last_queue_length,1}],
         running,<0.389.0>,
         [{statistics,{{{2016,9,19},{16,25,17}},
                       {reductions,9513561},
                       29528363,0}}],     %% 已压入 2900w+ 消息
         [{header,"Status for generic server rabbit_mgmt_db"},
          {data,[{"Status",running},
                 {"Parent",<0.389.0>},
                 {"Logged events",[]},
                 {"Queued messages",
                  {16,
                   [{0,{test,18507365}},
                    {0,{test,18507366}},
                    {0,{test,18507367}},
                    {0,{test,18507368}},
                    {0,{test,...}},
                    {0,{...}},
                    {0,...},
                    {...}|...]}}]},
          {data,[{"State",
                  {state,[{channel_stats,471133},
                          {connection_stats,467036},
                          {consumers_by_channel,479327},
                          {consumers_by_queue,475230},
                          {node_node_stats,487521},
                          {node_stats,483424},
                          {queue_stats,462939}],
                         491618,495715,499812,#Ref<0.0.3.25303>,
                         {{queue_stats,{resource,<<...>>,...}},
                          messages_unacknowledged},
                         [{exchange,#Fun<rabbit_exchange.lookup.1>},
                          {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                         10000,#Ref<0.0.3.14>,detailed}}]}]]}
(rabbit_2@sunfeideMacBook-Pro)64>

(rabbit_2@sunfeideMacBook-Pro)68> sys:get_status({global, rabbit_mgmt_db}).
{status,<0.394.0>,
        {module,gen_server2},
        [[{'$initial_call',{rabbit_mgmt_db,init,1}},
          {'$ancestors',[<0.389.0>,rabbit_mgmt_sup,
                         rabbit_mgmt_sup_sup,<0.367.0>]},
          {last_queue_length,7}],
         running,<0.389.0>,
         [{statistics,{{{2016,9,19},{16,25,17}},
                       {reductions,9513561},
                       33249279,0}}],       %% 已压入 3300w+ 消息
         [{header,"Status for generic server rabbit_mgmt_db"},
          {data,[{"Status",running},
                 {"Parent",<0.389.0>},
                 {"Logged events",[]},
                 {"Queued messages",
                  {271511,{summary,[{{0,{test,'_'}},271511}]}}}]},  %% 积压 27w+ 消息在优先级队列中，触发 web 页面黄色告警
          {data,[{"State",
                  {state,[{channel_stats,471133},
                          {connection_stats,467036},
                          {consumers_by_channel,479327},
                          {consumers_by_queue,475230},
                          {node_node_stats,487521},
                          {node_stats,483424},
                          {queue_stats,462939}],
                         491618,495715,499812,#Ref<0.0.3.25390>,
                         {{node_stats,'rabbit_2@sunfeideMacBook-Pro'},
                          io_sync_avg_time},
                         [{exchange,#Fun<rabbit_exchange.lookup.1>},
                          {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                         10000,#Ref<0.0.3.14>,detailed}}]}]]}
(rabbit_2@sunfeideMacBook-Pro)69>
(rabbit_2@sunfeideMacBook-Pro)69>
(rabbit_2@sunfeideMacBook-Pro)69> sys:get_status({global, rabbit_mgmt_db}).
{status,<0.394.0>,
        {module,gen_server2},
        [[{'$initial_call',{rabbit_mgmt_db,init,1}},
          {'$ancestors',[<0.389.0>,rabbit_mgmt_sup,
                         rabbit_mgmt_sup_sup,<0.367.0>]},
          {last_queue_length,0}],
         running,<0.389.0>,
         [{statistics,{{{2016,9,19},{16,25,17}},
                       {reductions,9513561},
                       37731650,0}}],
         [{header,"Status for generic server rabbit_mgmt_db"},
          {data,[{"Status",running},
                 {"Parent",<0.389.0>},
                 {"Logged events",[]},
                 {"Queued messages",
                  {8971,{summary,[{{0,{test,'_'}},8971}]}}}]},
          {data,[{"State",
                  {state,[{channel_stats,471133},
                          {connection_stats,467036},
                          {consumers_by_channel,479327},
                          {consumers_by_queue,475230},
                          {node_node_stats,487521},
                          {node_stats,483424},
                          {queue_stats,462939}],
                         491618,495715,499812,#Ref<0.0.2.45722>,
                         {{node_node_stats,{'rabbit_3@sunfeideMacBook-Pro','rabbit_1@sunfeideMacBook-Pro'}},
                          send_bytes},
                         [{exchange,#Fun<rabbit_exchange.lookup.1>},
                          {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                         10000,#Ref<0.0.3.14>,detailed}}]}]]}
(rabbit_2@sunfeideMacBook-Pro)70>

(rabbit_2@sunfeideMacBook-Pro)76> sys:get_status({global, rabbit_mgmt_db}).
{status,<0.394.0>,
        {module,gen_server2},
        [[{'$initial_call',{rabbit_mgmt_db,init,1}},
          {'$ancestors',[<0.389.0>,rabbit_mgmt_sup,
                         rabbit_mgmt_sup_sup,<0.367.0>]},
          {last_queue_length,0}],
         running,<0.389.0>,
         [{statistics,{{{2016,9,19},{16,25,17}},
                       {reductions,9513561},
                       43568879,0}}],          %% 已压入 4300w+ 消息
         [{header,"Status for generic server rabbit_mgmt_db"},
          {data,[{"Status",running},
                 {"Parent",<0.389.0>},
                 {"Logged events",[]},
                 {"Queued messages",
                  {69713,{summary,[{{0,{test,'_'}},69713}]}}}]},
          {data,[{"State",
                  {state,[{channel_stats,471133},
                          {connection_stats,467036},
                          {consumers_by_channel,479327},
                          {consumers_by_queue,475230},
                          {node_node_stats,487521},
                          {node_stats,483424},
                          {queue_stats,462939}],
                         491618,495715,499812,#Ref<0.0.2.45954>,
                         {{node_stats,'rabbit_2@sunfeideMacBook-Pro'},
                          io_sync_avg_time},
                         [{exchange,#Fun<rabbit_exchange.lookup.1>},
                          {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                         10000,#Ref<0.0.3.14>,detailed}}]}]]}
(rabbit_2@sunfeideMacBook-Pro)77>

(rabbit_2@sunfeideMacBook-Pro)97> sys:get_status({global, rabbit_mgmt_db}).
{status,<0.394.0>,
        {module,gen_server2},
        [[{'$initial_call',{rabbit_mgmt_db,init,1}},
          {'$ancestors',[<0.389.0>,rabbit_mgmt_sup,
                         rabbit_mgmt_sup_sup,<0.367.0>]},
          {last_queue_length,0}],
         running,<0.389.0>,
         [{statistics,{{{2016,9,19},{16,25,17}},
                       {reductions,9513561},
                       111022237,0}}],       %% 已压入 11100w+ 消息
         [{header,"Status for generic server rabbit_mgmt_db"},
          {data,[{"Status",running},
                 {"Parent",<0.389.0>},
                 {"Logged events",[]},
                 {"Queued messages",{0,[]}}]},
          {data,[{"State",
                  {state,[{channel_stats,471133},
                          {connection_stats,467036},
                          {consumers_by_channel,479327},
                          {consumers_by_queue,475230},
                          {node_node_stats,487521},
                          {node_stats,483424},
                          {queue_stats,462939}],
                         491618,495715,499812,#Ref<0.0.2.47860>,
                         {{queue_stats,{resource,<<...>>,...}},
                          messages_unacknowledged},
                         [{exchange,#Fun<rabbit_exchange.lookup.1>},
                          {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                         10000,#Ref<0.0.3.14>,detailed}}]}]]}
(rabbit_2@sunfeideMacBook-Pro)98>

```

![压入 10000w 消息后 web](https://github.com/moooofly/ImageCache/blob/master/Pictures/%E5%8E%8B%E5%85%A5%2010000w%20%E6%B6%88%E6%81%AF%E5%90%8E%20web.png "压入 10000w 消息后 web")

![压入 10000w 消息后 entop](https://github.com/moooofly/ImageCache/blob/master/Pictures/%E5%8E%8B%E5%85%A5%2010000w%20%E6%B6%88%E6%81%AF%E5%90%8E%20entop.png "压入 10000w 消息后 entop")

此时我的 MacPro 已经发热的不要不要的了～～

下面给出压测过程中的一些精彩瞬间

![2.9G时](https://github.com/moooofly/ImageCache/blob/master/Pictures/2.9G%208187.png "2.9G时")

![4.2G时](https://github.com/moooofly/ImageCache/blob/master/Pictures/4.2G%2011182.png "4.2G时")

![4.9G时](https://github.com/moooofly/ImageCache/blob/master/Pictures/4.9G%2022010.png "4.9G时")

![6G时](https://github.com/moooofly/ImageCache/blob/master/Pictures/6G%2011556.png "6G时")

![7.3G时](https://github.com/moooofly/ImageCache/blob/master/Pictures/7.3G%202411.png "7.3G时")

![10G时](https://github.com/moooofly/ImageCache/blob/master/Pictures/10G%205847.png "10G时")

![16G时](https://github.com/moooofly/ImageCache/blob/master/Pictures/16G%20%E6%9C%80%E5%90%8E%E7%8A%B6%E6%80%81.png "16G时")



### 实验结果说明

从上面的实验中可以得出以下几个结论：
- 只有统计数据库所在的 RabbitMQ 节点的内存占用随着压入消息量的增加而增大（从 56M 开始，最终达到 16G）；
- rabbit_mgmt_db 进程的 reductions 值从初始状态时的、同数量级值（3606267）增长到最终状态时的、高出其他所有进程两个数量级的值（3156198688）；注意：请忽略 <0.7856.0> 这个进程，因为这个是测试 shell 引入的，非 RabbitMQ 内部原生进程；
- 调用 sys:status/1 的瞬间 last_queue_length 的值并不一定与 "Queued messages" 中展示的值完全对应上，因为存在获取时差问题；web 页面上的黄色告警判定是基于 last_queue_length 的值（大于 1000  告警），从上面的信息中可以看到，存在 last_queue_length 为 0 ，而 "Queued messages" 积压上万消息的情况；同样，也存在 last_queue_length 值不为 0 ，而 "Queued messages" 为 0 的情况；（具体原因详见 `gen_server2.erl` 中 `drain/1` 消息搬移实现）



## 线上环境实际数据

信息获取
- last_queue_length : 44
- "Queued messages" : 0

```shell
[fei.sun@xg-napos-rmq-3 ~]$ sudo rabbitmqctl eval "sys:get_status(global:whereis_name(rabbit_mgmt_db))."
{status,<5548.2968.1472>,
        {module,gen_server2},
        [[{delegate,delegate_1},
          {'$ancestors',[<5548.2868.0>,rabbit_mgmt_sup,rabbit_mgmt_sup_sup,
                         <5548.2846.0>]},
          {last_queue_length,44},
          {'$initial_call',{gen,init_it,7}}],
         running,<5548.2868.0>,[],
         [{header,"Status for generic server rabbit_mgmt_db"},
          {data,[{"Status",running},
                 {"Parent",<5548.2868.0>},
                 {"Logged events",[]},
                 {"Queued messages",{0,[]}}]},
          {data,[{"State",
                  {state,[{channel_stats,2363484},
                          {connection_stats,2359387},
                          {consumers_by_channel,2371678},
                          {consumers_by_queue,2367581},
                          {node_node_stats,2379872},
                          {node_stats,2375775},
                          {queue_stats,2355287}],
                         2383969,2388066,2392086,#Ref<5548.0.106437.224150>,
                         {{channel_stats,<6332.21439.2856>},publish},
                         [{exchange,#Fun<rabbit_exchange.lookup.1>},
                          {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                         5000,#Ref<5548.0.1401.4079>,basic}}]}]]}
[fei.sun@xg-napos-rmq-3 ~]$
```


信息获取
- last_queue_length: 1031
- "Queued messages": 1107
    - channel_closed: 371
    - channel_created: 369
    - channel_stats: 307
    - connection_closed: 3
    - connection_created: 3
    - connection_stats: 49
    - queue_stats: 1
    - user_authentication_success: 3
    - augment_nodes: 1

```shell
[fei.sun@xg-napos-rmq-3 ~]$
[fei.sun@xg-napos-rmq-3 ~]$ sudo rabbitmqctl eval "sys:get_status(global:whereis_name(rabbit_mgmt_db))."
{status,<5548.2968.1472>,
    {module,gen_server2},
    [[{delegate,delegate_1},
      {'$ancestors',
          [<5548.2868.0>,rabbit_mgmt_sup,rabbit_mgmt_sup_sup,<5548.2846.0>]},
      {last_queue_length,1031},
      {'$initial_call',{gen,init_it,7}}],
     running,<5548.2868.0>,[],
     [{header,"Status for generic server rabbit_mgmt_db"},
      {data,
          [{"Status",running},
           {"Parent",<5548.2868.0>},
           {"Logged events",[]},
           {"Queued messages",
            {1107, 
             {summary,
                 [{{0,
                    {'$gen_cast',{event,{event,channel_closed,'_',none,'_'}}}},
                   371},
                  {{0,
                    {'$gen_cast',
                        {event,{event,channel_created,'_',none,'_'}}}},
                   369},
                  {{0,
                    {'$gen_cast',{event,{event,channel_stats,'_',none,'_'}}}},
                   307},
                  {{0,
                    {'$gen_cast',
                        {event,{event,connection_closed,'_',none,'_'}}}},
                   3},
                  {{0,
                    {'$gen_cast',
                        {event,{event,connection_created,'_',none,'_'}}}},
                   3},
                  {{0,
                    {'$gen_cast',
                        {event,{event,connection_stats,'_',none,'_'}}}},
                   49},
                  {{0,{'$gen_cast',{event,{event,queue_stats,'_',none,'_'}}}},
                   1},
                  {{0,
                    {'$gen_cast',
                        {event,
                            {event,user_authentication_success,'_',none,
                                '_'}}}},
                   3},
                  {{5,
                    {'$gen_call',
                        {'_','_'},
                        {augment_nodes,'_',
                            {no_range,no_range,no_range,no_range}}}},
                   1}]}}}]},
      {data,
          [{"State",
            {state,
                [{channel_stats,2363484},
                 {connection_stats,2359387},
                 {consumers_by_channel,2371678},
                 {consumers_by_queue,2367581},
                 {node_node_stats,2379872},
                 {node_stats,2375775},
                 {queue_stats,2355287}],
                2383969,2388066,2392086,#Ref<5548.0.106438.32815>,
                {{channel_stats,<6341.27963.3261>},publish},
                [{exchange,#Fun<rabbit_exchange.lookup.1>},
                 {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                5000,#Ref<5548.0.1401.4079>,basic}}]}]]}
[fei.sun@xg-napos-rmq-3 ~]$
```

信息获取
- last_queue_length: 79
- "Queued messages": 9
    - channel_closed: 3
    - channel_created: 1
    - channel_stats: 4
    - connection_stats: 1


```shell
[fei.sun@xg-napos-rmq-3 ~]$ sudo rabbitmqctl eval "sys:get_status(global:whereis_name(rabbit_mgmt_db))."
{status,<5548.2968.1472>,
    {module,gen_server2},
    [[{delegate,delegate_1},
      {'$ancestors',
          [<5548.2868.0>,rabbit_mgmt_sup,rabbit_mgmt_sup_sup,<5548.2846.0>]},
      {last_queue_length,79},
      {'$initial_call',{gen,init_it,7}}],
     running,<5548.2868.0>,[],
     [{header,"Status for generic server rabbit_mgmt_db"},
      {data,
          [{"Status",running},
           {"Parent",<5548.2868.0>},
           {"Logged events",[]},
           {"Queued messages",
            {9,
             [{0,
               {'$gen_cast',
                   {event,
                       {event,channel_closed,
                           [{pid,<6330.8200.3428>}],
                           none,1474194622633}}}},
              {0,
               {'$gen_cast',
                   {event,
                       {event,channel_stats,
                           [{pid,<6332.23792.4335>},
                            {transactional,false},
                            {confirm,false},
                            {consumer_count,0},
                            {messages_unacknowledged,0},
                            {messages_unconfirmed,0},
                            {messages_uncommitted,0},
                            {acks_uncommitted,0},
                            {prefetch_count,0},
                            {global_prefetch_count,0},
                            {state,closing},
                            {channel_queue_stats,[]},
                            {channel_exchange_stats,
                                [{{resource,<<"napos">>,exchange,
                                      <<"napos_luna_push_topic">>},
                                  [{publish,1}]}]},
                            {channel_queue_exchange_stats,
                                [{{{resource,<<"napos">>,queue,
                                       <<"napos_luna_push_queue">>},
                                   {resource,<<"napos">>,exchange,
                                       <<"napos_luna_push_topic">>}},
                                  [{publish,1}]}]}],
                           none,1474194622633}}}},
              {0,
               {'$gen_cast',
                   {event,
                       {event,channel_closed,
                           [{pid,<6332.23792.4335>}],
                           none,1474194622633}}}},
              {0,
               {'$gen_cast',
                   {event,
                       {event,channel_stats,
                           [{pid,<6332.16904.4254>},
                            {transactional,false},
                            {confirm,false},
                            {consumer_count,0},
                            {messages_unacknowledged,0},
                            {messages_unconfirmed,0},
                            {messages_uncommitted,0},
                            {acks_uncommitted,0},
                            {prefetch_count,0},
                            {global_prefetch_count,0},
                            {state,closing},
                            {channel_queue_stats,[]},
                            {channel_exchange_stats,
                                [{{resource,<<"napos">>,exchange,
                                      <<"napos.delivery">>},
                                  [{publish,1}]}]},
                            {channel_queue_exchange_stats,
                                [{{{resource,<<"napos">>,queue,
                                       <<"napos_delivery_order_queue">>},
                                   {resource,<<"napos">>,exchange,
                                       <<"napos.delivery">>}},
                                  [{publish,1}]}]}],
                           none,1474194622633}}}},
              {0,
               {'$gen_cast',
                   {event,
                       {event,connection_stats,
                           [{pid,<6330.11279.4042>},
                            {recv_oct,1526357},
                            {recv_cnt,171106},
                            {send_oct,696263},
                            {send_cnt,86018},
                            {send_pend,0},
                            {state,running},
                            {channels,0}],
                           none,1474194622634}}}},
              {0,
               {'$gen_cast',
                   {event,
                       {event,channel_closed,
                           [{pid,<6332.16904.4254>}],
                           none,1474194622633}}}},
              {0,
               {'$gen_cast',
                   {event,
                       {event,channel_created,
                           [{pid,<6332.25213.1465>},
                            {name,
                                <<"10.0.43.161:57192 -> 10.0.21.154:5672 (1)">>},
                            {connection,<6332.12472.5346>},
                            {number,1},
                            {user,<<"napos">>},
                            {vhost,<<"napos">>}],
                           none,1474194622633}}}},
              {0,
               {'$gen_cast',
                   {event,
                       {event,channel_stats,
                           [{pid,<6332.25213.1465>},
                            {transactional,false},
                            {confirm,false},
                            {consumer_count,0},
                            {messages_unacknowledged,0},
                            {messages_unconfirmed,0},
                            {messages_uncommitted,0},
                            {acks_uncommitted,0},
                            {prefetch_count,0},
                            {global_prefetch_count,0},
                            {state,starting},
                            {channel_queue_stats,[]},
                            {channel_exchange_stats,[]},
                            {channel_queue_exchange_stats,[]}],
                           none,1474194622633}}}},
              {0,
               {'$gen_cast',
                   {event,
                       {event,channel_stats,
                           [{pid,<6332.25213.1465>},
                            {transactional,false},
                            {confirm,false},
                            {consumer_count,0},
                            {messages_unacknowledged,0},
                            {messages_unconfirmed,0},
                            {messages_uncommitted,0},
                            {acks_uncommitted,0},
                            {prefetch_count,0},
                            {global_prefetch_count,0},
                            {state,running},
                            {channel_queue_stats,[]},
                            {channel_exchange_stats,[]},
                            {channel_queue_exchange_stats,[]}],
                           none,1474194622634}}}}]}}]},
      {data,
          [{"State",
            {state,
                [{channel_stats,2363484},
                 {connection_stats,2359387},
                 {consumers_by_channel,2371678},
                 {consumers_by_queue,2367581},
                 {node_node_stats,2379872},
                 {node_stats,2375775},
                 {queue_stats,2355287}],
                2383969,2388066,2392086,#Ref<5548.0.106440.254613>,
                {{channel_stats,<6330.9567.7066>},ack},
                [{exchange,#Fun<rabbit_exchange.lookup.1>},
                 {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                5000,#Ref<5548.0.1401.4079>,basic}}]}]]}
[fei.sun@xg-napos-rmq-3 ~]$
```

信息获取
- last_queue_length: 67
- "Queued messages": 501
    - channel_closed: 98
    - channel_created: 96
    - channel_stats: 291
    - connection_stats: 15
    - queue_stats: 1


```shell
[fei.sun@xg-napos-rmq-3 ~]$
[fei.sun@xg-napos-rmq-3 ~]$ sudo rabbitmqctl eval "sys:get_status(global:whereis_name(rabbit_mgmt_db))."
{status,<5548.2968.1472>,
    {module,gen_server2},
    [[{delegate,delegate_1},
      {'$ancestors',
          [<5548.2868.0>,rabbit_mgmt_sup,rabbit_mgmt_sup_sup,<5548.2846.0>]},
      {last_queue_length,67},
      {'$initial_call',{gen,init_it,7}}],
     running,<5548.2868.0>,[],
     [{header,"Status for generic server rabbit_mgmt_db"},
      {data,
          [{"Status",running},
           {"Parent",<5548.2868.0>},
           {"Logged events",[]},
           {"Queued messages",
            {501,
             {summary,
                 [{{0,
                    {'$gen_cast',{event,{event,channel_closed,'_',none,'_'}}}},
                   98},
                  {{0,
                    {'$gen_cast',
                        {event,{event,channel_created,'_',none,'_'}}}},
                   96},
                  {{0,
                    {'$gen_cast',{event,{event,channel_stats,'_',none,'_'}}}},
                   291},
                  {{0,
                    {'$gen_cast',
                        {event,{event,connection_stats,'_',none,'_'}}}},
                   15},
                  {{0,{'$gen_cast',{event,{event,queue_stats,'_',none,'_'}}}},
                   1}]}}}]},
      {data,
          [{"State",
            {state,
                [{channel_stats,2363484},
                 {connection_stats,2359387},
                 {consumers_by_channel,2371678},
                 {consumers_by_queue,2367581},
                 {node_node_stats,2379872},
                 {node_stats,2375775},
                 {queue_stats,2355287}],
                2383969,2388066,2392086,#Ref<5548.0.106473.143417>,
                {{connection_stats,<5548.23703.4077>},send_oct},
                [{exchange,#Fun<rabbit_exchange.lookup.1>},
                 {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                5000,#Ref<5548.0.1401.4079>,basic}}]}]]}
[fei.sun@xg-napos-rmq-3 ~]$
```


信息获取
- last_queue_length: 1224
- "Queued messages": 1264
    - channel_closed: 314
    - channel_created: 311
    - channel_stats: 583
    - connection_closed: 3
    - connection_created: 3
    - connection_stats: 43
    - queue_stats: 1
    - user_authentication_success: 3
    - augment_nodes: 2
    - augment_queues: 1

```shell
[fei.sun@xg-napos-rmq-3 ~]$ sudo rabbitmqctl eval "sys:get_status(global:whereis_name(rabbit_mgmt_db))."
{status,<5548.2968.1472>,
    {module,gen_server2},
    [[{delegate,delegate_1},
      {'$ancestors',
          [<5548.2868.0>,rabbit_mgmt_sup,rabbit_mgmt_sup_sup,<5548.2846.0>]},
      {last_queue_length,1224},
      {'$initial_call',{gen,init_it,7}}],
     running,<5548.2868.0>,[],
     [{header,"Status for generic server rabbit_mgmt_db"},
      {data,
          [{"Status",running},
           {"Parent",<5548.2868.0>},
           {"Logged events",[]},
           {"Queued messages",
            {1264,
             {summary,
                 [{{0,
                    {'$gen_cast',{event,{event,channel_closed,'_',none,'_'}}}},
                   314},
                  {{0,
                    {'$gen_cast',
                        {event,{event,channel_created,'_',none,'_'}}}},
                   311},
                  {{0,
                    {'$gen_cast',{event,{event,channel_stats,'_',none,'_'}}}},
                   583},
                  {{0,
                    {'$gen_cast',
                        {event,{event,connection_closed,'_',none,'_'}}}},
                   3},
                  {{0,
                    {'$gen_cast',
                        {event,{event,connection_created,'_',none,'_'}}}},
                   3},
                  {{0,
                    {'$gen_cast',
                        {event,{event,connection_stats,'_',none,'_'}}}},
                   43},
                  {{0,{'$gen_cast',{event,{event,queue_stats,'_',none,'_'}}}},
                   1},
                  {{0,
                    {'$gen_cast',
                        {event,
                            {event,user_authentication_success,'_',none,
                                '_'}}}},
                   3},
                  {{5,
                    {'$gen_call',
                        {'_','_'},
                        {augment_nodes,'_',
                            {no_range,no_range,no_range,no_range}}}},
                   2},
                  {{5,
                    {'$gen_call',
                        {'_','_'},
                        {augment_queues,'_',
                            {no_range,no_range,no_range,no_range},
                            basic}}},
                   1}]}}}]},
      {data,
          [{"State",
            {state,
                [{channel_stats,2363484},
                 {connection_stats,2359387},
                 {consumers_by_channel,2371678},
                 {consumers_by_queue,2367581},
                 {node_node_stats,2379872},
                 {node_stats,2375775},
                 {queue_stats,2355287}],
                2383969,2388066,2392086,#Ref<5548.0.106476.233231>,
                {{connection_stats,<5548.3082.6017>},send_oct},
                [{exchange,#Fun<rabbit_exchange.lookup.1>},
                 {queue,#Fun<rabbit_amqqueue.lookup.1>}],
                5000,#Ref<5548.0.1401.4079>,basic}}]}]]}
[fei.sun@xg-napos-rmq-3 ~]$
```














