
# Netflix

关于 Netflix 的介绍，可以直接看如下三篇文章：
- [【转载】跟花和尚学系统设计：明星公司之Netflix(上篇)](https://my.oschina.net/moooofly/blog/754413)
- [【转载】跟花和尚学系统设计：明星公司之Netflix(中篇)](https://my.oschina.net/moooofly/blog/754417)
- [【转载】跟花和尚学系统设计：明星公司之Netflix(下篇)](https://my.oschina.net/moooofly/blog/754426)

# Simian Army

![Simian Army](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/simian%20army.png)

# Chaos Monkey

![chaos monkey](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/chaos%20monkey.png)

## What is Chaos Monkey?

Chaos Monkey 是一种服务，用于将系统分组，并随机终止属于某个分组中的系统中的一部分；该服务运行在一定的受控时间段和时间间隔内（不会无故运行在周末和假日中，并且仅在上班时间内运行）；在大多数情况下，我们的应用设计要保证当某个 peer 下线时仍能继续工作，但是在那些特殊的场景下，我们需要确保有人在值守，以便解决问题，并从问题中进行经验学习；基于这个想法，Chaos Monkey 仅会在工作时间内被使用，以保证工程师能发现警告信息，并作出适当的回应；

## Why Run Chaos Monkey?

失效一定会发生，并且无法避免；如果你的应用无法容忍系统失效的情况，你是愿意在凌晨 3 点被叫醒，还是希望在办公室中用过 morning coffee 之后呢？即使你非常自信你所设计的架构本身能够容忍系统失效，但是你能确保其在下周仍能运转良好么，亦或是下个月呢？软件本身具有复杂和动态特性，因此你上周进行的 "simple fix" 可能会导致未知的后果；在系统失效期间，你的负载均衡器正确的探测和路由请求了么？你能可靠的重建你的系统么？可能发生某个工程师上周 "quick patched" 了一个在线系统，但是忘记提交变更到你到源码仓库的情况么？

请参考 [Quick start guide](https://github.com/Netflix/SimianArmy/wiki/Quick-Start-Guide) 以便开始安装和使用 Chaos Monkey ；


----------

# [Chaos Monkey](http://whatis.techtarget.com/definition/Chaos-Monkey)

> Chaos Monkey is a software tool that was developed by Netflix  engineers to test the resiliency and recoverability of their Amazon Web Services (AWS).

Chaos Monkey 测试 AWS 的弹性和恢复能力的工具；

> The software simulates failures of instances of services running within Auto Scaling Groups (ASG) by shutting down one or more of the virtual machines. According to the developers, Chaos Monkey was named for the way it wreaks havoc like a wild and armed monkey set loose in a data center.

AWS 的 service 运行在 ASG 中；
Chaos Monkey 通过关停一个或多个虚拟机来模拟 service 实例的失效；
Chaos Monkey 的名字来源于其工作的方式：如同一只野生的、武装了的猴子，在数据中心释放后，造成的严重破坏；

> Chaos Monkey works on the principle that the best way to avoid major failures is to fail constantly. However, unlike unexpected failures, which seem to occur at the worst possible times, the software is opt-out by default. It can also be configured for opt-in.

Chaos Monkey 的原则：**避免大多数失效的主要方式就是经常失效；**

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
Chaos Gorilla ，作为 Simian Army 中的成员之一，可以模拟整个区域的断电场景；

> Netflix engineers plan to add more monkeys to the army, some based on community suggestions.

Netflix 的工程师计划添加更多的 monkeys 进入到 army 中，其中一些就是基于社区的建议创造的；


----------

# [Netflix新放出来的开源工具Chaos Monkey](http://www.infoq.com/cn/news/2012/08/chaos-monkey)

- Chaos Monkey是一套用来故意把服务器搞下线的软件，可以测试云环境的恢复能力。
- Netflix专门开发的一系列捣乱工具，已经有不少被拿出来和技术社区自由分享，现在Chaos Monkey也加入了这个行列。
- Netflix团队让Chaos Monkey亮相的时间，最早是在2010年12月的一篇[官博文章](http://techblog.netflix.com/2010/12/5-lessons-weve-learned-using-aws.html)，文章内容是他们在AWS云上托管其热门视频流服务所得到的经验教训。文中总结了一点，叫做“避免失败的最好办法是经常失败”,反映Netflix通过主动破坏自身环境来发现弱点的做法。
- 我们的工程师在AWS上最早建立的系统之一叫Chaos Monkey。这猴子的工作是随机杀掉架构中的运行实例和服务。
- Netflix技术团队在2012年7月20日的[官博文章](http://techblog.netflix.com/2012/07/chaos-monkey-released-into-wild.html)上宣布，Chaos Monkey作为开源项目公开。文中解释了Chaos Monkey的设计意图和运营中的注意事项。
- Netflix声称软件可以成功运行在AWS以外的云上，主要给用户检测自身环境中的失败条件。
- Chaos Monkey可以被设定为只在支持人员现场待命，准备救灾的时候才运行。
- 服务具有可配置的执行计划，默认只在非假日的周一到周五上午9点至下午3点执行。我们设计只有在预料警报会被工程师发现并作出响应的有限时间段，才把Chaos Monkey放出来。
- 用户可以决定Chaos Monkey对新应用的攻击强度。
- 对于禁不起下线的应用，Chaos Monkey允许自主退出。Chaos Monkey中止实例的几率也是可调的。

----------

#[放到野外的Chaos Monkey](http://techblog.netflix.com/2012/07/chaos-monkey-released-into-wild.html)

- the best defense against major unexpected failures is to fail often. By frequently causing failures, we force our services to be built in a way that is more resilient. 
-  With this in mind, Chaos Monkey only runs within a limited set of hours with the intent that engineers will be alert and able to respond.
- Chaos Monkey allows for an **Opt-In** or an **Opt-Out** model. 
- there is a simple REST interface that allows you to query Chaos Monkey termination events. We keep records of what was terminated and when, so if something disappears, you can see if Chaos Monkey was responsible.

----------

# [我们使用AWS得到的5个教训](http://techblog.netflix.com/2010/12/5-lessons-weve-learned-using-aws.html)

1. Dorothy, you’re not in Kansas anymore.
2. Co-tenancy is hard.
3. The best way to avoid failure is to fail constantly.
4. Learn with real scale, not toy models.
5. Commit yourself.

大意如下：

- 在自己的数据中心好用的东西，到了云上不一定好用（软件设计/容量设计/网络延迟等）；
- 云设计的基础就是资源共享，面对的是多租户；而这会导致各个层面上的吞吐量变化；
- 我们常常将位于 AWS 上的 Netflix 软件架构戏称为“**兰博架构**”；因为我们要求每一个系统都能做到，即使自身依赖的全部系统全军覆没，仍能提供服务；
- 以真实的网络环境数据进行实际测试；

----------


# [Netflix继续开源，更多猴子进入视野](http://www.infoq.com/cn/news/2013/02/netflix-opensource)

除了上述提到的内容之外，自然少不了著名的猴子们，[Simian Army](https://github.com/Netflix/SimianArmy)是这群猴子的统称，除了Chaos Monkey和Janitor Monkey之外，更多的猴子也将在今年陆续开源。在Netflix的技术博客上有篇[文章](http://techblog.netflix.com/2011/07/netflix-simian-army.html)，详细介绍了Simian Army中的各位成员：

- **Chaos Monkey**，可以随机关闭生产环境中的实例，确保网站系统能够经受故障的考验，同时不会影响客户的正常使用。
- **Latency Monkey**，在RESTful服务的调用中引入人为的延时来模拟服务降级，测量上游服务是否会做出恰当响应。通过引入长时间延时，还可以模拟节点甚至整个服务不可用。
- **Conformity Monkey**，查找不符合最佳实践的实例，并将其关闭。例如，如果某个实例不在自动伸缩组里，那么就该将其关闭，让服务所有者能重新让其正常启动。
- **Doctor Monkey**，查找不健康实例的工具，除了运行在每个实例上的健康检查，还会监控外部健康信号，一旦发现不健康实例就会将其移出服务组。
- **Janitor Monkey**，查找不再需要的资源，将其回收，这能在一定程度上降低云资源的浪费。
- **Security Monkey**，这是Conformity Monkey的一个扩展，检查系统的安全漏洞，同时也会保证SSL和DRM证书仍然有效。
- **10-18 Monkey**，进行本地化及国际化的配置检查，确保不同地区、使用不同语言和字符集的用户能正常使用Netflix。
- **Chaos Gorilla**，Chaos Monkey的升级版，可以模拟整个Amazon Availability Zone故障，以此验证在不影响用户，且无需人工干预的情况下，能够自动进行可用区的重新平衡。

