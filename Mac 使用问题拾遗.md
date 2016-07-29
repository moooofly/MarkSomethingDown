

本文记录 Mac 使用过程中遇到的一些小问题；


----------


# Mac 快速锁屏

- 系统偏好设置 –> Mission Control –> 触发角；
- 活跃的屏幕角，选择一个角，设置成“将显示器置入睡眠状态”；
- 每次把鼠标移动到那个角上的时候，立即执行了该项动作，类似鼠标手势一样。

# Spotlight 搜索
Spotlight 是 Mac OS X 中非常实用的搜索功能，可以通过 `control+空格` 来快速搜索 Mac 中的内容。

# Mac 必备工具

- [seil](https://pqrs.org/osx/karabiner/seil.html.en)
- [iterm2](https://www.iterm2.com/)
- [tmux](https://tmux.github.io/)
- [powerline/fonts](https://github.com/powerline/fonts)
- [zsh]()


# Mac 上的“邮件”不断提示输入密码

解决：最终在 Mac 上通过 web 登录邮箱后，会弹出提示“是否允许其他应用使用该账户“，允许后就没有问题了；
Mac 官方解答：[这里](https://support.apple.com/zh-cn/HT204187)

# Mac 上的剪切复制粘贴

Command + 拖拽 = 剪切
Option + 拖拽 = 复制
Command + Option + 拖拽 = 创建快捷方式

在高版本 Mac OS X 中，复制还是 `Command + C` ，粘贴时用 `Command + Option + V` 可以产生剪切＋粘贴效果，也就是复制成功以后删掉原文件。

参考：[这里](http://www.baifeng.me/apple/macosx/2010/04/1295/)


# mac 下如何更新 locate 命令依赖的数据库

sudo /usr/libexec/locate.updatedb

如果执行上述命令出现权限问题，则切到根目录 / 下执行；

locate 使用的数据库在 /var/db/locate.database

.DS_Store 文件是什么？如何禁止
http://www.zhihu.com/question/20345704

.DS_Store是Mac OS保存文件夹的自定义属性的隐藏文件，如文件的图标位置或背景色，相当于Windows的desktop.ini。

1，禁止.DS_store生成：
打开 “终端” ，复制黏贴下面的命令，回车执行，重启Mac即可生效。

defaults write com.apple.desktopservices DSDontWriteNetworkStores -bool TRUE

2，恢复.DS_store生成：

defaults delete com.apple.desktopservices DSDontWriteNetworkStores

作者：Marsokit
链接：http://www.zhihu.com/question/20345704/answer/19471793
来源：知乎
著作权归作者所有，转载请联系作者获得授权。

https://discuss.circleci.com/t/brew-update-command-fails/5211/2

/usr/local/Library/brew.sh: line 32: /usr/local/Library/ENV/scm/git: No such file or directory

This is a confirmed issue with Homebrew, per them it should be fixed, but if not the following will correct.

cd "$(brew --repository)" && git fetch && git reset --hard origin/master

ref: https://github.com/Homebrew/brew/issues/55799