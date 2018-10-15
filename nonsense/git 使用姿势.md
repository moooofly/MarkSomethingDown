# git 使用姿势

## 图示

![Git Cheet Sheet](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Git%20Cheat%20Sheet.jpg "Git Cheet Sheet")

![Git 常用命令速查表](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Git%20%E5%B8%B8%E7%94%A8%E5%91%BD%E4%BB%A4%E9%80%9F%E6%9F%A5%E8%A1%A8.png "Git 常用命令速查表")


----------


## 目录

- [fatal: could not read Username for 'https://git.llsapp.com': terminal prompts disabled](#fatal-could-not-read-username-for-httpsgitllsappcom-terminal-prompts-disabled)
- [已经 push 到远端仓库后发现未 signoff](#已经-push-到远端仓库后发现未-signoff)
- [git 四个区和五种状态的切换](#git-四个区和五种状态的切换)
    - [git 的四种区](#git-的四种区)
    - [git 的五种状态](#git-的五种状态)
    - [秒懂 Git 的区和状态](#秒懂-git-的区和状态)
- [git diff 使用](#git-diff-使用)
- [git 常用的撤销操作](#git-常用的撤销操作)
- [拉取远程特定分支上的内容](#拉取远程特定分支上的内容)
- [如何撤销错误的 merge 结果](#如何撤销错误的-merge-结果)
- [如何在 git 中写出好的 commit 说明](#如何在-git-中写出好的-commit-说明)
    - [为什么要关注提交信息](#为什么要关注提交信息)
    - [commit 编写的基本要求](#commit-编写的基本要求)
    - [commit 中要回答哪些问题](#commit-中要回答哪些问题)
    - [tips](#tips)
- [git merge 与 git rebase](#git-merge-与-git-rebase)
    - [git merge](#git-merge)
    - [git rebase](#git-rebase)
    - [对比](#对比)
- [git tag 使用](#git-tag-使用)
- [git clone 时直接 rename](#git-clone-时直接-rename)
- [Create New Repository](#create-new-repository)
- [本地新建分支后 push 到远端仓库](#本地新建分支后-push-到远端仓库)
- [获取指定 tag 代码](#获取指定-tag-代码)
- [直接获取远端仓库指定 branch 代码](#直接获取远端仓库指定-branch-代码)
- [将远端仓库里指定 branch 拉取到本地](#将远端仓库里指定-branch-拉取到本地)
- [重命名 remote branch 名](#重命名-remote-branch-名)
- [删除不存在对应远程分支的本地分支](#删除不存在对应远程分支的本地分支)
- [fork 别人项目后如何同步其后续更新](#fork-别人项目后如何同步其后续更新)
- [定制化 git 全局配置](#定制化-git-全局配置)
- [基于 ssh 协议访问 git](#基于-ssh-协议访问-git)
- [其他](#其他)
    - [GitHub 秘籍](#github-秘籍)
    - [Git 简明指南](#git-简明指南)
    - [Git Cheat Sheet](#git-cheat-sheet)
    - [图解 Git](#图解-git)
    - [Think like a Git](#think-like-a-git)
    - [Git 社区参考书](#git-社区参考书)
- [Tools](#tools)
    - [Tower](#tower)


## fatal: could not read Username for 'https://git.llsapp.com': terminal prompts disabled

### 背景

该症状应该是创建了 GitLab Personal Access Tokens 后才出现的；

错误信息如下

```sh
[#51#root@ubuntu-1604 /go/src/git.llsapp.com/codelab]$go get -v git.llsapp.com/codelab/bazel_grpc
Fetching https://git.llsapp.com/codelab/bazel_grpc?go-get=1
Parsing meta tags from https://git.llsapp.com/codelab/bazel_grpc?go-get=1 (status code 200)
get "git.llsapp.com/codelab/bazel_grpc": found meta tag get.metaImport{Prefix:"git.llsapp.com/codelab/bazel_grpc", VCS:"git", RepoRoot:"https://git.llsapp.com/codelab/bazel_grpc.git"} at https://git.llsapp.com/codelab/bazel_grpc?go-get=1
git.llsapp.com/codelab/bazel_grpc (download)
# cd .; git clone https://git.llsapp.com/codelab/bazel_grpc.git /go/src/git.llsapp.com/codelab/bazel_grpc
Cloning into '/go/src/git.llsapp.com/codelab/bazel_grpc'...
fatal: could not read Username for 'https://git.llsapp.com': terminal prompts disabled
package git.llsapp.com/codelab/bazel_grpc: exit status 128
[#52#root@ubuntu-1604 /go/src/git.llsapp.com/codelab]$
```

### 原因分析

主要问题在于 "fatal: could not read Username for 'https://git.llsapp.com': terminal prompts disabled" ；

查到的原因有：

- "go get disable "terminal prompt" by default"
- "go get will refuse to authenticate on the command line. So you need to cache the credentials in git."
- "For me, I have 2FA enabled and thus could not use my password to auth. Instead I had to generate a personal access token to use in place of the password."

> 上面的说法都是正确的

### 解决

解决办法（土办法）：执行 `GIT_TERMINAL_PROMPT=1 go get -v git.llsapp.com/codelab/bazel_grpc`

```sh
[#52#root@ubuntu-1604 /go/src/git.llsapp.com/codelab]$GIT_TERMINAL_PROMPT=1 go get -v git.llsapp.com/codelab/bazel_grpc
Fetching https://git.llsapp.com/codelab/bazel_grpc?go-get=1
Parsing meta tags from https://git.llsapp.com/codelab/bazel_grpc?go-get=1 (status code 200)
get "git.llsapp.com/codelab/bazel_grpc": found meta tag get.metaImport{Prefix:"git.llsapp.com/codelab/bazel_grpc", VCS:"git", RepoRoot:"https://git.llsapp.com/codelab/bazel_grpc.git"} at https://git.llsapp.com/codelab/bazel_grpc?go-get=1
git.llsapp.com/codelab/bazel_grpc (download)
Username for 'https://git.llsapp.com': fei.sun
Password for 'https://fei.sun@git.llsapp.com':
package git.llsapp.com/codelab/bazel_grpc: no Go files in /go/src/git.llsapp.com/codelab/bazel_grpc
[#53#root@ubuntu-1604 /go/src/git.llsapp.com/codelab]$
```

上面的 Password 为 Personal Access Tokens 内容；


### 更好的解决办法

Ref: [go get results in 'terminal prompts disabled' error for github private repo](https://stackoverflow.com/questions/32232655/go-get-results-in-terminal-prompts-disabled-error-for-github-private-repo)


- before

```sh
[#56#root@ubuntu-1604 /go/src/git.llsapp.com/codelab]$cat ~/.gitconfig
...
[url "git@github.com:"]
    insteadOf = https://github.com/
[url "git@github.com:"]
    insteadOf = http://github.com/
[url "git@github.com:"]
    insteadOf = git://github.com/
[#57#root@ubuntu-1604 /go/src/git.llsapp.com/codelab]$
```

- after

```sh
[#60#root@ubuntu-1604 /go/src/git.llsapp.com/codelab]$git config --global --add url."git@git.llsapp.com:".insteadOf "https://git.llsapp.com/"

[#61#root@ubuntu-1604 /go/src/git.llsapp.com/codelab]$cat ~/.gitconfig
...
[url "git@github.com:"]
    insteadOf = https://github.com/
[url "git@github.com:"]
    insteadOf = http://github.com/
[url "git@github.com:"]
    insteadOf = git://github.com/
[url "git@git.llsapp.com:"]
    insteadOf = https://git.llsapp.com/
[#62#root@ubuntu-1604 /go/src/git.llsapp.com/codelab]$


[#62#root@ubuntu-1604 /go/src/git.llsapp.com/codelab]$go get -v git.llsapp.com/codelab/bazel_grpc
Fetching https://git.llsapp.com/codelab/bazel_grpc?go-get=1
Parsing meta tags from https://git.llsapp.com/codelab/bazel_grpc?go-get=1 (status code 200)
get "git.llsapp.com/codelab/bazel_grpc": found meta tag get.metaImport{Prefix:"git.llsapp.com/codelab/bazel_grpc", VCS:"git", RepoRoot:"https://git.llsapp.com/codelab/bazel_grpc.git"} at https://git.llsapp.com/codelab/bazel_grpc?go-get=1
git.llsapp.com/codelab/bazel_grpc (download)
package git.llsapp.com/codelab/bazel_grpc: no Go files in /go/src/git.llsapp.com/codelab/bazel_grpc
[#63#root@ubuntu-1604 /go/src/git.llsapp.com/codelab]$
```

### 其他

[Fix git error: Could not read Username for github.com](http://albertech.blogspot.com/2016/11/fix-git-error-could-not-read-username.html)


## 已经 push 到远端仓库后发现未 signoff

![DCO failed and resolve method](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/DCO%20failed%20and%20resolve%20method.png)

基本原理：

- 基于 rebase 命名添加 signoff
- 重新强制推送远端仓库

```
git rebase HEAD~n --signoff
git push --force <remote> <remote_branch>
```


## git 四个区和五种状态的切换

### git 的四种区

四种区：

- 工作区
- 暂存取
- 本地仓库
- 远程仓库

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Git%20%E5%90%84%E5%8C%BA%E5%92%8C%E7%8A%B6%E6%80%81%E5%8F%98%E6%9B%B4.png)

> 注意：图中最后两处的命令写反位置了

| 区域 | 命令 |
| -- | -- |
| 加入到**暂存区** | `git add .` (`git add <file>`) |
| 加入到**本地仓库** | `git commit -m "xxx"` |
| 加入到**远程仓库** | `git push origin master` |

### git 的五种状态

| 状态 | 描述 |
| -- | -- |
| Origin | 未修改 |
| Modified | 已修改 |
| Staged | 已暂存 |
| Committed | 已提交 |
| Pushed | 已推送 |

### 秒懂 Git 的区和状态

Ref: http://shengshui.com/?p=2582?zhihu

## git diff 使用

| 状态 | 命令 |
| -- | -- |
| 已 modify ，未 add  | `git diff` |
| 已 add ，未 commit  | `git diff --cached` |
| 已 commit ，未 push | `git diff <branch> origin/<branch>` |

## git 常用的撤销操作

| 状态 | 命令 |
| -- | -- |
| 已 modify ，未 add | 直接删文件就好了 |
| 已 add ，未 commit | `git reset HEAD <file>` 回到目标文件被 add 前的状态； <br> `git reset -–hard HEAD` 会产生目标文件被删除的效果；|
| 已 commit ，未 push | `git reset --hard origin/master` 基于远程仓库内容覆盖本地仓库内容；<br> `git reset -–hard HEAD` 不会产生效果，因为 HEAD 对应的就是最新的 commit ；<br> `git reset --hard HEAD^` 回退一个本地 commit ；<br> `git reset -–hard <commitID>` 回退一个本地指定 commit ；|
| 已 push | `git reset --hard <commitID>` 回退到本地指定的 commit ；<br> 如果后续需要对远程仓库的内容进行覆盖（因为本地经过修改后，commit tree 肯定和远程仓库不一致了），则必须通过 `git push -f` 进行强推) |

> http://www.netpi.me/uncategorized/gitrevoke/

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Git%20%E5%B8%B8%E7%94%A8%E7%9A%84%E6%92%A4%E9%94%80%E6%93%8D%E4%BD%9C.png)

基本状态标识

- `A-` = untracked 未跟踪
- `A` = tracked 已跟踪未修改
- `A+` = modified - 已修改未暂存
- `B` = staged - 已暂存未提交
- `C` = committed - 已提交未PUSH

各状态之间变化

- `A- -> B` : `git add <FILE>`
- `B -> A-` : `git rm --cached <FILE>`
- `B -> 删除不保留文件` : `git rm -f <FILE>`
- `A -> A-` : `git rm --cached <FILE>`
- `A -> A+` : 修改文件
- `A+ -> A` : `git checkout -- <FILE>`
- `A+ -> B` : `git add <FILE>`
- `B -> A+` : `git reset HEAD <FILE>`
- `B -> C` : `git commit`
- `C -> B` : `git reset --soft HEAD^`
- 修改最后一次提交: `git commit --amend`


## 拉取远程特定分支上的内容

有时候我们需要得到其它人的代码仓库，将别人的修改与自己的修改进行合并；

相关命令如下

```
git remote -v                                         -- 查看远程版本库信息
git remote add <name_this_remote> <remote_address>    -- 添加新的远程版本库
git checktout <target_branch>                         -- 切换到本地目标分支（即希望合并远程代码的分支）
git fetch <name_this_remote>                          -- 从新添加的远程版本库中获取最新的变更（no merge）
git merge <name_this_remote>/<target_branch>          -- 合并指定的远程分支到当前分支
```

> 注意：这里使用了 **先 fetch** + **后 merge** 的组合；似乎和直接使用 pull 没啥区别（待确认）

## 如何撤销错误的 merge 结果

假设你要 merge 远程分支 A 的内容，结果不小心 merge 了分支 B 上的内容；

```
git checktout <target_branch>
git fetch <name_this_remote>
git merge <name_this_remote>/B

# 此时可能会有一堆冲突要你解决，如果确定是一次错误 merge ，则按如下步骤操作

git reflog  # 确定错误 merge 之前当前 branch 所在 commit 对应 hash 值
git reset --hard <hash>
```

## 如何在 git 中写出好的 commit 说明

> 原文地址：[这里](https://ruby-china.org/topics/15737)

### 为什么要关注提交信息

- 加快 Reviewing Code 的过程；
- 帮助我们写好 release note ；
- 5 年后帮你快速想起来某个分支，`tag` 或者 `commit` 增加了什么功能，改变了哪些代码；
- 让其他的开发者在运行 `git blame` 的时候想跪谢；
- 总之一个好的提交信息，会帮助你提高项目的整体质量

### commit 编写的基本要求

- 第一行应该少于 50 个字。随后是一个空行；第一行题目也可以写成：`Fix issue #8976` ；
- 喜欢用 `vim` 的哥们把下面这行代码加入 `.vimrc` 文件中，来检查拼写和自动折行 `autocmd Filetype gitcommit setlocal spell textwidth=72`
- 永远不在 `git commit` 上增加 `-m <msg>` 或 `--message=<msg>` 参数，而单独写提交信息；

一个不好的例子 `git commit -m "Fix login bug"`

一个推荐的 commit message 应该是这样：

```
Redirect user to the requested page after login

https://trello.com/path/to/relevant/card

Users were being redirected to the home page after login, which is less
useful than redirecting to the page they had originally requested before
being redirected to the login form.

* Store requested path in a session variable
* Redirect to the stored location after successfully logging in the user
```

- 注释最好包含一个链接指向你们项目的 **issue**/**story**/**card**。一个完整的链接比一个 issue numbers 更好；
- 提交信息中包含一个简短的故事，能让别人更容易理解你的项目；

### commit 中要回答哪些问题

#### 为什么这次修改是必要的？

要告诉 Reviewers，你的提交包含什么改变。让他们更容易审核代码和忽略无关的改变。

#### 如何解决的问题？

这可不是说技术细节。看下面的两个例子：

```
Introduce a red/black tree to increase search speed
```

```
Remove <troublesome gem X>, which was causing <specific description of issue introduced by gem>
```

如果你的修改特别明显，就可以忽略这个。

#### 这些变化可能影响哪些地方？

这是你最需要回答的问题。因为它会帮你发现在某个 branch 或 commit 中的做了过多的改动。**一个提交尽量只做 1，2 个变化**。

你的团队应该有一个自己的行为规则，规定每个 commit 和 branch 最多能含有多少个功能修改。

### tips

- 使用 fix, add, change 而不是 fixed, added, changed ；
- 永远别忘了第 2 行是空行；
- 用 Line break 来分割提交信息，让它在某些软件里面更容易读；
- 请将每次提交限定于完成一次逻辑功能。并且可能的话，适当地分解为多次小更新，以便每次小型提交都更易于理解。

### 例子

```
Fix bug where user can't signup.

[Bug #2873942]

Users were unable to register if they hadn't visited the plans
and pricing page because we expected that tracking
information to exist for the logs we create after a user
signs up.  I fixed this by adding a check to the logger
to ensure that if that information was not available
we weren't trying to write it.
```

```
Redirect user to the requested page after login

https://trello.com/path/to/relevant/card

Users were being redirected to the home page after login, which is less
useful than redirecting to the page they had originally requested before
being redirected to the login form.

* Store requested path in a session variable
* Redirect to the stored location after successfully logging in the user
```

### 参考

- [A Note About Git Commit Messages](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html)
- [Writing good commit messages](https://github.com/erlang/otp/wiki/Writing-good-commit-messages)
- [A Better Git Commit](https://web-design-weekly.com/2013/09/01/a-better-git-commit/)
- [5 Useful Tips For A Better Commit Message](https://robots.thoughtbot.com/5-useful-tips-for-a-better-commit-message)
- [stopwritingramblingcommitmessages](http://stopwritingramblingcommitmessages.com/)
- [Git Commit Good Practice](https://wiki.openstack.org/wiki/GitCommitMessages)
- [AngularJS Git Commit Message Conventions](https://docs.google.com/document/d/1QrDFcIiPjSLDn3EL15IJygNPiHORgU1_OOAqWjiDU5Y/edit#heading=h.uyo6cb12dt6w)


## git merge 与 git rebase 

> 原文参考：[这里](http://blog.csdn.net/wh_19910525/article/details/7554489)

假设分支初始状态为

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/git%20merge%20vs%20git%20rebase%20-%201.jpg)

之后，我们在 mywork 分支上进行一些修改，生成两个提交 (commit) ；与此同时，其他人也在 origin 分支上做了一些修改，并且生成了两个提交；

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/git%20merge%20vs%20git%20rebase%20-%202.jpg)

这就意味着 origin 和 mywork 这两个分支各自"前进"了，即它们"分叉"了。

### git merge

通过 `git merge` 合并的效果如下：

```
$ git checkout mywork
$ git merge origin
```

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/git%20merge%20vs%20git%20rebase%20-%203.jpg)

实际效果为：使用 pull 命令把 origin 分支上的修改拉下来，并 mywork 分支上的修改进行合并；结果看起来就像一个新的“合并的提交” (merge commit) ；

> 在 mywork 分支历史中能看到 merge 相关信息；

### git rebase

通过 `git rebase` 合并的效果如下：

```
$ git checkout mywork
$ git rebase origin
```

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/git%20merge%20vs%20git%20rebase%20-%204.jpg)

实际效果为：rebase 命令会将 mywork 分支里的每个提交 (commit) 取消掉（和 origin 分叉之后的），并把它们临时保存为 ".git/rebase" 目录下补丁 (patch) ，之后将 mywork 分支更新为最新的 origin 分支，最后再把保存的那些补丁应用到 mywork 分支上。

当 mywork 分支更新完成之后，它会指向最新创建的 C6' 提交 (commit) ，而之前那些老的提交 (C5, C6) 会被丢弃。如果运行垃圾收集命令，之前那些被丢弃的提交会被删除；

> 让 mywork 分支历史看起来像没有经过任何合并一样；

### 对比

若执行 `git rebase` 时遇到冲突，则 git 会停止 rebase 过程并让你去解决冲突，在解决之后，需要通过 `git add` 命令更新相应内容的索引，之后无需执行 `git commit`，只需执行 `git rebase --continue` 继续应用余下的补丁；在任何时候，你都可以用 `git rebase --abort` 命令终止 rebase 行为，令 mywork 分支会回退到 rebase 开始前的状态。

若执行 `git merge` 时遇到冲突，则xxx；

**时间顺序问题**：

当我们使用 `git log` 查看 commit 信息时，会发现两种命令下 commit 顺序会有所不同。

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/git%20merge%20vs%20git%20rebase%20-%206.jpg)

假设

- C3 提交于 9:00AM
- C5 提交于 10:00AM
- C4 提交于 11:00AM
- C6 提交于 12:00AM

基于 `git merge` 合并后看到的 commit 顺序（从新到旧）是：`C7,C6,C4,C5,C3,C2,C1` ；

基于 `git rebase` 合并后看到的 commit 顺序（从新到旧）是：`C7,C6',C5',C4,C3,C2,C1` ；

因为 C6' 提交只是 C6 提交的克隆，C5' 提交只是 C5 提交的克隆；因此从用户角度看，就是 `C7,C6,C5,C4,C3,C2,C1` ；


## git tag 使用

- 显示已有标签

```shell
git tag
```

- 基于搜索模式列出符合条件的标签

```shell
git tag -l '<pattern>'
```

> 模式为正则表达式；

- 新建标签（基于当前提交点）

    > Git 使用的标签有两种类型：**轻量级的（lightweight）**和**含附注的（annotated）**；

    - 新建**含附注的标签**

        ```shell
        git tag -a <tag_name> -m "your comment"
        ```

    - 新建**基于 GPG 签署的标签**

        ```shell
        git tag -s <tag_name> -m "your comment"
        ```

        > 签署的目的是为了进行后续验证，防止篡改；

    - 新建**轻量级标签**

        ```shell
        git tag <tag_name>
        ```

- 查看标签相关信息

```shell
git show <tag_name>
```

- 基于 GPG 验证已经签署的标签

```shell
git tag -v <tag-name>
```

> 需要有签署者的公钥，存放在 keyring 中（即导入），才能验证；

- 后期（补）加注标签（只想某个特定提交）

```shell
git log --pretty=oneline                       # 查看提交历史，确定某次提交的哈希值
git tag -a <tag_name> <hashValue> -m "comment" # 可以只给出哈希值的前几位
```

- 将标签推送到远端仓库

    > 默认情况下，`git push` 并不会将标签推送到远端仓库；必须显式指定才行；

    - 推送指定标签

        ```shell
        git push origin <tag_name>
        ```

    - 推送所有标签

        ```shell
        git push origin --tags

        git push --tags      // 默认也是 push 到 origin
        ```

- 删除远端标签

```shell
git push origin --delete tag <tag_name>
git tag -d <tag_name>
git push origin :refs/tags/<tag_name>
```

- 获取远端标签

```shell
git fetch origin tag <tag_name>
```

> 尚未理解该命令的作用


## git clone 时直接 rename

> 可以直接指定路径，否则默认为 `./original_name`

```shell
git clone git@github.com:moooofly/original_name.git /path/to/new_name
```


## Create New Repository

> 以下内容取自 GitLab

Git global setup

```
git config --global user.name "moooofly"
git config --global user.email "centos.sf@gmail.com"
```

Create a new repository

```
git clone git@git.llsapp.com:fei.sun/scaffolding.git
cd scaffolding
touch README.md
git add README.md
git commit -m "add README"
git push -u origin master
```

Existing folder

```
cd existing_folder
git init
git remote add origin git@git.llsapp.com:fei.sun/scaffolding.git
git add .
git commit -m "Initial commit"
git push -u origin master
```

另

```
cd existing_folder
git init
git add .
git commit -m "first commit"
git remote add origin git@github.com:moooofly/your_project_name.git
(git pull origin master)
git push -u origin master
```

Existing Git repository

```
cd existing_repo
git remote rename origin old-origin
git remote add origin git@git.llsapp.com:fei.sun/scaffolding.git
git push -u origin --all
git push -u origin --tags
```


> 注1：执行 `git remote add` 前，首先要在 github 上创建一个名为 your_project_name 的 repo ；
> 
> 注2：上面用小括号括起来的命令的使用场景为：若在 github 上新建 repo 的时候，顺带创建了 README 或 .gitignore 或 LICENSE 等文件时，则需要先将上述文件拉取到本地；
> 
> 注3：上面的 git@github.com:moooofly/your_project_name.git 可以换成 https://github.com/moooofly/your_project_name.git ，还可以使用 https://github.com/moooofly/your_project_name
> 
> 注4：`git remote add` 后就可以进行 pull 了，但仍无法 push ；需要通过 `git push -u` 或 `git push --set-upstream` 命令，在 push 的同时，建立跟踪关系；
> 
> 注5：上面执行 `git pull origin master` 时可能会报 "fatal: refusing to merge unrelated histories" 错误，此时可以使用 `--allow-unrelated-histories` 选项解决，即 `git pull origin master --allow-unrelated-histories` ；详情参见[这里](https://stackoverflow.com/questions/37937984/git-refusing-to-merge-unrelated-histories/40107973#40107973?newreg=5095f8141c34479ba419f5e8b2d1b415)；


## 本地新建分支后 push 到远端仓库

基于当前所在本地分支创建并切换到新分支（新分支 hash 位置和源分支 hash 位置相同）

```shell
git checkout -b new_branch
```
此时分支信息**仅在本地存在**；

**若 github 上尚不存在 new_branch 分支**，则通过执行下面的命令，就可以将新建的本地 new_branch 分支中的内容 push 到 github 上对应的 new_branch 分支上（包含创建行为），并**建立跟踪关系**；

```shell
git push -u  # 下面命令的省略写法
或
git push --set-upstream origin new_branch
```

之后就可以在 github 页面中看到对应分支内容了；

**若 github 上已经存在 new_branch 分支**，那么你可能只是想要将本地的 new_branch 分支与其建立跟踪关系，可以执行如下命令；

```shell
git branch --set-upstream-to=origin/new_branch new_branch
```

## 获取指定 tag 代码

如果本地有代码仓库

```
git tag
git checkout <tag_name>
```

如果本地没有代码仓库

```
git clone git@xxx.xxx.xxx:/project_name.git
git tag
git checkout <tag_name>
```

> `git checkout <tag_name>` 也写作 `git checkout tags/<tag_name>`

注意：

- 通过上述办法切换到目标 tag 时，会提示：当前处于一个 "detached HEAD" 状态；我们知道每一个 tag 都只是代码仓库中的一个快照（commit），因此，如果你想基于此 tag 进行代码开发，上面的方法就太合适了；
- 若想基于 tag 快照进行代码开发，你需要把 tag 快照对应的代码拉取（关联）到一个新分支上，例如，我想基于 tag v1.6 的代码进行编辑，则需要创建一个分支 branch-v1.6 ，然后把 v1.6 的代码拉取（关联）到此分支上；

```
# git checkout <tag_name> -b <branch_related_to_tag_name>
git checkout v1.6 -b branch-v1.6 
```


## 直接获取远端仓库指定 branch 代码

`git clone` 默认会把远程仓库整个给 clone 下来，但默认只会在本地创建一个 master 分支；

`git clone` 可以使用 `-b` 参数指定目标分支。实际上，`-b` 参数不仅支持**分支名**，还支持 **tag** 和 **commit** 。

```
git clone <remote-addr:repo.git> -b <branch-or-tag-or-commit>

或

git clone <remote-addr:repo.git> -b <branch_name>
git clone <remote-addr:repo.git> -b <tag_name>
git clone <remote-addr:repo.git> -b <commit_hash>
```

## 将远端仓库里指定 branch 拉取到本地


> 前提：
> 
> - 本地不存在该分支；
> - 已经通过 `git clone` 获取远程仓库；

当想要从远端仓库里拉取一条本地不存在的分支时，可以通过如下命令查看远端仓库中包含哪些分支

```
# 查看所有分支（本地的和远端的）
git branch -av

或

# 功能同 -av 但输出内容可直接用于下面的命令
git branch -r
```

基于**远端特定分支**创建并切换到本地新分支，自动建立分支关联（新本地分支 hash 位置和源远端分支 hash 位置相同）

```
git checkout -b local_branch_name origin/remote_branch_name
```

如果出现提示：

```
fatal: Cannot update paths and switch to branch 'aaa' at the same time.
Did you intend to checkout 'origin/bbb' which can not be resolved as commit?
```

需要先执行

```
git fetch
```

再执行

```
git checkout -b local_branch_name origin/remote_branch_name
```


## 重命名 remote branch 名

当在本地执行过如下命令后，你将会创建一个本地分支 old_branch 并且关联到远程的 old_branch 分支上；

```shell
git checkout -b old_branch
git push -u
或
git push -u origin old_branch
```

此时若想更改远程分支的名字，则可以按如下方式进行操作：

- 修改本地分支名字

```shell
git branch -m old_branch new_branch 
```

- 删除远程待修改分支名（其实就是推送一个空分支到远程分支，以达到删除远程分支的目的）

```shell
git push origin :old_branch
```

在 Git v1.7.0 之后，可以使用这种语法删除远程分支

```shell
git push origin --delete <remote_branch_name>
```

> 删除本地分支命令为（需要切换到其他分支上执行该命令）
> ```shell
> # delete fully merged branch
> git branch -d <local_branch_name>
>
> # delete branch (even if not merged)
> git branch -D <local_branch_name>
> ```

- 将本地的新分支 push 到远程

```shell
git push -u
或
git push -u origin new_branch
```


分支信息查看说明：

> - 查看本地分支名字
>
> ```shell
> git branch
> ```
> 
> - 查看本地和远程分支名字（红色显示部分为远程分支）
> 
> ```shell
> git branch -a
> ```
> 
> - 查看本地和远程分支名字，会显示出本地和远程的 tracking 关系；
> 
> ```shell
> git branch -a -vv
> ```


## 删除不存在对应远程分支的本地分支

一种情况：提交 PR 后，远端 master 分支在 PR 合并完成后，一般会直接删除对应的 PR 分支，而提交 PR 的人在本地会看到如下提示信息；

```shell
➜  redis_dissector_for_wireshark git:(master)  git branch -a  -- 该命令看不出问题
* master
  revert
  remotes/origin/master
  remotes/origin/revert
➜  redis_dissector_for_wireshark git:(master)
➜  redis_dissector_for_wireshark git:(master) git remote show origin
* remote origin
  Fetch URL: git@github.com:moooofly/redis_dissector_for_wireshark.git
  Push  URL: git@github.com:moooofly/redis_dissector_for_wireshark.git
  HEAD branch: master
  Remote branches:
    master                     tracked
    refs/remotes/origin/revert stale (use 'git remote prune' to remove) -- 这里
  Local branch configured for 'git pull':
    master merges with remote master
  Local ref configured for 'git push':
    master pushes to master (local out of date)
➜  redis_dissector_for_wireshark git:(master)
```

上述信息表明：
> remote 分支 revert 处于 stale 状态（过时）

两种解决方法：
- 使用 `git remote prune origin` 将对应的分支关联信息从本地版本库中去除；
- 更简单的方法是使用 `git fetch -p` 命令，在 fetch 之后，自动删除掉没有与远程分支对应的本地分支；

输出结果如下

```shell
➜  redis_dissector_for_wireshark git:(master) git fetch -p
From github.com:moooofly/redis_dissector_for_wireshark
 x [deleted]         (none)     -> origin/revert
remote: Counting objects: 1, done.
remote: Total 1 (delta 0), reused 0 (delta 0), pack-reused 0
Unpacking objects: 100% (1/1), done.
   ea2ef49..193304e  master     -> origin/master
➜  redis_dissector_for_wireshark git:(master)
➜  redis_dissector_for_wireshark git:(master) git branch -a
* master
  revert
  remotes/origin/master
➜  redis_dissector_for_wireshark git:(master)
➜  redis_dissector_for_wireshark git:(master) git remote show origin
* remote origin
  Fetch URL: git@github.com:moooofly/redis_dissector_for_wireshark.git
  Push  URL: git@github.com:moooofly/redis_dissector_for_wireshark.git
  HEAD branch: master
  Remote branch:
    master tracked
  Local branch configured for 'git pull':
    master merges with remote master
  Local ref configured for 'git push':
    master pushes to master (local out of date)
➜  redis_dissector_for_wireshark git:(master)
```

最后还需要通过 `git branch -D revert` 删除 local 分支；
```shell
➜  redis_dissector_for_wireshark git:(master) git branch -D revert
Deleted branch revert (was be9b8d1).
➜  redis_dissector_for_wireshark git:(master)
```


## fork 别人项目后如何同步其后续更新

先从自己的 github 中 clone 一份内容
```shell
git clone git@github.com:moooofly/sre.git
```

此时只有名为 origin 的 remote
```shell
git remote -v
```

新增名为 eleme_sre 的 remote ，即当前 repo 的始祖；
```shell
git remote add eleme_sre https://github.com/eleme/sre.git
```

再次查看时，已经增加了名为 eleme_sre 的 remote ；
```shell
git remote -v
```

从 eleme_sre 的 master 分支上拉取内容到本地；
```shell
git pull eleme_sre master
```

将拉取到本地的内容推到名为 origin 的 master 分支，即推到自己的 github 仓库中；
```shell
git push origin master
```

> 上述命令最好写完整，否则容易引起混乱或歧义；


## 定制化 git 全局配置

取自：[SRE 团队 git 配置参考](https://github.com/eleme/sre/blob/master/git.md)

```shell
[color]

    ui = auto

[alias]

    lg1 = log --graph --all --format=format:'%C(bold blue)%h%C(reset) - %C(bold green)(%ar)%C(reset) %C(white)%s%C(reset) %C(bold white)— %an%C(reset)%C(bold yellow)%d%C(reset)' --abbrev-commit --date=relative

    lg2 = log --graph --all --format=format:'%C(bold blue)%h%C(reset) - %C(bold cyan)%aD%C(reset) %C(bold green)(%ar)%C(reset)%C(bold yellow)%d%C(reset)%n''          %C(white)%s%C(reset) %C(bold white)— %an%C(reset)' --abbrev-commit

[core]

    editor = vim

    safecrlf = true

    excludesfile = ~/.gitignore

[push]

    default = current

[rerere]

    enabled = 1

    autoupdate = 1

[user]

    name = your-name

    email = your-email

[merge]

    tool = vimdiff

[url "git@github.com:"]

    insteadOf = https://github.com/

[url "git@github.com:"]

    insteadOf = http://github.com/

[url "git@github.com:"]

    insteadOf = git://github.com/
```



## 基于 ssh 协议访问 git

```shell
vagrant@vagrant-ubuntu-trusty:~$ ll .ssh

total 16
drwx------ 2 vagrant vagrant 4096 Jun  6 09:30 ./
drwxr-xr-x 7 vagrant vagrant 4096 Jun  6 09:38 ../
-rw------- 1 vagrant vagrant  389 Jun  6 08:18 authorized_keys
-rw-r--r-- 1 vagrant vagrant  884 Jun  6 09:30 known_hosts
vagrant@vagrant-ubuntu-trusty:~$
```


创建 RSA 密钥对
```shell
vagrant@vagrant-ubuntu-trusty:~$ ssh-keygen -t rsa -b 4096 -C "aaa@bbb.com"
Generating public/private rsa key pair.
Enter file in which to save the key (/home/vagrant/.ssh/id_rsa):
Enter passphrase (empty for no passphrase):               -- 最好不使用密码
Enter same passphrase again:

Your identification has been saved in /home/vagrant/.ssh/id_rsa.
Your public key has been saved in /home/vagrant/.ssh/id_rsa.pub.

The key fingerprint is:
f1:a9:8a:a3:23:18:64:3e:5b:3a:e9:39:54:32:23:83 aaa@bbb.com

The key's randomart image is:
+--[ RSA 4096]----+
|                 |
|                 |
|.       .        |
|E* .     o .     |
|=.=     S o      |
|.+ .     .       |
|o.*     .        |
|o*o .. .         |
|.++o...          |
+-----------------+
vagrant@vagrant-ubuntu-trusty:~$
```

```shell
vagrant@vagrant-ubuntu-trusty:~$ ll .ssh

total 24
drwx------ 2 vagrant vagrant 4096 Jun  6 09:40 ./
drwxr-xr-x 7 vagrant vagrant 4096 Jun  6 09:38 ../
-rw------- 1 vagrant vagrant  389 Jun  6 08:18 authorized_keys
-rw------- 1 vagrant vagrant 3243 Jun  6 09:40 id_rsa
-rw-r--r-- 1 vagrant vagrant  745 Jun  6 09:40 id_rsa.pub
-rw-r--r-- 1 vagrant vagrant  884 Jun  6 09:30 known_hosts
vagrant@vagrant-ubuntu-trusty:~$
```

```shell
vagrant@vagrant-ubuntu-trusty:~/workspace/eleme_project$ eval "$(ssh-agent -s)"
Agent pid 1371
vagrant@vagrant-ubuntu-trusty:~/workspace/eleme_project$ ssh-add ~/.ssh/id_rsa
Identity added: /home/vagrant/.ssh/id_rsa (/home/vagrant/.ssh/id_rsa)
vagrant@vagrant-ubuntu-trusty:~/workspace/eleme_project$
```

最后将 id_rsa.pub 文件中的内容添加到 github 账户中


## 其他

### GitHub 秘籍

Ref: https://snowdream86.gitbooks.io/github-cheat-sheet/content/zh/index.html

### Git 简明指南

Ref: http://rogerdudler.github.io/git-guide/index.zh.html

适合用作新手的 PPT

### Git Cheat Sheet

Ref: http://rogerdudler.github.io/git-guide/files/git_cheat_sheet.pdf

### 图解 Git

Ref: http://marklodato.github.io/visual-git-guide/index-zh-cn.html

### Think like a Git

Ref: http://think-like-a-git.net/

### Git 社区参考书

Ref: https://book.git-scm.com/

## Tools

### Tower

Ref: https://www.git-tower.com/mac?source=rd


----------


**[⬆ 返回顶部](#目录)**
