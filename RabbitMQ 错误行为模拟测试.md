





----------

# handshake_timeout 模拟


```shell
➜  ~ nc 11.11.11.15 5672
（自行退出）
➜  ~
```

```shell
=INFO REPORT==== 26-Jul-2016::05:57:17 ===
accepting AMQP connection <0.30457.1> (11.11.11.1:57752 -> 11.11.11.15:5672)

=ERROR REPORT==== 26-Jul-2016::05:57:17 ===
closing AMQP connection <0.30457.1> (11.11.11.1:57752 -> 11.11.11.15:5672):
{handshake_timeout,handshake}
```
> 上面两条日志是同时输出的，即时间是一样的；这个与 error_logger 刷盘方式有关；

```shell
05:57:07.527784 IP 11.11.11.1.57752 > 11.11.11.15.amqp: Flags [S], seq 1831788832, win 65535, options [mss 1460,nop,wscale 5,nop,nop,TS val 1240822351 ecr 0,sackOK,eol], length 0
05:57:07.527840 IP 11.11.11.15.amqp > 11.11.11.1.57752: Flags [S.], seq 3830321900, ack 1831788833, win 28960, options [mss 1460,sackOK,TS val 18905480 ecr 1240822351,nop,wscale 6], length 0
05:57:07.528284 IP 11.11.11.1.57752 > 11.11.11.15.amqp: Flags [.], ack 1, win 4117, options [nop,nop,TS val 1240822351 ecr 18905480], length 0

05:57:17.530495 IP 11.11.11.15.amqp > 11.11.11.1.57752: Flags [R.], seq 1, ack 1, win 453, options [nop,nop,TS val 0 ecr 1240822351], length 0
```
> 可以看到在三次握手后的第 10 秒，服务器发送 RST 给 client ； 


----------

# connection_closed_with_no_data_received 模拟


> 该日志信息的含义：
>> The connection was closed before any packet was received. It's probably a load-balancer healthcheck: don't consider this a failure.


```shell
➜  ~ nc 11.11.11.15 5672
^C
➜  ~
```
> 通过 `ctrl + c` 终止

```shell
=INFO REPORT==== 26-Jul-2016::06:22:19 ===
accepting AMQP connection <0.7776.0> (11.11.11.1:57820 -> 11.11.11.15:5672)

=INFO REPORT==== 26-Jul-2016::06:22:19 ===
closing AMQP connection <0.7776.0> (11.11.11.1:57820 -> 11.11.11.15:5672):
connection_closed_with_no_data_received
```
> 注意：
> - 若 `closing AMQP connection` 对应的 **REPORT** 级别为 **INFO** ，则说明此为正常连接关闭处理；
> - 一般来讲，也不会看到 `accepting AMQP connection` 出现在 **REPORT** 级别 **DEBUG** 下，这种类型日志对应了 TCP healthcheck ，例如 HAProxy 的保活机制；


```shell
06:22:14.003771 IP 11.11.11.1.57820 > 11.11.11.15.amqp: Flags [S], seq 1007623008, win 65535, options [mss 1460,nop,wscale 5,nop,nop,TS val 1242321795 ecr 0,sackOK,eol], length 0
06:22:14.003794 IP 11.11.11.15.amqp > 11.11.11.1.57820: Flags [S.], seq 1833358307, ack 1007623009, win 28960, options [mss 1460,sackOK,TS val 19282099 ecr 1242321795,nop,wscale 6], length 0
06:22:14.003935 IP 11.11.11.1.57820 > 11.11.11.15.amqp: Flags [.], ack 1, win 4117, options [nop,nop,TS val 1242321795 ecr 19282099], length 0
06:22:19.220908 IP 11.11.11.1.57820 > 11.11.11.15.amqp: Flags [F.], seq 1, ack 1, win 4117, options [nop,nop,TS val 1242326997 ecr 19282099], length 0
06:22:19.221244 IP 11.11.11.15.amqp > 11.11.11.1.57820: Flags [R.], seq 1, ack 2, win 453, options [nop,nop,TS val 0 ecr 1242326997], length 0
```
> 可以看到，对于 RabbitMQ 而言，正常关闭序列为

```sequence
client->RabbitMQ: FIN
RabbitMQ->client: RST
```

> 对于使用这种关闭序列的原因，官方解释如下
>>  We don't call gen_tcp:close/1 here since it waits for pending output to be sent, which results in unnecessary delays. We could just terminate - the reader is the controlling process and hence its termination will close the socket. However, to keep the file_handle_cache accounting as accurate as possible we ought to close the socket w/o delay before termination.


----------


```shell
➜  ~ nc 11.11.11.15 5672
^D
➜  ~
```
> 通过 `ctrl + d` 终止


```shell
=INFO REPORT==== 26-Jul-2016::07:06:11 ===
accepting AMQP connection <0.30621.0> (11.11.11.1:57919 -> 11.11.11.15:5672)

=WARNING REPORT==== 26-Jul-2016::07:06:13 ===
closing AMQP connection <0.30621.0> (11.11.11.1:57919 -> 11.11.11.15:5672):
client unexpectedly closed TCP connection
```
> 概率出现上述错误；有时会出现 connection_closed_with_no_data_received 错误；


