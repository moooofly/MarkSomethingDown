

# Disk Alarms

RabbitMQ 能够在磁盘空闲空间低于某个阈值时阻塞 producer ；该实现方式很有必要，因为即使是非持久消息，也会在某些时候被 page out 到磁盘上；而磁盘空间的耗尽可能会导致服务的崩溃；默认情况下，RabbitMQ 会在空闲磁盘空间低于 50MB 时阻塞住 producer ，并阻止驻留内存的消息被 page out 到磁盘；这种行为能够减少，但无法消除由于磁盘空间耗尽导致崩溃的可能性；尤其需要注意的是，如果消息被 page out 到磁盘过快，则非常有可能耗尽磁盘空间，并在两次运行磁盘空间监控之间发生崩溃；一种更加保守的方式为设置阈值限制等于系统中已安装的内存量 ；

全局流控会在空闲磁盘空间低于配置的阈值时触发；drive 或 partition 上能够被 broker 的数据库使用的空闲空间量，会以至少 10 秒一次的频率被监控，从而决定是否触发或清除告警；

监控进程是伴随 broker 启动而启动的，并会在 broker 的日志文件中输出如下内容：

```shell
=INFO REPORT==== 23-Jun-2012::14:52:41 ===
Disk free limit set to 953MB
```

在无法识别的平台上，监控进程将不被启动，并在日志中输出如下内容：

```shell
=WARNING REPORT==== 23-Jun-2012::15:45:29 ===
Disabling disk free space monitoring
```

当在 cluster 中运行 RabbitMQ 时，磁盘告警为 cluster 范围的；只要 cluster 中有一个节点发生了低于阈值的情况，所有节点都会阻塞连接；

RabbitMQ 会周期性检查空闲磁盘空间量；检查的频率和上次检查时的空闲量有关（这是为了保证磁盘告警能够在空闲空间耗尽时被及时触发）；通常情况下，磁盘空间检查为 10 一次，但该频率会随着阈值的接近而增加；当非常接近阈值时，RabbitMQ 会以每秒 10 次的频率进行检查，此时会对系统负载有所影响；

# Configuring the Disk Free Space Limit

磁盘空闲空间阈值可以通过 `disk_free_limit` 进行配置；默认要求在数据库所在分区上只要存在 50MB 的空闲空间；如下配置文件设置磁盘空间空闲阈值为 1GB ：
```shell
[{rabbit, [{disk_free_limit, 1000000000}]}].
```

或者你也可以使用具体的内存单位进行设置（kB, kiB, MB, MiB, GB, GiB 等）：
```shell
[{rabbit, [{disk_free_limit, "1GB"}]}].
```

还可以设置相对于 RAM 的磁盘空闲阈值；如下配置文件中设置了磁盘空闲阈值和机器中的 RAM 相同：
```shell
[{rabbit, [{disk_free_limit, {mem_relative, 1.0}}]}].
```

该阈值可以通过 `rabbitmqctl set_disk_free_limit disk_limit` 或 `rabbitmqctl set_disk_free_limit mem_relative fraction` 命令进行运行时调整；变更的效果在 broker 关闭前一直有效；若想配置一直有效，请写入配置文件；






