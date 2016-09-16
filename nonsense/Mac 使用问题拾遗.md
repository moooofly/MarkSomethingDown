

本文记录 Mac 使用过程中遇到的一些小问题；


----------

# pip 安装问题

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


# Mac 上的“邮件”不断提示输入密码

最终在 Mac 上通过 web 登录邮箱后，会弹出提示“是否允许其他应用使用该账户“，允许后就没有问题了；    
Mac 官方解答：[这里](https://support.apple.com/zh-cn/HT204187)；

# Mac 上的剪切复制粘贴

- Command + 拖拽 = 剪切
- Option + 拖拽 = 复制
- Command + Option + 拖拽 = 创建快捷方式

在高版本 Mac OS X 中，复制还是 `Command + C` ，粘贴时用 `Command + Option + V` 可以产生剪切＋粘贴效果，也就是复制成功以后删掉原文件；

参考：[这里](http://www.baifeng.me/apple/macosx/2010/04/1295/)；


# mac 下如何更新 locate 命令依赖的数据库

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


# Homebrew 使用姿势


## 简介

以下内容取自：[这里](http://brew.sh/index_zh-cn.html)

- OS X 不可或缺的套件管理器；
- Homebrew 将套件安装到独立目录，并将文件软链接至 /usr/local ；
- 与 Homebrew 相关的所有文件均会被安装到预定义目录下，无需操心 Homebrew 的安装位置问题；
- Homebrew 以 git, ruby 为其筋骨；因此，可以借助您的相关知识进行自由修改；可以方便地撤回您的修改或者合并上游更新；
- Homebrew 的程式都是简单的 Ruby 脚本；

## 安装

（安装过程中需要用到 root 权限）

```shell
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

## 卸载

（卸载过程中会提示有些内容需要手动清理）

```shell
ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/uninstall)"
```



# 执行 brew update 时报错如何解决

当执行 brew update 时可能会出现如下错误信息
```shell
/usr/local/Library/brew.sh: line 32: /usr/local/Library/ENV/scm/git: No such file or directory
```

网上查到的[解答](https://discuss.circleci.com/t/brew-update-command-fails/5211/2)为：
> This is a confirmed issue with Homebrew, per them it should be fixed, but if not the following will correct.    
> `cd "$(brew --repository)" && git fetch && git reset --hard origin/master`    
> ref: https://github.com/Homebrew/brew/issues/55799    

另外，在有些时候，如果第一次 `brew update` 失败了，再执行一次可能就会成功（原因未知）；
