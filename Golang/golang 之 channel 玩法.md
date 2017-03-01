# golang 之 channel 玩法

标签（空格分隔）： golang

---

## [Go-简洁的并发](http://www.thinksaas.cn/topics/0/594/594568.html)

本文从**并发模式**的角度讲述了 channel 和 goroutine 的使用模型；主要涉及一下内容：

- **生成器**（让函数返回“服务”，用于数据读取，生成 ID，实现定时器等）；
- **多路复用**（将若干个相似的小“服务”整合成一个大“服务”）；
- **Furture 技术**（实现函数调用和函数参数准备两个过程的完全解耦，提供更大的自由度和并发度；还可以用来实现 pipe filter）；
- **并发循环**（利用多核提升性能，解决 CPU 热点）；
- Chain Filter 技术（属于 pipe filter 的特殊情况）；
- **共享变量**（一般来说，协程之间不推荐使用共享变量来交互，但是在一些场合，使用共享变量也是可取的）；
- **协程泄漏**（一般而言，协程执行结束后就会销毁；协程也会占用内存；只有两种情况会导致协程无法结束：一种情况是协程想从一个通道读数据，但无人往这个通道写入数据，或许这个通道已经被遗忘了；还有一种情况是程想往一个通道写数据，可是由于无人监听这个通道，该协程将永远无法向下执行）；


----------


> 基于 goroutine 服务化的生成器

```golang
// 函数 rand_generator_2，返回 通道(Channel)
func rand_generator_2() chan int {
	// 创建通道
	out := make(chan int)
	// 创建协程
	go func() {
		for {
			//向通道内写入数据，如果无人读取会等待
			out <- rand.Int()
		}
	}()
	return out
}

func main() {
	// 生成随机数作为一个服务
	rand_service_handler := rand_generator_2()
	// 从服务中读取随机数并打印
	fmt.Printf("%dn", <-rand_service_handler)
}
```

> 多路复用

```golang
// 函数 rand_generator_3 ，返回通道(Channel)
func rand_generator_3() chan int {
	// 创建两个随机数生成器服务
	rand_generator_1 := rand_generator_2()
	rand_generator_2 := rand_generator_2()

	//创建通道
	out := make(chan int)

	//创建协程
	go func() {
		for {
			//读取生成器1中的数据，整合
			out <- <-rand_generator_1
		}
	}()
	go func() {
		for {
			//读取生成器2中的数据，整合
			out <- <-rand_generator_2
		}
	}()
	return out
}
```

> future 技术示例

```golang
//一个查询结构体
type query struct {
	//参数Channel
	sql chan string
	//结果Channel
	result chan string
}

//执行Query
func execQuery(q query) {
	//启动协程
	go func() {
		//获取输入
		sql := <-q.sql
		//访问数据库，输出结果通道
		q.result <- "get " + sql
	}()

}

func main() {
	//初始化Query
	q :=
		query{make(chan string, 1), make(chan string, 1)}
	//执行Query，注意执行的时候无需准备参数
	execQuery(q)

	//准备参数
	q.sql <- "select * from table"
	//获取结果
	fmt.Println(<-q.result)
}
```

> 并发循环

```golang
//建立计数器
sem := make(chan int, N); 
//FOR循环体
for i,xi := range data {
	//建立协程
    go func (i int, xi float) {
        doSomething(i,xi);
		//计数
        sem <- 0;
    } (i, xi);
}
// 等待循环结束
for i := 0; i < N; ++i { <-sem }
```

> 基于 Chain-Filter 生成素数

```golang
// A concurrent prime sieve

package main

// Send the sequence 2, 3, 4, ... to channel 'ch'.
func Generate(ch chan<- int) {
	for i := 2; ; i++ {
		ch <- i // Send 'i' to channel 'ch'.
	}
}

// Copy the values from channel 'in' to channel 'out',
// removing those divisible by 'prime'.
func Filter(in <-chan int, out chan<- int, prime int) {
	for {
		i := <-in // Receive value from 'in'.
		if i%prime != 0 {
			out <- i // Send 'i' to 'out'.
		}
	}
}

// The prime sieve: Daisy-chain Filter processes.
func main() {
	ch := make(chan int) // Create a new channel.
	go Generate(ch)      // Launch Generate goroutine.
	for i := 0; i < 10; i++ {
		prime := <-ch
		print(prime, "n")
		ch1 := make(chan int)
		go Filter(ch, ch1, prime)
		ch = ch1
	}
}
```

