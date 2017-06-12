# GIT 写出好的 commit message

> 原文地址：[这里](https://ruby-china.org/topics/15737)

## 为什幺要关注提交信息

- 加快 Reviewing Code 的过程；
- 帮助我们写好 release note ；
- 5 年后帮你快速想起来某个分支，`tag` 或者 `commit` 增加了什么功能，改变了哪些代码；
- 让其他的开发者在运行 `git blame` 的时候想跪谢；
- 总之一个好的提交信息，会帮助你提高项目的整体质量

## 基本要求

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

## 注释要回答如下信息

### 为什么这次修改是必要的？

要告诉 Reviewers，你的提交包含什么改变。让他们更容易审核代码和忽略无关的改变。

### 如何解决的问题？

这可不是说技术细节。看下面的两个例子：

```
Introduce a red/black tree to increase search speed
```

```
Remove <troublesome gem X>, which was causing <specific description of issue introduced by gem>
```

如果你的修改特别明显，就可以忽略这个。

### 这些变化可能影响哪些地方？

这是你最需要回答的问题。因为它会帮你发现在某个 branch 或 commit 中的做了过多的改动。**一个提交尽量只做 1，2 个变化**。

你的团队应该有一个自己的行为规则，规定每个 commit 和 branch 最多能含有多少个功能修改。

### 小提示

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

### 本文参考阅读

- [A Note About Git Commit Messages](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html)
- [Writing good commit messages](https://github.com/erlang/otp/wiki/Writing-good-commit-messages)
- [A Better Git Commit](https://web-design-weekly.com/2013/09/01/a-better-git-commit/)
- [5 Useful Tips For A Better Commit Message](https://robots.thoughtbot.com/5-useful-tips-for-a-better-commit-message)
- [stopwritingramblingcommitmessages](http://stopwritingramblingcommitmessages.com/)
- [Git Commit Good Practice](https://wiki.openstack.org/wiki/GitCommitMessages)
- [AngularJS Git Commit Message Conventions](https://docs.google.com/document/d/1QrDFcIiPjSLDn3EL15IJygNPiHORgU1_OOAqWjiDU5Y/edit#heading=h.uyo6cb12dt6w)




