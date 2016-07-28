

# git clone 时直接 rename

```shell
git clone git@github.com:moooofly/aaa.git bbb
```


先在 github 上创建一个名为 your_project_name 的 repo


## create a new repository on the command line

```shell
cd /path/to/project/dir/
git init
git add .
git commit -m "first commit"
git remote add origin git@github.com:moooofly/your_project_name.git
(git pull origin master)
git push -u origin master
```
上面用小括号括起来的命令的使用场景为：若在 github 上新建 repo 的时候，顺带创建了 README 或 .gitignore 或 LICENSE 等文件时，则需要先将上述文件拉取到本地；

## push an existing repository from the command line

```shell
git remote add origin git@github.com:moooofly/your_project_name.git
(git pull origin master)
git push -u origin master
```

 PS：上面的 git@github.com:moooofly/rabbitmq-server-3.6.1.git 可以换成 https://github.com/moooofly/rabbitmq-server-3.6.1.git ，还可以使用 https://github.com/moooofly/rabbitmq-server-3.6.1


# 本地创建分支后关联到 github


创建并切换分支
```shell
git checkout -b new_branch
```
此时分支信息仅在本地存在；

执行下面的命令后就可以将 new_branch 中的内容 push 到 github 上对应的 new_branch 分支中，并建立跟踪关系；
```shell
git push -u
```

