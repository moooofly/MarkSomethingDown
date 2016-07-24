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

由于 RabbitMQ 允许的并发连接数目存在上限，而上述问题理论上会导致相应的连接被长时间占用，因此，应该会导致特定时间内，可建立连接数目下降；

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

两者其实存在本质上的区别：对 RabbitMQ 而言，HAProxy 的方式没有完成 TCP 三次握手，sdfsfsf
sfsfsfsdfsf


（据说此方式为老版本的实现，新版本已经和 haproxy 实现方式一致）

## 影响范围

目前看来只是造成了日志中出现大量相关错误信息，浪费了部分可用连接数量；

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



## 监督树总体结构


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
                          |                             |
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
                   +------+---------+
                   |                |
               help_sup        rabbit_reader
```

其中 

- **topman** - 维护 pubsub 进程和 Topic 的映射关系； 
- **pubsub** - 关联特定  Topic 的进程 ；维护所有订阅到该 Topic 的进程信息； 
- **janus_acceptor** - 处理来自网络的 TCP 连接；动态创建 transport 和 client_proxy 进程，以处理后续协议交互； 
- **transport** - 针对某个 TCP 连接上的数据处理； 
- **client_proxy** - 实际处理订阅，取消订阅，以及消息推送的模块； 
- **mapper** -  提供轻量级进程注册管理功能； 


## 源码分析

## 影响范围


----------

# goproxy agent 探测 RabbitMQ 活性问题

## RabbitMQ 自身 heartbeat 保活方式

- 业务以 2.5s 时间间隔发送 heartbeat 给 RMQ
- RMQ 以 5s 时间间隔发送 heartbeat 给业务
似乎均为请求，没有应答

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



## 影响范围










