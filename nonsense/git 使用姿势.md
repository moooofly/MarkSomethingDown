

# git tag 操作

- 显示已有标签

```shell
git tag
```

- 基于搜索模式列出符合条件的标签

```shell
git tag -l '<pattern>'
```

> 模式为正则表达式；

- 新建标签

    > Git 使用的标签有两种类型：轻量级的（lightweight）和含附注的（annotated）；

    - 新建含附注的标签

        ```shell
        git tag -a <tag_name> -m "your comment"
        ```

    - 新建基于 GPG 签署的标签

        ```shell
        git tag -s <tag_name> -m "your comment"
        ```

        > 签署的目的是为了进行后续验证，防止篡改；

    - 新建轻量级标签

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

- 后期（补）加注标签

```shell
git log --pretty=oneline           # 查看提交历史，确定某次提交的哈希值
git tag -a <tag_name> <hashValue>  # 可以只给出哈希值的前几位
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



# 删除不存在对应远程分支的本地分支

一种情况：提交 PR 后，master 分支在合并完成后会删除对应的 PR 分支，而提交 PR 的人的本地分支会看到如下提示信息；

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
- 使用 `git remote prune origin` 可以将其从本地版本库中去除；
- 更简单的方法是使用 `git fetch -p` 命令，在 fetch 之后删除掉没有与远程分支对应的本地分支；

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

# 定制化 git 全局配置

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



# 基于 SSH 协议访问 git

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


----------


# 问题



