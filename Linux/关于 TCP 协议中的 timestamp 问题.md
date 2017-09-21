# 关于 TCP 协议中的 timestamp 问题

之前和同事聊抓包分析过程中，进行应用层 request-response 耗时计算时，时间戳从哪里获取的问题时，扯出了“TCP 协议是否带有时间戳”的问题；我认为 TCP 协议已经带有 timestamp（因为我记得看到过），而另外一个人认为没有；

今天在看到另外一篇文章时，看到和 timestamp 相关的信息，遂成此文；


## RFC 793

### Options

> An Option field may contain several options, and each option may be several octets in length.  The options are used primarily in testing situations; for example, to carry **timestamps**.  Both the Internet Protocol and TCP provide for options fields.

## RFC 1137

解释了当 TIME-WAIT 状态不足时将会发生什么；

## RFC 1323

实现了 TCP 拓展规范，以保证网络繁忙状态下的高可用。除此之外，另外，它定义了一个新的 TCP 选项，即两个四字节的 timestamp fields 时间戳字段，第一个是 TCP 发送方的当前时钟时间戳，而第二个是从远程主机接收到的最新时间戳；


## wiki：TCP timestamps

**TCP timestamps**, defined in **RFC 1323**, can help TCP determine in which order packets were sent. TCP timestamps are not normally aligned to the system clock and start at some random value. Many operating systems will increment the timestamp for every elapsed millisecond; however the RFC only states that the ticks should be proportional.

There are two timestamp fields:

- a 4-byte **sender timestamp** value (my timestamp)
- a 4-byte **echo reply timestamp** value (the most recent timestamp received from you).

TCP timestamps are used in an algorithm known as Protection Against Wrapped Sequence numbers, or **PAWS** (see RFC 1323 for details). **PAWS is used when the receive window crosses the sequence number wraparound boundary**. In the case where a packet was potentially retransmitted it answers the question: "Is this sequence number in the first 4 GB or the second?" And the timestamp is used to break the tie.

Also, the Eifel detection algorithm ([RFC 3522](https://tools.ietf.org/html/rfc3522)) uses TCP timestamps to determine if retransmissions are occurring because packets are lost or simply out of order.


## sysctl 配置

linux 上 tcp_timestamps 默认开启

```
root@vagrant-ubuntu-trusty:~] $ sysctl -a|grep tcp_timestamps
net.ipv4.tcp_timestamps = 1
```


## 参考

- [RFC 793](http://www.rfc-editor.org/rfc/rfc793.txt)
- [RFC 1137](http://www.rfc-editor.org/rfc/rfc1137.txt)
- [RFC 1323](http://www.rfc-editor.org/rfc/rfc1323.txt)
- [wiki: Transmission Control Protocol](https://en.wikipedia.org/wiki/Transmission_Control_Protocol)