> 实现一个协程安全的共享变量

```golang
//共享变量有一个读通道和一个写通道组成
type sharded_var struct {
	reader chan int
	writer chan int
}

//共享变量维护协程
func sharded_var_whachdog(v sharded_var) {
	go func() {
		//初始值
		var value int = 0
		for {
			//监听读写通道，完成服务
			select {
			case value = <-v.writer:
			case v.reader <- value:
			}
		}
	}()
}

func main() {
	//初始化，并开始维护协程
	v := sharded_var{make(chan int), make(chan int)}
	sharded_var_whachdog(v)

	//读取初始值
	fmt.Println(<-v.reader)
	//写入一个值
	v.writer <- 1
	//读取新写入的值
	fmt.Println(<-v.reader)
}
```

> 使用**超时**避免读堵塞，使用**缓冲**避免写堵塞

```golang
func never_leak(ch chan int) {
	//初始化timeout，缓冲为1
	timeout := make(chan bool, 1)
	//启动timeout协程，由于缓存为1，不可能泄露
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()
	//监听通道，由于设有超时，不可能泄露
	select {
	case <-ch:
		// a read from ch has occurred
	case <-timeout:
		// the read from ch has timed out
	}
}
```


----------


## [Curious Channels](https://dave.cheney.net/2013/04/30/curious-channels)

