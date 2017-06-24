# root 权限之 su 与 sudo

## su

> 以下内容取自 man 手册

`su` 命令用于变更 user ID 或者称为超级用户（superuser）；

- `su` 命令用于在会话登录期间切换成另外一个用户身份；
- 如果没有指定切换的目标用户名，则 `su` 默认切换成超级用户；
- 可选参数 `-` 可以用来提供与用户直接登录时（login shell）相似的环境；
- 可以在用户名后面继续提供额外的参数，而这些参数会将提供给用户登录 shell 使用；特别地，若指定 `-c` 参数，则其后的参数将被当作 command 被命令行解析器识别；该 command 将被 `/etc/passwd` 中针对目标用户配置的 shell 所执行；
- 可以使用 `--` 参数将 `su` 的选项和供 shell 使用的参数分开；
- 用户会被提供要求输入密码，当然是在需要的情况下（在 root 用户下执行 `su` 命令就不需要提供密码了）；
- 所有的 `su` 切换，无论是成功切换，还是失败切换，该行为都会被记录下来以便探测到系统是否被乱用；
- （在通过 `su` 切换用户时）当前环境会被传递给新 shell 中；但针对 `$PATH` 的内容：针对普通用户来说，重置为 `/bin:/usr/bin`；针对超级用户来说，重置为 `/sbin:/bin:/usr/sbin:/usr/bin`；该行为可以通过 `/etc/login.defs` 文件中的 `ENV_PATH` 和 `ENV_SUPATH` 进行变更；
- A subsystem login is indicated by the presence of a "*" as the first character of the login shell. The given home directory will be used as the root of a new file system which the user is actually logged into.

常用参数：

- **`-, -l, --login`**

提供与“用户直接时的登录环境（login shell）”相似的环境；当使用 `-` 形式时，其只能作为 su 的最后一个选项使用；而其他形式（-l 和 --login）则没有此限制；

### 实验

> 注：以下执行过程是连续进行的；

- 基于普通用户登录

```shell
vagrant@vagrant-ubuntu-trusty:~$ echo $$
1031
vagrant@vagrant-ubuntu-trusty:~$
vagrant@vagrant-ubuntu-trusty:~$ echo $PATH
/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games
vagrant@vagrant-ubuntu-trusty:~$
vagrant@vagrant-ubuntu-trusty:~$ go version
The program 'go' can be found in the following packages:
 * gccgo-5
 * golang-go
Try: sudo apt-get install <selected package>
vagrant@vagrant-ubuntu-trusty:~$
```

- 基于 `su - root` 切换超级用户

```shell
vagrant@vagrant-ubuntu-trusty:~$ su - root
Password:
root@vagrant-ubuntu-trusty:~#
root@vagrant-ubuntu-trusty:~# echo $$
1210
root@vagrant-ubuntu-trusty:~#
root@vagrant-ubuntu-trusty:~# echo $PATH
/usr/local/go/bin:/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/mysql/bin
root@vagrant-ubuntu-trusty:~#
root@vagrant-ubuntu-trusty:~# go version
go version go1.7.4 linux/amd64
root@vagrant-ubuntu-trusty:~#
```

- 基于 `su` 切换超级用户（当前已经是 root）

```shell
root@vagrant-ubuntu-trusty:~# su
root@vagrant-ubuntu-trusty:~#
root@vagrant-ubuntu-trusty:~# echo $$
1238
root@vagrant-ubuntu-trusty:~#
root@vagrant-ubuntu-trusty:~# echo $PATH
/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games
root@vagrant-ubuntu-trusty:~#
root@vagrant-ubuntu-trusty:~# go version
The program 'go' can be found in the following packages:
 * gccgo-5
 * golang-go
Try: apt-get install <selected package>
root@vagrant-ubuntu-trusty:~#
```

- 基于 `su -` 切换超级用户（当前已经是 root）

```shell
root@vagrant-ubuntu-trusty:~# exit
root@vagrant-ubuntu-trusty:~#
root@vagrant-ubuntu-trusty:~# echo $$
1210
root@vagrant-ubuntu-trusty:~#
root@vagrant-ubuntu-trusty:~# su -
root@vagrant-ubuntu-trusty:~#
root@vagrant-ubuntu-trusty:~# echo $$
1251
root@vagrant-ubuntu-trusty:~#
root@vagrant-ubuntu-trusty:~# echo $PATH
/usr/local/go/bin:/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/mysql/bin
root@vagrant-ubuntu-trusty:~#
root@vagrant-ubuntu-trusty:~# go version
go version go1.7.4 linux/amd64
root@vagrant-ubuntu-trusty:~#
```

### 小结

- `su -` 和 `su - root` 是等价的，成功切换需要 root 密码；
- `su` 和 `su -` 的差别在于：`su` 具有重置行为（按 non-login shell 处理），而 `su -` 按 login shell 处理；


