# brew 和 brew cask

## Homebrew 与 Homebrew-Cask 的关系

> - Remember that Homebrew-Cask is an independent project from Homebrew.
> - Homebrew-Cask is implemented as a subcommand of Homebrew.

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

[Homebrew](https://github.com/Homebrew/brew) 介绍（[官方主页](http://brew.sh/)）：

> The missing package manager for macOS.

知乎上关于这个问题的讨论：[这里](https://www.zhihu.com/question/22624898)。
网友提供的很好的介绍：[这里](https://aaaaaashu.gitbooks.io/mac-dev-setup/content/Homebrew/Cask.html)。


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


## 提 bug 前的确认工作

在你认为遇到 bug 但找不到解决办法前 make sure you have the latest versions of Homebrew, Homebrew-Cask, and all Taps by running the following commands.

```shell
$ cd $(brew --repo); git fetch; git reset --hard origin/master
$ brew untap phinze/cask; brew untap caskroom/cask; brew uninstall --force brew-cask
$ brew cleanup; brew cask cleanup; brew update
```

补充说明：


```shell
$ cd $(brew --repo); git fetch; git reset --hard origin/master`
```
> If Homebrew was updated on **Aug 10-11th 2016** and `brew update` always says *Already up-to-date*. you need to run these commands above.

```shell
brew uninstall --force brew-cask
```
> To uninstall all versions of a Cask, use `--force`


```shell
$ brew cleanup; brew cask cleanup; brew update
```
> `cleanup` — cleans up cached downloads


----------


## 提 issue 的套路

> Start by searching for your issue before posting a new one. If you find an open issue and have any new information not reported in the original, please add your insights. If you find a closed issue, try the solutions there. If the issue is still not solved, open a new one with your new information and a link back to the old related issue.



