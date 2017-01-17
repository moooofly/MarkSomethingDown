# shell 使用常见问题

标签（空格分隔）： shell

---

## 如何解决 ubuntu 默认不提供 root 用户问题

Ubuntu 默认安装时，并没有给 root 用户设置口令，也没有启用 root 帐户；

设置 root 用户密码（也用于重置密码），即开启 root 帐号

```shell
$ sudo passwd root    # 设置 root 用户密码
$ su - root           # 切换成 root 用户
```

如果要再次禁用 root 帐号（锁定帐户），可以执行

```shell
# sudo passwd -l root
```

如果要再次启动 root 账号（需要以前锁定过，否则无效），可以执行

```shell
sudo passwd -u root
```


## bash v.s. dash

- dash — command interpreter (shell)

> `dash` is the standard command interpreter for the system.  The current
> version of `dash` is in the process of being changed to conform with the
> POSIX 1003.2 and 1003.2a specifications for the shell.  This version
> has many features which make it appear similar in some respects to the
> **Korn** shell, but it is not a **Korn** shell clone (see ksh(1)).  Only fea‐
> tures designated by POSIX, plus a few Berkeley extensions, are being
> incorporated into this shell.


- bash - GNU Bourne-Again SHell

> `Bash` is an `sh`-compatible command language interpreter that executes
> commands read from the standard input or from a file. `Bash` also
> incorporates useful features from the **Korn** and **C** shells (`ksh` and
> `csh`).
>
> `Bash` is intended to be a conformant implementation of the Shell and
> Utilities portion of the IEEE POSIX specification (IEEE Standard
> 1003.1). `Bash` can be configured to be POSIX-conformant by default.

如今 Ubuntu 系统中，`/bin/sh` 默认已经指向 `dash` ，这是一个不同于 `bash` 的 shell ，它主要是为了执行脚本而出现，而不是交互；它速度更快，但功能相比 `bash` 要少很多，语法严格遵守 POSIX 标准；

```shell
root@vagrant-ubuntu-trusty:~# ll /bin/sh
lrwxrwxrwx 1 root root 4 Jan 17 07:14 /bin/sh -> dash*
root@vagrant-ubuntu-trusty:~#
```

如果想切换 `/bin/sh` 的默认指向（不建议修改），则执行
```shell
sudo dpkg-reconfigure dash
```

> 想要进行修改的理由：脚本兼容性问题；因为 `bash` 和 `dash` 对 shell 脚本语法对兼容性有所不同；


## shell 脚本执行问题

> 参考：[这里](http://stackoverflow.com/questions/31473298/i-am-puzzled-with-the-differences-between-sh-xxx-sh-and-xxx-sh-for-running/31473406)以及[这里](https://www.cyberciti.biz/faq/run-execute-sh-shell-script/)

一般情况下，在编写好 shell 脚本后，会按照如下方式执行脚本

```shell
# chmod +x script-name-here.sh
# 如下任选一种
# ./script-name-here.sh        -- 1
# sh script-name-here.sh       -- 2
# bash script-name-here.sh     -- 3
```

- 1 的成功执行，要求 `.sh` 脚本必须有可执行权限；
- 1 和 2 的区别在于：1 执行时，解释器使用的是 "shebang" 行上指定的内容；2 执行时，解释器使用的是 `sh` 所链接到的 shell ；
- 2 和 3 的区别在于：`sh` 可能链接到 `bash` 上，也可能链接到其他 shell 上；



### 使用 `sudo bash xxx.sh` 的原因

一般情况下，通过 shell 脚本安装应用程序时，你会需要 root 权限才能保证安装过程的顺利进行；而 root 权限在许多 Linux 和 UNIX like 系统中默认是不允许直接使用的；

通过 `sudo bash xxx.sh` 方式，可以基于用户自身密码获取 root 权限，以 root shell 运行脚本；

当然，如果知道 root 密码的话，也可以如下操作（不建议这么使用）

```shell
$ su - root
# bash xxx.sh
```


## shell 切换问题

> 参考：[Change Shell To Bash](https://www.cyberciti.biz/faq/how-to-change-shell-to-bash/)

### 确定指定 user 的默认 shell

```shell
$ grep <username> /etc/passwd
```

### 确认当前系统中提供了哪些 shell

```shell
$ cat /etc/shells
```

### 进行不同 shell 切换

简单一句话就是：想使用哪种 shell 就直接输入其名字；

例如从 bash 切到 dash 时，输入

```shell
# dash
```

### 改变指定 user 的默认 shell

可用命令形式

```shell
chsh
chsh -s /bin/bash
chsh -s /bin/bash <username>
```

实际操作

```shell
vagrant@vagrant-ubuntu-trusty:~$ cat /etc/passwd|grep "vagrant"
vagrant:x:1000:1000:,,,:/home/vagrant:/bin/bash
vagrant@vagrant-ubuntu-trusty:~$
vagrant@vagrant-ubuntu-trusty:~$ chsh
Password:
Changing the login shell for vagrant
Enter the new value, or press ENTER for the default
	Login Shell [/bin/bash]: /bin/dash
vagrant@vagrant-ubuntu-trusty:~$
vagrant@vagrant-ubuntu-trusty:~$ cat /etc/passwd|grep "vagrant"
vagrant:x:1000:1000:,,,:/home/vagrant:/bin/dash
vagrant@vagrant-ubuntu-trusty:~$
```

> 若想确定效果，则需要先 logout 再 login ；

### 查找指定 shell 的完整路径

```shell
# type -a bash
```

或者

```shell
# which bash
```


