




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


----------


#[放到野外的Chaos Monkey](http://techblog.netflix.com/2012/07/chaos-monkey-released-into-wild.html)

We have found that the best defense against major unexpected failures is to fail often. By frequently causing failures, we force our services to be built in a way that is more resilient. We are excited to make a long-awaited announcement today that will help others who embrace this approach.

We have written about our [Simian Army](http://techblog.netflix.com/2011/07/netflix-simian-army.html) in the past and we are now proud to announce that the source code for the founding member of the Simian Army, Chaos Monkey, [is available to the community](https://github.com/Netflix/SimianArmy).

Do you think your applications can handle a troop of mischievous monkeys loose in your infrastructure? Now you can find out.

## What is Chaos Monkey?

Chaos Monkey is a service which runs in the Amazon Web Services (AWS) that seeks out Auto Scaling Groups (ASGs) and terminates instances (virtual machines) per group. The software design is flexible enough to work with other cloud providers or instance groupings and can be enhanced to add that support. The service has a configurable schedule that, by default, runs on non-holiday weekdays between 9am and 3pm. In most cases, we have designed our applications to continue working when an instance goes offline, but in those special cases that they don't, we want to make sure there are people around to resolve and learn from any problems. With this in mind, Chaos Monkey only runs within a limited set of hours with the intent that engineers will be alert and able to respond.

## Why Run Chaos Monkey?

Failures happen and they inevitably happen when least desired or expected. If your application can't tolerate an instance failure would you rather find out by being paged at 3am or when you're in the office and have had your morning coffee? Even if you are confident that your architecture can tolerate an instance failure, are you sure it will still be able to next week? How about next month? Software is complex and dynamic and that "simple fix" you put in place last week could have undesired consequences. Do your traffic load balancers correctly detect and route requests around instances that go offline? Can you reliably rebuild your instances? Perhaps an engineer "quick patched" an instance last week and forgot to commit the changes to your source repository?
There are many failure scenarios that Chaos Monkey helps us detect. Over the last year Chaos Monkey has terminated over 65,000 instances running in our production and testing environments. Most of the time nobody notices, but we continue to find surprises caused by Chaos Monkey which allows us to isolate and resolve them so they don't happen again.

## Auto Scaling Groups

The default instance groupings that Chaos uses for selection is Amazon's Auto Scaling Group (ASG). Within an ASG, Chaos Monkey will select an instance at random and terminate it. The ASG should detect the instance termination and automatically bring up a new, identically configured, instance. If you are not using Auto Scaling Groups that should be the first step to making your application handle these isolated instance failure scenarios. Check out [Asgard](http://techblog.netflix.com/2012/06/asgard-web-based-cloud-management-and.html) to make deploying and managing ASGs easy. There are many [great features for ASGs](http://techblog.netflix.com/2012/01/auto-scaling-in-amazon-cloud.html) beyond replacing terminated instances, like enabling the use of Amazon's Elastic Load Balancers (ELBs) to distribute traffic to all instances in your application. Netflix has a best-practice where all instances should be run within an ASG and we have Janitor Monkey to remind us by terminating all instances not following this best-practice.

## Configuration

Chaos Monkey allows for an Opt-In or an Opt-Out model. At Netflix, we use the Opt-Out model, so if an application owner does nothing, Chaos Monkey will be acting on their application. For your organization, you have the option to choose what is right for you. This allows you to "test the water" and try out Chaos Monkey on a specific application to see how it reacts. Not every application can trivially handle an instance going offline.  Sometimes it takes a human to manually recover instances, perhaps exercising backups to bring them back. Ideally, engineers work towards making that process easier and faster and eventually automatic. For those applications, there is the ability to Opt-Out of Chaos Monkey. There is also a tunable "probability" that Chaos Monkey uses to control the chance of a termination.  A probability of 1 (or 100%) will terminate one instance per day per ASG.  If instance recovery is difficult and you only want a termination weekly, you can reduce the probability to 0.2 or 20% (daily is 100%, it runs 5 work days per week, so weekly is 20%). Note that this is still a probability and only meaningful when sampled multiple times. With a 20% probability, Chaos Monkey would terminate one instance a week on average. In practice, it might be 2 days in a row followed by 2 weeks of no terminations, but given a large enough sample it will terminate weekly on average. For an environment as large as Netflix, the configuration can get a bit tricky to manage and for this we have developed a dashboard to help that we hope to open source soon. You can read more about how to configure Chaos Monkey on the [documentation wiki](https://github.com/Netflix/SimianArmy/wiki/Configuration).

## REST

Currently, there is a simple [REST interface](https://github.com/Netflix/SimianArmy/wiki/REST) that allows you to query Chaos Monkey termination events. We keep records of what was terminated and when, so if something disappears, you can see if Chaos Monkey was responsible. You could use this API to get notifications of terminations, but we encourage you to use a more general application monitoring solution like servo to discover what is happening to your applications at runtime.

## Costs

The termination events are stored in an Amazon SimpleDB table by default. There could be associated costs with Amazon SimpleDB but the activity of Chaos Monkey should be small enough to fall within Amazon's Free Usage Tier. Ultimately the costs associated with running Chaos Monkey are your responsibility.
Cost references: http://aws.amazon.com/simpledb/pricing/

## More Monkey Business

We have a long line of simians waiting to be released.  The next likely candidate will be Janitor Monkey which helps keep your environment tidy and your costs down.  Stay tuned for more announcements.
If building tools to automate the operations and improve the reliability of the cloud sounds exciting, we're always looking for new members to join the team.  Take a look at jobs.netflix.com for current openings or contact [@atseitlin](https://twitter.com/atseitlin).

### Chaos Monkey
- [Home Page](https://github.com/Netflix/SimianArmy/wiki)
- [Quick Start Guide](https://github.com/Netflix/SimianArmy/wiki/Quick-Start-Guide)
- [REST API](https://github.com/Netflix/SimianArmy/wiki/REST)
- [Source Code](https://github.com/Netflix/SimianArmy)

### Netflix Cloud Platform
- The Netflix Simian Army
- Asgard Web Based Cloud Management
- 5 Lessons We’ve Learned Using AWS
- Netflix Open Source Projects
- Auto Scaling in the Amazon Cloud
- Servo (Publish application metrics for auto scaling)
- @NetflixOSS Twitter Feed

### Amazon Web Services
- Auto Scaling
- Elastic Load Balancing (ELB)
- SimpleDB Costs

----------

# [我们使用AWS得到的5个教训](http://techblog.netflix.com/2010/12/5-lessons-weve-learned-using-aws.html)

In my [last post](http://techblog.netflix.com/2010/12/four-reasons-we-choose-amazons-cloud-as.html) I talked about some of the reasons we chose AWS as our computing platform. We’re about one year into our transition to AWS from our own data centers. We’ve learned a lot so far, and I thought it might be helpful to share with you some of the mistakes we’ve made and some of the lessons we’ve learned.

## Dorothy, you’re not in Kansas anymore.

If you’re used to designing and deploying applications in your own data centers, you need to be
prepared to unlearn a lot of what you know. Seek to understand and embrace the differences operating in a cloud environment.

Many examples come to mind, such as hardware reliability. In our own data centers, session-based memory management was a fine approach, because any single hardware instance failure was rare. Managing state in volatile memory was reasonable, because it was rare that we would have to migrate from one instance to another. I knew to expect higher rates of individual instance failure in AWS, but I hadn’t thought through some of these sorts of implications.

Another example: in the Netflix data centers, we have a high capacity, super fast, highly reliable
network. This has afforded us the luxury of designing around chatty APIs to remote systems. AWS networking has more variable latency. We’ve had to be much more structured about “over the wire” interactions, even as we’ve transitioned to a more highly distributed architecture.

## Co-tenancy is hard.

When designing customer-facing software for a cloud environment, it is all about managing down expected overall latency of response. AWS is built around a model of sharing resources; hardware, network, storage, etc. Co-tenancy can introduce variance in throughput at any level of the stack. You’ve got to either be willing to abandon any specific subtask, or manage your resources within AWS to avoid co-tenancy where you must.

Your best bet is to build your systems to expect and accommodate failure at any level, which introduces the next lesson.

## The best way to avoid failure is to fail constantly.

We’ve sometimes referred to the Netflix software architecture in AWS as our Rambo Architecture. Each system has to be able to succeed, no matter what, even all on its own. We’re designing each distributed system to expect and tolerate failure from other systems on which it depends.

If our recommendations system is down, we degrade the quality of our responses to our customers, but we still respond. We’ll show popular titles instead of personalized picks. If our search system is intolerably slow, streaming should still work perfectly fine.

One of the first systems our engineers built in AWS is called the Chaos Monkey. The Chaos Monkey’s job is to randomly kill instances and services within our architecture. If we aren’t constantly testing our ability to succeed despite failure, then it isn’t likely to work when it matters most – in the event of an unexpected outage.

## Learn with real scale, not toy models.

Before we committed ourselves to AWS, we spent time researching the platform and building test systems within it. We tried hard to simulate realistic traffic patterns against these research projects.

This was critical in helping us select AWS, but not as helpful as we expected in thinking through our architecture. Early in our production build out, we built a simple repeater and started copying full customer request traffic to our AWS systems. That is what really taught us where our bottlenecks were, and some design choices that had seemed wise on the white board turned out foolish at big scale.

We continue to research new technologies within AWS, but today we’re doing it at full scale with real data. If we’re thinking about new NoSQL options, for example, we’ll pick a real data store and port it full scale to the options we want to learn about.

## Commit yourself.

When I look back at what the team has accomplished this year in our AWS migration, I’m truly amazed. But it didn’t always feel this good. AWS is only a few years old, and building at a high scale within it is a pioneering enterprise today. There were some dark days as we struggled with the sheer size of the task we’d taken on, and some of the differences between how AWS operates vs. our own data centers.

As you run into the hurdles, have the grit and the conviction to fight through them. Our CEO, Reed Hastings, has not only been fully on board with this migration, he is the person who motivated it! His commitment, the commitment of the technology leaders across the company, helped us push through to success when we could have chosen to retreat instead.

AWS is a tremendous suite of services, getting better all the time, and some big technology companies are running successfully there today. You can too! We hope some of our mistakes and the lessons we’ve learned can help you do it well.

----------
