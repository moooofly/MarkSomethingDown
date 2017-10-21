# Mac 新装

## [安装 homebrew](https://brew.sh/)

> 此步骤可以在安装 iTerm2 后执行；

```
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

建议安装 cask（详见《[Homebrew 相关](https://github.com/moooofly/MarkSomethingDown/blob/master/nonsense/Homebrew%20%E7%9B%B8%E5%85%B3.md)》中的说明）；

```
$ brew install cask
```

## [下载安装 iTerm2](http://iterm2.com/)

安装方法：

- 直接从官网上下载，然后拖放到**应用程序**目录中（推荐）；
- 通过 brew 进行安装；

```
$ brew install cask
$ brew cask install iterm2
```


## ZSH 和 iTerm2 设置

- 调整 zsh 的 ZSH_THEME 配置（默认即可）；
- 调整 iTerm2 配色方案 Color Presets ；
- 调整 Font 和 Non-ASCII Font ；

常见组合：**agnoster 主题 + Solarized 配色方案 + Powerline 字体**


## 针对 iTerm2 的其他调整

### 取消 iTerm2 响铃声

```
Preferences -> Profiles -> Terminal -> silence bell
```


## [安装 zsh](https://github.com/robbyrussell/oh-my-zsh/wiki/Installing-ZSH)

Most versions of macOS ship `zsh` by default, but it's normally an older version. Try `zsh --version` before installing it from `Homebrew`. If it's newer than 4.3.9 you might be OK. Preferably newer than or equal to 5.0.

macOS 默认支持并已安装 zsh ；

```
sunfei@sunfeideMacBook-Pro:~|⇒  cat /etc/shells
# List of acceptable shells for chpass(1).
# Ftpd will not allow users to connect who are not using
# one of these shells.

