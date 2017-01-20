# 小女孩也能看懂的插画版 Kubernetes 指南

标签（空格分隔）： Kubernetes

---

> Matt Butcher 是 Deis 的平台架构师，热爱哲学，咖啡和精雕细琢的代码。有一天女儿走进书房问他什么是 Kubernetes，于是就有了这本插画版的 Kubernetes 指南，讲述了勇敢的 Phippy（一个 PHP 应用），在 Kubernetes 的冒险故事，满满的父爱有木有！

----------

- 某一天

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/kubernetes-illustrated-guide-illustration-1.png)

> 有一天，女儿走进书房问我：『亲爱的爸爸，什么是 Kubernetes 呢？』

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/kubernetes-illustrated-guide-diagram-1.png)

> 我回答她：『Kubernetes 是一个开源的 Docker 容器编排系统，它可以调度计算集群的节点，动态管理上面的作业，保证它们按用户期望的状态运行。通过使用「labels」和「pods」的概念，Kubernetes 将应用按逻辑单元进行分组，方便管理和服务发现。』
>
> 女儿更疑惑了……于是就有了这个故事。

- 给孩子的插画版 Kubernetes 指南

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/kubernetes-illustrated-guide-illustration-2.png)

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/kubernetes-illustrated-guide-illustration-3.png)

> 很久很久以前，有一个叫 Phippy 的 PHP 应用，她很单纯，只有一个页面。她住在一个托管服务里，周围还有很多可怕的应用，她都不认识，也不想去认识，但是他们却要共享这里的环境。所以，她一直都希能有一个属于自己的环境：一个可以称作 home 的 webserver。

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/kubernetes-illustrated-guide-diagram-2.png)

*每个应用的运行都要依赖一个环境，对于一个 PHP 应用来说，这个环境包括了一个 webserver，一个可读的文件系统和 PHP 的 engine。*

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/kubernetes-illustrated-guide-illustration-4.png)

> 有一天，一只可爱的鲸鱼拜访了 Phippy，他建议 Phippy 住在容器里。Phippy 听从了鲸鱼的建议搬家了，虽然这个容器看起来很好，但是……怎么说呢，就像是漂浮在海上的一个小房间一样，还是没有家的感觉。

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/kubernetes-illustrated-guide-diagram-3.png)

*不过容器倒是为应用提供了隔离的环境，在这个环境里应用就能运行起来。但是这些相互隔离的容器需要管理，也需要跟外面的世界沟通。共享的文件系统，网络，调度，负载均衡和资源分配都是挑战。*

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/kubernetes-illustrated-guide-illustration-5.png)

> 『抱歉……孩子……』鲸鱼耸耸肩，一摇尾消失在了海平面下…… Phippy 还没有来得及失望，就看到远方驶来一艘巨轮，掌舵的老船长非常威风。这艘船乍一看就是大了点，等到船走近了，Phippy 才发现船体两边挂满了皮筏。
> 
> 老船长用充满智慧的语气对 Phippy 说：『你好，我是 Kube 船长』。

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/kubernetes-illustrated-guide-diagram-4.png)

*『Kubernetes』是希腊语中的船长，后来的『Cybernetic』和『Gubernatorial』这两个词就是从 Kubernetes 衍生来的。Kubernetes 项目由 Google 发起，旨在为生产环境中成千上万的容器，构建一个健壮的平台。*

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/kubernetes-illustrated-guide-illustration-6.png)

> 『您好，我是 Phippy。』
> 
> 『很高兴认识你。』船长边说，边在 Phippy 身上放了一个 name tag。

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/kubernetes-illustrated-guide-diagram-5.png)

*Kubernetes 使用 label 作为『nametag』来区分事物，还可以根据 label 来查询。label 是开放式的：可以根据角色，稳定性或其它重要的特性来指定。*

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/kubernetes-illustrated-guide-illustration-7.png)

