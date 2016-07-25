[toc]


----------

# handshake_timeout 错误

## 错误日志

```shell
=ERROR REPORT==== 21-Jul-2016::19:21:26 ===
closing AMQP connection <0.11270.4127> (10.0.133.159:38511 -> 10.0.21.154:5672):
{handshake_timeout,handshake}
```

## 出错频率

n sec ～ n min    
多数为 10+ minute 级别，偶尔会 1 minute 中多次；     


## 原因总结

业务 client 与 RabbitMQ 成功建立 TCP 连接后，在 10s 内（默认值）未发送 AMQP handshake 协议包（即 Protocol-Header 0-9-1 包）；此问题属于业务 client 侧 bug ；

## 影响范围

由于 RabbitMQ 用于并行处理 client 连接的 ranch_acceptor 进程数目存在上限（默认值为 10），而上述问题理论上会导致相应的连接被长时间占用，因此，应该会导致特定时间内，可被及时处理的client 连接数目下降；

## 源码分析

在 `rabbit_reader.erl` 中
```erlang
...
init(Parent, HelperSup, Ref, Sock) ->
    %% 成功获取到 client socket 所有权
    rabbit_net:accept_ack(Ref, Sock),
    Deb = sys:debug_options([]),
    %% 开始 TCP 连接处理
    start_connection(Parent, HelperSup, Deb, Sock).
...
%% 正式开始 TCP + AMQP 协议处理
start_connection(Parent, HelperSup, Deb, Sock) ->
	...
    %% 允许的最长 AMQP 0-9-1 handshake 时间，即在成功建立 socket 连接后最长等待时间
    %% 以毫秒为单位，默认为 10000ms ，即 10s
    {ok, HandshakeTimeout} = application:get_env(rabbit, handshake_timeout),
	...
    %% 启动超时定时器，针对 handshake 状态机下未收到任何协议包的情况
    erlang:send_after(HandshakeTimeout, self(), handshake_timeout),
    ...
    try
	    %% 进入数据接收＋处理循环
        run({?MODULE, recvloop,
             [Deb, [], 0, switch_callback(rabbit_event:init_stats_timer(
                                            State, #v1.stats_timer),
                                          handshake, 8)]}),
        %% 循环正常退出时，即 TCP 连接正常关闭情况下，会输出如下日志
        log(info, "closing AMQP connection ~p (~s)~n", [self(), Name])
    catch
		%% 针对异常情况的处理
		%% handshake_timeout 异常会在这里被捕获
        Ex ->
          log_connection_exception(Name, Ex)
    after
	    %% 无论正常关闭，或者异常关闭 TCP 连接，服务器侧均采用如下方式进行 socket 关闭处理
        rabbit_net:fast_close(Sock),
		...
    end,
    done.
...
log_connection_exception(Name, Ex) ->
	%% 根据捕获到的异常信息类型，确定异常级别
    Severity = case Ex of
                   connection_closed_with_no_data_received -> debug;
                   connection_closed_abruptly              -> warning;
                   %% handshake_timeout 异常输出如下级别
                   _                                       -> error
               end,
    log_connection_exception(Severity, Name, Ex).
...
log_connection_exception(Severity, Name, Ex) ->
    log(Severity, "closing AMQP connection ~p (~s):~n~p~n",
        [self(), Name, Ex]).
...
recvloop(Deb, Buf, BufLen, State = #v1{sock = Sock, recv_len = RecvLen})
  when BufLen < RecvLen ->
    case rabbit_net:setopts(Sock, [{active, once}]) of
        ok              -> mainloop(Deb, Buf, BufLen,
                                    State#v1{pending_recv = true});
        {error, Reason} -> stop(Reason, State)
    end;
...
mainloop(Deb, Buf, BufLen, State = #v1{sock = Sock,
                                       connection_state = CS,
                                       connection = #connection{
                                         name = ConnName}}) ->
    %% 从进程邮箱中获取数据，其中包括如下内容
    %% 1. 从 TCP 连接上获取到的数据报文；
    %% 2. 从 TCP 连接上获取到的 FIN or RST 报文；
    %% 3. 当前 TCP 连接异常通知消息；
    %% 4. 由超时定时器发送的 handshake_timeout 消息；
    Recv = rabbit_net:recv(Sock),
    ...
    case Recv of
		...
        {other, Other}  ->  %% 其它错误处理
            case handle_other(Other, State) of
                stop     -> ok;
                NewState -> recvloop(Deb, Buf, BufLen, NewState)
            end
    end.
...
handle_other(handshake_timeout, State)
  when ?IS_RUNNING(State) orelse ?IS_STOPPING(State) ->
    State;
handle_other(handshake_timeout, State) ->
    maybe_emit_stats(State),
    %% 这里抛出 handshake_timeout 异常
    throw({handshake_timeout, State#v1.callback});
...

```