/bin/bash
/bin/csh
/bin/ksh
/bin/sh
/bin/tcsh
/bin/zsh
sunfei@sunfeideMacBook-Pro:~|⇒
sunfei@sunfeideMacBook-Pro:~|⇒  where zsh
/bin/zsh
sunfei@sunfeideMacBook-Pro:~|⇒  /bin/zsh --version
zsh 5.0.8 (x86_64-apple-darwin15.0)
sunfei@sunfeideMacBook-Pro:~|⇒
```

理论上讲，默认的版本已经足够，如果一定想要升级为最新版本，则可以通过如下命令完成：

```
brew install zsh zsh-completions
```

注意：基于这种方式安装后，新装到 zsh 位于 `/usr/local/bin/zsh` ；此时系统中存在两个版本的 zsh ；

```
sunfei@sunfeideMacBook-Pro:~|⇒  where zsh
/bin/zsh
/usr/local/bin/zsh
```

若想设置 `/usr/local/bin/zsh` 为默认 shell ，首先必须要将其保存到已授权 shell 列表 `/etc/shells` 中；

最后，需要将 zsh 设置为默认 shell ；

```
chsh -s /usr/local/bin/zsh
```

## [oh-my-zsh](https://github.com/robbyrussell/oh-my-zsh/)
```
$ sh -c "$(curl -fsSL https://raw.github.com/robbyrussell/oh-my-zsh/master/tools/install.sh)"
$ sh -c "$(wget https://raw.github.com/robbyrussell/oh-my-zsh/master/tools/install.sh -O -)"
```


## [开启 zsh 插件](https://github.com/robbyrussell/oh-my-zsh/wiki/Plugins-Overview)

找到 `~/.zshrc` 有一行 `plugins=(git)`，想加什么插件就把名字放里面就是了，比如 `plugins=(rails git ruby)` 就开启了 rails，git 和 ruby 三个插件。

更多默认自带插件请进入 `~/.oh-my-zsh/plugins/` 文件夹探索（自定义插件位于 `~/.oh-my-zsh/custom/plugins/` 目录），也可以看看 wiki 里的 **Plugins Overview** ，每个人的需求不一样，里面有一些比较神奇的插件，比如

- 敲两下 esc 它会给你自动加上 sudo 的 sudo 插件；
- 让复制显示进度条的 cp 插件；
- 解压用的 extract 插件（有没有觉得在命令行下敲一大堆选项才能解压有点奇怪？）；
- vi 粉的 vi-mode ；
- 等等...

> TODO: 需要研究一下

配置示例

```
plugins=(bundler git git-flow gnu-utils osx ruby gem perl rails rvm mercurial svn macports osx virtualenvwrapper django pip) 
```

我的配置

```
plugins=(git zsh-autosuggestions z)
```


### Z

嗯，这也是个自带的但是没有开启的插件。为什么单独把它拿出来讲呢？因为太好用了，没有它我根本就不想用命令行。简直就是可以**无脑跳跃**，比如你经常进入 `~/Documents` 目录，按下 **z doc** 一般它就可以跳进去了（当然首先你得用一段时间让它积累一下数据才能用）。类似的插件还有好几个，比如 `autojump` ，`fasd` ，这类东西好像叫 FS Jumping ；

### [zsh-autosuggestions](https://github.com/tarruda/zsh-autosuggestions)

官方介绍："Fish-like fast/unobtrusive autosuggestions for zsh. It suggests commands as you type, based on command history." ；

没错，这是模仿 `fish shell` 的一个插件，作用基本上就是**根据历史记录即时提示**。没有这个东西让我感觉自己很盲目。没有用过 fish 的同学可能觉得它有点奇怪，但是一旦适应它以后就会发现它会大幅度的提高效率（按 ctrl+E 是正确姿势）；

安装方法：

- 手动安装

假定下载到 `~/.zsh/zsh-autosuggestions` ；

```
git clone git://github.com/zsh-users/zsh-autosuggestions ~/.zsh/zsh-autosuggestions
```

在 `.zshrc` 中添加

```
source ~/.zsh/zsh-autosuggestions/zsh-autosuggestions.zsh
```

- 基于 oh-my-zsh 使用

将目标仓库 clone 到 `$ZSH_CUSTOM/plugins` (默认为 `~/.oh-my-zsh/custom/plugins`)

```
git clone git://github.com/zsh-users/zsh-autosuggestions $ZSH_CUSTOM/plugins/zsh-autosuggestions
```

添加该插件添加到 oh-my-zsh 的插件列表中以便加载：

```
plugins=(zsh-autosuggestions)
```

------

其它：

- [My Extravagant Zsh Prompt](http://stevelosh.com/blog/2010/02/my-extravagant-zsh-prompt/)


## Vim 安装

安装

```
brew install vim
```

查看

```
$ where vim
/usr/local/bin/vim
/usr/bin/vim

