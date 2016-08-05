
# Memory Alarms

在启动了 RabbitMQ 后，若执行过 `rabbitmqctl set_vm_memory_high_watermark fraction` 命令，RabbitMQ server 会探测计算机上已安装 RAM 总量；默认情况下，当 RabbitMQ server 使用了超过 40% 的 RAM 内存时，会触发内存告警并阻塞住所有 connection ；一旦告警被解除（例如当 server 将消息 page out 到磁盘时，或者将消息 delivery 到客户端时）常规服务能力就能够恢复了；

内存阈值默认被设置为安装 RAM 的 **40%** ；

> ⚠️ 该值并不会真正阻止 RabbitMQ server 使用超过 40% 的内存，而只是一个会令 publisher 开始被阻塞的点；在最坏的情况下，由于 Erlang 垃圾回收器的原因，可能会导致内存使用量被 double（默认情况下为 RAM 的 80%）；因此强烈建议开启操作系统本身的 swap 或 page 功能；

在 32-bit 体系架构中，每个进程可用内存的限制为 2GB ；而常规实现的 64-bit 体系架构（即AMD64 和 Intel EM64T）中仅允许每个进程使用区区 256TB ；64-bit Windows 更进一步将其限制为8TB ；另外需要注意的是，就算在 64-bit 操作系统中，一个 32-bit 进程通常也只会使用最大 2GB 的地址空间；

## Configuring the Memory Threshold

可以通过配置文件调整触发**流控机制**的内存阈值；下面的示例中将阈值设置成默认的 0.4 ：

```shell
[{rabbit, [{vm_memory_high_watermark, 0.4}]}].
```

默认值 0.4 代表了 40% 的已安装 RAM 或者 40% 的可用虚拟地址空间（比前者更小）；例如在一个32-bit 平台上，如果你安装了 4GB 的 RAM ，那么 40% 的 4GB 为 1.6GB ，但是在 32-bit 的 Windows 上，通常会限制进程只能使用 2GB 内存，因此这里实际得到的阈值为 2GB 的 40% ，即 820MB ；

另一种方案为，直接设置节点可用内存阈值为一个具体数值；下面的例子中设置阈值为 1073741824 字节（1024 MB） ：

```shell
[{rabbit, [{vm_memory_high_watermark, {absolute, 1073741824}}]}].
```

与上面相同的例子，但使用了内存单位 MiB：

```shell
[{rabbit, [{vm_memory_high_watermark, {absolute, "1024MiB"}}]}].
```

如果上面设置的绝对数值超过了实际安装的 RAM 大小，或者可用的虚拟地址空间大小，阈值会被自动调整为两者中较小的那个值；

在 RabbitMQ server 启动时，内存使用限制信息会输出到 `RABBITMQ_NODENAME.log` 文件中：

```shell
=INFO REPORT==== 29-Oct-2009::15:43:27 ===
Memory limit set to 2048MB.
```

可以通过 `rabbitmqctl status` 命令查询内存阈值的具体数值；

可以在 broker 处于运行状态时进行阈值的修改，只需执行 `rabbitmqctl set_vm_memory_high_watermark fraction` 或者 `rabbitmqctl set_vm_memory_high_watermark absolute memory_limit` 命令；

可以在上述命令中直接使用内存单位（如 MiB）；

通过上述命令进行变更后，修改在 broker 停止运行前一直有效；若想 broker 重启后仍然有效，需要将相应的配置写入到配置文件中；

在具有 hot-swappable RAM 的系统中，内存限制会有所不同，when this command is executed without altering the threshold, due to the fact that the total amount of system RAM is queried.

### Disabling all publishing

若设置 `set_vm_memory_high_watermark` 为 0 则会立刻触发内存告警 ，并且令所有的 publishing 行为立即停止（这对于希望实现全局范围内进行 publish 停止功能来说非常有用）；

设置命令为
```shell
rabbitmqctl set_vm_memory_high_watermark 0
```

## Limited Address Space

当在 64 bit 的操作系统（或者 32 bit 带有 PAE 的操作系统）上将 RabbitMQ 运行在 32 bit 的 Erlang VM 中时，可访问的内存是受限的；服务器会检测到这种情况，并记录如下日志信息：

```shell
=WARNING REPORT==== 19-Dec-2013::11:27:13 ===
Only 2048MB of 12037MB memory usable due to limited address space.
Crashes due to memory exhaustion are possible - see
http://www.rabbitmq.com/memory.html#address-space
```

memory alarm 系统是不完美的；尽管停止 publishing 行为通常会阻止任何后续内存使用，但仍有可能存在其他的东东继续内存消耗；当发生这种情况时，通常物理内存会被耗尽，之后操作系统会开始进行 swap 操作；但是当运行在受限地址空间的情况下，内存使用超限将会导致 VM 崩溃；

因此，强烈建议在  64 bit 操作系统上只使用 64 bit Erlang VM ；


## Configuring the Paging Threshold

在 broker 真正得到内存使用上限并阻塞 publish 行为前，会尝试通过 page out queue 中内容到磁盘的方式释放内存占用；page out 行为同时针对 persistent 和 transient 消息（persistent 消息已经存在于磁盘上了，上述操作只是将其从内存中清除干净）；

默认情况下，上述行为开始于 broker 内存使用达到了高水位 50% 的时候；（也就是说，当默认的高水位设置为 0.4 时，当内存使用量达到 20% 后就会触发）；可以通过修改 `vm_memory_high_watermark_paging_ratio` 配置项进行调整，默认值为 0.5 ：

```shell
[{rabbit, [{vm_memory_high_watermark_paging_ratio, 0.75},
         {vm_memory_high_watermark, 0.4}]}].
```

上述配置会在内存使用达到 30% 时开始进行 page 操作，达到 40% 时阻塞 publisher ；

设置 `vm_memory_high_watermark_paging_ratio` 的值大于 1.0 是可以的；在这种情况下，queues 将不会将其内容 page 到磁盘上；在这种情况下，如果 memory alarm 被触发了，producer  会照上面所说的被阻塞；


## Unrecognised platforms

如果 RabbitMQ server 无法识别你的系统，其会在 RABBITMQ_NODENAME.log 文件中附加如下告警信息；并假定系统中安装了 1GB 的 RAM ：

```shell
=WARNING REPORT==== 29-Oct-2009::17:23:44 ===
Unknown total memory size for your OS {unix,magic_homebrew_os}. Assuming memory size is 1024MB.
```

在这种情况下，`vm_memory_high_watermark` 配置值被用作 scale 假定的 1GB RAM 的乘数；若 `vm_memory_high_watermark` 被设置为 0.4 ，RabbitMQ 的内存阈值将被设置为 410MB ，即无论何时 RabbitMQ 使用了超过 410MB 的内存，都会导致 producer 被阻塞；也就是说，当 RabbitMQ 无法识别你的平台时，如果你实际安装了 8GB RAM ，并且你想让 RabbitMQ 在内存使用超过 3GB 时阻塞 producer ，你就可以设置 `vm_memory_high_watermark` 为 3 ；


关于推荐 RAM 水位设置，可以参考 [Production Checklist](http://www.rabbitmq.com/production-checklist.html) ；


----------

官网原文：[这里](http://www.rabbitmq.com/memory.html)