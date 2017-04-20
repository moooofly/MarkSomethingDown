# PF_RING 杂七杂八

## [How to capture from a bond interface using PF_RING ZC?](http://www.ntop.org/support/faq/how-to-capture-from-a-bond-interface-using-pf_ring-zc/)

由于 PF_RING ZC 是一种内核旁路（kernel-bypass）技术，并且应用程序会直接访问网卡，因此不可能从 bond 设备上进行捕获操作，然而你能够在 ZC 模式下直接从多个接口上聚合（aggregate）traffic ，详见示例 `zbalance_ipc -i zc:ethX,zc:ethY` ；

## [ZC is faster than DNA because ZC does no copy, but DNA does?](http://www.ntop.org/support/faq/zc-is-faster-than-dna-because-zc-does-no-copy-but-dna-does/)

性能是相近的，但基于 Zero-Copy 的 ZC 提供了更高的灵活性（提供了 sw queues, multi-process 和 multi-vm 支持等等）；

## [Do DNA and ZC have any relationship , dependence, or they are completely isolated technologies?](http://www.ntop.org/support/faq/do-dna-and-zc-have-any-relationship-dependence-or-they-are-completely-isolated-technologies-2/)

完全不相关；

## [What is the PF_RING ZC distribution format?](http://www.ntop.org/support/faq/why-is-the-pf_ring-dna-distribution-format/)

PF_RING ZC 由两部分构成：kernel drivers 和 user-space library；内核驱动以源码格式作为 PF_RING 的一部分发布，而用户空间库以二进制格式发布，并且要求 per-MAC licenses ；


## [Modern Packet Capture and Analysis: Multi-Core, Multi-Gigabit, and Beyond](http://www.ntop.org/pf_ring/modern-packet-capture-and-analysis-multi-core-multi-gigabit-and-beyond/)

- 对于每一个想要使用 PF_RING 获取 packet 捕获加速能力的人来说，都应该阅读这个 [pdf 文档](http://luca.ntop.org/IM2009_Tutorial.pdf)；
- 如今，packet 捕获过程中的最大成本耗费已经受限于 packet 分析速度了；
- 由于上述原因，您应该将 PF_RING 作为框架用以创建简单、但强大的 traffic 监控应用程序；