$ where vi
/usr/bin/vi
```

在 `.zshrc` 中通过 alias 解决 vi 问题引用老版本问题：

```
alias vi="/usr/local/bin/vim"
```

在安装 vim 的时候自动安装 perl/libyaml/ruby/python 几个依赖项；

> 在安装 python 相关内容时会自动安装 `pip` 和 `setuptools` ；建议顺便执行 `pip install --upgrade pip setuptools` 进行升级（也可以后续执行）；

## golang 安装

```
brew install go
```

安装后，在 `.zshrc` 中设置好 `GOPATH` 变量（默认为 `$HOME/go` 目录）和 `PATH` 变量；

常见设置为（linux 系统中具有 root 权限时的玩法）

```
export GOPATH="/go"
export PATH="$GOPATH/bin:/usr/local/bin:$PATH"
```

在 mac 环境下，由于 /go 需要 root 权限才能访问，所以建议设置为：

```
export GOPATH="$HOME/go"
export PATH="$GOPATH/bin:/usr/local/bin:$PATH"
```

另外，在 go1.8 版本中，针对 GOPATH 设置问题有了新变化，详见 [The default GOPATH](https://rakyll.org/default-gopath/) ；

简单摘录如下：

- Go 1.8 will set a **default GOPATH** if the GOPATH env variable is not set.
- **Default GOPATH** is:
    - `$HOME/go` on Unix-like systems
    - `%USERPROFILE%\go` on Windows
- Users still need to add `$GOPATH/bin` into their `PATH` to run binaries installed by `go get` and `go install`.
- The users who are developing with the Go language still need to understand that the presence of `GOPATH`, its location and its structure.
- If your `GOROOT` (the location where you checkout the Go’s source code) is the **default GOPATH** and if you don’t have a `GOPATH` set, the tools will reject to use the **default GOPATH** not to corrupt your GOROOT.
- You still may prefer to [set a custom GOPATH](https://github.com/golang/go/wiki/SettingGOPATH) if the default is not working for you.



## 基于 Vim 构建 golang 开发环境

详见：[Golang 开发环境搭建之 Vim 篇](https://my.oschina.net/moooofly/blog/1036706)

可能需要翻墙安装的内容（时好时坏）

```
go get golang.org/x/tools/cmd/goimports
go get golang.org/x/tools/cmd/guru
go get golang.org/x/tools/cmd/gorename
go get github.com/golang/lint/golint
go get github.com/kisielk/errcheck
go get github.com/zmb3/gogetdoc
go get github.com/josharian/impl
go get -u github.com/haya14busa/gopkgs/cmd/gopkgs
go get -u github.com/dominikh/go-tools/cmd/keyify
```

## .vimrc 配置

参考：

- [Golang 开发环境搭建之 Vim 篇](https://my.oschina.net/moooofly/blog/1036706)
- [.vimrc 配置梳理](https://my.oschina.net/moooofly/blog/983252)


## 重映射 Caps Lock 键

### 为何要重映射 Caps Lock 键

**A popular `caps lock` remap for `vi` users is `escape`.** Many people find it awkward to constantly reach for the upper left corner of the keyboard for the `escape` key. 

### 重映射 Caps Lock 键的最简单办法

系统偏好设置 => 键盘 => 修饰键 => 将**大写锁定键**改为 Escape

> 采用下面 Karabiner-Elements 进行设置的好处在于，可以进行更多个性化设置（比如外接键盘的使用）；

----------


### Seil

- github 地址：[这里](https://github.com/tekezo/Seil/)
- 官方地址：[这里](https://pqrs.org/osx/karabiner/seil.html)

This application is for **older** macOS. Please use ·Karabiner-Elements· since macOS Sierra.

> 已过时


### Karabiner

- Karabiner github 地址：[这里](https://github.com/tekezo/Karabiner)
- Karabiner-Elements github 地址：[这里](https://github.com/tekezo/Karabiner-Elements)
- 官方地址：[这里](https://pqrs.org/osx/karabiner/)

A powerful and stable keyboard customizer for OS X.

- We made new `Karabiner` as `Karabiner-Elements-11.0.0` from scratch due to kernel architecture changes in macOS Sierra.
- Prior to version 9.3.0, `Karabiner` was called `KeyRemap4MacBook`.
- `Karabiner` and `KeyRemap4MacBook` work with all Mac products, including the MacBook family, iMac, Mac mini, and Mac Pro.

注意：

- Karabiner 对应的是 KeyRemap4MacBook ，用于老版本 macOS
- Karabiner-Elements 用于新版本 macOS


其他参考：

- [Remap Caps Lock](http://wiki.c2.com/?RemapCapsLock)
- [Remapping Caps Lock](http://www.drbunsen.org/remapping-caps-lock/)
- [Map Keys in Mac](http://www.legendu.net/en/blog/map-keys-in-mac/)
- [9 Enhancements to Shell and Vim Productivity](https://danielmiessler.com/blog/enhancements-to-shell-and-vim-productivity/#gs.nhdi9CI)
- [提高 Vim 和 Shell 效率的 9 个建议](http://blog.jobbole.com/86809/)


## 重映射 Esc 键

貌似不处理这个也 ok ；

## meta 键

为什么需要 Meta 键？

> 在 Emacs Meta 键的使用非常频繁，而 OSX 系统中却没有提供 Meta 键。在 iTerm2 中可以选择将左右两个 Option 键中的一个作为 Meta 键。官方推荐使用左边的 Option 键作为 Meta ，右边的 Option 键依然保留 OSX 的默认功能（输入特殊字符）。


----------

如何将 `Option` 键配置成 Meta 键？

Q: How do I make the `option`/`alt` key act like `Meta` or send **escape codes**?

A: Go to `Preferences->Profiles` tab. Select your **profile** on the left, and then open the **Keyboard** tab. At the bottom is a set of buttons that lets you select the behavior of the `Option` key. For most users, `+Esc` will be the best choice.


----------


在 Mac 下默认是没有 Meta 键的，可通过如下方法进行修改：

- 若是系统自带的 terminal ，在`设置 -> 键盘选项`中，将“使用 option 键作为 Meta 键”选中；
- 若是在 iTerm2 下，需要在 `Profiles -> Keys` 中，将 "Left option Key" 的单选按钮修改为 "+Esc" 即可。


## 安装 chrome 浏览器

```
$ brew cask info google-chrome
google-chrome: latest
https://www.google.com/chrome/
Not installed
From: https://github.com/caskroom/homebrew-cask/blob/master/Casks/google-chrome.rb
==> Name
Google Chrome
==> Artifacts
Google Chrome.app (app)

