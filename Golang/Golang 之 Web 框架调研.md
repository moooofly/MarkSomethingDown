# Golang 之 Web 框架调研

> 本文是关于 golang Web 框架信息的学习和汇总；

## [超全的Go Http路由框架性能比较](http://colobu.com/2016/03/23/Go-HTTP-request-router-and-web-framework-benchmark/)


- `net/http` 自己的 [default request multiplexer (i.e. mux)](https://golang.org/pkg/net/http/#ServeMux)
    - 简单、功能有限，很容易扩展实现自定义路由器；
    - 很多第三方路由库本质上就是在其基础上的扩展；
- [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter) 在 `net/http` 默认的 mux 基础上，实现了 routing pattern 的变量支持，并且更具扩展性；
    - 增加了一些贴心功能（路径末尾的'\'、path 自动修正等）；
    - 衍生出了一个专门用于测试 routing 性能的测试框架 [julienschmidt/go-http-routing-benchmark](https://github.com/julienschmidt/go-http-routing-benchmark)；
- 测试框架模拟了静态路由、Github API、Goolge+ API、Parse API 的各种情况；

> 完整的 Web 框架的特性包括：
> 
> - 上下文管理；
> - Session 维护；
> - 模版的处理；
> - ORM ；
>
> ------ 
>
> 选择一个框架的理由:
> 
> - 灵活性；
> - 扩展性；
> - API 友好程度；
> - 文档详细程度；
> - 项目活跃度；
> - 社区活跃度；
> - 性能和内存占用；



----------



## [iris 真的是最快的Golang 路由框架吗?](http://colobu.com/2016/04/01/Is-iris-the-fastest-golang-router-library/)

- 发现 [julienschmidt/go-http-routing-benchmark](https://github.com/julienschmidt/go-http-routing-benchmark) 这个测试框架存在着一个严重的问题，就是 Handler 的业务逻辑实现的非常简单，各个框架的 handler 类似，这和生产环境中的情况严重不符；
- 如果加上业务逻辑的处理时间，Go 内置的路由功能要远远好于 [kataras/iris](https://github.com/kataras/iris) ，甚至可以说 iris 的路由根本无法应用的有业务逻辑的产品中，随着业务逻辑的时间耗费加大，iris 的吞吐量急剧下降；
- 而对于 Go 的内置路由来说，业务逻辑的时间耗费加大，单个 client 会等待更长的时间，但是并发量大的网站来说，吞吐率不会下降太多；
- Go http server 实现的是每个 request 对应一个 goroutine (goroutine per request) ，考虑到 Http Keep-Alive 的情况，更准确的说是每个连接对应一个 goroutine (goroutine per connection) ；
- 提供了一个用来获取 goroutine 的 Id 的函数；
- 从代码实现层面分析对比了 Go http server 和 iris 的性能差异原因：
    - 对于 Go http server 来说，业务逻辑的时间花费会影响单个 goroutine 的执行时间，反映到客户的浏览器是是延迟时间 latency 增大了，如果并发量足够多，影响的是系统中的 goroutine 的数量以及它们的调度，吞吐率不会剧烈影响；
    - 对于 iris 来说，其为了提高性能，缓存了 context ，并且对于相同的请求 URL 和 Method ，会从缓存中使用相同的 context 使用；在并发量较大的时候，对于相同的请求（即 req.URL.Path 和 Method 相同的请求），会进入排队的状态，导致性能低下；



----------


## [谁是最快的Go Web框架](http://colobu.com/2016/04/06/the-fastest-golang-web-framework/)

- 发现了 [julienschmidt/go-http-routing-benchmark](https://github.com/julienschmidt/go-http-routing-benchmark) 这个测试框架的问题；
    - 只测试了 web 框架的路由功能（包括路径中参数的解析）；
    - 并未测试完整的 web 框架处理过程（接受连接、路由、Handler 处理）；
    - 其 Handler 实现非常的简单，不能反映实际产品的业务处理；
- 发现之前的测试结论存在问题（即“ iris 是最快的 HTTP 路由框架”的结论）；
- （倒逼）现在 iris 已经改成了 `valyala/fasthttp` 实现，性能超级好（基本上 iris 的高性能来源于 `valyala/fasthttp` 和 `julienschmidt/httprouter` 等一些框架的努力）；
- 促使重新实现了 [smallnest/go-web-framework-benchmark](https://github.com/smallnest/go-web-framework-benchmark) 这个测试框架；
    - 为每个 web 框架实现了 `/hello` 的 Http Get 服务，它返回 `hello world` 字符串。所有的 web 框架的实现都是一致的；
    - 可以指定业务处理的时间，如 10 毫秒，100 毫秒，500 毫秒等；
    - 自动化测试；
- `valyala/fasthttp` 表现非常的好，唯一需要考虑的是：如果选它做 web 框架，你的代码将难以迁移到别的框架上，因为它实现了和标准库`net/http` 不一样的接口。如果开启 http pipelining ，`valyala/fasthttp` 会远远好于基于 `net/http` 实现的框架。


> 基本上 go web 框架分为两个门派：
> 
> - 基于标准库 `net/http` 的框架；
> - 基于 `valyala/fasthttp` 库的框架；
>
> ------ 
>
> 可能处理速度非常慢的业务逻辑：
> 
> - 从一个网络连接中读取数据；
> - 写数据到硬盘中；
> - 访问数据库；
> - 访问缓存服务器；
> - 调用其它服务，等待服务结果的返回；


----------


## [Go Web 框架性能比拼 2017 春季版](http://colobu.com/2017/04/07/go-webframework-benchmark-2017-Spring/)


- 基于 Go v1.8.0 ；
- HTTP pipelining 是将多个 HTTP 请求（request）整批提交的技术，而在发送过程中不需先等待服务端的回应。请求结果管线化使得 HTML 网页加载时间动态提升，特别是在具有高延迟的连接环境下，如卫星上网。在宽带连接中，加速不是那么显著的；因为需要服务器端应用 HTTP/1.1 协议：服务器端必须按照客户端的请求顺序恢复请求，这样整个连接还是先进先出的，对头阻塞（HOL blocking）可能会发生，造成延迟。未来的 HTTP/2.0 或者 SPDY 中的异步操作将会解决这个问题。因为它可能将多个 HTTP 请求填充在一个 TCP 数据包内，HTTP 管线化需要在网络上传输较少的 TCP 数据包，减少了网络负载；
- iris 存在争议，已经从测试结果中被禁掉了；


> 业务不同处理时间的现实意义：
> 
> - **0ms** 模拟理想的业务处理，每个请求基本只耗费小于 1 毫秒的处理时间，这是理想的极端的情况，比如访问内存中缓存的对象就返回；
> - **10ms** 模拟比较好的业务处理，服务器只需在极短的情况下处理完请求，如果业务不是太复杂，没有访问本地磁盘、数据库或者其它远程服务的情况下，比较符合这种测试；
> - **100ms** 模拟一般的业务处理，一般接收请求后，可能访问本地磁盘上的文件、数据库或者调用一个或多个远程服务，在这种情况下，完成一次请求可能要花费较少的时间；


----------


## [Go语言的几个Web开发框架](http://blog.sina.com.cn/s/blog_68f1adf70101cz4w.html)

### [revel/revel](https://github.com/revel/revel)

是一个高效的 Go 语言 full-stack Web 开发框架，其思路完全来自 Java 的 Play Framework ；

特点：

- 热编译
- 简单可选
- 同步（每个请求都创建自己的 goroutine 来处理）

### [astaxie/beego](https://github.com/astaxie/beego)

beego 是一个用 Go 开发的应用框架，思路来自于 tornado ，路由设计来源于 Sinatra ，作者是 [build-web-application-with-golang](https://github.com/astaxie/build-web-application-with-golang/) 电子书的作者；

支持如下特性：

- MVC
- REST
- 智能路由
- 日志调试
- 配置管理
- 模板自动渲染
- layout 设计
- 中间件插入逻辑
- 方便的 JSON/XML 服务

### [golangers/framework](https://github.com/golangers/framework)

Golanger 是一个轻量级的 Web 应用框架，使用 Go 语言编写。

Golanger 框架主要实现了 MVC (Model-View-Controller) 模式（三层架构模式），把软件系统分为三个基本部分：模型、视图和`net/http`；

Golanger 约定的命名规则：

- 控制器 Controller ：存放在 controllers 目录中，负责转发请求，对请求进行处理；
- 模型 Model ：存放在 models 目录中，程序员编写程序应有的功能（实现算法等等）、数据管理和数据库设计（可以实现具体的功能）；
- 视图 View ：存放在 views 目录中，界面设计人员进行图形界面设计；
- 静态文件放在 static 目录中；
- add-on 存放第三方库文件，默认是把 GOPATH 设置为这个目录；

### [QLeelulu/goku](https://github.com/QLeelulu/goku)

国人开发的 Go web MVC 框架，仿照 ASP.NET MVC ，简单而且强大。

基本功能：

- mvc (Lightweight model)
- 路由
- 多模板引擎和布局
- 简单数据库 API
- 表单验证
- 控制器或 Action 的过滤
- 中间件


----------


## [Go语言的Web框架](http://fuxiaohei.me/2014/3/13/go-web-framework.html)

### [revel/revel](https://github.com/revel/revel)

- 最早的 Go 语言 Web 框架，借鉴的 java 和 scala 语言的 play 框架的很多想法；
- 带有和 play 一样的毛病，舍弃了原有的标准完全自己来；revel 完全不理 Go 标准库的一套，全部是自己的概念；

### [astaxie/beego](https://github.com/astaxie/beego)

- 国内最火热的、比较中型的框架；
- 除了基础的 MVC 结构外，还带有 Cache ，ORM ，Session 等多个库的支持；
- 用的人很多，文档也很齐全（更新不太及时），社区和 Q 群也很活跃；

###  [Martini](https://github.com/go-martini/martini) 

- 概念非常不错的新锐框架；
- 微型框架：只带有简单的核心，路由功能和依赖注入容器 inject ，因此很多东西需要自己写，比如 view ，session 等；
- martini 营造的是一种组件生态（martini-contrib），即 nodejs 中的 expressjs 在做的事情；
- 基于其 DI 实现，第三方库很容易改造为 martini 规范的中间件；
- 由于依赖注入的实现依赖 reflect 反射，而 Go 语言的反射库效率很差，因此过多的中间件肯定会拖慢整体的速度（取决于 Go 本身的发展情况）；



----------


## [Go语言的Web框架比较](http://www.jdon.com/47016)

英文原文：[这里](https://medium.com/square-corner-blog/a-comparison-of-go-web-frameworks-f47804cf86f6)


- 推荐使用 [`net/http`](https://golang.org/pkg/net/http/) 作为入门起步的标准库；
- 如果你需要**路由**方面功能，可使用 [Gorilla](http://www.gorillatoolkit.org/) 和 [Gocraft/web](https://github.com/gocraft/web) ；
- [Revel](http://revel.github.io/) 和 [Martini](https://github.com/go-martini/martini) 有太多的依赖注入和其他魔术让人感觉不舒服；
- 上述所有的 Web 框架都是基于 `net/http` 包构建的；
- 本文从路由功能、数据绑定、控制器、中间件，以及杂类几个方面进行了比较，值得深入思考；
    - **路由**是一种将 Web 请求映射到一个处理器函数的机制，比较应该从灵活性和实现是否直接的角度考虑；
    - **数据绑定**是将请求参数转换成处理器使用参数的机制，会涉及到反射和依赖注入相关实现；
    - **控制器**或上下文（context）是用于维护每个请求的状态的东东；
    - **中间件**是一个跨处理器提供通用功能的技术；
    - **杂类**主要涉及设计哲学，灵感来源等问题；


----------


## [Golang Web Framework Comparsion](https://github.com/diyan/go-web-framework-comparsion)

This suite aims to compare the public API of various Go web frameworks and routers.


----------

## [What are the best web frameworks for Go?](https://www.slant.co/topics/1412/~web-frameworks-for-go)

> 这是一个根据投票给出的排行榜，指出了一些优点和缺点，可以看看；


----------

## [A Survey of 5 Go Web Frameworks](https://thenewstack.io/a-survey-of-5-go-web-frameworks/)

> 文章是 2014 年的，内容有点老；

本文从 COMMUNITY 和 BEST FEATURE 两个方面对如下 web 框架进行了比较：

- Beego
- Martini
- Gorilla
- GoCraft
- `net/http`


----------


## [Why I Don’t Use Go Web Frameworks](https://medium.com/code-zen/why-i-don-t-use-go-web-frameworks-1087e1facfa4)


此文是 [jochasinga/requests](https://github.com/jochasinga/requests) 作者用于说明为何不要使用框架的战斗檄文；有助于从相反的角度审视一下框架使用问题；

推荐阅读！


Checklist for (Not) Using a Framework

Frameworks are especially tempting for newcomers considering most of the time they have a task and/or requirements to work on and not just an interest and time to invest in Go. Here is a simple checklist for you to run down before you consider using a framework:

- You understand interfaces in Go thoroughly like @rob_pike does.
- You understand context thoroughly or want to deal with it instead of responses and requests.
- You intend to build a simple REST web service which might only handle JSON tasks of which a framework makes simple.
- You work alone or it is unlikely someone will work on your code in the future.
- You do not often consult online docs and resources when programming.
- You do not intend to use other packages outside of a framework’s functionalities.

If most of your answers are false, then

Stick to the bare metal until you don’t have to.


----------


## [Go 开发 HTTP](http://fuxiaohei.me/2016/9/20/go-and-http-server.html?nsukey=iIAIpUMuTIu8B9MXcgDI51F9rGZ6s3Zpy3Qc6ruvD%2BF3D65%2FEWyw9vEqDxhPahbDyrLUHILxsUfpM55bCZZbXAAA5LOSOz56ljdc6kLZkw5hFIRS3GDiBQ1eATDVIutbz6jrYDM6G3GpvlXQgJ50WCWi8z2PbHoDYJ7rWkqx45Z1blqfSC%2F0ShtXuX%2F9Dml7)

本文完整介绍了基于标准库 `net/http` 开发 Web 服务时的使用细节，推荐阅读！


----------


## [HTTP Response Snippets for Go](http://www.alexedwards.net/blog/golang-response-snippets)

本文给出了如何返回 JSON, XML 数据，以及如何渲染模板内容等问题的代码示例程序。


