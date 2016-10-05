




# Netflix

关于 Netflix 的介绍，可以直接看如下三篇文章：
- [【转载】跟花和尚学系统设计：明星公司之Netflix(上篇)](https://my.oschina.net/moooofly/blog/754413)
- [【转载】跟花和尚学系统设计：明星公司之Netflix(中篇)](https://my.oschina.net/moooofly/blog/754417)
- [【转载】跟花和尚学系统设计：明星公司之Netflix(下篇)](https://my.oschina.net/moooofly/blog/754426)

# Simian Army

![Simian Army](http://mmbiz.qpic.cn/mmbiz/edSnflBH738kheia7KNvP6p2BDVnFx8bAOWdqyZPU689Duib1w6Jiaxb2t7dZyQ8zpZRrsEH7d6EPMvqxhxx5TwhA/640?wx_fmt=png&tp=webp&wxfrom=5)


# [Chaos Monkey](https://github.com/Netflix/SimianArmy/wiki/Chaos-Monkey)

![chaos monkey](http://mmbiz.qpic.cn/mmbiz/edSnflBH738kheia7KNvP6p2BDVnFx8bApdJyEJTXVJUic13umHxm7rtvvibuyviah9VDv2VDAUlYcHtdIicO39RUOQ/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1)

## What is Chaos Monkey?

Chaos Monkey 是一种服务，用于将系统分组，并随机终止属于某个分组中的系统中的一部分；该服务运行在一定的受控时间段和时间间隔内（不会无故运行在周末和假日中，并且仅在上班时间内运行）；在大多数情况下，我们的应用设计要保证当某个 peer 下线时仍能继续工作，但是在那些特殊的场景下，我们需要确保有人在值守，以便解决问题，并从问题中进行经验学习；基于这个想法，Chaos Monkey 仅会在工作时间内被使用，以保证工程师能发现警告信息，并作出适当的回应；

## Why Run Chaos Monkey?

失效一定会发生，并且无法避免；如果你的应用无法容忍系统失效的情况，你是愿意在凌晨 3 点被叫醒，还是希望在办公室中用过 morning coffee 之后呢？即使你非常自信你所设计的架构本身能够容忍系统失效，但是你能确保其在下周仍能运转良好么，亦或是下个月呢？软件本身具有复杂和动态特性，因此你上周进行的 "simple fix" 可能会导致未知的后果；在系统失效期间，你的负载均衡器正确的探测和路由请求了么？你能可靠的重建你的系统么？可能发生某个工程师上周 "quick patched" 了一个在线系统，但是忘记提交变更到你到源码仓库的情况么？

Refer to the [Quick start guide](https://github.com/Netflix/SimianArmy/wiki/Quick-Start-Guide) to get started setting up and using Chaos Monkey

----------

# [Netflix新放出来的开源工具Chaos Monkey](http://www.infoq.com/cn/news/2012/08/chaos-monkey)

----------

#[放到野外的Chaos Monkey](http://techblog.netflix.com/2012/07/chaos-monkey-released-into-wild.html)

----------

# [我们使用AWS得到的5个教训](http://techblog.netflix.com/2010/12/5-lessons-weve-learned-using-aws.html)

----------


# [Netflix继续开源，更多猴子进入视野](http://www.infoq.com/cn/news/2013/02/netflix-opensource)

除了上述提到的内容之外，自然少不了著名的猴子们，[Simian Army](https://github.com/Netflix/SimianArmy)是这群猴子的统称，除了Chaos Monkey和Janitor Monkey之外，更多的猴子也将在今年陆续开源。在Netflix的技术博客上有篇[文章](http://techblog.netflix.com/2011/07/netflix-simian-army.html)，详细介绍了Simian Army中的各位成员：

- Chaos Monkey，可以随机关闭生产环境中的实例，确保网站系统能够经受故障的考验，同时不会影响客户的正常使用。
- Latency Monkey，在RESTful服务的调用中引入人为的延时来模拟服务降级，测量上游服务是否会做出恰当响应。通过引入长时间延时，还可以模拟节点甚至整个服务不可用。
- Conformity Monkey，查找不符合最佳实践的实例，并将其关闭。例如，如果某个实例不在自动伸缩组里，那么就该将其关闭，让服务所有者能重新让其正常启动。
- Doctor Monkey，查找不健康实例的工具，除了运行在每个实例上的健康检查，还会监控外部健康信号，一旦发现不健康实例就会将其移出服务组。
- Janitor Monkey，查找不再需要的资源，将其回收，这能在一定程度上降低云资源的浪费。
- Security Monkey，这是Conformity Monkey的一个扩展，检查系统的安全漏洞，同时也会保证SSL和DRM证书仍然有效。
- 10-18 Monkey，进行本地化及国际化的配置检查，确保不同地区、使用不同语言和字符集的用户能正常使用Netflix。
- Chaos Gorilla，Chaos Monkey的升级版，可以模拟整个Amazon Availability Zone故障，以此验证在不影响用户，且无需人工干预的情况下，能够自动进行可用区的重新平衡。

----------

# [Chaos Monkey](http://whatis.techtarget.com/definition/Chaos-Monkey)

> Chaos Monkey is a software tool that was developed by Netflix  engineers to test the resiliency and recoverability of their Amazon Web Services (AWS).

Chaos Monkey 测试 AWS 的弹性和恢复能力的工具；

> The software simulates failures of instances of services running within Auto Scaling Groups (ASG) by shutting down one or more of the virtual machines. According to the developers, Chaos Monkey was named for the way it wreaks havoc like a wild and armed monkey set loose in a data center.

AWS 的 service 运行在 ASG 中；
Chaos Monkey 通过关停一个或多个虚拟机来模拟 service 实例的失效；
Chaos Monkey 的名字来源于其工作的方式：如同一只野生的、武装了的猴子，在数据中心释放后，造成的严重破坏；

> Chaos Monkey works on the principle that the best way to avoid major failures is to fail constantly. However, unlike unexpected failures, which seem to occur at the worst possible times, the software is opt-out by default. It can also be configured for opt-in.

Chaos Monkey 的原则：避免大多数失效的主要方式就是经常失效；
然而，又与一般性的、发生在预期外的、最差时机上的失败不同，该软件默认就提供了退出机制，同时也可以配置合适进入；

> Chaos Monkey has a configurable schedule that allows simulated failures to occur at times when they can be closely monitored.  In this way, it’s possible to prepare for major unexpected errors rather than just waiting for catastrophe to strike and seeing how well you can manage.

Chaos Monkey 可针对调度规划进行配置，即允许模拟出来的失败情况发生在我们的密切监督下；
因此，针对主要的非预期错误进行准备是可能的，而不会仅仅观察在灾难的冲击下我们到底能干什么；

> Chaos Monkey was the original member of Netflix’s Simian Army, a collection of software tools designed to test the AWS infrastructure. The software is open source to allow other cloud services users to adapt it for their use. 

Chaos Monkey 属于 Netflix 公司的 Simian Army 产品中的一员；
Simian Army 由一组软件工具构成，用于测试 AWS 基础设施；
该软件开源，可用于其他云服务用户进行相应测试使用；

> Other Simian Army members have been added to create failures and check for abnormal conditions, configurations and security issues.  Chaos Gorilla, another member of the Simian Army, simulates outages for entire regions. 

其他 Simian Army 成员可用于创建失败场景，检查异常条件，检查配置问题和安全问题；
Chaos Gorilla ，Simian Army 中的成员之一，可以模拟整个区域的断电场景；

Netflix engineers plan to add more monkeys to the army, some based on community suggestions.



