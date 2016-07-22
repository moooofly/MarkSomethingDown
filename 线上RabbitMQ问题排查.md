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

## 源码分析

在 rabbit_reader.erl 中
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
    {ok, HandshakeTimeout} = application:get_env(rabbit, handshake_timeout),
	...
    %% 启动超时定时器，针对 handshake 状态机下未收到任何协议包的情况
    erlang:send_after(HandshakeTimeout, self(), handshake_timeout),
    ...
    try
        run({?MODULE, recvloop,
             [Deb, [], 0, switch_callback(rabbit_event:init_stats_timer(
                                            State, #v1.stats_timer),
                                          handshake, 8)]}),
        log(info, "closing AMQP connection ~p (~s)~n", [self(), Name])
    catch
        Ex ->
          log_connection_exception(Name, Ex)
    after
        rabbit_net:fast_close(Sock),
		...
    end,
    done.
...
log_connection_exception(Name, Ex) ->
    Severity = case Ex of
                   connection_closed_with_no_data_received -> debug;
                   connection_closed_abruptly              -> warning;
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
    Recv = rabbit_net:recv(Sock),
    case CS of
        pre_init when Buf =:= [] ->
            %% We only log incoming connections when either the
            %% first byte was received or there was an error (eg. a
            %% timeout).
            %%
            %% The goal is to not log TCP healthchecks (a connection
            %% with no data received) unless specified otherwise.
            log(case Recv of
                  closed -> debug;
                  _      -> info
                end, "accepting AMQP connection ~p (~s)~n",
                [self(), ConnName]);
        _ ->
            ok
    end,
    case Recv of
        {data, Data} ->
            recvloop(Deb, [Data | Buf], BufLen + size(Data),
                     State#v1{pending_recv = false});
        closed when State#v1.connection_state =:= closed ->
            ok;
        closed when CS =:= pre_init andalso Buf =:= [] ->
            stop(tcp_healthcheck, State);
        closed ->
            stop(closed, State);
        {other, {heartbeat_send_error, Reason}} ->
            %% The only portable way to detect disconnect on blocked
            %% connection is to wait for heartbeat send failure.
            stop(Reason, State);
        {error, Reason} ->
            stop(Reason, State);
        {other, {system, From, Request}} ->
            sys:handle_system_msg(Request, From, State#v1.parent,
                                  ?MODULE, Deb, {Buf, BufLen, State});
        {other, Other}  ->
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
    throw({handshake_timeout, State#v1.callback});
...

```


## 影响范围

----------

# enotconn (socket is not connected) 错误

## 错误日志

```shell
=ERROR REPORT==== 21-Jul-2016::19:22:03 ===
Error on AMQP connection <0.21024.4027>: enotconn (socket is not connected)
```

## 出错频率


n sec ～ n min
多数为 1~2 minute 级别，偶尔会 1 minute 中多次； 

## 源码分析



## 影响范围


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