> Kube 船长建议 Phippy 可以把她的容器搬到船上的 pod 里，Phippy 很高兴地接受了这个提议，把容器搬到了 Kube 的大船上。Phippy 感觉自己终于有家了。

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/kubernetes-illustrated-guide-diagram-6.png)

*在 Kubernetes 中，pod 代表着一个运行着的工作单元。通常，每个 pod 中只有一个容器，但有些情况下，如果几个容器是紧耦合的，这几个容器就会运行在同一个 pod 中。Kubernetes 承担了 pod 与外界环境通信的工作。*

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/kubernetes-illustrated-guide-illustration-8.png)

> Phippy 对这一切都感到很新奇，同时她也有很多与众不同的关注点：『如果我想要复制自己该怎么做呢？按需的……任意次数的可以吗？』
> 
> 『很简单。』船长说道，接着就给 Phippy 介绍起了 replication controller。

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/kubernetes-illustrated-guide-diagram-7.png)

*Replication controller 提供了一种管理任意数量 pod 的方式。一个 replication controller 包含了一个 pod 模板，这个模板可以被不限次数地复制。通过 replication controller，Kubernetes 可以管理 pod 的生命周期，包括扩/缩容，滚动部署和监控等功能。*

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/kubernetes-illustrated-guide-illustration-9.png)

> Phippy 就这样在船上和自己的副本愉快地生活了好多天。但是每天只能面对自己的副本，这样的生活也太孤单了。
> 
> Kube 船长慷慨地笑道：『我有好东西给你。』
> 
> 说着，Kube 船长就在 Phippy 的 replication controller 和船上其它地方之间建了一个隧道：『就算你们四处移动，这个隧道也会一直待在这里，它可以帮你找到其它 pod，其它 pod 也可以找到你。』

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/kubernetes-illustrated-guide-diagram-8.png)

*service 可以和 Kubernetes 环境中其它部分（包括其它 pod 和 replication controller）进行通信，告诉它们你的应用提供什么服务。Pod 可以四处移动，但是 service 的 IP 地址和端口号是不变的。而且其它应用可以通过 Kubernetes 的服务发现找到你的 service。*

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/kubernetes-illustrated-guide-illustration-10.png)

> 有了 service，Phippy 终于敢去船上其它地方去玩了，她很快就有了新朋友 Goldie。有一天，Goldie 送了 Phippy 一件礼物，没想到 Phippy 只看了一眼就哭了。
> 
> 『你怎么哭了？』Goldie 问道。
> 
> 『我太喜欢这个礼物了，可惜没地儿放……』Phippy 都开始抽泣了。Goldie 一听原来是这么回事，马上就告诉 Phippy：『为什么不放在一个 volume 里呢？』

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/kubernetes-illustrated-guide-diagram-9.png)

*Volume 代表了一块容器可以访问和存储信息的空间，对于应用来说，volume 是一个本地的文件系统。实际上，除了本地存储，Ceph、Gluster、Elastic Block Storage 和很多其它后端存储都可以作为 volume。*

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/kubernetes-illustrated-guide-illustration-11.png)

> Phippy 渐渐地爱上了船上的生活，她很享受和新朋友的相处（Goldie 的每个 pod 副本也都很 nice）。但是回想起以前的生活，她又在想是不是可以有一点点私人空间呢？
> 
> Kube 船长很理解：『看起来你需要 namespace。』

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/kubernetes-illustrated-guide-diagram-10.png)

*Namespace 是 Kubernetes 内的分组机制。Service，pod，replication controller 和 volume 可以很容易地和 namespace 配合工作，但是 namespace 为集群中的组件间提供了一定程度的隔离。*

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/kubernetes-illustrated-guide-illustration-12.png)

> 于是，在 Kube 船长的船上，Phippy 和她的朋友们开始了海上的历险，最重要的是，Phippy 找到了自己的家。
> 
> 从此，Phippy 过上了幸福的生活。

