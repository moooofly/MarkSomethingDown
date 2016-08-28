# Linux Man Pages in Dash
Oct 21st, 2013

> Dash works great as a man page browser, but I sometimes get requests to make extra docsets containing the man pages of various flavours of Linux.

Dash 作为 man 手册浏览工具已经做的很好了，但是还是经常会收到针对不同种类 Linux 添加额外 docsets 的需求；

> I’ve decided not to pursue these requests, because:

但我最终决定不处理类似的请求，因为：

> - Updates for these docsets would be a nightmare, as man pages change a lot, individually.
> - I’d have to choose which man pages to include and which not to. I’d never be able to guess which obscure man page a user might want.


- man 手册经常会发生变更，且都是独立变更；
- 不得不选择包含哪些 man 页面，以及不包含哪些 man 页面；我无法猜测哪些费解的 man 页面是一个用户可能需要的；

> The current `Man Pages` docset solves these issues by indexing the man pages that are actually on your Mac.

当前提供的 `Man Pages` docset 通过索引你的 Mac 机器上实际安装的 man 手册的方式解决了上述问题；

## The workaround
> You can copy the man pages from any Linux box to your Mac and Dash will index them as part of the regular Man Pages docset.

你可以将任意 Linux 中的 man 手册拷贝到你的 Mac 上，之后 Dash 就能够将其作为 `Man Pages` docset 的一部分进行索引了；

> Step by step instructions:

操作步骤如下：

> - Log into your Linux box
> - Run man -w to list the folders that contain man pages
> - Copy these folders to your Mac
> - Optional, but highly recommended: use a batch renamer to rename all of the man page files to have a common prefix. This will help you differentiate between the default macOS man pages and the Linux ones
> - Move the man pages anywhere on your MANPATH, or any folder from man -w

- 登录到你的 Linux 机器上；
- 运行 `man -w` 查看包含 man 手册页面的文件夹；
- 拷贝这些文件夹到你的 Mac 中；
- 可选的，但却是极力推荐的方式：使用一个批量重命名工具对所有 man 手册文件进行加前缀重命名；这样你才能区分出哪些是 macOS 默认的 man 手册，哪些是 Linux 中的 man 手册；
- 将 man 手册移动到 MANPATH 环境变量中指定的路径里，或者移动到 man -w 命令输出的文件夹里；


That’s it!