在 `rabbit_net.erl` 中
```erlang
recv(Sock) when is_port(Sock) ->
    recv(Sock, {tcp, tcp_closed, tcp_error}).

recv(S, {DataTag, ClosedTag, ErrorTag}) ->
    receive
        {DataTag, S, Data}    -> {data, Data};    %% 接收到 TCP 数据包
        {ClosedTag, S}        -> closed;          %% TCP 连接关闭（FIN or RST）
        {ErrorTag, S, Reason} -> {error, Reason}; %% TCP 连接相关错误
        Other                 -> {other, Other}   %% 其它错误处理（如 handshake_timeout）
    end.
```


----------

# enotconn (socket is not connected) 错误

## 错误日志

```shell
=ERROR REPORT==== 21-Jul-2016::19:22:03 ===
Error on AMQP connection <0.21024.4027>: enotconn (socket is not connected)
```

## 出错频率

n sec ～ n min ；   
多数为 1~2 minute 级别，偶尔会 1 minute 中多次；     

## 原因总结

应该是由于业务 client 或 goproxy agent 发起的保活探测 TCP 序列导致的，通过抓包能够看到如下 TCP 序列每隔 2s 重复一次；

```sequence
goproxy agent->RabbitMQ: SYN
RabbitMQ->goproxy agent: SYN,ACK
goproxy agent->RabbitMQ: ACK
goproxy agent->RabbitMQ: RST,ACK
```

而 HAProxy 的健康监测 TCP 序列如下

```sequence
HAProxy->RabbitMQ: SYN
RabbitMQ->HAProxy: SYN,ACK
HAProxy->RabbitMQ: RST,ACK
```

两者其实存在一定区别：差异分析还需要查阅一些资料，后续补上；


## 影响范围

目前看来只是造成了日志中出现大量相关错误信息，浪费了部分可用连接数量；    
据说上述行为逻辑属于老版本的实现，新版本已经和 haproxy 实现方式一致；建议升级 goproxy 和 goproxy agent 为最新版本；

## 源码分析

在 `rabbit_reader.erl` 中

```erlang
...
%% 正式开始 TCP + AMQP 协议处理
start_connection(Parent, HelperSup, Deb, Sock) ->
	...
	%% 获取当前 TCP 连接两端的 ip 和 port 信息，拼接成连接信息字符串
    Name = case rabbit_net:connection_string(Sock, inbound) of
               {ok, Str}         -> Str;
               %% 在获取时，触发 socket 错误，认为当前连接已不存在
               %% 注意：这里没有输出异常日志
               {error, enotconn} -> rabbit_net:fast_close(Sock),
                                    exit(normal);
               {error, Reason}   -> socket_error(Reason),
                                    rabbit_net:fast_close(Sock),
                                    exit(normal)
           end,
	...
    %% 输出异常日志的地方
    %% 关键：上面 rabbit_net:connection_string 中调用的就是 rabbit_net:socket_ends
    %% 同样的代码上面没有报错，而此处会报错，说明在两段临近代码的执行间发生了 TCP 连接断开！
    {PeerHost, PeerPort, Host, Port} =
        socket_op(Sock, fun (S) -> rabbit_net:socket_ends(S, inbound) end),
    ...
    done.
...
socket_op(Sock, Fun) ->
    case Fun(Sock) of
        {ok, Res}       -> Res;
        {error, Reason} -> %% 输出错误日志
					       socket_error(Reason),
					       %% 关闭 TCP 连接
                           rabbit_net:fast_close(Sock),
                           %% 正常退出 rabbit_reader 进程
                           exit(normal)
    end.
...
socket_error(Reason) when is_atom(Reason) ->
    log(error, "Error on AMQP connection ~p: ~s~n",
        [self(), rabbit_misc:format_inet_error(Reason)]);
...
```

在 `rabbit_net.erl` 中

