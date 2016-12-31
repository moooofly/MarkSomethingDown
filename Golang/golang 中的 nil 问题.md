# golang 中的 nil 问题

标签（空格分隔）： golang nil interface

---

## 问题代码

如下代码在函数 `checkError` 中实现了针对 err 是否为 `nil` 的判定，但是此处判定代码的编写是否正确呢？如果不正确，这样编写存在什么问题？

```golang
func checkError (err error) {
    if err != nil {   // 问题点 2
		panic(err)
    }
}

type Error struct {
    errCode uint8
}
func (e *Error) Error() string {
        switch e.errCode {
        case 1:
                return "file not found"
        case 2:
                return "time out"
        case 3:
                return "permission denied"
        default:
                return "unknown error"
         }
}

func main() {
    var e *Error  // 问题点 1
    checkError(e)
}
```

上述代码的本意是：在未对 `e` 进行显式初始化时，`e` 将被默认初始化为零值，将零值传入 `checkError()` 中与 `nil` 进行想等判定，应该不会触发 `panic(err)` 调用；

## 原因说明

事实上，上述代码总是会触发 `panic(err)` 调用；原因总结如下：

- 在通过 `var e *Error` 定义 `e` 的时候，运行时系统会将其默认初始化为对应类型的零值，即 (*Error)nil ；
- 当 `e` 被传递到 `checkError()` 函数中，通过 error 接口拿到后，产生到效果为 `interface{Error() string}((*Error)nil)` ；
- 当通过 `if err != nil` 进行相等判定时，由于两者类型不同，因此永远得到 false 结果，所以总是会触发 `panic()` ；


----------


其实，若想要深刻理解上述错误原因，实际上要正确理解如下问题：

- `nil` 是什么？如何使用？
- `error` 是什么？接收入参后如何工作？

下面给出一些能够解答上述问题的相关材料；

## nil 

在 golang 源码 `libexec/src/builtin/builtin.go` 中有

```golang
// nil is a predeclared identifier representing the zero value for a
// pointer, channel, func, interface, map, or slice type.
var nil Type // Type must be a pointer, channel, func, interface, map, or slice type

// Type is here for the purposes of documentation only. It is a stand-in
// for any Go type, but represents the same type for any given function
// invocation.
type Type int
```

从上面的注释中可以看出：`nil` 是一个预定义**标识符**（identifier），其代表（用作）一些类型的**零值**；这些类型包括：pointer, channel, func, interface, map, slice ；

> 扩展问题：`nil` 是有类型的，还是无类型的？identifier 和 value 的关系是？
> 回答：虽然 `nil` 本身在底层是基于 int 定义的，但应该将其理解成**无类型**的值，且仅作为其他静态类型的**零值**使用；


## error

在 golang 源码 `libexec/src/builtin/builtin.go` 中有

```golang
// The error built-in interface type is the conventional interface for
// representing an error condition, with the nil value representing no error.
type error interface {
	Error() string
}
```

可以看出：error 是一个内置接口类型，以 `nil` 标记没有错误；

## [Why is my nil error value not equal to nil?](https://golang.org/doc/faq#nil_error)

从底层实现上讲，interface 是由两个元素构成的，即 **type** 和 **value** ；其中，**value** 由一个任意的具体值表示，称作 interface 的 dynamic value ；而 **type** 则对应该 value 的类型（即 dynamic type）；对于 `int` 类型值 3 来说，interface 的值大致等价于 `(int, 3)` ；

如果我们说 interface 的值为 `nil` ，则（只可能）对应的是其内部 value 和 type 均未设置的情况，即 `(nil, nil)` ；尤其需要知道的是：`nil` interface 和 `nil` type 是对应的；如果我们保存了一个 type 为 `*int` 的 `nil` 指针到一个 interface 中，那么其内部 type 则为 `*int` ，而不管指针的具体值是什么：`(*int, nil)`；对于这样的 interface 值来说，虽然其内部的 value 为 `nil` 但我们仍旧认为 interface 本身是 `non-nil` 的；

这种令人困惑的情况总是会出现在 `nil` 值被保存到 interface 中的时候，比如函数返回 `error` interface 的情况：

