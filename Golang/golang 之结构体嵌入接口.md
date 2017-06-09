# golang 之结构体嵌入接口

## 一篇文章

> 参考：[Golang embedded interface on parent struct](https://stackoverflow.com/questions/38043678/golang-embedded-interface-on-parent-struct)


One thing you seem to be missing is **how embedding interfaces affects a structure in Go**. See, embedding promotes all of the methods of the embedded type (`struct` or `interface`, doesn't matter) to be methods of the parent type, but called using the embedded object as the receiver.

The practical **side effect** of this is that **embedding an interface into a structure guarantees that that structure fulfills the interface it is embedding**, because it by definition has all of the methods of that interface. **Trying to call any of those methods without defining something to fill that interface field in the struct, however, will panic**, as that interface field defaults to `nil`.

Embedding of interfaces, rather than being for `duck-typing`, is more about `polymorphism`. 

Do not mix **`interface` embedding** with **`struct` embedding**.

If you embed interfaces on a struct, you are actually adding new fields to the struct with the name of the interface so if you don't init those, you will get panics because they are nil.


## 实验

### struct 内嵌 interface 但不实现任何 method

```
package main

import (
	"fmt"
)

type Human interface {
    eat()
    sleep()
}

type Tom struct {
    age int
    Human
}

func main() {	
	tom := Tom{}
	fmt.Printf("tom => %v   tom.Human => %v\n", tom, tom.Human)  // tom => {0 <nil>}   tom.Human => <nil>
	
	var human Human
	human = &tom     // main.go:22: human declared and not used
	//tom.eat()    // panic: runtime error: invalid memory address or nil pointer dereference
}
```

小节：

- Human 接口的匿名嵌入导致同名成员变量的产生；
- Human 接口的匿名嵌入导致其定义中的 method **自动提升**，产生的实际副作用如下：
    - （由于自动提升的缘故）Tom 在未实现任何 Human 接口方法情况下，就可以对其进行赋值（上面的错误信息是说 human 接口声明后未使用，说明赋值本身是合法的，虽然赋值后啥也干不了）；
    - 但 Human 被提升的所有 method 实际都不可调用（可以按如下方式理解：Tom 拿到了 Human 指定的“虚方法”，但由于没有具体实现，因此不能被调用）；


### struct 内嵌 interface 但只实现部分 method

```
package main

import (
	"fmt"
)

type Human interface {
    eat()
    sleep()
}

type Tom struct {
    age int
    Human
}

// 产生覆盖 Human 自动提升的 eat 方法的效果
func (t *Tom) eat() {
    fmt.Println("I'm eating!")
}

/*
func (t *Tom) sleep() {
    fmt.Println("I'm sleeping!")
}
*/

func main() {	
	tom := Tom{}
	fmt.Printf("tom => %v   tom.Human => %v\n", tom, tom.Human)  // tom => {0 <nil>}   tom.Human => <nil>
	
	var human Human
	human = &tom
	
	human.eat()       // - 1
	//human.sleep()   // - 2  panic: runtime error: invalid memory address or nil pointer dereference
	
	tom.eat()         // - 3
	//tom.sleep()     // - 4  panic: runtime error: invalid memory address or nil pointer dereference
}
```

小节：

- 位置 1 调用的是 Tom 实现的 `eat` 方法（覆盖了 Human 处提升的“虚方法”）；
- 位置 2 调用的是 Tom 从 Human 处得到的 `sleep` “虚方法”；
- 位置 3 调用的是 Tom 实现的 `eat` 方法（覆盖了 Human 处提升的“虚方法”）；
- 位置 4 调用的是 Tom 从 Human 处得到的 `sleep` “虚方法”；


----------


## 复杂组合

> 以下代码取自 `beats` ；

```golang
type elasticsearchOutput struct {
	index    outil.Selector  // 具名组合 struct
	beatName string
	pipeline *outil.Selector

	mode mode.ConnectionMode // 具名组合 interface
	topology                 // 匿名组合 struct

	template      map[string]interface{}
	template2x    map[string]interface{}
	templateMutex sync.Mutex
}
```

在 `select.go` 中有

```
type Selector struct {
	sel SelectorExpr   // 具名组合 interface
}

type SelectorExpr interface {
	sel(evt common.MapStr) (string, error)
}
```

在 `mode.go` 中有

```golang
type ConnectionMode interface {
	Close() error

	PublishEvents(sig op.Signaler, opts outputs.Options, data []outputs.Data) error

	PublishEvent(sig op.Signaler, opts outputs.Options, data outputs.Data) error
}

type Connectable interface {
	Connect(timeout time.Duration) error

	Close() error
}

type ProtocolClient interface {
	Connectable       // 匿名组合 interface

	PublishEvents(data []outputs.Data) (nextEvents []outputs.Data, err error)

	PublishEvent(data outputs.Data) error
}
```

在 `topology.go` 中有

```golang
type topology struct {
	clients []mode.ProtocolClient

	TopologyMap atomic.Value
}
```

## 问题

### 为何需要在 struct 中嵌入 interface 呢？

提出这个问题的原因在于，golang 本身推崇的是 `duck-typing` ；既然如此，嵌入 interface 似乎是没有必要的，而事实不是这样；

上面给出的引文中有一句说明非常到位：

> Embedding of interfaces, rather than being for `duck-typing`, is more about `polymorphism`. 

用自己的话说就是：**如果你在 struct 中进行 interface 嵌入，那么你要的其实是 `polymorphism`（多态），而不是 `duck-typing`** ；更进一步，由于目标不是 `duck-typing` ，因此上面**实验**中给出的使用方法虽然在语言角度上没问题，但却是错误用法；

另外补充一点：**在 struct 嵌入 interface 的场景中，是否匿名嵌入似乎并不重要了（虽然大部分情况下会选择具名形式），因为效果是一样的**；

至于“多态”，完全和其他面向对象语言里的用法相同，即当需要对多种不同具体实现进行抽象时，可以通过接口方式实现统一调用；

比如底层实现了 TCP 和 UDP 两种协议类型，但上层只想按照统一的方式进行接口调用；


### interface 之间组合只使用匿名？

接口之间只会使用**匿名组合**来进行扩展（没见过具名的）；