```erlang
...
connection_string(Sock, Direction) ->
	%% 获取 socket 连接两端 ip 和 port 信息
    case socket_ends(Sock, Direction) of
        {ok, {FromAddress, FromPort, ToAddress, ToPort}} ->
            {ok, rabbit_misc:format(
                   "~s:~p -> ~s:~p",
                   [maybe_ntoab(FromAddress), FromPort,
                    maybe_ntoab(ToAddress),   ToPort])};
        Error ->
            Error
    end.

socket_ends(Sock, Direction) ->
    {From, To} = sock_funs(Direction),
    %% 获取 tcp 通信两端的 ip 和 port
    case {From(Sock), To(Sock)} of
        {{ok, {FromAddress, FromPort}}, {ok, {ToAddress, ToPort}}} ->
            {ok, {rdns(FromAddress), FromPort,
                  rdns(ToAddress),   ToPort}};
        {{error, _Reason} = Error, _} ->
            Error;
        {_, {error, _Reason} = Error} ->
            Error
    end.
...
```



----------


# rabbit_channel_sup_sup 的 shutdown 错误 
## 错误日志

```shell
=SUPERVISOR REPORT==== 22-Jul-2016::15:21:21 ===
     Supervisor: {<0.25278.357>, rabbit_channel_sup_sup}
     Context:    shutdown_error
     Reason:     shutdown
     Offender:   [{nb_children,1},
                  {name,channel_sup},
                  {mfargs,{rabbit_channel_sup,start_link,[]}},
                  {restart_type,temporary},
                  {shutdown,infinity},
                  {child_type,supervisor}]
```


## 出错频率

n sec ～ n min    
业务高峰时段 1 second  中多次，平时 minute 级别    


## RabbitMQ 监督树结构

若想理解上述错误报告的含义，首先需要正确理解 RabbitMQ 内部进程的组织形式；

```
                                                      |
                                                  rabbit_sup
                                                      |
                                                      | (one_for_all)
                                               tcp_listener_sup
                                                      |
                                                      | (one_for_all)
                                   +------------------+--------------------+
                                   |                                       |
                                   |                                       |
                          ranch_listener_sup                        tcp_listener
                                   |
                                   | (rest_for_one)
                          +--------+--------------------+
                          |(1)                          |(2)
                          |                             |
                   ranch_conns_sup             ranch_acceptors_sup
                          |                             |
                          |                             | (one_for_one)
              +-----------+-----------+        +--------+--------+
              |           |           |        |        |        |
             ...          |          ...      ...       |       ...
                rabbit_connection_sup             ranch_acceptor
                          |
                          | (one_for_all)
              +-----------+-----------+
              |                       |
  rabbit_connection_helper_sup    rabbit_reader
              |
```

其中 

- **ranch_acceptors_sup** - 创建监听 socket ；默认监听所有 interface 上的 `5672` 端口；
- **ranch_acceptor** - N 个 ranch_acceptor 进程共享同一个监听 socket ，并获取来自 client 的 TCP 连接；ranch_acceptor 进程数目可通过配置项 num_tcp_acceptors 进行配置，默认为 10 ；
- **ranch_conns_sup** - 负责创建以 rabbit_connection_sup 为根的 AMQP connection 相关进程树；采取同步的方式从 ranch_acceptor 获取 client socket 控制权，并最终转移给 rabbit_reader 进程；负责控制最大并发连接；
- **rabbit_connection_sup** - 与每一条 TCP 连接对应的进程树结构的根； 
- **rabbit_reader** - 负责 TCP 协议和 AMQP 协议的相关处理； 


```
                                                |
                                   rabbit_connection_helper_sup
                                                |
                                                | (one_for_one)
                               +-----------+----+-----+--------------+
                               |           |          |              |
                               |           |          |              |
                               |     heartbeat_sender |              |
                               |                heartbeat_receiver   |
                     rabbit_channel_sup_sup             rabbit_queue_collector
                               |
                               | (simple_one_for_one)
                   +-----------+--------------+
                   |           |              |
                  ...          |             ...
                        rabbit_channel_sup
                               |
                               | (one_for_all)
               +---------------+---------------+
               |               |               |
               |               |               |
         rabbit_channel   rabbit_writer   rabbit_limiter
```

其中 


- **heartbeat_sender** - 负责 heartbeat 发送处理；
- **heartbeat_receiver** - 负责 heartbeat 接收处理；
- **rabbit_queue_collector** - 负责处理具有 exclusive 属性的 queue ；
- **rabbit_channel_sup** - 收到来自 client 的 channel.open 信令时，会在 rabbit_channel_sup_sup 下创建以 rabbit_channel_sup 为根的进程树，对应一个 channel 的处理；
- **rabbit_channel** - 对应 AMQP 0-9-1 中的 channel 实现；
- **rabbit_writer** - 负责发送 frame 给 client ；
- **rabbit_limiter** - 负责与 QoS 和流控相关的 channel 处理；



