# Homebrew 相关

## Homebrew 安装和卸载

> 以下内容取自：[这里](http://brew.sh/index_zh-cn.html)

- OS X 不可或缺的套件管理器；
- Homebrew 将套件安装到独立目录，并将文件软链接至 `/usr/local` ；
- 与 Homebrew 相关的所有文件均会被安装到预定义目录下，无需操心 Homebrew 的安装位置问题；
- Homebrew 以 `git`, `ruby` 为其筋骨；因此，可以借助您的相关知识进行自由修改；可以方便地撤回您的修改或者合并上游更新；
- Homebrew 的程式都是简单的 Ruby 脚本；

### 安装

（安装过程中需要用到 root 权限）

```shell
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
```

### 卸载

（卸载过程中会提示有些内容需要手动清理）

```shell
ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/uninstall)"
```


----------


## Homebrew 相关术语

| Term           | Description                                                | Example                                                         |
|----------------|------------------------------------------------------------|-----------------------------------------------------------------|
| **Formula**    | package 定义文件                                     | `/usr/local/Homebrew/Library/Taps/homebrew/homebrew-core/Formula/foo.rb` |
| **Keg**        | The installation prefix of a **Formula**（可理解成版本号）                  | `/usr/local/Cellar/foo/0.1`                                     |
| **opt prefix** | A symlink to the active version of a **Keg**               | `/usr/local/opt/foo `                                           |
| **Cellar**     | All **Kegs** are installed here                            | `/usr/local/Cellar`                                             |
| **Tap**        | 存放 **Formulae** 和/或 commands 的其它可选 Git 仓库  | `/usr/local/Homebrew/Library/Taps/homebrew/homebrew-versions`   |
| **Bottle**     | 用于替代源码构建而预先构建好的 **Keg**      | `qt-4.8.4.mavericks.bottle.tar.gz`                              |
| **Cask**       | 用于安装 macOS native apps 的 [homebrew 扩展](https://github.com/caskroom/homebrew-cask)  | `/Applications/MacDown.app/Contents/SharedSupport/bin/macdown` |
| **Brew Bundle**| 用于描述依赖关系的 [homebrew 扩展](https://github.com/Homebrew/homebrew-bundle)   | `brew 'myservice', restart_service: true` |


----------


## Homebrew 命令简介

- 查看是否存在可更新的版本

```
brew outdated
```

- 针对所有或特定的 formulae 移除 cellar 中的老旧版本（因为 upgrade 默认保留老旧版本）

```
brew cleanup [formulae]
```

- 升级版本（若指定 cleanup 则直接移除之前安装的老旧版本）

```
brew upgrade [--cleanup] [formulae]
```

- 从 github 上获取最新版本的 Homebrew 以及全部 formulae ，并进行必要的迁移

```
brew update [--merge] [--force]
```

- 在不同版本间进行切换（实际切换的是符号连接，必须存在多个版本，即不能被 cleanup 掉）

```
brew switch name version
```

- 在编辑器中打开 formula（即 xxx.rb 文件）

```
brew edit formula
```

- 从 Homebrew prefix 对应的目录下移除针对特定 formula 的符号链接

```
unlink [--dry-run] formula
```


----------


## Homebrew 与 Homebrew-Cask 的关系

> - Remember that Homebrew-Cask is an independent project from Homebrew.
> - Homebrew-Cask is implemented as a subcommand of Homebrew.

[Homebrew](https://github.com/Homebrew/brew) 介绍（[官方主页](http://brew.sh/)）：

> The missing package manager for macOS.

[Homebrew-Cask](https://github.com/caskroom/homebrew-cask) 介绍（[官方主页](https://caskroom.github.io/)）：

> “To install, drag this icon…” no more!
> 
> **Homebrew-Cask** extends [Homebrew](http://brew.sh/) and brings its elegance, simplicity, and speed to the installation and management of **GUI macOS applications** such as Google Chrome and Adium.
> 
> We do this by providing a friendly Homebrew-style CLI workflow for the administration of macOS applications distributed **as binaries**.
> 
> It’s implemented as a `homebrew` [external command](https://github.com/Homebrew/brew/blob/master/docs/External-Commands.md) called `cask`.
>
> To start using Homebrew-Cask, you just need Homebrew installed. 
>> As of December 2015, it is now unnecessary to install cask manually as it is now part of homebrew's installation. So after updating homebrew via `brew update`, you are set to use `brew cask`.

知乎上关于这个问题的讨论：[这里](https://www.zhihu.com/question/22624898)；   
网友提供的很好的介绍：[这里](https://aaaaaashu.gitbooks.io/mac-dev-setup/content/Homebrew/Cask.html)；



----------

## Homebrew-Cask v.s. Mac App Store

Homebrew-Cask 和 Mac App Store 相比，目前还有很多优势：

- 安装软件体验非常一致、简洁优雅；
- 对常用软件支持更全面，例如 MPlayerX 已经宣布不在更新 Mac App Store上 的版本；
- 软件更新速度快，体验好。例如 Alfred 2.0 已经出了很久（目前都 3.0 版本了），但在 Mac App Store 上还是 1.2 版本，QQ 也是这样的情况；

Mac App Store 生态圈远不完善，审核流程过长，限制太多，维护成本过高，让很多应用开发者被迫离开。虽然我个人很喜欢 Homebrew-Cask，但还是希望 Apple 尽快完善 Mac App Store ；



----------



## Homebrew 安装指定版本的软件

简单的讲，有如下几种办法：

- 如果之前已经基于 brew 安装了多个版本，即在 `/usr/local/Cellar/xxx/` 下存在多个版本号目录，则可以直接通过 `brew switch xxx <version_num>` 进行切换；
- 如果目标软件是基于 `brew install xxx` 安装的，则此时只会存在唯一一个版本号目录，则默认安装的是最新版本；通过 `brew edit xxx` 进入针对 xxx.rb 的编辑模式（不要通过任何其他方式操作），其中包含了当前版本 url 和 sha256 信息（用于获取对应版本以及进行校验）；因此，只需在 github 上查找对应目标版本的 `xxx.rb` 文件的修改历史，将其中的 url 和 sha256 内容拷贝到当前文件即可；修改后只需要执行 `brew unlink xxx; brew install xxx` 进行安装；
- 还有一种方式是自行手动下载对应版本的包（可以根据 xxx.rb 文件中 url 地址确定），解压到相应的 `/usr/local/Cellar/xxx/<version_num>/` 目录下，之后就可以通过 `brew switch xxx <version_num>` 进行切换了；


网上可以搜索到一些关于 brew 如何安装特定版本的说法：

- `brew versions` (已经废弃了)
- `brew tap homebrew/versions`（已经废弃了）

参考：

- [brew安装指定版本的软件](http://www.jianshu.com/p/aadb54eac0a8)


----------


## 执行 brew update 时报错

当执行 `brew update` 时可能会出现如下错误信息

```shell
/usr/local/Library/brew.sh: line 32: /usr/local/Library/ENV/scm/git: No such file or directory
```

网上的[解答](https://discuss.circleci.com/t/brew-update-command-fails/5211/2)：

> This is a confirmed issue with Homebrew, per them it should be fixed, but if not the following will correct.    
> `cd "$(brew --repository)" && git fetch && git reset --hard origin/master`    
> ref: https://github.com/Homebrew/brew/issues/55799    

另外，在有些时候，如果第一次 `brew update` 失败了，再执行一次可能就会成功（原因未知）；


----------


## 执行 brew 命令提示 “GitHub API Error”

执行 brew 时会输出如下错误信息：


```shell
...
Error: GitHub API Error: API rate limit exceeded for 103.215.2.69. (But here's the good news: Authenticated requests get a higher rate limit. Check out the documentation for more details.)
Try again in 11 minutes 48 seconds, or create a personal access token:
  https://github.com/settings/tokens/new?scopes=gist,public_repo&description=Homebrew
and then set the token as: export HOMEBREW_GITHUB_API_TOKEN="your_new_token"
```

在浏览器中打开上述 URL 并创建 token ；

```
Personal access tokens

Tokens you have generated that can be used to access the GitHub API.

Make sure to copy your new personal access token now. You won’t be able to see it again!

f7c**********(略)*********ccb

Personal access tokens function like ordinary OAuth access tokens. They can be used instead of a password for Git over HTTPS, or can be used to authenticate to the API over Basic Authentication.
```

将生成的 token 值添加到对应 shell 的 rc 文件中，例如 .zshrc 文件
```shell
if [ -f /usr/local/bin/brew ]; then
    export HOMEBREW_GITHUB_API_TOKEN=f7c**********(略)*********ccb
fi
```

最后通过 source 命令使其立即生效；
```shell
source ~/.zshrc
```


----------


## 可通过 brew cask 安装的常用软件

基于下面的指令可以实现传说中的一键装机

```
brew cask install iterm2
brew cask install sublime-text
brew cask install qq
brew cask install thunder
brew cask install google-chrome
brew cask install alfred
brew cask install skitch
brew cask install dropbox
brew cask install skype
brew cask install mplayerx
brew cask install evernote
brew cask install mou
brew cask install virtualbox
brew cask install visual-studio-code
```

一些常用开发软件

```
brew install cask
brew install lua
brew install luajit
brew install go
brew install openssl
brew install vim
brew install emacs
brew install wget
brew install zsh
brew install zsh-completions
brew install elasticsearch
brew install sqlite
brew install readline
brew install gettext
brew install ruby
brew install pcre
brew install perl
brew install gdbm
brew install libyaml
brew install pkg-config
brew install geoip c-ares cmake glib dbus gnutls libgcrypt libsmi qt
brew install wireshark --with-qt --with-lua --with-libsmi --with-headers
brew cask install wireshark-chmodbpf
brew install libevent
brew install tmux
```

可能需要运行的命令

```
brew install caskroom/cask
```




