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