$ brew cask install google-chrome
```


## 安装 sublime-text

> 注意：sublime 和 sublime-text 不是同一个东西；

```
$ brew cask info sublime-text
sublime: latest
https://www.salukistudios.com/sublime/
Not installed
From: https://github.com/caskroom/homebrew-cask/blob/master/Casks/sublime.rb
==> Name
Sublime
==> Artifacts
Sublime.app (app)

$ brew cask install sublime-text
```

成功安装后，可以通过如下命令进行确认：

```
$ subl --version
Sublime Text Build 3126
```

## 安装 Alfred

```
brew cask install alfred
```

问题：

- 该版本是否有限制
- powerpack 是否需要安装
- 默认设置调整


## 安装 Wireshark

```
brew install lua
brew install wireshark --with-qt5 --with-lua --with-libsmi --with-headers
brew cask install wireshark-chmodbpf
```

补充说明：

```
...

--with-headers
	Install Wireshark library headers for plug-in development
--with-libsmi
	Build with libsmi support
--with-lua
	Build with lua support
--with-qt
	Build the wireshark command with Qt (can be used with or without either GTK option)

...

If your list of available capture interfaces is empty
(default macOS behavior), try installing ChmodBPF from homebrew cask:

  brew cask install wireshark-chmodbpf

This creates an 'access_bpf' group and adds a launch daemon that changes the
permissions of your BPF devices so that all users in that group have both
read and write access to those devices.
```

## 安装 ctags

mac 上默认存在一个 ctags ，但不符合要求（连版本信息都没有）

```
➜  ~ where ctags
/usr/bin/ctags
➜  ~ ctags --version
/Library/Developer/CommandLineTools/usr/bin/ctags: illegal option -- -
usage: ctags [-BFadtuwvx] [-f tagsfile] file ...
➜  ~ ctags -V
/Library/Developer/CommandLineTools/usr/bin/ctags: illegal option -- V
usage: ctags [-BFadtuwvx] [-f tagsfile] file ...
```

因此，需要通过如下命令安装个最新版本

```
➜  ~ brew info ctags
ctags: stable 5.8 (bottled), HEAD
Reimplementation of ctags(1)
https://ctags.sourceforge.io/
Not installed
From: https://github.com/Homebrew/homebrew-core/blob/master/Formula/ctags.rb
==> Caveats
Under some circumstances, emacs and ctags can conflict. By default,
emacs provides an executable `ctags` that would conflict with the
executable of the same name that ctags provides. To prevent this,
Homebrew removes the emacs `ctags` and its manpage before linking.

