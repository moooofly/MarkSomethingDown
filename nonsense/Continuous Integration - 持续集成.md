# Continuous Integration - 持续集成

## [使用 Travis 进行持续集成](http://www.liaoxuefeng.com/article/0014631488240837e3633d3d180476cb684ba7c10fda6f6000)

持续集成：Continuous Integration，简称 CI ，意思是，在一个项目中，任何人对代码库的任何改动，都会触发 CI 服务器自动对项目进行构建，自动运行测试，甚至自动部署到测试环境。这样做的好处就是，随时发现问题，随时修复。因为修复问题的成本随着时间的推移而增长，越早发现，修复成本越低。

Travis CI 是在线托管的 CI 服务，用 Travis 来进行持续集成，不需要自己搭服务器，在网页上点几下就好，用起来更方便。最重要的是，它对开源项目是免费的。

因为 GitHub 和 Travis 是一对好基友，不用 GitHub 虽然也能用 Travis ，但是配置起来太麻烦。而且，作为开源项目，为什么不用 GitHub ？

需要编写一个 `.travis.yml` 文件来告诉 Travis 一些项目信息；

> 是不是用了 CI 代码质量就有保证了？

这个问题的答案是：否。如果 CI 能提高代码质量，那软件公司只需要招实习生写代码外加 CI 就可以了，招那么贵的高级工程师浪费钱干啥？

> 是不是用了 Travis 就实现了 CI ？

这个问题的答案还是：否。CI 是解决问题的手段而不是目的。问题是如何提高代码质量。我见过很多公司的项目，编译一次半小时（不是编译 Linux 内核那种），测试一次几个小时。不能在短时间内完成编译、测试的代码都有问题，导致 CI 形同虚设。这里的“短时间”是指 5 分钟以内。


## [使用 Jenkins 进行持续集成](http://www.liaoxuefeng.com/article/001463233913442cdb2d1bd1b1b42e3b0b29eb1ba736c5e000)

开源项目如何利用 Travis CI 进行持续集成；如果你的项目不是开源项目，可以自己搭建 CI 环境，即利用 Jenkins CI 进行持续集成；

Jenkins官方网站：https://jenkins.io/

虽然 Jenkins 提供了 Windows、Linux、OS X 等各种安装程序，但是，这些安装程序都没有 war 包好使。只需要运行命令：

```
java -jar jenkins.war
```

浏览器登录：http://localhost:8080/


## 其他

- [Core Concepts for Beginners](https://docs.travis-ci.com/user/for-beginners)
- [Getting started](https://docs.travis-ci.com/user/getting-started/)
- [Building a Go Project](https://docs.travis-ci.com/user/languages/go)
- [Customizing the Build](https://docs.travis-ci.com/user/customizing-the-build/)