### 其他

提及 login shell 就涉及到用户登录时，都需要执行到哪些脚本的问题，因为知道在何种情况下执行哪些脚本，以及各个脚本的作用，才能准确的将自定义内容添加到合适的位置；

> 参考：《[Execution sequence for .bash_profile, .bashrc, .bash_login, .profile and .bash_logout](http://www.thegeekstuff.com/2008/10/execution-sequence-for-bash_profile-bashrc-bash_login-profile-and-bash_logout/)》

- 登录 interactive login shell 时（伪代码）

```
execute /etc/profile
IF ~/.bash_profile exists THEN
    execute ~/.bash_profile
ELSE
    IF ~/.bash_login exist THEN
        execute ~/.bash_login
    ELSE
        IF ~/.profile exist THEN
            execute ~/.profile
        END IF
    END IF
END IF
```
> 注意：
> - 在不同 LINUX 操作系统下，可能存在 ~/.bash_profile、~/.bash_login、~/.profile 中的一种或几种；如果存在几种的话，那么执行的顺序便如上述伪代码；
> - 若存在多种配置文件，则要将自定义内容放入优先被处理的文件中，否则无法生效；

- 登出 interactive login shell 时（伪代码）

```
IF ~/.bash_logout exists THEN
    execute ~/.bash_logout
END IF
```

- 登录 interactive non-login shell 时（伪代码）

```
IF ~/.bashrc exists THEN
    execute ~/.bashrc
END IF
```

另外，Mac 的终端默认开启为 Login Shell ，而 Ubuntu 的 Gnome Terminal 默认开启的是 Non-Login Shell ；


----------


## sudo

> 以下内容取自 man 手册

以其他用户身份执行命令；


- sudo 允许授权用户使用超级用户或其他用户身份执行命令；
- sudo 基于插件架构实现安全策略和 input/output logging 功能；第三方能够开发和分发自身实现的策略和 I/O logging 插件，以无缝对接 `sudo` 前端；
- 默认的安全策略为 `sudoers` ，其通过文件 /etc/sudoers 进行配置；也可以采用 LDAP 实现；
- 安全策略决定了一个用户能够以何种特权执行 `sudo` 命令；策略本身可能会要求用户对自身进行密码鉴权，或者基于其他方式进行鉴权；
-  如果鉴权是必须的，那么在预先配置的时间约束内，用户没有输入密码，则 `sudo` 将会主动推出；该时间约束是策略相关的；默认情况下，基于 sudoers 安全策略的密码提示超时时间是无限的；
- 安全策略可以支持密码缓存以便允许用户在一定时间内再次使用 `sudo` 时无需再次鉴权；
- 基于 sudoers 的策略，密码缓存时间为 15 分钟，除非在 sudoers(5) 中进行改写；
- 若使用 `-v` 选项运行 `sudo` 命令，用户可以直接更新缓存的密码，而不用真正允许运行命令；
- 当使用 `sudoedit` 启用 `sudo` 命令，隐含使用了 `-e` 选项；
- 安全策略会记录调用 `sudo` 命令时的成功和失败情况；如果配置了 I/O 插件，则运行命令本身的输入和输出同样会被进行日志记录；


常用参数：

- **`-i, --login`**

以 login shell 运行由目标用户的密码库入口（`/etc/passwd`）所指定的 shell ；这意味着登录相关的资源文件，如 `.profile` 或 `.login` ，将被 shell 所读取；

如果指定了具体的 command ，则会传递给该 shell 并通过 shell 的 `-c` 选项执行；如果没有指定具体 command ，那么直接启动一个 interactive shell ；

`sudo` 会在运行 shell 前，尝试变更到相应用户的家目录下；
该命令被调用时得到的运行环境，类似于一个用户在刚登录时所得到的环境； 

在 sudoers(5) 手册中能够看到 `-i` 选项是如何在 sudoers 安全策略下影响 command 运行的环境的； 


- **`-s, --shell`**

当设置了 SHELL 环境变量的话，则运行由其指定的 shell ；否则将运行由用户密码库入口（通常为 /etc/passwd）指定的 shell ；如果指定了 command ，则会将 command 传递给 shell 并通过 shell 的 `-c` 选项执行；如果未指定 command ，则启动交互式 shell ；

> 执行 `sudo bash` 等价于 执行 `sudo -s` 并且当前默认 shell 为 `bash` 对情况；

- **`-u user, --user=user`**

以指定用户身份运行 command ，而不是基于默认目标用户身份（通常为 root 用户）；

user 内容可以为用户名或者带有 '#' 前缀的 user ID (UID)（例如 #0 代表 UID 0）

当基于 UID 运行 command 时，许多 shell 都要求对 '#' 进行转义（使用 '\'）；


- **`-v, --validate`**

更新用户的缓存密码，在必要的情况下需要对用户身份进行鉴权；对于 sudoers 插件来说，该选项能够将 sudo 的默认 15 分钟超时进行再次扩展，而无需运行任何具体 command ；需要注意的时，并不是所有的安全策略都支持密码缓存功能；

- **`--`**

选项 `--` 表明 sudo 应该停止处理命令行参数（即其余参数属于 command）；


## 试验

```shell
vagrant@vagrant-ubuntu-trusty:~$ echo $$
1009
vagrant@vagrant-ubuntu-trusty:~$ echo $PATH
/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games
vagrant@vagrant-ubuntu-trusty:~$
vagrant@vagrant-ubuntu-trusty:~$ go version
The program 'go' can be found in the following packages:
 * gccgo-5
 * golang-go
Try: sudo apt-get install <selected package>
vagrant@vagrant-ubuntu-trusty:~$
vagrant@vagrant-ubuntu-trusty:~$ sudo -i
root@vagrant-ubuntu-trusty:~#
root@vagrant-ubuntu-trusty:~# echo $$
1417
root@vagrant-ubuntu-trusty:~#
root@vagrant-ubuntu-trusty:~# echo $PATH
/usr/local/go/bin:/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/local/mysql/bin
root@vagrant-ubuntu-trusty:~#
root@vagrant-ubuntu-trusty:~# go version
go version go1.7.4 linux/amd64
root@vagrant-ubuntu-trusty:~#
```

## 小结

- `su - root` 切换时需要 root 密码；`sudo -i` 切换 root 时不需要密码；
- 推荐使用 sudo 的理由：
    - su - root 需要将 root 密码提供给别人；
    - 可以在 /etc/sudoers 文件中限制哪些用户能够临时获得 root 权限；
    - 基于 sudo 容易实现“以普通用户身份运行服务程序”的功能（降低权限，防攻击）；


----------

## 几个令人困惑的问题

### [su vs sudo -s vs sudo -i vs sudo bash](http://unix.stackexchange.com/questions/35338/su-vs-sudo-s-vs-sudo-i-vs-sudo-bash)

- `sudo` 会记录你所运行的命令，并与你自身用户身份相关联（可用于事后问责）；
- `sudo` 的权限控制更加弹性；`su` 的权限控制为 "all or nothing" ；
- 设置 `sudo` 规则的时候，若未对 `sudo -s` 或 `sudo bash` 或 `sudo vi` 进行限制，则理论上使用者能够进行越权操作；
- `sudo` 能够进行权限隔离的原因在于：`sudo` 使用每个 sudoers 自身的密码，而非 root 密码；
- 使用 `sudo` 时不存在“ root 密码变更后需要通知所有基于 su 切换 root ”问题；
- `sudo bash` 和 `sudo -s` 的差别在于：
    - `sudo bash` 使用确定的 shell ，即 bash ；而 `sudo -s` 使用配置相关的 shell ；
    - `sudo -s command` 基本上等价于 `sudo $SHELL -c command` ；
    - `sudo -s` 可以基于 shell 的 stdin 接收 command ；例如 `sudo -s < my-shell-script` ，或者 `sudo -s heredoc` ；通过这种方式，可以达到仅输入一次 `sudo` 就执行多条 command 的目的；
- `sudo -s` 和 `sudo -i` 的区别在于：前者会继承当前用户环境；后者会获得类似于 login shell 的干净环境；这也是 `sudo -i` 要比 `sudo -s` 更安全的原因；
- Roughly speaking, `sudo -i` is to `sudo -s` as `su -` is to `su`: it resets all but a few key environment variables and sends you back to your user's home directory.


### [Difference between sudo -i and sudo su](http://unix.stackexchange.com/questions/98531/difference-between-sudo-i-and-sudo-su)

`sudo -i` 的几大优点：

- So the big take away is that the method used to become root `sudo -i` is advantages over the others because you use your own password to do so, protecting root's password from needing to be given out.
- There is logging when you became root, vs. mysteriously some one becoming root via `su` or `su -`.
- `sudo -i` gives you a better user experience over either `su`'s because it protects your `$DISPLAY` and `$TERM`.
- `sudo -i` provides some protection to the system when user's become root, by limiting the environment which they are given.

不要使用 `sudo su` 的原因：

> When you run `sudo su` the `sudo` command masks the effects of the `su` and so much of the environment that you'd get from a regular `su` is lost.


### ['sudo su -' vs 'sudo -i' vs 'sudo /bin/bash' - when does it matter which is used, or does it matter at all?](http://askubuntu.com/questions/376199/sudo-su-vs-sudo-i-vs-sudo-bin-bash-when-does-it-matter-which-is-used)


the difference between **login**, **non-login**, **interactive** and **non-interactive** shells:

- **login shell**: A login shell logs you into the system as a specified user, necessary for this is a username and password. 
- **non-login shell**: A shell that is executed without logging in, necessary for this is a currently logged-in user.
- **interactive shell**: A shell (login or non-login) where you can interactively type or interrupt commands.
- **non-interactive shell**: A (sub)shell that is probably run from an automated process. You will see neither input nor output.


`sudo su` Calls `sudo` with the command `su`. Bash is called as **interactive** **non-login** shell. So bash only executes `.bashrc`. You can see that after switching to root you are still in the same directory.

`sudo su -` This time it is a login shell, so `/etc/profile`, `.profile` and `.bashrc` are executed and you will find yourself in root's home directory with root's environment.

`sudo -i` It is nearly the same as sudo `su -`. The `-i` (simulate initial login) option runs the shell specified by the password database entry of the target user as a login shell. This means that login-specific resource files such as `.profile`, `.bashrc` or `.login` will be read and executed by the shell.

`sudo /bin/bash` This means that you call `sudo` with the command `/bin/bash`. `/bin/bash` is started as non-login shell so all the dot-files are not executed, but bash itself reads `.bashrc` of the calling user. Your environment stays the same. Your home will not be root's home. So you are root, but in the environment of the calling user.

`sudo -s` reads the `$SHELL` variable and executes the content. If `$SHELL` contains `/bin/bash` it invokes `sudo /bin/bash` (see above).


To check if you are in a login shell or not (works only in `bash` because `shopt` is a builtin command):

```shell
shopt -q login_shell && echo 'Login shell' || echo 'No login shell'
```

### [Difference between sudo su and sudo -s](https://ubuntuforums.org/showthread.php?t=983645&p=6188826#post6188826)

Here are the differences I found:

- "sudo -s" 
```
HOME=/home/applic
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/X11R6/bin
reads $USER's ~/.bashrc
```

- "sudo su"
```
HOME=/root
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games
reads /etc/environment
reads /root/.bashrc
```

Being root and having $HOME set to the normal user's home can cause problems. 

Here is a summary:

| | HOME=/root | uses root's PATH | corrupted by user's env vars |
---|---|---|---
sudo -i | Y	| Y[2] | N
sudo -s	| N	| Y[2] | Y
sudo bash | N | Y[2] | Y
sudo su | Y | N[1] | Y

[1] `PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games` probably set by `/etc/environment`
[2] `PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/X11R6/bin`

To start a **root shell** (i.e. a command window where you can run root commands), starting root's environment and login scripts, use:
```shell
sudo -i     (similar to sudo su - , gives you roots environment configuration)
```
To start a **root shell**, but keep the current shell's environment, use:
```
sudo -s     (similar to sudo su)
```

底线就是："`sudo -i`" 是，当你想要获得 root shell 的同时，又不想“污染”用户自身环境的情况下，最合适的命令；

The quick easy answer was

> Don't do "`sudo su`" or "`sudo -i`" to get a "true" root login. It open up your system to potential harm, and may cause viruses, trojan's etc to infect all the windows systems in the network. If you insist on using it, sit in front of the terminal, and pull the network cable out of the wall!

## [RootSudo](https://help.ubuntu.com/community/RootSudo)

> 此文是 ubuntu 官方给出的使用建议，其中的信息非常有价值；

- 默认情况下，Ubuntu 上 root 账户密码是锁定的；
> This means that you cannot login as root directly or use the su command to become the root user. However, since the root account physically exists it is still possible to run programs with root-level privileges. This is where sudo comes in - it allows authorized users to run certain programs as root without having to know the root password.

- 由于使用 sudo 引发的不便
> Redirecting the output of commands run with sudo requires a different approach. For instance consider `sudo ls > /root/somefile` will not work since it is the shell that tries to write to that file. You can use `ls | sudo tee -a /root/somefile` to append, or `ls | sudo tee /root/somefile` to overwrite contents. You could also pass the whole command to a shell process run under sudo to have the file written to with root permissions, such as `sudo sh -c "ls > /root/somefile"`.

- 开启 root 账户（不建议）
> To enable the root account (i.e. set a password) use:
> 
> `sudo passwd root`
>
> Enabling the root account is rarely necessary. Almost everything you need to do as administrator of an Ubuntu system can be done via `sudo` or `gksudo`. If you really need a persistent root login, **the best alternative is to simulate a root login shell** using the following command...
> 
> `sudo -i` （后文又补充说明：即使是这种命令也不建议使用）

- 免密码使用 sudo （不建议）
>If you disable the sudo password for your account, you will seriously compromise the security of your computer. Anyone sitting at your unattended, logged in account will have complete root access, and remote exploits become much easier for malicious crackers.







