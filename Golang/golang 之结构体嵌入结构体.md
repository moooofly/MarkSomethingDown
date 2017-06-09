# golang 之结构体嵌入结构体

> 本文仅针对**结构体嵌入结构体**的情况进行讨论；

## 问题

- 为什么要使用匿名组合，组合是否匿名的差别在哪里；
- 匿名组合时是否使用指针类型的差别在哪里；
- 是不是应该全都使用匿名组合；

----------

## （成员方法）继承和覆盖

Go 语言采用**组合**文法提供了**继承**功能；

```golang
type Base struct {
    Name string
}

func (base *Base) Foo() {...}
func (base *Base) Bar() {...}

type Foo struct {
    Base    // 匿名组合
    ...     // 其他成员
}

func (foo *Foo) Bar() {
    foo.Base.Bar()
    ...
}
```

若“派生类” Foo 没有通过定义同名方法（签名可不同）**覆盖（改写或扩展）**“基类” Base 的成员方法时，则 Base 的成员方法就会被“继承”；

例如上面例子中，调用 `foo.Foo()` 和调用 `foo.Base.Foo()` 效果一致。

若定义同名方法（签名可不同），则产生**覆盖**效果；


## （成员变量）继承和覆盖

> 组合中的名字冲突问题

- Case 1

```golang
type X struct {
    Name string
}
type Y struct {
    X
    Name string
}
```

组合的类型和被组合的类型都包含一个 Name 成员，会不会有问题呢？答案是否定的。所有针对 Y 类型的 Name 成员的访问，都只会访问到最外层的那个 Name 变量，即 X.Name 变量相当于被隐藏起来了。

- Case 2

```golang
type Logger struct {
    level int
}
type Y struct {
    *Logger
    Name string
    *log.Logger  // 以指针方式匿名组合log.logger结构
}
```

这里显然会有问题，因为匿名组合类型相当于以其类型名称（去掉包名部分）作为成员变量的名字。按此规则，Y 类型中相当于存在两个名为 Logger 的成员，虽然类型不同。因此，我们预期会收到编译错误。有意思的是，这个编译错误并不是一定会发生的。假如这两个 Logger 在定义后再也没有被用过，那么编译器将直接忽略掉这个冲突问题，直至开发者开始使用其中的某个 Logger 。

> 重要补充：
> - 对于结构体类型的多层嵌入来说，可以在被嵌入的结构体类型的值是那个像调用其自己的字段或方法那样调用**任意深度**的嵌入类型值的字段或方法，只要其未被隐藏；
> - 被嵌入的结构体类型的字段或方法可以隐藏**任意深度**的嵌入类型的同名字段或方法；本质上，任何较浅层次的嵌入类型的字段或方法都会隐藏较深层次嵌入类型的同名字段或方法；
> - 上述隐藏是可以**交叉进行**的，即字段可以隐藏方法，方法也可以隐藏字段，之遥它们的名字相同即可；

## 匿名组合 & 指针类型

可以使用指针类型进行匿名组合：

```golang
type Foo struct {
    *Base
    ...
}
```

该代码仍然有“派生”的效果，但基于 Foo 类型创建实例的时候，需要提供一个 Base 类型实例的指针；

关于匿名组合是否**使用指针类型**的示例：

```golang
type S struct {
    T1        // 字段名自动为 T1
    *T2       // 字段名自动为 T2
    P.T3      // 字段名自动为 T3
    *P.T4     // 字段名自动为 T4
    x, y int  // 非匿名字段 x, y
}
```

重要结论：

- **如果在结构体 S 中包含一个匿名字段 T** ，那么 S 和 *S 的方法集合中都会包含接收者类型为 T 的方法；除此之外，*S 的方法集合中还会包含接收者类型为 *T 的方法；
- **如果在结构体 S 中包含一个匿名字段 \*T** ，那么 S 和 *S 的方法集合中都会包含接收者类型为 T 或 *T 的所有方法；

## 内存布局

与其他语言不同，Go 语言很清晰地告诉你的内存布局是怎样的，你还可以随心所欲地修改内存布局，如：

```golang
type Foo struct {
    ...    // 其他成员
    Base
}
```

这段代码从语义上来说，和上面的例子并无不同，但内存布局发生了改变，“基类” Base 的数据放在了“派生类” Foo 的最后。


## 价值

- Case 1

```golang
type Job struct {
    Command string
    *log.Logger   // 以指针方式匿名组合log.logger结构
}
```

在合适的赋值后，在 Job 类型的所有成员方法中，可以很舒适地借用所有 `log.Logger` 提供的方法。比如：

```golang
func (job *Job) Start() {
    job.Log("starting now...")
    ...//做一些事情
    job.Log("started.")
}
```

当我们想要精确控制 Logger 的方法时，可以

```golang
func (job *Job) Logf(format string, args ...interface{}) {
    job.Logger.Logf("%q: %s", job.Command, fmt.Sprintf(format, args...))
}
```

对于 Job 的实现者来说，他甚至根本就不用意识到 `log.logger` 类型的存在，这就是匿名组合的魅力所在。在实际工作中，只有合理利用才能最大发挥这个功能的价值。

注意：不管是**非匿名的类型组合**还是**匿名组合**，被组合的类型所包含的方法虽然都升级了外部这个组合类型的方法，但其实它们被组合的方法**调用时接收者并没有改变**。比如上面这个 Job 例子，即使组合后调用的方式变成了 `job.Log(...)` ，但 Log 函数的接收者仍然是 `log.Logger` 指针，因此在 Log 中不可能访问到 job 的其他成员方法和变量。

- Case 2

在 `io` 包中有

```golang
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type ReadWriter interface {
    Reader
    Writer
}
```

在 `bufio` 包有两个结构体类型，`bufio.Reader` 和 `bufio.Writer` 各自实现了来自 `io` 包中的对应接口。进一步定义如下结构体：

```golang
// ReadWriter 包含了指向一个 Reader 和一个 Writer 的指针
// ReadWriter 实现了 io.ReadWriter 接口
type ReadWriter struct {
    *Reader  // 即 *bufio.Reader
    *Writer  // 即 *bufio.Writer
}
```

这里就有个问题需要回答，为什么要定义成匿名形式？而不是如下具名形式？

```golang
type ReadWriter struct {
    reader *Reader
    writer *Writer
}
```

问题关键在于：**具名组合的情况下，若要使成员变量的方法提升（promote）为结构体的方法（以满足 `io` 接口），就需要提供转换（forwarding）方法**，如下所示：

```golang
func (rw *ReadWriter) Read(p []byte) (n int, err error) {
    return rw.reader.Read(p)
}
```

而通过直接嵌入结构体（即匿名组合），就可以不必这么繁琐，因为**嵌入类型的方法会被自动继承（进而满足一些接口定义）**；

补充：使用**具名嵌入**时能够对成员变量的可见性进行控制（似乎是个鸡肋功能），例如，Foo 结构具名嵌入 Base 结构时，可以指定名字 base ，因此，其他包将无法直接访问 base 成员变量相关内容；


## 参考

- [golang中的匿名组合](http://studygolang.com/wr?u=http%3a%2f%2f11317783.blog.51cto.com%2f11307783%2f1877972)
- [匿名字段和内嵌结构体](https://github.com/Unknwon/the-way-to-go_ZH_CN/blob/master/eBook/10.5.md)
- [Effective Go](https://golang.org/doc/effective_go.html#embedding)