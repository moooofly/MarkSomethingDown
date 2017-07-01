
本文记录 Mac 使用过程中遇到的一些小问题；

----------

# 升级 Mac vim7.3 到 >7.4

原文地址：[这里](https://xjliao.me/2015/06/16/mac-vim73-update-vim74.html)

结论如下：

- 直接升级原生的 Vim7.3（问题：会覆盖原生 vim 及其配置信息）
- 使用 MacVim（问题：不能直接在终端中使用 vim ，使用时会单独开启一个窗口）
- 自己编译（问题：编译过程中需要解决一堆问题，系统中会存在两个版本）
- 通过 brew 进行安装（问题：系统中会存在两个版本）

综上，第四种办法最简单，安装成功后，只需要通过别名处理一下就好；

# pip 安装问题

结论：不必通过 `brew install pip` 安装 `pip` ，因为 `brew install python` 时会直接提供 `pip` 供使用；但一般也不建议通过 `brew install python` 安装 `python` ，因此会导致出现两个 `python` 版本（需要自行调整）；

```shell
➜  ~ brew install pip
Error: No available formula with the name "pip"
Homebrew provides pip via: `brew install python`. However you will then
have two Pythons installed on your Mac, so alternatively you can install
pip via the instructions at:
  https://pip.readthedocs.io/en/stable/installing/
➜  ~
```

通过 brew 安装 pip 提示说，只要安装 python 就会附带 pip ；同时告知系统中已经存在 python 了，如果通过 brew 安装 python 会导致同时存在两个版本；

按照上面链接中的方式，执行如下操作
```shell
wget https://bootstrap.pypa.io/get-pip.py
sudo python get-pip.py
```

之后就可以正常使用 pip 命令安装各种功能 python 包了；

stackoverflow:[installing-pip-on-mac-os-x](http://stackoverflow.com/questions/17271319/installing-pip-on-mac-os-x)


# Mac 快速锁屏

- 系统偏好设置 –> Mission Control –> 触发角；
- 活跃的屏幕角，选择一个角，设置成“将显示器置入睡眠状态”；
- 每次把鼠标移动到那个角上的时候，立即执行了该项动作，类似鼠标手势一样；

# Spotlight 搜索
Spotlight 是 Mac OS X 中非常实用的搜索功能，可以通过 `control+空格` 来快速搜索 Mac 中的内容；

# Mac 必备工具

- [seil](https://pqrs.org/osx/karabiner/seil.html.en)
- [iterm2](https://www.iterm2.com/)
- [tmux](https://tmux.github.io/)
- [powerline/fonts](https://github.com/powerline/fonts)
- [zsh]()



# Mac 上的 root 权限的使用（sudo 和密码问题）

默认情况下，OS X 中的 root 用户处于停用状态。如果需要，请按照[本文](https://support.apple.com/zh-cn/HT204012)中的步骤启用并使用 root 用户。


# Mac 中启用 sshd 服务

```shell
sudo launchctl load -w /System/Library/LaunchDaemons/ssh.plist
```
参考：[这里](https://segmentfault.com/a/1190000001732729)

# Mac 系统环境变量加载顺序

## bash

- /etc/profile
- /etc/paths
- ~/.bash_profile
- ~/.bash_login
- ~/.profile
- ~/.bashrc

当然 /etc/profile 和 /etc/paths 是**系统级别**的，系统启动就会加载；后面几个是当前**用户级**的环境变量。

~/.bashrc 没有上述规则，它是 bash shell 打开的时候载入的。

## zsh

参考《[zsh 使用](https://github.com/moooofly/MarkSomethingDown/blob/master/nonsense/zsh%20%E4%BD%BF%E7%94%A8.md)》；

# Mac 上进行 shell 切换

查看系统当前可用的 shell

```shell
➜  ~ cat /etc/shells
# List of acceptable shells for chpass(1).
# Ftpd will not allow users to connect who are not using
# one of these shells.

/bin/bash
/bin/csh
/bin/ksh
/bin/sh
/bin/tcsh
/bin/zsh
➜  ~
```

确定当前使用的 shell

```shell
➜  ~ echo $SHELL
/bin/zsh
➜  ~
```

进行 shell 切换

```shell
➜  ~ chsh -s /bin/bash
Changing shell for sunfei.
Password for sunfei:
➜  ~
➜  ~ echo $SHELL
/bin/zsh
➜  ~
```

> ⚠️ 此时虽然成功切换，但在当前窗口中不会有效果，需要新开一个窗口才能看到变化！

在新窗口中进行确认

```shell
sunfeideMacBook-Pro:~ sunfei$ echo $SHELL
/bin/bash
sunfeideMacBook-Pro:~ sunfei$
```

# Mac 上的“邮件”不断提示输入密码

最终在 Mac 上通过 web 登录邮箱后，会弹出提示“是否允许其他应用使用该账户“，允许后就没有问题了；    
Mac 官方解答：[这里](https://support.apple.com/zh-cn/HT204187)；

# Mac 上的剪切复制粘贴

- Command + 拖拽 = 剪切
- Option + 拖拽 = 复制
- Command + Option + 拖拽 = 创建快捷方式

在高版本 Mac OS X 中，复制还是 `Command + C` ，粘贴时用 `Command + Option + V` 可以产生剪切＋粘贴效果，也就是复制成功以后删掉原文件；

参考：[这里](http://www.baifeng.me/apple/macosx/2010/04/1295/)；


# Mac 下通过 brew 安装 wireshark

安装命令

```shell
brew install lua
brew install wireshark --with-qt --with-lua --with-libsmi --with-headers
brew cask install wireshark-chmodbpf
```

安装后

```shell
➜  ~ Wireshark -v
Wireshark 2.2.1 (Git Rev Unknown from unknown)

Copyright 1998-2016 Gerald Combs <gerald@wireshark.org> and contributors.
License GPLv2+: GNU GPL version 2 or later <http://www.gnu.org/licenses/old-licenses/gpl-2.0.html>
This is free software; see the source for copying conditions. There is NO
warranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

Compiled (64-bit) with Qt 5.7.0, with libpcap, without POSIX capabilities, with
GLib 2.50.2, with zlib 1.2.5, with SMI 0.4.8, with c-ares 1.12.0, with Lua
5.2.4, with GnuTLS 3.4.16, with Gcrypt 1.7.3, with MIT Kerberos, with GeoIP,
with QtMultimedia, without AirPcap.

Running on Mac OS X 10.11.6, build 15G31 (Darwin 15.6.0), with locale
zh_CN.UTF-8, with libpcap version 1.5.3 - Apple version 54, with GnuTLS 3.4.16,
with Gcrypt 1.7.3, with zlib 1.2.5.
Intel(R) Core(TM) i5-5257U CPU @ 2.70GHz (with SSE4.2)

Built using clang 4.2.1 Compatible Apple LLVM 7.3.0 (clang-703.0.31).
➜  ~
➜  ~ tshark -v
TShark (Wireshark) 2.2.1 (Git Rev Unknown from unknown)

Copyright 1998-2016 Gerald Combs <gerald@wireshark.org> and contributors.
License GPLv2+: GNU GPL version 2 or later <http://www.gnu.org/licenses/old-licenses/gpl-2.0.html>
This is free software; see the source for copying conditions. There is NO
warranty; not even for MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.

Compiled (64-bit) with libpcap, without POSIX capabilities, with GLib 2.50.2,
with zlib 1.2.5, with SMI 0.4.8, with c-ares 1.12.0, with Lua 5.2.4, with GnuTLS
3.4.16, with Gcrypt 1.7.3, with MIT Kerberos, with GeoIP.

Running on Mac OS X 10.11.6, build 15G31 (Darwin 15.6.0), with locale
zh_CN.UTF-8, with libpcap version 1.5.3 - Apple version 54, with GnuTLS 3.4.16,
with Gcrypt 1.7.3, with zlib 1.2.5.
Intel(R) Core(TM) i5-5257U CPU @ 2.70GHz (with SSE4.2)

Built using clang 4.2.1 Compatible Apple LLVM 7.3.0 (clang-703.0.31).
➜  ~
```


# Mac 下如何更新 locate 命令依赖的数据库

在 CentOS 系统上，更新 `locate` 命令依赖的数据库，只需要执行 `updatedb` 命令；

在 Mac OS X 系统中，则需要执行如下命令
```shell
sudo /usr/libexec/locate.updatedb
```
> 注意：如果在执行上述命令时出现权限问题，则可以尝试切到根目录 / 下执行该命令；

locate 命令依赖的数据库位于 `/var/db/locate.database` ；


# .DS_Store 文件是干什么的？如何禁止？

.DS_Store 文件是 Mac OS 中保存文件夹自定义属性的隐藏文件，如文件的图标位置或背景色，相当于 Windows 中的 desktop.ini 。

若想要禁止 .DS_store 文件生成，可以打开“终端”，复制黏贴下面的命令，回车执行，重启 Mac 即可生效。
```shell
defaults write com.apple.desktopservices DSDontWriteNetworkStores -bool TRUE
```
或者执行
```shell
defaults write com.apple.finder AppleShowAllFiles FALSE; killall Finder;
```

若想恢复 .DS_store 文件的生成，则需执行
```shell
defaults delete com.apple.desktopservices DSDontWriteNetworkStores
```

删除系统中存在的所有 .DS_Store 文件
```shell
sudo find / -name ".DS_Store" -depth -exec rm {} \;
```

参考：[这里](http://www.zhihu.com/question/20345704)；


# Bash Completion on OS X With Brew

I live and breath OS X on a daily basis, with a large portion of my work revolving around the command line using mostly tools which I’ve installed with brew. Bash completion has likely saved me days worth of time over the past decade or so. Little did I know, up until recently, however, that there is an official tap with completion scripts (in addition to the ones which come with individual recipes such as git) which can be installed for tools like docker, vagrant and grunt. Using it is dog simple. To start with you’ll want to go ahead and install bash-completion (if you haven’t already) and then tap homebrew/completion to gain access to these additional formulae:

```shell
$ brew install bash-completion
$ brew tap homebrew/completions
```

After you run that first command, in typical brew fashion, it will request that you add the following tidbit to your ~/.bash_profile. Don’t forget this part. It’s critical!

```shell
if [ -f $(brew --prefix)/etc/bash_completion ]; then
    . $(brew --prefix)/etc/bash_completion
fi
```

Once you’ve done this, you’ll be able to install the additional completion scripts. You can find a complete list of these [over here](https://github.com/Homebrew/homebrew-completions) on GitHub. Happy tabbing!

> UPDATED

[Homebrew/homebrew-completions](https://github.com/Homebrew/homebrew-completions) supports `bash` and `fish` completions (only for `brew-cask` right now). Install them with:
```shell
$ brew install homebrew/completions/brew-cask-completion
```

For zsh completion support, simply run:
```shell
$ brew install zsh-completions
```

# Mac 下的 Homebrew 使用

详见《[Homebrew 相关](https://github.com/moooofly/MarkSomethingDown/blob/master/nonsense/Homebrew%20%E7%9B%B8%E5%85%B3.md)》