However, if you install emacs with the `--keep-ctags` option, then
the `ctags` emacs provides will not be removed. In that case, you
won't be able to install ctags successfully. It will build but not
link.
➜  ~
➜  ~ brew install ctags
➜  ~
➜  ~ ctags --version
Exuberant Ctags 5.8, Copyright (C) 1996-2009 Darren Hiebert
  Compiled: Sep 13 2016, 04:58:37
  Addresses: <dhiebert@users.sourceforge.net>, http://ctags.sourceforge.net
  Optional compiled features: +wildcards, +regex
➜  ~
```

具体使用

```
# ctags -R .
```


## 安装 fzf

fzf 项目由以下组件构成：

- fzf 可执行程序
- fzf-tmux 脚本用于在 tmux pane 中启动 fzf
- Shell 扩展
    - Key bindings (`CTRL-T`, `CTRL-R`, and `ALT-C`) (`bash`, `zsh`, `fish`)
    - Fuzzy auto-completion (`bash`, `zsh`)
- Vim/Neovim plugin


直接通过 git 下载安装：

```
git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf
~/.fzf/install
```

在 OS X 上可以通过 Homebrew 安装：

```
brew install fzf

# Install shell extensions
/usr/local/opt/fzf/install
```

还可以作为 Vim 插件安装

```
# 手动添加相应的目录到 &runtimepath 中

" If installed using git
set rtp+=~/.fzf

