# golang 之 context 使用

标签（空格分隔）： golang

---

## [golang context初探](http://www.wxdroid.com/index.php/2073.html)

### 什么是 context

- 从 go1.7 开始，`golang.org/x/net/context` 包正式作为 `context` 包进入了标准库；
- 官方的文档说明：*Package context defines the Context type, which carries deadlines, cancelation signals, and other request-scoped values across API boundaries and between processes.*
- 通过 `context` ，我们可以方便地针对由同一个请求所产生的 goroutines 进行约束管理：可以设定超时时间、deadline，甚至是取消该请求相关的所有 goroutines ；

### 如何使用 context

示例代码

```golang
package main

import (
    "context"
    "log"
    "net/http"
    _ "net/http/pprof"
    "time"
)

func main() {
    go http.ListenAndServe(":8989", nil)
    ctx, cancel := context.WithCancel(context.Background())
    go func() {
        time.Sleep(3 * time.Second)
        cancel()
    }()
    log.Println(A(ctx))
    select {}
}

func C(ctx context.Context) string {
    select {
    case <-ctx.Done():
        return "C Done"
    }
    return ""
}

func B(ctx context.Context) string {
    ctx, _ = context.WithCancel(ctx)
    go log.Println(C(ctx))
    select {
    case <-ctx.Done():
        return "B Done"
    }
    return ""
}

func A(ctx context.Context) string {
    go log.Println(B(ctx))
    select {
    case <-ctx.Done():
        return "A Done"
    }
    return ""
}
```

> 问题：如何查看 pprof 的输出信息？

### context 的使用规范

在最新的 1.8 beta 版本中，很多相关的包都加入了 `context` ，比如 `database` 包。那么，在使用 context 的时候有哪些需要注意呢？

- 不要把 context 存储在结构体中，而是要显式地进行传递；
- 把 context 作为第一个参数，并且一般都把变量命名为 ctx ；
- 就算是程序允许，也不要传入一个 nil 的 context ；如果不知道要使用哪种 context 的话，用 context.TODO() 来替代；
- context.WithValue() 只用来传递请求范围的值（即和请求相关的元数据），不要用它来传递可选参数；
- 同一个 Context 可以用来传递到不同的 goroutine 中，Context 在多个 goroutine 中是安全的


----------


## [Golang之Context的使用](http://www.nljb.net/default/Golang%E4%B9%8BContext%E7%9A%84%E4%BD%BF%E7%94%A8/)

在 golang 中的创建一个新的线程（即 goroutine）时，并不会返回类似 C 语言中 pid 的东东，所以我们无法从外部直接杀死某个线程（goroutine）；因此，若想要在目标 goroutine 自行结束前，终止其执行，就需要采用 **channel + select** 的方式来解决；然而在有些场景下，基于这种方式的实现比较麻烦；例如，由一个请求衍生出多个 goroutine ，并且 goroutine 之间需要满足一定的约束关系，以实现诸如有效期、中止线程树、传递请求全局变量之类功能的场景；于是 google 为我们提供一个解决方案：基于 `context` 包实现上下文功能；

> 此文中给出的多个示例程序比较有阅读价值；


----------

## [使用context实现多个goroutine的依赖管理](http://studygolang.com/articles/2790)

> 此文将《[Go Concurrency Patterns: Context](https://blog.golang.org/context)》中的示例改造成可以本地运行的东东了～（可以进一步深入研究下）

----------


## [Go Concurrency Patterns: Context](https://blog.golang.org/context)

关于 context 的权威文章（官方＋英文）；

## [Go语言并发模型：使用 context](http://studygolang.com/articles/9048)

上述文章的中文翻译；

----------

## [My Go Resolutions for 2017](https://research.swtch.com/go2017)

```
...
Context & best practices

We added the new context package in Go 1.7 for holding request-scoped information like timeouts, cancellation state, and credentials. An individual context is immutable (like an individual string or int): it is only possible to derive a new, updated context and pass that context explicitly further down the call stack or (less commonly) back up to the caller. The context is now carried through APIs such as database/sql and net/http, mainly so that those can stop processing a request when the caller is no longer interested in the result. Timeout information is appropriate to carry in a context, but—to use a real example we removed—database options are not, because they are unlikely to apply equally well to all possible database operations carried out during a request. What about the current clock source, or logging sink? Is either of those appropriate to store in a context? I would like to try to understand and characterize the criteria for what is and is not an appropriate use of context.
...
```


----------


## 基于 go tool vet 分析 context 问题

使用 vet 命令分析本文最开始的 context 演示代码

```shell
root@vagrant-ubuntu-trusty:/go/src# go tool vet -lostcancel context_usage.go
context_usage.go:31: the cancel function returned by context.WithCancel should be called, not discarded, to avoid a context leak
root@vagrant-ubuntu-trusty:/go/src#
```

可以看到，分析结果认为调用 context.WithCancel 得到的 cancel 函数不应该被忽略掉；