```golang
func returnsError() error {
	var p *MyError = nil
	if bad() {
		p = ErrBad
	}
	return p // Will always return a non-nil error.
}
```

在一切正常的情况下，函数会返回值为 `nil` 的 p 指针，此时作为返回值的 `error` interface 中包含的内容为 `(*MyError, nil)` ；如果调用者将该返回值与 `nil` 进行比较，则看起来好像总是存在错误，因此永远都不会相等（==），虽然事实上代码中没有任何错误发生；为了返回一个正确的 `nil` 给调用者，函数必须按照如下代码返回 `nil` ：

```golang
func returnsError() error {
	if bad() {
		return ErrBad
	}
	return nil
}
```

**最佳实践**：总是使用函数签名中指定的 `error` 类型进行返回，而不是返回一个具体的类型，例如 `*MyError` ，以确保作为返回值所创建的 error 能被正确使用；例如，`os.Open` 要么返回 `nil` ，要么返回类型为 `*os.PathError` 的 error ；

无论何时，只要 interface 被使用，就可能会遇到上面的问题；因此在使用 interface 时要时刻谨记，只要存在具体值（除了 `nil` 以外的值）被保存到 interface 中的情况，那么 interface 的值将不会是 `nil` ；更多信息，可以查看 [The Laws of Reflection](https://blog.golang.org/laws-of-reflection)


## [Understanding Go's `nil` value](http://www.gmarik.info/blog/2016/understanding-golang-nil-value/)

下面这段代码用于输出各种指定静态类型下的 `nil` 值；

```golang
package main

import "fmt"

func main() {

	fmt.Println("-----------------------------------------------------------------------------------------------")
	fmt.Println("Type                       default-format-value(%v)\tgo-style-type(%T)\tgo-style-value(%#v)")
	fmt.Println("-----------------------------------------------------------------------------------------------")
	fmt.Printf("Func type nil          ==> %v\t\t\t%T\t\t\t%#v\n", (func())(nil), (func())(nil), (func())(nil))
	fmt.Printf("Map type nil           ==> %v\t\t\t%T\t%#v\n", map[string]string(nil), map[string]string(nil), map[string]string(nil))
	fmt.Printf("Interface{} type nil   ==> %v\t\t\t%T\t\t\t%#v\n", nil, nil, nil)
	fmt.Printf("Interface{} type nil   ==> %v\t\t\t%T\t\t\t%#v\n", interface{}(nil), interface{}(nil), interface{}(nil))
	fmt.Printf("Channel type nil       ==> %v\t\t\t%T\t\t%#v\n", (chan struct{})(nil), (chan struct{})(nil), (chan struct{})(nil))
	fmt.Printf("Pointer type nil       ==> %v\t\t\t%T\t\t%#v\n", (*struct{})(nil), (*struct{})(nil), (*struct{})(nil))
	fmt.Printf("Pointer type nil       ==> %v\t\t\t%T\t\t\t%#v\n", (*int)(nil), (*int)(nil), (*int)(nil))
	fmt.Printf("Slice type nil         ==> %v\t\t\t\t%T\t\t%#v\n", []string(nil), []string(nil), []string(nil))
}
```

从输出中可以看到 `%v` 和 `%#v` 在特定静态类型零值输出上的差别，以及查看零值当前属于哪种静态类型（`%T`）；

```shell
root@vagrant-ubuntu-trusty:~/workspace/CODE/Golang# go run nil_issue_analysis.go
-----------------------------------------------------------------------------------------------
Type                       default-format-value(%v)	go-style-type(%T)	go-style-value(%#v)
-----------------------------------------------------------------------------------------------
Func type nil          ==> <nil>			func()			(func())(nil)
Map type nil           ==> map[]			map[string]string	map[string]string(nil)
Interface{} type nil   ==> <nil>			<nil>			<nil>
Interface{} type nil   ==> <nil>			<nil>			<nil>
Channel type nil       ==> <nil>			chan struct {}		(chan struct {})(nil)
Pointer type nil       ==> <nil>			*struct {}		(*struct {})(nil)
Pointer type nil       ==> <nil>			*int			(*int)(nil)
Slice type nil         ==> []				[]string		[]string(nil)
root@vagrant-ubuntu-trusty:~/workspace/CODE/Golang#
```

----------


# [变量初始化问题](https://golang.org/ref/spec#Variables)

一个变量，代表的是值的存储位置；而变量的类型，决定了其允许保存的值的集合；

> A variable is a storage location for holding a value. The set of permissible values is determined by the variable's type.

会为命名变量预留存储位置的地方：

- 变量声明；
- 函数参数和返回值；
- 用于函数声明的签名；
- 函数字面量；

> A variable declaration or, for function parameters and results, the signature of a function declaration or function literal reserves storage for a named variable. Calling the built-in function new or taking the address of a composite literal allocates storage for a variable at run time. Such an anonymous variable is referred to via a (possibly implicit) pointer indirection.

结构型变量（array/slice/struct）包含的元素和域可以被单独访问，行为同变量；

> Structured variables of array, slice, and struct types have elements and fields that may be addressed individually. Each such element acts like a variable.

变量的**静态类型**是指：

- 变量声明中指定的类型；
- 在 `new` 调用中指定的类型；
- 复合字面量对应的类型；
- 结构型变量中的元素的类型；

接口类型的变量还具有一个独立的**动态类型**，对应的是运行时赋值给接口变量的值的**静态类型**（除非该值为预定义的标识符 `nil`，此时是无类型的）；

动态类型在执行过程中可能会发生变化，但保存在接口变量中的值总是可以被赋值给相应静态类型的变量；

> The **static type** (or just type) of a variable is the type given in its declaration, the type provided in the `new` call or composite literal, or the type of an element of a structured variable. Variables of interface type also have a distinct **dynamic type**, which is the concrete type of the value assigned to the variable at run time (unless the value is the predeclared identifier `nil`, which has no type). The dynamic type may vary during execution but values stored in interface variables are always assignable to the static type of the variable.

（下面这几行代码能够说明很多问题）

```golang
var x interface{}  // x is nil and has static type interface{}
var v *T           // v has value nil, static type *T
x = 42             // x has value 42 and dynamic type int
x = v              // x has value (*T)(nil) and dynamic type *T
```

**变量值的提取**是通过在表达式中引用相应的变量实现的；变量的值总是最近一次赋值的内容；如果一个变量未曾被赋予过值，那么其值被赋予其类型的**零值**；

> A variable's value is retrieved by referring to the variable in an expression; it is the most recent value assigned to the variable. If a variable has not yet been assigned a value, its value is the zero value for its type.


# [零值问题](https://golang.org/ref/spec#The_zero_value)

当一个变量被分配了存储空间时，即

- 或者通过直接声明；
- 或者通过 `new` 调用；

或一个新（变量）值被创建时，即

- 或者通过复合字面量；
- 或者通过 `make` 调用；

若满足未显式初始化这个条件，那么相应的变量或值将被赋予默认值，即相应类型的零值；

- boolean 类型变量的零值为 `false` ；
- integer 类型变量的零值为 `0` ；
- float 类型变量的零值为 `0.0` ；
- string 类型变量的零值为 `""` ；
- pointer, function, interface, slice, channel 和 map 类型变量的零值为 `nil` ；

上述初始化行为是递归处理的（对于复合类型由意义）；

> When storage is allocated for a variable, either through a declaration or a call of `new`, or when a new value is created, either through a composite literal or a call of `make`, and no explicit initialization is provided, the variable or value is given a default value. Each element of such a variable or value is set to the zero value for its type: false for booleans, 0 for integers, 0.0 for floats, "" for strings, and `nil` for pointers, functions, interfaces, slices, channels, and maps. This initialization is done recursively, so for instance each element of an array of structs will have its fields zeroed if no value is specified.

These two simple declarations are equivalent:

```
var i int
var i int = 0
```

After
```
type T struct { i int; f float64; next *T }
t := new(T)
```
the following holds:
```
t.i == 0
t.f == 0.0
t.next == nil
```
The same would also be true after
```
var t T
```