## 源码分析

```erlang
...
%% 终止当前监督者进程
terminate(_Reason, State) ->
    terminate_children(State#state.children, State#state.name).
...
%% 终止当前监督者进程下的所有子进程
terminate_children(Children, SupName) ->
    terminate_children(Children, SupName, []).
...
terminate_children([Child | Children], SupName, Res) ->
    NChild = do_terminate(Child, SupName),
    terminate_children(Children, SupName, [NChild | Res]);
...
%% 终止特定子进程
do_terminate(Child, SupName) when is_pid(Child#child.pid) ->
	%% 按照子进程规范 shutdown 子进程
    case shutdown(Child#child.pid, Child#child.shutdown) of
        ok ->
            ok;
        {error, normal} when not ?is_permanent(Child#child.restart_type) ->
            ok;
        {error, OtherReason} ->
	        %% 在关闭时发生错误
            report_error(shutdown_error, OtherReason, Child, SupName)
    end,
    Child#child{pid = undefined};
...
report_error(Error, Reason, Child, SupName) ->
    ErrorMsg = [{supervisor, SupName},      %% 所属监督者
		{errorContext, Error},              %% 发生错误时的上下文
		{reason, Reason},                   %% 发生错误的原因
		{offender, extract_child(Child)}],  %% 发生问题的进程
	%% 输出错误信息到 SASL 日志中
    error_logger:error_report(supervisor_report, ErrorMsg).

extract_child(Child) when is_list(Child#child.pid) ->
    [{nb_children, length(Child#child.pid)},
     {name, Child#child.name},
     {mfargs, Child#child.mfargs},
     {restart_type, Child#child.restart_type},
     {shutdown, Child#child.shutdown},
     {child_type, Child#child.child_type}];
...

```

## 原因总结

通过源码可知，rabbit_channel_sup_sup 进程的创建对应了 rabbit_reader 进程收到来自 client 的 AMQP connection.open 信令；而 rabbit_channel_sup 和其下子进程的创建对应了 rabbit_reader 进程收到来自 client 的 AMQP channel.open 信令；    
从 SASL 日志中看到：Supervisor: {<0.25278.357>, rabbit_channel_sup_sup} 中的 <0.25278.357> 即 rabbit_channel_sup_sup 进程 pid 不断变化，没有重复（可以通过过滤进行确认），说明以 rabbit_channel_sup_sup 为根的进程树在不断的销毁和创建；    
在正常的连接关闭序列下，应该不会报上述错误日志（后续进行试验验证），因此，该问题应该和业务的连接关闭处理逻辑有关；

## 影响范围

目前看来，该问题除了会导致 RabbitMQ 进程的不断创建和销毁外（增加了一定开销），未造成其它直接影响；

----------

# goproxy agent 探测 RabbitMQ 活性问题

## RabbitMQ 自身 heartbeat 保活方式

- 业务以 2.5s 时间间隔发送 heartbeat 给 RMQ
抓包信息如下：
![业务每 2.5 秒发送一次 heartbeat 包](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/%E4%B8%9A%E5%8A%A1%E6%AF%8F2.5%E7%A7%92%E5%8F%91%E9%80%81%E4%B8%80%E6%AC%A1heartbeat%E5%8C%85.png "业务发送 heartbeat 包情况")


- RMQ 以 5s 时间间隔发送 heartbeat 给业务

## haproxy 健康检查方式

```sequence
HAProxy->RabbitMQ: SYN
RabbitMQ->HAProxy: SYN,ACK
HAProxy->RabbitMQ: RST,ACK
```

## goproxy agent 健康检查方式

（据说此方式为老版本的实现，新版本已经和 haproxy 实现方式一致）

```sequence
goproxy agent->RabbitMQ: SYN
RabbitMQ->goproxy agent: SYN,ACK
goproxy agent->RabbitMQ: ACK
goproxy agent->RabbitMQ: RST,ACK
```


## 源码分析

略

## 影响范围

- goproxy 中实现的健康检查方式建议和 HAProxy 一致；
- 理论上讲，业务通过 AMQP 协议中的 heartbeat 功能就能够实现可靠监测；个人认为使用这种方式更合理；
- 若想要同时使用 goproxy 和 heartbeat 两种方式进行监测，则建议：根据实际情况，调整 heartbeat 监测超时时间；目前抓包显示，业务使用了 2.5s 的时间间隔，会导致 RabbitMQ 处理大量心跳消息，理论上讲，会导致常规业务消息的处理被拖慢；







