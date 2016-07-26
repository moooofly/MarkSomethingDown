





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
