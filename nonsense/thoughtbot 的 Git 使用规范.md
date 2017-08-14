# [Git Protocol](https://github.com/thoughtbot/guides/tree/master/protocol/git)

A guide for programming within version control.

## Maintain a Repo

* 避免在源码控制中包含特定于开发机器或 process 的文件；
* 在 merge 完成后，请删除 local 和 remote 上的 feature branches ；
* 在 feature branch 开展开发工作；
* 经常 rebase 以保证及时将 upstream 上发生的变更纳入到本地；
* 基于 [pull request] 进行 code reviews ；

[pull request]: https://help.github.com/articles/using-pull-requests/

## Write a Feature

- 基于 master 分支（最新代码）创建一个本地 feature branch ；

```
git checkout master
git pull
git checkout -b <branch-name>
```

- 经常 rebas码，以保证与 upstream 上的代码一致；

```
git fetch origin
git rebase origin/master
```

解决冲突；

- 当 feature 开发完成、测试通过后，stage 所有变更；

```
git add --all
```

之后进行 commit ；

```
git status
git commit --verbose
```

编写一份 [good commit message] ，示例格式如下：

```
    Present-tense summary under 50 characters

    * More information about commit (under 72 characters).
    * More information about commit (under 72 characters).

    http://project.management-system.com/ticket/123
```

- 如果在 feature 开发中创建了不止一个 commit ，则 [use `git rebase` interactively](https://help.github.com/articles/about-git-rebase/)
以 squash 这些 commits ，使之成为更加内聚的 (cohesive)、更具可读性的 commits ：

```
git rebase -i origin/master
```

- 推送该分支到你自己的远程仓库；

```
git push origin <branch-name>
```

- 提交 [GitHub pull request] ；

在项目的 chat room 中寻求 code review ；

[good commit message]: http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html
[GitHub pull request]: https://help.github.com/articles/using-pull-requests/

## Review Code

除了 author 外的其他组员也要进行 pull request 的 review ；所有人都需要遵循 [Code Review](https://github.com/thoughtbot/guides/tree/master/code-review) 准则以免 miscommunication ；

所有人应该基于 GitHub 的 web 接口直接在代码行上进行评论和提问，也可以在项目专门的 chat room 中进行沟通；

针对 review 者自己进行变更的情况，可以

```
git checkout <branch-name>
./bin/setup
git diff staging/master..HEAD
```

在当前分支上进行小变更，进行 feature 测试，运行 tests，之后 commit 和 push ；

当满足要求后，在 pull request 时直接填写 `Ready to merge.` ；

## Merge

**Rebase interactively**；将类似 "Fix whitespace" 这种 commits 进行 Squash ，将其变成一条或少量有价值的 commit(s) ；编辑 commit messages 以正确显示意图；运行测试；

```
git fetch origin
git rebase -i origin/master
```

**Force push your branch**. 当你的 commit(s) 被 push 到 master 上时，这允许 GitHub 自动关闭你的 pull
request 并将其标记为 merged ；这同样令 [find the pull request] 成为可能；

```
git push --force-with-lease origin <branch-name>
```

查看新 commits 列表；查看发生变更的文件；Merge 分支内容到 master ；

```
git log origin/master..<branch-name>
git diff --stat origin/master
git checkout master
git merge <branch-name> --ff-only
git push
```

删除你的远程 feature branch ；

```
git push origin --delete <branch-name>
```

删除你的本地 feature branch ；

```
git branch --delete <branch-name>
```

[find the pull request]: http://stackoverflow.com/a/17819027


----------


- [Git 使用规范流程](http://www.ruanyifeng.com/blog/2015/08/git-use-process.html)
- [Git Interactive Rebase, Squash, Amend and Other Ways of Rewriting History](https://robots.thoughtbot.com/git-interactive-rebase-squash-amend-rewriting-history)