```shell
07:06:11.418649 IP 11.11.11.1.57919 > 11.11.11.15.amqp: Flags [S], seq 1497952456, win 65535, options [mss 1460,nop,wscale 5,nop,nop,TS val 1244951533 ecr 0,sackOK,eol], length 0
07:06:11.418672 IP 11.11.11.15.amqp > 11.11.11.1.57919: Flags [S.], seq 1733377893, ack 1497952457, win 28960, options [mss 1460,sackOK,TS val 19941452 ecr 1244951533,nop,wscale 6], length 0
07:06:11.418819 IP 11.11.11.1.57919 > 11.11.11.15.amqp: Flags [.], ack 1, win 4117, options [nop,nop,TS val 1244951533 ecr 19941452], length 0
07:06:11.990055 IP 11.11.11.1.57919 > 11.11.11.15.amqp: Flags [P.], seq 1:2, ack 1, win 4117, options [nop,nop,TS val 1244952101 ecr 19941452], length 1
07:06:11.990085 IP 11.11.11.15.amqp > 11.11.11.1.57919: Flags [.], ack 2, win 453, options [nop,nop,TS val 19941595 ecr 1244952101], length 0
07:06:12.161104 IP 11.11.11.1.57919 > 11.11.11.15.amqp: Flags [P.], seq 2:3, ack 1, win 4117, options [nop,nop,TS val 1244952270 ecr 19941595], length 1
07:06:12.161122 IP 11.11.11.15.amqp > 11.11.11.1.57919: Flags [.], ack 3, win 453, options [nop,nop,TS val 19941638 ecr 1244952270], length 0
07:06:12.339042 IP 11.11.11.1.57919 > 11.11.11.15.amqp: Flags [P.], seq 3:4, ack 1, win 4117, options [nop,nop,TS val 1244952447 ecr 19941638], length 1
07:06:12.339062 IP 11.11.11.15.amqp > 11.11.11.1.57919: Flags [.], ack 4, win 453, options [nop,nop,TS val 19941683 ecr 1244952447], length 0
07:06:12.524796 IP 11.11.11.1.57919 > 11.11.11.15.amqp: Flags [P.], seq 4:5, ack 1, win 4117, options [nop,nop,TS val 1244952632 ecr 19941683], length 1
07:06:12.524813 IP 11.11.11.15.amqp > 11.11.11.1.57919: Flags [.], ack 5, win 453, options [nop,nop,TS val 19941729 ecr 1244952632], length 0
07:06:13.160406 IP 11.11.11.1.57919 > 11.11.11.15.amqp: Flags [F.], seq 5, ack 1, win 4117, options [nop,nop,TS val 1244953263 ecr 19941729], length 0
07:06:13.160660 IP 11.11.11.15.amqp > 11.11.11.1.57919: Flags [R.], seq 1, ack 6, win 453, options [nop,nop,TS val 0 ecr 1244953263], length 0
```



----------

# bad_header 和 bad_version 错误

> 这两种错误一般不会遇到，此处仅贴出来供参考；

```shell
➜  ~ telnet 11.11.11.15 5672
Trying 11.11.11.15...
Connected to 11.11.11.15.
Escape character is '^]'.
amqp0091
AMQP	Connection closed by foreign host.
➜  ~

...

➜  ~
➜  ~ telnet 11.11.11.15 5672
Trying 11.11.11.15...
Connected to 11.11.11.15.
Escape character is '^]'.
AMQP0091
AMQP	Connection closed by foreign host.
➜  ~
```


```shell
=INFO REPORT==== 26-Jul-2016::07:24:34 ===
accepting AMQP connection <0.7395.1> (11.11.11.1:57970 -> 11.11.11.15:5672)

=ERROR REPORT==== 26-Jul-2016::07:24:34 ===
closing AMQP connection <0.7395.1> (11.11.11.1:57970 -> 11.11.11.15:5672):
{bad_header,<<"amqp0091">>}

...

=INFO REPORT==== 26-Jul-2016::07:24:51 ===
accepting AMQP connection <0.7533.1> (11.11.11.1:57971 -> 11.11.11.15:5672)

=ERROR REPORT==== 26-Jul-2016::07:24:51 ===
closing AMQP connection <0.7533.1> (11.11.11.1:57971 -> 11.11.11.15:5672):
{bad_version,{48,48,57,49}}
```

抓包信息就不贴出来了，错误原因参考如下代码：

```erlang
handle_input(handshake, <<"AMQP", A, B, C, D, Rest/binary>>, State) ->
    {Rest, handshake({A, B, C, D}, State)};
handle_input(handshake, <<Other:8/binary, _/binary>>, #v1{sock = Sock}) ->
    refuse_connection(Sock, {bad_header, Other});
...
handshake({0, 0, 9, 1}, State) ->
    start_connection({0, 9, 1}, rabbit_framing_amqp_0_9_1, State);
...
handshake(Vsn, #v1{sock = Sock}) ->
    refuse_connection(Sock, {bad_version, Vsn}).
...
refuse_connection(Sock, Exception, {A, B, C, D}) ->
    ok = inet_op(fun () -> rabbit_net:send(Sock, <<"AMQP",A,B,C,D>>) end),
    throw(Exception).

-ifdef(use_specs).
-spec(refuse_connection/2 :: (rabbit_net:socket(), any()) -> no_return()).
-endif.
refuse_connection(Sock, Exception) ->
    refuse_connection(Sock, Exception, {0, 0, 9, 1}).
```