" If installed using Homebrew
set rtp+=/usr/local/opt/fzf
```

若使用了 Vundle 插件管理器，则在 `.vimrc` 中指定 `Plugin 'junegunn/fzf'` ，之后通过 `:PluginInstall` 安装后，再执行如下命令

```
cd /root/.vim/bundle/fzf
./install --all
cp /root/.vim/bundle/fzf/bin/* /bin  # 该命令不执行也可以
```

## docker

详见《[Docker for Mac](https://github.com/moooofly/MarkSomethingDown/blob/master/Docker/Docker%20for%20Mac.md)》

## vagrant

```
brew cask install vagrant
```

## virtualbox

```
brew cask install virtualbox
```

## IntelliJ IDEA


## tmux

- [A tmux Crash Course](https://robots.thoughtbot.com/a-tmux-crash-course)
- [Love, Hate, & tmux](https://robots.thoughtbot.com/love-hate-tmux)
- [A tmux Primer](https://danielmiessler.com/study/tmux/)


### 为什么要选择 tmux

- 掉线仍然能保证工作现场；
- tmux 完全使用键盘控制窗口，实现窗口的切换；

```
brew install tmux
```

### 快捷键

bind-key    -T prefix "                split-window
bind-key    -T prefix %                split-window -h
bind-key    -T prefix s                split-window -v
bind-key    -T prefix v                split-window -h
bind-key    -T root   M-e              split-window -v
bind-key    -T root   M-q              split-window -h


```
Ctrl+b  // 激活控制台；

系统操作 
?   // 列出所有快捷键；按q返回 
d   // 脱离当前会话；这样可以暂时返回Shell界面，输入tmux attach能够重新进入之前的会话 
D   // 选择要脱离的会话；在同时开启了多个会话时使用 
Ctrl+z  // 挂起当前会话 
r   // 强制重绘未脱离的会话 
s   // 选择并切换会话；在同时开启了多个会话时使用 
:   // 进入命令行模式；此时可以输入支持的命令，例如kill-server可以关闭服务器 
[   // 进入复制模式；此时的操作与vi/emacs相同，按q/Esc退出 
~   // 列出提示信息缓存；其中包含了之前tmux返回的各种提示信息 

窗口操作 
c   // 创建新窗口 
&   // 关闭当前窗口 
数字键 // 切换至指定窗口 
p   // 切换至上一窗口 
n   // 切换至下一窗口 
l   // 在前后两个窗口间互相切换 
w   // 通过窗口列表切换窗口 
,   // 重命名当前窗口；这样便于识别 
.   // 修改当前窗口编号；相当于窗口重新排序 
f   // 在所有窗口中查找指定文本 

面板操作 
”   // 将当前面板平分为上下两块 
%   // 将当前面板平分为左右两块 
x   // 关闭当前面板 
!   // 将当前面板置于新窗口；即新建一个窗口，其中仅包含当前面板 
Ctrl+方向键    // 以1个单元格为单位移动边缘以调整当前面板大小 
Alt+方向键 // 以5个单元格为单位移动边缘以调整当前面板大小 
Space   // 在预置的面板布局中循环切换；依次包括even-horizontal、even-vertical、main-horizontal、main-vertical、tiled 
q   // 显示面板编号 
o   // 在当前窗口中选择下一面板 
方向键 // 移动光标以选择面板 
{   // 向前置换当前面板 
}   // 向后置换当前面板 
Alt+o   // 逆时针旋转当前窗口的面板 
Ctrl+o  // 顺时针旋转当前窗口的面板
```


Now a `Ctrl-b` options reference:

- **Basics**
    - `?` get help
- **Session management**
    - `s` list sessions
    - `$` rename the current session
    - `d` detach from the current session
- **Windows**
    - `c` create a new window
    - `,` rename the current window
    - `w` list windows
    - `%` split horizontally
    - `"` split vertically
    - `n` change to the next window
    - `p` change to the previous window
    - `0` to `9` select windows 0 through 9
- **Panes**
    - `%` create a horizontal pane
    - `"` create a vertical pane
    - `h` move to the left pane. *
    - `j` move to the pane below *
    - `l` move to the right pane *
    - `k` move to the pane above *
    - `q` show pane numbers
    - `o` toggle between panes
    - `}` swap with next pane
    - `{` swap with previous pane
    - `!` break the pane out of the window
    - `x` kill the current pane
- **Miscellaneous**
    - `t` show the time in current pane

> 更详细的 tmux cheetsheet 见[这里](https://www.tmuxcheatsheet.com/)；


char  digraph   hex     dec     official name
←       <-      2190    8592    LEFTWARDS ARROW
↑       -!      2191    8593    UPWARDS ARROW
→       ->      2192    8594    RIGHTWARDS ARROW
↓       -v      2193    8595    DOWNWARDS ARROW




### 脚本

其中 xdev 为自行实现的脚本

```
[#1#root@dockermonitor ~]$cat `which xdev`
tmux new-session -d -n 'vim' -s $1

tmux new-window -n 'runner' -t $1:3
tmux new-window -n 'tester' -t $1:4
tmux new-window -n 'logger' -t $1:5

tmux attach-session -t $1
[#2#root@dockermonitor ~]$
```



## 其他

- 安装迅雷

```
brew cask install thunder
```

- 安装 QQ

```
brew cask install qq
```

- 安装 unrar

```
brew install unrar
```

- 安装 lrzsz

```
brew install lrzsz
```

使用方法：[moooofly/iterm2-zmodem](https://github.com/moooofly/iterm2-zmodem)

- 安装 mtr

```
brew install mtr
```

安装后查看

```
➜  ~ brew list mtr
/usr/local/Cellar/mtr/0.92/sbin/mtr
/usr/local/Cellar/mtr/0.92/sbin/mtr-packet
/usr/local/Cellar/mtr/0.92/share/bash-completion/completions/mtr
/usr/local/Cellar/mtr/0.92/share/man/ (2 files)
➜  ~
```

发现 mtr 命令没有在 `/usr/local/bin/` 目录下软链接，手动创建如下

```
➜  ~ ln -s /usr/local/Cellar/mtr/0.92/sbin/mtr /usr/local/bin/mtr
➜  ~ ln -s /usr/local/Cellar/mtr/0.92/sbin/mtr-packet /usr/local/bin/mtr-packet
```

之后可以使用了

```
sudo mtr github.com
```




----------


