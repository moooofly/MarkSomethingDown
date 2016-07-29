

本文记录 Mac 使用过程中遇到的一些小问题；


----------


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


# 执行 brew update 时报错如何解决

https://discuss.circleci.com/t/brew-update-command-fails/5211/2
当执行 brew update 时可能会出现如下错误信息
```shell
/usr/local/Library/brew.sh: line 32: /usr/local/Library/ENV/scm/git: No such file or directory
```

网上查到的解答为：
> This is a confirmed issue with Homebrew, per them it should be fixed, but if not the following will correct.
> `cd "$(brew --repository)" && git fetch && git reset --hard origin/master`
> ref: https://github.com/Homebrew/brew/issues/55799