


# 浏览目录

直接执行 `vim <目录名>`
进入 vim 后执行 `:E`

两种方式存在一定区别：
第一种是基于 NERD tree 实现的；
第二种是基于 Netrw Directory Listing 实现的；


----------


# 缓冲区 buffer 相关命令

查看位于缓冲区中的文件 `:ls` ；
切换缓冲区文件显示 `:buffer <num>` 或 `:buffer <Path/to/FileName>`；

切换缓冲区中下一个文件 `:bnext` 缩写 `:bn` ；
切换缓冲区中前一个文件 `:bprevious` 缩写 `:bp` ；
切换缓冲区中最后一个文件 `:blast` 缩写 `:bl` ；
切换缓冲区中第一个文件 `:bfirst` 缩写 `:bf` ；

buffer 中用于标示文件状态的标记：

```shell
– （非活动的缓冲区）
a （当前被激活缓冲区）
h （隐藏的缓冲区）
% （当前的缓冲区）
# （交换缓冲区）
= （只读缓冲区）
+ （已经更改的缓冲区）
```

注意：
vim 命令支持 Tab 键补全；
vim 命令支持命令缩写（同 gdb 模式）；


----------

# 窗口分割（分屏）

### 命令行分屏启动 vim

垂直分屏 `vim -On file1 file2 file3`
水平分屏 `vim -on file1 file2 file3`
其中 n 为分屏数目，不指定 n 则根据文件数目进行分屏；

> 注意：
执行 `vim file1 file2 file3` 时，只打开了一个窗口，file1 显示在窗口中，file1, file2 和 file3 均存在于 buffer 中；
执行 `vim -O file1 file2 file3` 时，会以分屏方式打开 3 个窗口，分别显示 3 个文件，与此同时，3 个文件也存在于 buffer 中；

`ctrl+W C` 关闭光标所在分屏窗口（注意：窗口中显示的文件仍旧存在于 buffer 中）；

### 基于文件进行分屏

`ctrl+W S` 上下分割当前打开的文件；
`ctrl+W V` 左右分割当前打开的文件；
`:sp filename` 上下分割，并打开一个新的文件；
`:vsp filename` 左右分割，并打开一个新的文件；


### 基于目录进行分屏

在下边分屏浏览目录 `:Hexplore` 缩写 `:He`
在上边分屏浏览目录 `:Hexplore!` 缩写 `:He!`
在左边分屏浏览目录 `:Vexplore` 缩写 :`Ve`
在右边分屏浏览目录 `:Vexplore!` 缩写 :`Ve!`

### 在分屏间移动光标

若想在分屏间移动光标，先按 `ctrl + W` ，再按方向键：
h 向左侧屏移动
j 向下侧屏移动
k 向上侧屏移动
l 向右侧屏移动
w 向下一个屏移动

### 分屏同步移动

要让两个分屏中的文件同步移动，需要到待同步移动的两个屏中都输入 `:set scb` 命令，解开同步移动输入 `:set scb!` 命令；
注：`:set scb` 是 `:set scrollbind` 的缩写；


### 分屏尺寸调整

`ctrl + W <` 调整宽度（需要较新版本支持）
`ctrl + W >` 调整宽度（需要较新版本支持）
`ctrl + W =` 让所有屏一样高度
`ctrl + W +` 增加光标所在屏高度
`ctrl + W -` 减少光标所在屏高度

### 分屏转 Tab 页

`ctrl + W T` 将当前分屏窗口转成 Tab 页；

> 注意：
>> - 上面的 T 为大写字母；
>> - 上面的命令，只能将光标所在分屏窗口转成 Tab 页，若存在多个分屏窗口，则需要执行多次上述命令；
>> - 标签转换为分屏[据说](http://stackoverflow.com/questions/14688536/move-adjacent-tab-to-split)没有内置的解法；


----------


# Tab 页浏览目录


### 分 Tab 页启动 vim




若不喜欢上面的分屏方式，还可以采取分页方式：
`:Texplorer` 缩写 `:Te` ；

若想在多个 Tab 页中切换，只需在 normal 模式下执行如下命令（注意，下述命令不需要加冒号）：
切换到下一个 Tab 页 `gt` ；
切换到前一个 Tab 页 `gT` ；
切换到指定 Tab 页 `{i} gt` ，其中 i 是指定页对应的数字，从 1 开始计数；例如 5 gt 就是到第 5 页；
移动 Tab 页到目标位置 `:tabmove {i}` ，其中 i 是指 Tab 页目标位置，位置从 0 开始计数；
查看已打开 Tab 页情况 `:tabs`
关闭指定 Tab 页 `:tabclose {i}` ，若未指定 i 则关闭当前页；
在 Shell 命令行下，可以使用 vim 的 -p 参数以 Tab 页方式打开多个文件：`vim -p a.c b.c c.c`
将 buffer 中的文件全部转成 Tab 页中的文件执行 `:bufdo tab split`


----------

# 保存会话


----------

# Quickfix


`:cp` 跳到上一个错误
`:cn` 跳到下一个错误
`:cl` 列出所有错误
`:cc` 显示错误详细信息


----------

# 关键字补全

在 insert 模式下，可以按 `Ctrl + N`（Vim 会开始搜索当前目录下的代码，搜索完成后会出现一个下拉列表），再按 `Ctrl + P`，之后就可以通过上下光标键来进行选择；

更多的补齐，都在 Ctrl +X 下面：

Ctrl + X 和 Ctrl + D 宏定义补齐
Ctrl + X 和 Ctrl + ] 是Tag 补齐
Ctrl + X 和 Ctrl + F 是文件名 补齐
Ctrl + X 和 Ctrl + I 也是关键词补齐，但是关键后会有个文件名，告诉你这个关键词在哪个文件中
Ctrl + X 和 Ctrl +V 是表达式补齐
Ctrl + X 和 Ctrl +L 这可以对整个行补齐

----------


# 小技巧


`guu` 或 `Vu` 可以把一整行字符变成全小写；
`gUU` 或 `VU` 可以把一整行字符变成全大写；
按 `v` 键进入 virtual 模式，然后移动光标选择目标文本，再按 `u` 转小写，或者按 `U` 转大写；
`ga` 查看光标处字符的 **ascii** 码；
`g8` 查看光标处字符的 **utf-8** 编码；
`gf` 打开光标处所指的文件（可用于打开 #include 头文件，但必须有路径）
`*` 或 `#` 在当前文件中向下或向上搜索光标位置的单词；

`Ctrl + O` 向后回退光标移动；
`Ctrl + I` 向前追赶光标移动；

----------


# 参考


[无插件Vim编程技巧](http://coolshell.cn/articles/11312.html)
[Vim的分屏功能](http://coolshell.cn/articles/1679.html)
[]()
[]()
