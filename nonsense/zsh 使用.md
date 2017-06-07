# zsh 使用

## [Startup Files](http://zsh.sourceforge.net/Intro/intro_3.html)

> 取自《[An Introduction to the Z Shell](http://zsh.sourceforge.net/Intro/intro_toc.html#SEC3)》；

存在 5 个 startup 文件供 `zsh` 进行命令读取：

```
$ZDOTDIR/.zshenv
$ZDOTDIR/.zprofile
$ZDOTDIR/.zshrc
$ZDOTDIR/.zlogin
$ZDOTDIR/.zlogout
```

若 `ZDOTDIR` 未设置，则使用 `HOME` ；这也是最常见情况；

- 只要 shell 被启动，'`.zshenv`' 就会被 source ，除非指定了 `-f` 选项；该文件中应该包含用于设置**命令搜索路径**的命令，以及其它重要的环境变量；'`.zshenv`' 中不应该包含产生输出内容的命令，也不应该假定 shell 被附着在（attached）在 tty 上；
- 当启动的是**交互式 shell** 时，'`.zshrc`' 就会被 source ；该文件中应该包含用于设置 **aliases**, **functions**, **options**, **key bindings** 等内容的命令；
- 当启动的是**登录 shell** 时，'`.zlogin`' 就会被 source ；因此该文件中应该包含只在登录 shell 中才会执行的命令；
- '`.zlogout`' 在退出登录 shell 时被 source ；
- '`.zprofile`' 和 '`.zlogin`' 类似，除了其会在 '`.zshrc`' 之前被 source ；'`.zprofile`' 存在的意义在于对 `ksh` 粉来说，其可以作为 '`.zlogin`' 的一种等价替代；这两者不应被同时使用，尽管确实可以这么做；
- '`.zlogin`' 中不应该被放入 alias definitions, options, environment variable settings 等内容；作为一个通用规则，不应该通过该文件改变任何 shell environment ；更进一步，它应被用于设置 terminal 类型，以及运行一系列外部命令；


----------


## [The Zsh Startup Files](https://www-s.acm.illinois.edu/workshops/zsh/startup_files.html)

> 取自《[Zsh Workshop: Table of Contents](https://www-s.acm.illinois.edu/workshops/zsh/toc.html)》；

理解处理顺序很重要；理解什么条件下哪些文件内容会被忽略很重要；

> Like most shells, `zsh` processes a number of **system** and **user** startup files. It is very important to understand the order in which each file is read, and what conditions cause various files to be ignored. Otherwise, you may be entering commands and configurations into your startup files that aren't getting seen or executed by `zsh`.

### The Startup Process

In the below description, `zsh` looks for user startup files in the user's **home directory by default**. To make it look in another directory, set the parameter `ZDOTDIR` to where you'd like `zsh` to look.

When `zsh` starts up, the following files are read (in order):

- First, read `/etc/zshenv`
    If the `RCS` option is unset in this system file, all other startup files are skipped. (Can you say 'B O F H' ? I knew you could.)

- Next, read `~/.zshenv`
- Next, IF the shell is a **login shell**, read `/etc/zprofile`, followed by `~/.zprofile`
- Next, IF the shell is **interactive**, read `/etc/zshrc`, followed by `~/.zshrc`
- Finally, if the shell is a **login shell**, read `/etc/zlogin`, followed by `~/.zlogin`

### Logging Out

When a user logs out, `/etc/zlogout` is read, followed by `ZDOTDIR/.zlogout`.

### What do the terms mean?

A **login shell** is generally one that is spawned at login time. (IE, by either `/bin/login` or some other daemon that handles incoming connections). If you `telnet`, `rlogin`, `rsh`, or `ssh` to a host, chances are you have a **login shell**.

An **interactive shell** is one in which prompts are displayed and the user types in commands to the shell. (IE, a `tty` is associated with the shell)

For example, if I run the command

```
ssh SOME_HOST some_command
```

I will be running (presumably) a non-interactive program called `some_command`. This means that `zsh` will not be an **interactive shell**, and ignore the corresponding files. `Zsh` should, however, be a **login shell**, and read the appropriate files.

### Note

Another directory besides `/etc` can be used for the global files. This is determined during the installation of `zsh`.


----------


## [What should/shouldn't go in .zshenv, .zshrc, .zlogin, .zprofile, .zlogout?](http://unix.stackexchange.com/questions/71253/what-should-shouldnt-go-in-zshenv-zshrc-zlogin-zprofile-zlogout)

Here is a non-exclusive list of what each file tends to contain:

- Since `.zshenv` is always sourced, it often contains exported variables that should be available to other programs. For example, `$PATH`, `$EDITOR`, and `$PAGER` are often set in `.zshenv`. Also, you can set `$ZDOTDIR` in `.zshenv` to specify an alternative location for the rest of your `zsh` configuration.
- `.zshrc` is for **interactive shell** configuration. You set options for the interactive shell there with the `setopt` and `unsetopt` commands. You can also load shell modules, set your history options, change your prompt, set up zle and completion, et cetera. You also set any variables that are only used in the interactive shell (e.g. `$LS_COLORS`).
- `.zlogin` is sourced on the start of a **login shell**. This file is often used to start X using `startx`. Some systems start X on boot, so this file is not always very useful.
- `.zprofile` is basically the same as `.zlogin` except that it's sourced directly before `.zshrc` is sourced instead of directly after it. According to the `zsh` documentation, "`.zprofile` is meant as an alternative to `.zlogin' for `ksh` fans; the two are not intended to be used together, although this could certainly be done if desired."
- `.zlogout` is sometimes used to clear and reset the terminal.

You should go through [the configuration files of random Github users](https://github.com/search?q=zsh+dotfiles&ref=commandbar) to get a better idea of what each file should contain.


----------


## [【转】Mac OS X 中 Zsh 下 PATH 环境变量的正确设置](http://www.cnblogs.com/sdlypyzq/p/5001037.html)

```
+--------------+-------------+-------------+-----------------+-----------------+
|              |    login    |  non-login  |      login      |    non-login    |
|              | interactive | interactive | non-interactive | non-interactive |
+--------------+-------------+-------------+-----------------+-----------------+
|/etc/zshenv   |     A       |      A      |        A        |        A        |
+--------------+-------------+-------------+-----------------+-----------------+
|~/.zshenv     |     B       |      B      |        B        |        B        |
+--------------+-------------+-------------+-----------------+-----------------+
|/etc/zprofile |     C       |             |        C        |                 |
+--------------+-------------+-------------+-----------------+-----------------+
|~/.zprofile   |     D       |             |        D        |                 |
+--------------+-------------+-------------+-----------------+-----------------+
|/etc/zshrc    |     E       |      C      |                 |                 |
+--------------+-------------+-------------+-----------------+-----------------+
|~/.zshrc      |     F       |      D      |                 |                 |
+--------------+-------------+-------------+-----------------+-----------------+
|/etc/zlogin   |     G       |             |        E        |                 |
+--------------+-------------+-------------+-----------------+-----------------+
|~/.zlogin     |     H       |             |        F        |                 |
+--------------+-------------+-------------+-----------------+-----------------+
|              |             |             |                 |                 |
+--------------+-------------+-------------+-----------------+-----------------+
|              |             |             |                 |                 |
+--------------+-------------+-------------+-----------------+-----------------+
|~/.zlogout    |     I       |             |        G        |                 |
+--------------+-------------+-------------+-----------------+-----------------+
|/etc/zlogout  |     J       |             |        H        |                 |
+--------------+-------------+-------------+-----------------+-----------------+
```


----------


## [oh-my-zsh小记](https://segmentfault.com/a/1190000004695131?hmsr=toutiao.io)

> 此文讲述如何从 bash 转为 zsh ；


- 安装 zsh
- 安装 oh-my-zsh

```
# via wget
sh -c "$(wget https://raw.githubusercontent.com/robbyrussell/oh-my-zsh/master/tools/install.sh -O -)"

# via curl
sh -c "$(curl -fsSL https://raw.githubusercontent.com/robbyrussell/oh-my-zsh/master/tools/install.sh)"
```

- 配置信息搬移

我之前一直使用的是 bash ，有部分配置保存在相应的配置文件中，我的是在 `~/.bash_profile` 中，使用 `vim ~/.bash_profile` 进入编辑模式，把里面的个人配置拷贝出来粘贴到 `~/.zshrc` 的末尾即可。

- 字体安装

oh-my-zsh 最直观和 `bash` 不一样的地方要数它丰富的主题了，也是一开始吸引我使用它的地方。不过在配置主题之前最好先保证系统有丰富的字体，因为部分主题依赖于[这些字体](https://github.com/powerline/fonts)，按照说明安装即可，安装完成后在 shell 偏好设置里面选择，我使用的是 Meslo for Powerline 系列。

- 切换主题

切换主题只需要编辑 `~/.zshrc` 文件，找到下面这段文本：

```
# Set name of the theme to load.
# Look in ~/.oh-my-zsh/themes/
# Optionally, if you set this to "random", it'll load a random theme each
# time that oh-my-zsh is loaded.
ZSH_THEME="random"
```

我自己配置的是随机 random ，自带的主题在 `~/.oh-my-zsh/themes` 里面，想用哪个就把相应的名字替换进去就可以了，这是一些自带主题的截图⇒[我带你去看](https://github.com/robbyrussell/oh-my-zsh/wiki/themes)。有可能你不满足于这些，没关系，还有[扩展主题](https://github.com/robbyrussell/oh-my-zsh/wiki/External-themes)，每个主题都有详细的配置方法，照做就 OK 了。

- 插件

oh-my-zsh 另外一个强大的地方是插件，自带插件在 `~/.oh-my-zsh/plugins` 目录下，想了解各个插件的功能和使用方法，阅读各个插件目录下的 `*.plugin.zsh` 就可以了，比如在终端中输入 `vim ~/.oh-my-zsh/plugins/git/git.plugin.zsh` ，你可以看到：

```
# Query/use custom command for `git`.
zstyle -s ":vcs_info:git:*:-all-" "command" _omz_git_git_cmd
: ${_omz_git_git_cmd:=git}

#
# Functions
#
...
alias g='git'

alias ga='git add'
alias gaa='git add --all'
alias gapa='git add --patch'
...
```

贴心好用到哭有木有😭😭😭

除了自带插件外，还有一大票扩展插件，我目前只安装了一个 [zsh-completions](https://github.com/zsh-users/zsh-completions) 。安装方式很简单，把插件用 git 命令克隆到 `~/.oh-my-zsh/custom/plugins` ，然后在配置文件中按插件说明配置保存、重启就可以了。

有些插件在按照说明安装使用后会报类似于如下的错误：

```
_arguments:450: _vim_files: function definition file not found
```

我的解决办法是，直接删除 `~/.zcompdump` 文件，但是不知道会不会影响其他插件的功能。


----------


## [技术手札：如何全新安装 Mac OS X El Capitan](http://www.jianshu.com/p/fa45988bb270)

通过实验，得到 `zsh` 的配置文件的加载规律：

```
#
# A: /etc/zshenv   B: ~/.zshenv   C: /etc/zprofile   D: ~/.zprofile
# E: /etc/zshrc    F: ~/.zshrc    G: /etc/zlogin     H: ~/.zlogin
# I: ~/.zlogout    J: /etc/zlogout
#+-------------------+-------------------------------------------+
#|                   |                   login                   |
#|                   +------------------------------+------------+
#|                   |              yes             |     no     |
#+-------------+-----+------------------------------+------------+
#|             | yes | A->B->C->D->E->F->G->H->I->J | A->B->E->F |
#| interactive |-----+------------------------------+------------+
#|             | no  | A->B->C->D->      G->H->I->J | A->B       |
#+-------------+-----+------------------------------+------------+
#
```

从加载顺序中可以看出来，`.zshenv` 文件是能保证被第一个加载的。

另外，OS X El Capitan 系统中，**有两个 `zsh` 的默认配置文件**，其中内容如下：

在 `/etc/zprofile` 中有：

```
# system-wide environment settings for zsh(1)
if [ -x /usr/libexec/path_helper ]; then
    eval `/usr/libexec/path_helper -s`
fi
```

在 `/etc/zshrc` 中有

```
# Correctly display UTF-8 with combining characters.
if [ "$TERM_PROGRAM" = "Apple_Terminal" ]; then
    setopt combiningchars
fi
```

我们发现，`/etc/zprofile` 引用了一个可执行文件
`/usr/libexec/path_helper`，那这个文件的作用是什么呢？

原来，苹果使用一套新的机制希望来替换传统的直接修改环境变量的方式：`path_helper`。

`path_helper` 命令只是用来输出一个 shell 语句，例如：

```
export $PATH=<...>
export $MANPATH=<...>
```

而本身并不执行任何修改。因此，可使用 `eval` 命令执行修改。`-s` 参数的作用，是只生成 `$PATH` 的 `export` 语句。

而执行 `path_helper` 命令的时候，它会按照以下次序依次添加路径：

- `/etc/paths` 文件中的路径
- `/etc/paths.d` 目录下所有文件中的路径
- 当前 `$PATH` 变量

其中，重复路径不再添加。

现在我们来推测一下：当系统加载 `zsh` 环境的时候，`$PATH` 环境变量到底发生了什么？

由于 OS X El Capitan 系统中**默认不存**在 `/etc/zshenv` 文件，所以 zsh 加载的第一个文件是 `.zshenv`。加载 `.zshenv` 后，`rvm`、`nvm.sh` 等环境配置脚本被执行，此时 `$PATH` 是理想的状态；

当系统执行 `/etc/zprofile` 文件的时候，文件中的 `path_helper` 指令对 `$PATH` 变量中所有的路径重新做了一个排序，系统默认的 `/bin` 路径自动排到了最前面，元凶终于找到了：）

解决方案：

所以，原则上，将在 `$PATH` 中添加前置路径的脚本，从 `.zshenv` 移到 `.zprofile` 和 `.zshrc` 中加载，即可。

其余的，具体情况具体分析。



----------

## 其它

- [UNIX shell differences and how to change your shell (Monthly Posting)](http://www.faqs.org/faqs/unix-faq/shell/shell-differences/)
- [理解 bashrc 和 profile](https://wido.me/sunteya/understand-bashrc-and-profile)
