# [IRQ Balancing](http://www.ntop.org/pf_ring/irq-balancing/)

在 Linux 中，中断（interrupts）由内核自动进行处理；尤其需要说明的是，还存在一个名为 `irqbalancer` 的进程负责在多个处理器中进行中断的均衡；不幸的是，该默认行为是让全部处理器共同进行中断处理，而结果却是总体性能并非最优，尤其在多核系统中；这是因为拥有多 RX queues 的现代 NICs 在满足缓存一致性（cache coherency）原则前提下，能够工作于最佳状态；这意味着属于 ethX-RX queue Y 的中断必须被发送给 1 个核心或者最多 1 个核心及其相应的 Hyper-Threaded (HT) ；如果多处理器同时针对了相同的 RX queue 进行处理，则会导致缓存（cache）失效，性能下降；由于这个原因，IRQ balancing 成为了性能的关键；这也就是我为何建议针对相同的中断仅使用 1 或 2 个核心（对应 HT 情况）进行处理；由于在 Linux 上中断（默认）通常被发送给全部处理器的现状，即对应 `/prox/irq/X/smp_affinity` 被设置为 `ffffffff` 的情况，正如我上面所说的，最好避免这种全部处理器参与处理所有中断的方式；例如

```shell
~# grep eth /proc/interrupts
191: 0 0 3 0 0 0 2 310630615 454 0 0 0 0 0 2 0 PCI-MSI-edge eth5-rx-3
192: 0 3 0 0 0 2 0 314774529 0 0 0 0 0 2 0 0 PCI-MSI-edge eth5-rx-2
193: 3 0 0 0 2 309832652 454 0 0 0 0 0 2 0 0 0 PCI-MSI-edge eth5-rx-1
194: 0 0 0 2 0 314283930 0 0 0 0 0 2 0 0 0 3 PCI-MSI-edge eth5-TxRx-0
195: 0 0 1 0 1 0 0 0 0 0 1 0 0 0 1 0 PCI-MSI-edge eth5
196: 0 3 0 311226806 0 0 0 0 0 2 0 0 0
```

其中

```
# cat /proc/irq/19[12345]/smp_affinity
00008080
00008080
00002020
00002020
ffffffff
```

该设置允许最大化性能，尤其在使用了 **PF_RING** 和 **TNAPI** 时；当手动调整中断设置时，请去使能 `irqbalancer` ，因为 `irqbalancer` 会按照其自己希望的方式恢复中断处理行为，进而“危害”你所做的努力；

进一步阅读：

- [SMP affinity and proper interrupt handling in Linux](http://www.alexonlinux.com/smp-affinity-and-proper-interrupt-handling-in-linux)
- [Why interrupt affinity with multiple cores is not such a good thing](http://www.alexonlinux.com/why-interrupt-affinity-with-multiple-cores-is-not-such-a-good-thing)