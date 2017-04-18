# insmod 和 modprobe


`insmod` 与 `modprobe` 都是用于载入 kernel 模块的，差别在于 `modprobe` 能够根据配置文件自动处理 module 间的载入依赖问题。

比方你要载入 a module，但是 a module 要求系统要先载入 b module 时，直接用 `insmod` 挂入通常都会出现错误讯息，不过 `modprobe` 倒是能够知道先载入 b module  后才载入 a module，如此依赖关系就会满足。

不过 `modprobe` 获取 module 之间依赖关系是通过读取 `/lib/modules/$(uname -r)/modules.dep` 文件得知的，而该文件的内容是通过 `depmod` 所建立的。若在载入过程中发生错误，在 `modprobe` 会卸载整组的模块。

`modprobe` 在加载模块的时候，会检查模块里是否存在一些 symbols 在内核里没有定义，如果有这样的 symbols ，`modprobe` 函数会搜索其他模块，看其他模块里有没有相关的 symbols ；如果有，则将此模块也一起加载，这样的话，就算模块里有一些没有定义的 symbols 也能成功加载。但如果用 `insmod` 去加载的话，遇到这种情况就会加载失败。会出现 "unresolved symbols" 信息。


在用法上：

- `insmod` 一次只能加载特定的一个设备驱动，且需要驱动的具体地址，如 `insmod /path/to/drv.ko` ；
- `modprobe` 则可以一次将有依赖关系的驱动全部加载到内核。无需指定驱动的具体地址，如 `modprob drv` ，但需要在安装时按照 `make modues_install` 的方式进行驱动模块的安装，即驱动模块被安装在 `/lib/modules/$(uname -r)/xxx` 下。


lsmod 命令可以用来查看当前已经被加载的模块，它是通过读出 `/proc/modules` 这个虚拟文件来实现上述功能的。关于当前已加载模块的更多信息可以在 `/sys/module/<module name>/xxx` 中查看；