> 中文翻译版本：[这里](https://mikespook.com/2013/05/%E7%BF%BB%E8%AF%91%E7%BB%9D%E5%A6%99%E7%9A%84-channel/)

阅读建议：主要看英文版，中文翻译版本有些地方翻译的不知所以；

本文要点：

- A `closed` channel **never** blocks (Once a channel has been closed, you cannot send a value on this channel, but you can still receive from the channel)
- A `nil` channel **always** blocks (a channel value that has not been initalised, or has been set to nil will always block)


----------


> closed channel 不会阻塞读取操作（closed channel 可以被无限读取）

```golang
package main

import "fmt"

func main() {
        ch := make(chan bool, 2)
        ch <- true
        ch <- true
        close(ch)

        for i := 0; i < cap(ch) +1 ; i++ {
                // first var is the zero value for that channel’s type
                // second var is the open state of the channel
                v, ok := <- ch
                fmt.Println(v, ok)
        }
}
```

> detect if your channel is closed in the `range` idiom （用法说明：通过 close 关闭 channel 时，能够被 range 检测出来）

```golang
package main

import "fmt"

func main() {
        ch := make(chan bool, 2)
        ch <- true
        ch <- true
        close(ch)

        for v := range ch {
                fmt.Println(v) // called twice
        }
}
```

> 基于 sync.WaitGroup + select + quit channel （不完美方案）

```golang
package main

import (
        "fmt"
        "sync"
        "time"
)

func main() {
        finish := make(chan bool)
        var done sync.WaitGroup
        done.Add(1)
        go func() {
                select {
                case <-time.After(1 * time.Hour):
                case <-finish:
                }
                done.Done()
        }()
        t0 := time.Now()
        finish <- true // send the close signal
        done.Wait()    // wait for the goroutine to stop
        fmt.Printf("Waited %v for goroutine to stop\n", time.Since(t0))
}
```

这种实现可能遇到的问题：

1. 由于 finish 这个 channel 是 unbuffered 的，所以，一旦在目标 goroutine 中忘记了读取 finish channel ，则会导致发送端阻塞；
2. 即使将 finish channel 改成 buffered 的，其实也只能算治标不治本的办法，因为只解决了发送端的阻塞问题，无法发现接收端是否存在泄漏；另外，确定 buffer 的大小也是一个问题，有些情况下，根本无法确定；
3. 如果仅从解决发送端阻塞的角度，将 finish <- true 封装在 select 语句中，明显也是一个二逼方案；

>  基于 sync.WaitGroup + select + closed channel （好方案）

```golang
package main

import (
        "fmt"
        "sync"
        "time"
)

func main() {
        const n = 100
        finish := make(chan bool)
        var done sync.WaitGroup
        for i := 0; i < n; i++ { 
                done.Add(1)
                go func() {
                        select {
                        case <-time.After(1 * time.Hour):
                        case <-finish:
                        }
                        done.Done()
                }()
        }
        t0 := time.Now()
        close(finish)    // closing finish makes it ready to receive
        done.Wait()      // wait for all goroutines to stop
        fmt.Printf("Waited %v for %d goroutines to stop\n", time.Since(t0), n)
}
```

这种实现的好处是：利用 channel 实现信号通知任意数量的 goroutines ，无需知道具体数目，无需担心死锁问题；


> 基于 sync.WaitGroup + select + closed channel + chan struct{} （推荐方案）

```golang
package main

import (
        "fmt"
        "sync"
        "time"
)

func main() {
        finish := make(chan struct{})
        var done sync.WaitGroup
        done.Add(1)
        go func() {
                select {
                case <-time.After(1 * time.Hour):
                case <-finish:
                }
                done.Done()
        }()
        t0 := time.Now()
        close(finish)
        done.Wait()
        fmt.Printf("Waited %v for goroutine to stop\n", time.Since(t0))
}
```

> channel 未初始化，则默认初始化为对应的 nil 值

- 阻塞写入

```golang
package main

func main() {
        var ch chan bool
        ch <- true // blocks forever
}
```

- 阻塞读取

```golang
package main

func main() {
        var ch chan bool
        <- ch // blocks forever
}
```

> use the closed channel idiom to wait for multiple channels to close （利用 nil channel 的阻塞特性）

```golang
package main

import (
        "fmt"
        "time"
)

func WaitMany(a, b chan bool) {
        for a != nil || b != nil {
                select {
                case <-a:
                        a = nil 
                case <-b:
                        b = nil
                }
        }
}

func main() {
        a, b := make(chan bool), make(chan bool)
        t0 := time.Now()
        go func() {
                close(a)
                close(b)
        }()
        WaitMany(a, b)
        fmt.Printf("waited %v for WaitMany\n", time.Since(t0))
}
```

----------


## [Go并发模式：管道和取消](http://air.googol.im/post/go-concurrency-patterns-pipelines-and-cancellation/)

> 英文原本：[这里](https://blog.golang.org/pipelines)

> 另外一个中文版本：[Go语言并发模型：像Unix Pipe那样使用channel](https://segmentfault.com/a/1190000006261218)（该版本有一些译者自己的补充说明）

阅读建议：中文翻译版本还不错，可以适当对比英文原文；

本文要点：

- There's no formal definition of a `pipeline` in Go; it's just one of many kinds of concurrent programs. 
- Multiple functions can read from the same channel until that channel is closed; this is called `fan-out`. This provides a way to distribute work amongst a group of workers **to parallelize CPU use and I/O**.
- A function can read from multiple inputs and proceed until all are closed by multiplexing the input channels onto a single channel that's closed when all the inputs are closed. This is called `fan-in`.
- Sends on a closed channel panic, so it's important to ensure all sends are done before calling close. 
- **"Stopping short" 原则（“尽快终止”原则）**：
    - stages close their outbound channels when all the send operations are done. 
    - stages keep receiving values from inbound channels until those channels are closed.

    效果：This pattern allows each receiving stage to be written as a range loop and ensures that all goroutines exit once all values have been successfully sent downstream.
    
- **资源泄漏问题**：goroutines consume memory and runtime resources, and heap references in goroutine stacks keep data from being garbage collected. Goroutines are not garbage collected; they must exit on their own.
- **done channel 的使用（显式取消）**；a receive operation on a closed channel can always proceed immediately, yielding the element type's zero value.
- **构建 pipeline 的指导原则**：
    - stages close their outbound channels when all the send operations are done.
    - stages keep receiving values from inbound channels until those channels are closed or the senders are unblocked.
    - Pipelines unblock senders either by ensuring there's enough buffer for all the values that are sent or by explicitly signalling senders when the receiver may abandon the channel.

本文最后给出一个对目录下所有文件计算 md5 值的程序，讲解了如何基于 pipeline 实现并行计算，值得学习（建议看英文版）；






