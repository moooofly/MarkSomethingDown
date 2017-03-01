# Go 语言 Q&A

标签（空格分隔）： golang

---

## 使用 %q 和 %v 输出 []string{"a", "b", "c"} 的区别

- **`%q`**
a single-quoted character literal safely escaped with Go syntax.

- **`%v`**
the value in a default format
when printing structs, the plus flag (**`%+v`**) adds field names

测试代码

```golang
package main
import "fmt"

func main() {
    s := []string{"a", "b", "d"}
    fmt.Printf("output(%%q) => %q\n", s)
    fmt.Printf("output(%%v) => %v\n", s)
}
```

输出结果

```shell
root@vagrant-ubuntu-trusty:/go/src# ./format_output
output(%q) => ["a" "b" "d"]
output(%v) => [a b d]
root@vagrant-ubuntu-trusty:/go/src#
```

### fmt 包打印格式说明

官方文档：[这里](https://github.com/golang/go/blob/9af83462c6f432b77a846a24b4d8efae9bdf0567/src/fmt/doc.go)；

## 类型转换/类型开关/类型断言的英文分别是什么

- 类型转换 => type conversion
- 类型开关 => type switch
- 类型断言 => type assertion

## 影子变量的使用

```golang
switch x := x.(type) {
    case int:
        fmt.Printf("value = %d     type = %T", )
    ...
}
```

> TODO: 意义

## 如何直接构造可被 `json.Unmarshal()` 操作的数据

- []byte('xxx') 是**错误**的；
- []byte(\`xxx\`) 是**正确**的；

## channel 相关问题

一篇整理文章：《[golang 之 channel 玩法](https://github.com/moooofly/MarkSomethingDown/blob/master/Golang/golang%20%E4%B9%8B%20channel%20%E7%8E%A9%E6%B3%95.md)》

### chan struct{} 的使用

一句话概括：使用 `chan struct{}` 作为跨 goroutines 的通信旗语（semaphores）再合适不过了；

> semaphores 也可以翻译成**信号量**；

使用旗语（semaphores）的好处在于：我们无需关心实际传入的值是什么；因此，在这里使用 0 字节长度的类型是有意义的；

额外的好处：针对这种特定的使用 case ，编译器能够进行优化；详见[这里](https://docs.google.com/document/d/1yIAYmbvL3JxOKOjuCyon7JhW4cSv1wy5hC0ApeGMV9s/pub)；

另外一种使用：在某些场景中会使用 Exit 函数，其中可能会做很多用于回收资源的关闭操作，而被关闭的对象很可能就是内嵌无限循环 goroutines 的实体；那么如何把退出信号传递给每一个 goroutine 呢？通常的做法是定义一个 chan struct{} ，当调用 close() 函数关闭 channel 时，所有的 <-chan 操作会同时触发执行，这样就实现了将退出信号传递给每一个 goroutine 的功能；而各个 goroutine 都可以通过 select（默认 linux 下 select 调用的是 epoll 进行轮询）进行轮询，当从 <-chan 上获取到数据时，break 出无限循环；

stackoverflow 上到一篇[讨论](http://stackoverflow.com/questions/20793568/golang-anonymous-struct-and-empty-struct)如下（摘录关键信息）：

- `done <- struct{}{}` 这么奇怪的东西是干啥的

    `struct{}{}` 实际上是类型为 `struct{}` 的复合字面量（composite literal）；

- empty struct 的 size (`struct {}`) 是 0

    你可以创建任意多个类型为 `struct{}` 的东东 (`struct{}{}`) 并压入你的 channel ：但内存消耗并不会发生改变；因此，使用其进行跨 goroutines 的信号通信再合适不过了；

### struct {} 的用途

在 [The Go Programming Language Specification](https://golang.org/ref/spec#Struct_types) 中我们知道：`struct {}` 即 empty struct ；

除了上面用于定义 `chan struct{}` 之外，使用 `struct {}` 还有哪些好处呢？

由于 `struct {}` 也是一种 struct ，因此其保留了相应的特性：

- 可以对其定义 methods (that type can be a method receiver)
- 可以实现 interface (with said methods you just define on your empty struct)
- 可以用于实现单例模式（singleton）

    > 在 Go 中，你可以使用一个 empty struct ，然后将所有数据保存在全局变量中；这样可以保证，只有一个该类型的实例，因为全部的 empty structs 都是可交换的（interchangeable）；

PS：补充一篇专门讨论 [empty struct](https://dave.cheney.net/2014/03/25/the-empty-struct) 的文章；对应的中文翻译在[这里](http://www.golangtc.com/t/575442b8b09ecc02f7000057)；

### 针对 channel 调用 close

> 原文地址：[这里](https://golang.org/ref/spec#Close)

对于 channel `c` 来说，对其调用内置函数 `close()` ，表明不会再有更多值被发送到该 channel 上；如果 c 是 receive-only channel 则会触发错误；向一个已关闭的 channel 发送数据或者再次关闭，则会导致运行时 panic ；针对 nil channel 调用 `close()` 同样会导致运行时 panic ；在调用 `close()` 后，并且在任何之前已发送的值被接收后，接收操作将会返回 channel 类型对应的零值（和表明 channel 已关闭的 false），而不会发生阻塞（无论对其调用多少次接收操作）；若执行的是 multi-valued 接收操作，则会返回一个接收值和一个 indication ，（后者用于）以表明该 channel 是否已关闭；


## expvar 使用

参见《[golang 之 expvar 使用](https://github.com/moooofly/MarkSomethingDown/blob/master/Golang/golang%20%E4%B9%8B%20expvar%20%E4%BD%BF%E7%94%A8.md)》；

## 是否使用匿名字段的背后考虑？是否使用指针类型的背后考虑？

```golang
type stream struct {
    applayer.Stream
    parser   parser
    tcptuple *common.TCPTuple
}
```

## 适合放在 init() 中执行的操作有哪些？

- 获取初始调用时间

```golang
func init() {
	startTime = time.Now()
}
```

- 设置命令行选项

```golang
func init() {
	// Adds logging specific flags: -v, -e and -d.
	verbose = flag.Bool("v", false, "Log at INFO level")
	toStderr = flag.Bool("e", false, "Log to stderr and disable syslog/file output")
	debugSelectorsStr = flag.String("d", "", "Enable certain debug selectors")
}
```

- 初始化随机数生成器的 seed

```golang
func init() {
	// Initialize runtime random number generator seed using global, shared
	// cryptographically strong pseudo random number generator.
	//
	// On linux Reader might use getrandom(2) or /udev/random. On windows systems
	// CryptGenRandom is used.
	n, err := cryptRand.Int(cryptRand.Reader, big.NewInt(math.MaxInt64))
	var seed int64
	if err != nil {
		// fallback to current timestamp on error
		seed = time.Now().UnixNano()
	} else {
		seed = n.Int64()
	}

	rand.Seed(seed)
}
```

## map[string]struct{} 的用途

```golang
// Keep sorted for future command addition
var redisCommands = map[string]struct{}{
    "APPEND":           {},
    "AUTH":             {},
    "BGREWRITEAOF":     {},
    "BGSAVE":           {},
    "BITCOUNT":         {},
    "BITOP":            {},
    "BITPOS":           {},
...
```

- 利用 map 的特性，方便判断是否某个 key 值存在；
- 使用 `struct{}` 的目的：可以不占用内存空间，在不关心 value 时使用；


## 创建 goroutine 的方式

```golang
    go func() {
        server.NewHttpServer(services)
    }()
```

和 

```golang
    go server.NewHttpServer(services)
```

> TODO：区别？选择依据？

## select 相关

### [The Go Programming Language Specification](https://golang.org/ref/spec#Select_statements)

A "select" statement chooses which of a set of possible send or receive operations will proceed. It looks similar to a "switch" statement but with the cases all referring to communication operations.

A case with a RecvStmt may assign the result of a RecvExpr to one or two variables, which may be declared using a short variable declaration. The RecvExpr must be a (possibly parenthesized) receive operation. There can be at most one default case and it may appear anywhere in the list of cases.

`select` 语句的执行过程：

- For all the cases in the statement, the channel operands of receive operations and the channel and right-hand-side expressions of send statements are evaluated exactly once, in source order, upon entering the "select" statement. The result is a set of channels to receive from or send to, and the corresponding values to send. Any side effects in that evaluation will occur irrespective of which (if any) communication operation is selected to proceed. Expressions on the left-hand side of a RecvStmt with a short variable declaration or assignment are not yet evaluated.
- If one or more of the communications can proceed, a single one that can proceed is chosen via a uniform pseudo-random selection. Otherwise, if there is a default case, that case is chosen. If there is no default case, the "select" statement blocks until at least one of the communications can proceed.
- Unless the selected case is the default case, the respective communication operation is executed.
- If the selected case is a RecvStmt with a short variable declaration or an assignment, the left-hand side expressions are evaluated and the received value (or values) are assigned.
- The statement list of the selected case is executed.

Since communication on nil channels can never proceed, a select with only nil channels and no default case blocks forever.

```golang
var a []int
var c, c1, c2, c3, c4 chan int
var i1, i2 int
select {
case i1 = <-c1:
	print("received ", i1, " from c1\n")
case c2 <- i2:
	print("sent ", i2, " to c2\n")
case i3, ok := (<-c3):  // same as: i3, ok := <-c3
	if ok {
		print("received ", i3, " from c3\n")
	} else {
		print("c3 is closed\n")
	}
case a[f()] = <-c4:
	// same as:
	// case t := <-c4
	//	a[f()] = t
default:
	print("no communication\n")
}

for {  // send random sequence of bits to c
	select {
	case c <- 0:  // note: no statement, no fallthrough, no folding of cases
	case c <- 1:
	}
}

select {}  // block forever
```


----------

### [Go语言并发模型：使用 select](https://segmentfault.com/a/1190000006815341)

一个[代码片段](https://talks.golang.org/2012/concurrency.slide#32)：

```golang
select {
case v1 := <-c1:
    fmt.Printf("received %v from c1\n", v1)
case v2 := <-c2:
    fmt.Printf("received %v from c2\n", v1)
case c3 <- 23:
    fmt.Printf("sent %v to c3\n", 23)
default:
    fmt.Printf("no one was ready to communicate\n")
}
```

上面这段代码中，select 语句有四个 case 子语句，前两个是 receive 操作，第三个是 send 操作，最后一个是默认操作。

代码执行到 select 时，**case 语句会按照源代码的顺序被评估，且只评估一次**，评估的结果会出现下面这几种情况：

- 除 default 外，如果只有一个 case 语句评估通过，那么就执行这个case里的语句；
- 除 default 外，如果有多个 case 语句评估通过，那么通过伪随机的方式随机选一个；
- 如果 default 外的 case 语句都没有通过评估，那么执行default 里的语句；
- 如果没有 default，那么代码块会被阻塞，指导有一个 case 通过评估；否则一直阻塞

select 的使用场景：

- 为请求设置超时时间（再没有引入 context 情况下的常见用法）；
- 用于 done channel 或 quit channel 检测；


----------

### 代码片段

- 当 channel 的 buffer 已满，（故意）直接丢弃 event ，避免阻塞

```golang
    select {
    case p.trans <- event:
        return true
    default:
        // drop event if queue is full
        return false
    }
```

- 以阻塞的方式调用 select ，当 done channel 上收到信号后，丢弃后续 event

```golang
    select {
    case p.flows <- event:
        return true
    case <-p.done:
        // drop event, if worker has been stopped
        return false
    }
```

- 通过 `filepath.Walk` 遍历目录下文件时的一种控制（实现只要关闭 done channel 就停止遍历的行为）

```golang
    // 如果done被关闭了，停止walk
    select {
    case <-done:
        return errors.New("walk canceled")
    default:
        return nil
    }
```

----------


### select {} 问题

参考：[这里](http://stackoverflow.com/questions/18661602/what-does-an-empty-select-do)

在 `net/http/httptest` 中，有如下代码：

```golang
go s.Config.Serve(s.Listener)
if *serve != "" {
  fmt.Fprintln(os.Stderr, "httptest: serving on", s.URL)
  select {}
}
```

因此，经常有人会问：`select {}` 是干什么用的？

一个空的 `select {}` 语句能够实现无限（永远）阻塞的功能；该语句类似于，或者说，在实践中等价于一个空的 `for {}` 语句；

据我所知，一个空的 `for {}` 将会造成 cpu 的空转（至少在其它语言中是这样的）；在 OP 的示例程序中，通过 `select {}` 的使用已阻止主 goroutine 的退出； 

从 CSP 的角度来讲，空的 `select {}` 语句更像是 **STOP** 行为，即进程不再向前运行；比较像一种 ***self-deadlock*** ，虽然我不认为这种方式用的地方很多，但上述代码确实代表了一种有趣的、真实使用场景；一个空的 `for {}` 循环对应的是 ***self-livelock*** ，两者有所不同，因为后者消耗 cpu 资源；

#### 实验验证

- 测试一（`select {}`）

```golang
package main

import "fmt"

func main() {
    fmt.Println("vim-go")
    select {}
}
```

运行时发现会报 "*fatal error: all goroutines are asleep - deadlock!*" 错误；

```shell
root@vagrant-ubuntu-trusty:~/workspace/CODE/Golang# go run main.go
vim-go
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [select (no cases)]:
main.main()
	/root/workspace/CODE/Golang/main.go:8 +0xd1
exit status 2
root@vagrant-ubuntu-trusty:~/workspace/CODE/Golang#
```

> 经常在 go 标准库中能看到 `select {}` 的使用（一些开源代码也这么用），说明这种方式是好的；上述代码的问题在于：没有一个 goroutine 处于运行状态；

- 测试二（`for {}`）

```golang
package main

import "fmt"

func main() {
    fmt.Println("vim-go")
    for {
    }
}
```

运行时会卡住，没有出现之前的错误信息；

```shell
root@vagrant-ubuntu-trusty:~/workspace/CODE/Golang# go run main.go
vim-go
(卡住)
```

通过 `top` 可以看到，用户态 CPU 占用为 95.7%（证明了基于 `for {}` 方式阻塞 goroutine 的方式是不好的）

```shell
root@vagrant-ubuntu-trusty:~/workspace/CODE/Golang# ps aux|grep "for_usage"
root      4895 51.1  0.1   3460   768 pts/1    Rl   11:35   0:03 ./for_usage
root@vagrant-ubuntu-trusty:~/workspace/CODE/Golang#
root@vagrant-ubuntu-trusty:~/workspace/CODE/Golang# /usr/local/bin/pstack 4895
Thread 1 (process 4895):
#0  main.main () at /root/workspace/CODE/Golang/for_usage.go:7
root@vagrant-ubuntu-trusty:~/workspace/CODE/Golang#
root@vagrant-ubuntu-trusty:~/workspace/CODE/Golang# top -Hp 4895
top - 11:37:30 up 2 days, 33 min,  2 users,  load average: 1.98, 1.42, 1.18
Threads:   3 total,   1 running,   2 sleeping,   0 stopped,   0 zombie
%Cpu(s): 95.7 us,  4.0 sy,  0.0 ni,  0.0 id,  0.0 wa,  0.3 hi,  0.0 si,  0.0 st
KiB Mem:    501832 total,   445524 used,    56308 free,    48820 buffers
KiB Swap:   522236 total,        0 used,   522236 free.   294312 cached Mem

  PID USER      PR  NI    VIRT    RES    SHR S %CPU %MEM     TIME+ COMMAND
 4895 root      20   0    3460    772    472 R 48.6  0.2   1:01.51 for_usage
 4896 root      20   0    3460    772    472 S  1.0  0.2   0:01.24 for_usage
 4897 root      20   0    3460    772    472 S  0.0  0.2   0:00.00 for_usage
```

- 测试三（基于 channel 的阻塞特性）

```golang
package main

import "fmt"

func main() {
    fmt.Println("vim-go")
    ch := make(chan struct{})
    <-ch
}
```

运行时发现会报 "*fatal error: all goroutines are asleep - deadlock!*" 错误；

```shell
root@vagrant-ubuntu-trusty:~/workspace/CODE/Golang# go run main.go
vim-go
fatal error: all goroutines are asleep - deadlock!

goroutine 1 [chan receive]:
main.main()
	/root/workspace/CODE/Golang/main.go:10 +0x108
exit status 2
root@vagrant-ubuntu-trusty:~/workspace/CODE/Golang#
```

> 能够看到错误信息和使用 `select {}` 时基本相同；基于 channel 的方式利用了未初始化 channel 的默认值为其类型的零值，而针对 nil channel 的接收会导致阻塞；

#### 补充说明

关于 "*fatal error: all goroutines are asleep - deadlock!*" 错误信息的理解：其实这个错误信息并不表示真正出现了死锁；由于对于死锁的检测非常麻烦，go 就采用了这种比较简单粗暴的方法，当所有 goroutine 都在睡觉时，go 就认为是死锁了；

一个能够直接说明错误的含义的例子；

```golang
package main

import (
    "fmt"
)

func f1 () {
    fmt.Println("f1")
}

func f2 () {
    for {    // 无限循环，可以保证指定 f2 的 goroutine 永远不会处于 asleep 状态
        fmt.Println("f2")
    }
}

func main() {
    go f1()
    go f2()   // 若没有执行该语句，则会报上述错误信息
    ch := make(chan int)
    <-ch
}
```

## 字符串用反引号括起来

```golang
    `
    SELECT hostname, name, config FROM collectors
    INNER JOIN hosts ON hosts.id = collectors.host_id
    INNER JOIN collector_types ON collector_types.id = collectors.collector_type_id
    WHERE collector_types.name="corvus";
    `
```

## go run 编译出来的东西放在了哪里？

可以通过 `go run -work xxx.go` 看到临时目录相关信息；

```golang
➜  pkginit git:(master) go run -work initpkg_demo.go
WORK=/var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/go-build574909512
Map: map[1:A 2:B 3:C]
OS: darwin, Arch: amd64
➜  pkginit git:(master)
```

可以通过 `go run -x xxx.go` 看到命令执行期间用到的所有命令；

```golang
➜  pkginit git:(master) go run -x initpkg_demo.go
WORK=/var/folders/wg/w5bqgv311fx878j0swrntqg40000gn/T/go-build887487701
mkdir -p $WORK/command-line-arguments/_obj/
mkdir -p $WORK/command-line-arguments/_obj/exe/
cd /Users/sunfei/workspace/GiT/IdeaProjects/src/github.com/g0hacker/goc2p/src/basic/pkginit
/usr/local/Cellar/go/1.7.1/libexec/pkg/tool/darwin_amd64/compile -o $WORK/command-line-arguments.a -trimpath $WORK -p main -complete -buildid 8dd480db8fc9c54e2d73b883e36a24069936e235 -D _/Users/sunfei/workspace/GiT/IdeaProjects/src/github.com/g0hacker/goc2p/src/basic/pkginit -I $WORK -pack ./initpkg_demo.go
cd .
/usr/local/Cellar/go/1.7.1/libexec/pkg/tool/darwin_amd64/link -o $WORK/command-line-arguments/_obj/exe/initpkg_demo -L $WORK -w -extld=clang -buildmode=exe -buildid=8dd480db8fc9c54e2d73b883e36a24069936e235 $WORK/command-line-arguments.a
$WORK/command-line-arguments/_obj/exe/initpkg_demo
Map: map[1:A 2:B 3:C]
OS: darwin, Arch: amd64
➜  pkginit git:(master)
```

## `sync.WaitGroup` 的使用

原文：[这里](http://stackoverflow.com/documentation/go/376/concurrency/2490/waiting-for-goroutines#t=201607280838310884759)

**当 `main` 函数（或者说 main goroutine）运行结束时，整个 Go 程序就会终止**；因此，常规实践方案就是在 main goroutine 中等待所有其它 goroutines 运行完成；而具体的解决办法就是使用 `sync.WaitGroup` 完成相应的控制（同步）行为；

示例程序如下

```golang
package main

import (
    "fmt"
    "sync"
)

var wg sync.WaitGroup // 1

func routine(i int) {
    defer wg.Done() // 3
    fmt.Printf("routine %v finished\n", i)
}

func main() {
    for i := 0; i < 10; i++ {
        wg.Add(1) // 2
        go routine(i) // *
    }
    wg.Wait() // 4
    fmt.Println("main finished")
}
```

`sync.WaitGroup` 在执行顺序控制中的使用：

- **声明全局变量**：定义为全局变量是将其暴露给全部函数和方法的最简单方式；
- **增加计数器的值**：必须在 main goroutine 中进行调用，因为无法确保一个新启动的 goroutine 一定能够在步骤 4 之前得到执行（由于内存模型 guarantees 的缘故）
- **减少计数器的值**：必须在 goroutine 退出时执行；通过使用 defer 调用，可以确保对应函数一定会在函数结束时被调用，无论函数是如何结束的；
- **等待计数器的值归零**：必须在 main goroutine 中进行调用，以便阻止在所有 goroutines 程序执行完之前退出；

> **Parameters 应该在启动新 goroutine 前进行求值**（这其实是 golang 语言中的“closure 问题”）；因此，在调用 `wg.Add(1)` 之前，显式定义 Parameters 的值是非常有必要的；这样做的好处是，可能触发 panic 的代码将不会导致计数器值的增加；在这里，变量 i 的值是在 for 循环中定义的；


（更一般的情况）正确的实现：

```golang
param := f(x)
wg.Add(1)
go g(param)
```

会导致问题的实现：

```golang
wg.Add(1)
go g(f(x))
```

## closure 问题

此问题属于 golang 语言中比较奇怪的地方（即不要在闭包中引用指针），需要注意；

> 其实上面的代码实现中，在 `for i := 0; i < 10; i++ {...}` 里调用 `go routine(i)` 也属于 closure ，但不会出“closure 问题”；而在下面的代码中，通过 range 遍历得到相应值，在调用 `go func(...)` 的 closure 实现，则存在如下两种实现方式；

- 方式一

```golang
package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {
	values := []string{"a", "b", "c", "d"}
	for _, val := range values {
		wg.Add(1)
		go func(val interface{}) {   // 1. 传参数
			defer wg.Done()
			fmt.Println(val)         // 2.
		}(val)                       // 3. 传参数
	}

	wg.Wait()
}
```

- 方式二

```golang
package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {
	values := []string{"a", "b", "c", "d"}

	for _, val := range values {
		val := val             // 1. 没有这句是不行的
		wg.Add(1)
		go func() {            // 2. 无参数
			defer wg.Done()
			fmt.Println(val)   // 3.
		}()                    // 4. 无参数
	}
	
	wg.Wait()
}
```

> 问题：为啥不行？这又涉及到 golang 语言的一个奇怪问题，详见《[Iterators in Go](https://ewencp.org/blog/golang-iterators/)》和《[lavalamp/The Three Go Landmines.markdown](https://gist.github.com/lavalamp/4bd23295a9f32706a48f)》

## 如何控制 goroutine 的退出

针对 goroutine 退出的控制问题，分为两个层面：
 
- main goroutine 等待所有 sub goroutine 运行完成后，再终止自身运行（通过 `sync.WaitGroup` 进行同步控制）；
- main goroutine 需要在一定条件下，主动结束一个或一组和某任务相关的 sub goroutine ，再终止自身运行（通过定时器和 channel 进行同步控制，或者直接使用 `context.Context`）；

> 上述表达中的 main goroutine 在更一般的情况下，应该描述为 father goroutine ；

示例程序如下（可以理解成没有 context 前的最佳实现）

```golang
package main

import (
    "log"
    "sync"
    "time"
)

func main() {
    // The WaitGroup lets the main goroutine wait for all other goroutines
    // to terminate. However, this is no implicit in Go. The WaitGroup must
    // be explicitely incremented prior to the execution of any goroutine 
    // (i.e. before the `go` keyword) and it must be decremented by calling
    // wg.Done() at the end of every goroutine (typically via the `defer` keyword). 
    wg := sync.WaitGroup{}

    // The stop channel is an unbuffered channel that is closed when the main
    // thread wants all other goroutines to terminate (there is no way to 
    // interrupt another goroutine in Go). Each goroutine must multiplex its
    // work with the stop channel to guarantee liveness.
    stopCh := make(chan struct{})


    for i := 0; i < 5; i++ {
        // It is important that the WaitGroup is incremented before we start
        // the goroutine (and not within the goroutine) because the scheduler
        // makes no guarantee that the goroutine starts execution prior to 
        // the main goroutine calling wg.Wait().
        wg.Add(1)
        go func(i int, stopCh <-chan struct{}) {
            // The defer keyword guarantees that the WaitGroup count is 
            // decremented when the goroutine exits.
            defer wg.Done()

            log.Printf("started goroutine %d", i)

            select {
            // Since we never send empty structs on this channel we can 
            // take the return of a receive on the channel to mean that the
            // channel has been closed (recall that receive never blocks on
            // closed channels).   
            case <-stopCh:
                log.Printf("stopped goroutine %d", i)
            }
        }(i, stopCh)
    }

    time.Sleep(time.Second * 5)
    close(stopCh)
    log.Printf("stopping goroutines")
    wg.Wait()
    log.Printf("all goroutines stopped")
}
```

## 开发环境配置（仅在 root 权限下允许 go 被使用）

- 普通用户设置

```shell
vagrant@vagrant-ubuntu-trusty:~$ vi .profile
# ~/.profile: executed by the command interpreter for login shells.
# This file is not read by bash(1), if ~/.bash_profile or ~/.bash_login
# exists.
# see /usr/share/doc/bash/examples/startup-files for examples.
# the files are located in the bash-doc package.

# the default umask is set in /etc/profile; for setting the umask
# for ssh logins, install and configure the libpam-umask package.
#umask 022

# if running bash
if [ -n "$BASH_VERSION" ]; then
    # include .bashrc if it exists
    if [ -f "$HOME/.bashrc" ]; then
        . "$HOME/.bashrc"
    fi
fi

# set PATH so it includes user's private bin if it exists
if [ -d "$HOME/bin" ] ; then
    PATH="$HOME/bin:$PATH"
fi
```

- root 用户设置

```
root@vagrant-ubuntu-trusty:~# vi .profile
# ~/.profile: executed by Bourne-compatible login shells.

if [ "$BASH" ]; then
  if [ -f ~/.bashrc ]; then
    . ~/.bashrc
  fi
fi

mesg n

TZ='Asia/Shanghai'; export TZ

export GOROOT=/usr/local/go
export GOPATH=/go
export PATH=$GOROOT/bin:$GOPATH/bin:$PATH:/usr/local/mysql/bin
```


----------


## go get 和 go get -u 的区别

xx

## import "github.com/elastic/beats/libbeat/beat" 与版本控制问题

xx

## []string{} 和 make([]string, 0) 的区别

xx

## make(chan int) 和 make(chan int, 1) 的选择

xx

## gorountine 和操作系统线程之间的关系

gorountine是Go语言里很重要的新概念，有点类似线程，但消耗的资源比线程少很多，而且gorountine只是Go内部的概念，不会在操作系统层面有对应的实现。在Go里启动的各个gorountine之间是并行的，每个gorountine可能会映射到一个系统线程，也可能多个gorountine共用一个线程，如果是多核的机器，不同的gorountine会自动分配到不同的核心。gorountine间的切换也由Go来控制，不需要程序员操心。gorountine占用的内存远小于系统线程或进程，gorountine间的切换成本也很低。程序里可以轻易创建数万个gorountine做并行，而不用担心会占用过多的系统资源。

## for 用法

```
for {}                     // 相当于C语言里的while 1 {}
for i := 0; i < xx; i++ {} // 相当于C语言里for (int i=0; i<xx; i++) {}
for i > 0 {}               // 相当于C语言里的while (i>0) {}
for index, item := range array {} // 相当于Python里的foreach，index是循环序号
```

## defer 说明

defer 指定的函数不会立刻执行，而是在当前函数退出时才执行。defer 主要是用来做一些清扫类的工作，比如常见的关闭文件、释放缓存；

## 作用域问题

- x 的作用域被限制在 if 语句块内

```golang
if x := recover(); x != nil { ... }
```

- x 的作用域不受 if 语句块限制

```golang
x := recover()
if x != nil { ... }
```

## struct 中的 tag

Go 的 struct 的成员可以定义 Tag 这种描述信息，程序里没有直接的作用，但可以通过 reflect 访问到；虽然目前 Go 的标准库里只用 Tag 来做 json 编码/解码时对应域的名字，不过似乎可以有更有趣的用法。

```golang
type A struct {
	i int `tag:"abc"`
}

v := reflect.TypeOf(new(A))
f, _ := v.FieldByName("i")
f.Tag.Get("tag") == "abc"
```

虽然目前 Go 的标准库里只用 Tag 来做 json 编码/解码时对应域的名字，不过似乎可以有更有趣的用法。

在做 json 解码时，如果对应的 struct 没有特殊声明，会把 json 的域写入对应名字首字母大写的 struct 的成员里。而 json 编码时默认会以 struct 域的名字作为 json 域名，如果要对应到小写，需要声明相关域的 Tag ；



## golang 中调用 C

[Golang里调用C](http://air.googol.im/post/call-c-in-golang/)


## recover 是一个特殊函数

能猜到下面的程序会不会崩溃么？

```golang
package main

func main() {
	defer recover()
	panic(1)
}
```

程序依旧会崩溃。

原因是recover虽然看起来是个函数，但其实是编译器有特殊处理，可以当做一个关键字看待。

正确的写法是把recover放到一个函数里：

```golang
package main

func main() {
	defer func() { recover() }()
	panic(1)
}
```

## slice 的引用特性

依旧是先给出程序，猜输出：

```golang
package main

import (
	"fmt"
)

func main() {
	array := make([]int, 0, 3)
	array = append(array, 1)
	a := array
	b := array
	a = append(a, 2)
	b = append(b, 3)
	fmt.Println(a)
}
```

答案揭晓，输出是[1 3]。

就我的理解，slice是一个_{**指向内存的指针**，**当前已有元素的长度**，**内存最大长度**}_的结构体，其中只有**指向内存的指针**一项是真正具有引用语义的域，另外两项都是每个slice自身的值。因此，对slice做赋值时，会出现两个slice指向同一块内存，但是又分别具有各自的元素长度和最大长度。程序里把array赋值给a和b，所以a和b会同时指向array的内存，并各自保存一份当前元素长度1和最大长度3。之后对a的追加操作，由于没有超出a的最大长度，因此只是把新值2追加到a指向的内存，并把a的“当前已有元素的长度”增加1。之后对b进行追加操作时，因为a和b各自拥有各自的“当前已有元素的长度”，因此b的这个值依旧是1，追加操作依旧写在b所指向内存的偏移为1的位置，也就复写了之前对a追加时写入的2。

为了让slice具有引用语义，同时不增加array的实现负担，又不增加运行时的开销，似乎也只能忍受这个奇怪的语法了。



## go 语言设计哲学

go选择了CSP的内核实现，并在此上实现了goroutine和channel。这个选择最明显的一个结果，就是go的os库里没有其他语言os库的select/poll/kqueue的api。为什么？因为完全不需要。go里类似的功能可以通过对channel的select来实现，不需要暴露os的功能给使用者，而且也不需要让用户选择到底是用select还是poll。基于此，所有go的库在处理io时都是异步非阻塞的（除非同步等待channel或者用sync实现同步）。

利用编译器实现CSP还有另外的优势：动态栈。传统的编译语言，由于没有考虑并发的机制，整个程序是运行在同一个栈上；使用os的线程/进程库，会依赖os对线程/进程的栈的分配，这种分配和程序本身无关，是os根据经验设置的数值，对程序来说经常会出现要么栈效率不高，要么栈不够用溢出；使用库实现的纤程/线程/异步，都是在使用当前进程的栈，每个纤程/线程/异步模块没有自己的独立栈，导致很多纤程实现都是不依赖栈，而全靠堆来解决每个纤程内的分配问题，对内存造成了很大的压力。而go通过编译器实现的goroutine，可以在编译期知晓更多的goroutine信息，并在运行时动态分配每个goroutine的栈大小，做到既不浪费，也不溢出。

----------


注意interface只能接收一个实例的指针，而不能直接接收实例作为参数。