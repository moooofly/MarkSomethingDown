

# git clone 时直接 rename

```shell
git clone git@github.com:moooofly/aaa.git bbb
```


----------



# 基于本地项目创建 github repo

```shell
cd /path/to/project/dir/
git init
git add .
git commit -m "first commit"
git remote add origin git@github.com:moooofly/your_project_name.git
(git pull origin master)
git push -u origin master
```

> 注1：在执行 `git remote add xxx` 前，需要先在 github 上创建一个名为 your_project_name 的 repo ；
> 
> 注2：上面用小括号括起来的命令的使用场景为：若在 github 上新建 repo 的时候，顺带创建了 README 或 .gitignore 或 LICENSE 等文件时，则需要先将上述文件拉取到本地；
> 
> 注3：上面的 git@github.com:moooofly/your_project_name.git 可以换成 https://github.com/moooofly/your_project_name.git ，还可以使用 https://github.com/moooofly/your_project_name


# 本地新建分支后 push 到 github repo


创建并切换分支

```shell
git checkout -b new_branch
```
此时分支信息**仅在本地存在**；

执行下面的命令后，就可以将新建的本地 new_branch 分支中的内容 push 到 github 上对应的 new_branch 分支中，并**建立跟踪关系**；

```shell
git push -u
```

> 上面命令隐含了 `--set-upstream-to` 动作；

之后就可以在 github 页面中看到对应分支内容了；

若 github 上已经存在 new_branch 分支，想要将本地的某个分支与其建立跟踪关系，可以执行如下命令；

```shell
git branch --set-upstream-to=origin/new_branch new_branch
```

# 重命名 github repo 中的远程分支名


当在本地执行过如下命令后，你将会创建一个本地分支 old_branch 并且关联到远程的 old_branch 分支上；

```shell
git checkout -b old_branch
git push -u
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
> git branch -d <local_branch_name>
> ```

- 将本地的新分支 push 到远程

```shell
git push -u
```


分支信息查看说明：

> ```shell
> git branch
> ```
> 查看本地分支名字
> ```shell
> git branch -a
> ```
> 查看本地和远程分支名字（红色显示部分为远程分支）
> ```shell
> git branch -a -vv
> ```
> 查看本地和远程分支名字，会显示出本地和远程的 tracking 关系；








# fork 别人项目后如何同步其后续更新

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
git push
```


----------


# 问题

- 是否设置 --set-upstream-to 有何区别？
- 

