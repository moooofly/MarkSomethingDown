





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
> 上面两条日志是同时输出的，可以看到时间是一样的；

```shell
05:57:07.527784 IP 11.11.11.1.57752 > 11.11.11.15.amqp: Flags [S], seq 1831788832, win 65535, options [mss 1460,nop,wscale 5,nop,nop,TS val 1240822351 ecr 0,sackOK,eol], length 0
05:57:07.527840 IP 11.11.11.15.amqp > 11.11.11.1.57752: Flags [S.], seq 3830321900, ack 1831788833, win 28960, options [mss 1460,sackOK,TS val 18905480 ecr 1240822351,nop,wscale 6], length 0
05:57:07.528284 IP 11.11.11.1.57752 > 11.11.11.15.amqp: Flags [.], ack 1, win 4117, options [nop,nop,TS val 1240822351 ecr 18905480], length 0

05:57:17.530495 IP 11.11.11.15.amqp > 11.11.11.1.57752: Flags [R.], seq 1, ack 1, win 453, options [nop,nop,TS val 0 ecr 1240822351], length 0
```
> 可以看到在三次握手后的第 10 秒，服务器发送 RST 给 client ； 


----------

# connection_closed_with_no_data_received 模拟




```shell
➜  ~ nc 11.11.11.15 5672
^C
➜  ~
```

```shell
=INFO REPORT==== 26-Jul-2016::06:22:19 ===
accepting AMQP connection <0.7776.0> (11.11.11.1:57820 -> 11.11.11.15:5672)

=INFO REPORT==== 26-Jul-2016::06:22:19 ===
closing AMQP connection <0.7776.0> (11.11.11.1:57820 -> 11.11.11.15:5672):
connection_closed_with_no_data_received
```