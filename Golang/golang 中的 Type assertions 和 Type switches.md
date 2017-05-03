
> 此文实际为官方文档翻译；

# [Type assertions](https://golang.org/ref/spec#Type_assertions)

给定类型为接口的表达式 x 和实际类型 T ，下面的表达式

```golang
x.(T)
```

能够断言 x （是否）非 `nil` ，以及保存在 x 中的值为 T 类型；表达式 `x.(T)` 就被称作 **type assertion** ；

更准确的说，**如果 T 不是接口类型**，那么 `x.(T)` 将断言 x 的 dynamic type 是否和类型 T 相同；在这种情况下，T 必须实现 x 的（接口）类型，否则类型断言将会失败；因为在这种情况下 x 将无法保存 T 类型的值；**如果 T 是接口类型**，那么 `x.(T)` 将断言 x 的 dynamic type 是否实现了接口 T（可以理解成将 x 的 dynamic type 转换成 T 接口的类型）；

**如果 type assertion 成立**，那么表达式的值就是保存在 x 中的值，并且类型为 T ；**如果 type assertion 不成立**，那么将触发 run-time panic ；换句话说，即使 x 的 dynamic type 只能在 run time 时才知道，但 `x.(T)` 的类型在编写程序的时候应该明确知道是 T ；（is known to be T in a correct program.）

```golang
var x interface{} = 7  // x has dynamic type int and value 7
i := x.(int)           // i has type int and value 7

type I interface { m() }

func f(y I) {
	s := y.(string)        // illegal: string does not implement I (missing method m)
	r := y.(io.Reader)     // r has type io.Reader and the dynamic type of y must implement both I and io.Reader
	…
}
```

若 type assertion 被用在赋值（assignment）或初始化（initialization）语句中，则具有如下特殊形式：

```golang
v, ok = x.(T)
v, ok := x.(T)
var v, ok = x.(T)
var v, ok T1 = x.(T)
```

此时需要增加一个额外的无类型 boolean 变量；变量 ok 为 true 则表明断言成立；否则，表明断言不成立，此时 v 的值将被赋值为类型 T 的零值；**在这种情况下，不会触发运行时崩溃**；


# [Type switches](https://golang.org/ref/spec#Switch_statements)

type switch 是针对 types 进行的比较，而非针对 values ；其非常类似于 switch 表达式；可以认为它是一种特殊的 switch 表达式，具有 type assertion 的形式，使用了保留字 type 而非实际的具体 type ：

```golang
switch x.(type) {
// cases
}
```

每一个 Case 都会使用实际类型 T 去和表达式 x 的 dynamic type 去比较；和 type assertions 一样，x 必须为 interface 类型，同时每一个列出在 case 中的 non-interface 类型的 T 都必须实现 x 的类型；在 type switch 中每一个 case 中列出的类型都必须是完全不同的；

case 中的 type 可能为 nil ；最多只能有一个 nil case ；

假定表达式 x 的类型为 interface{} 则有如下 type switch ：

```golang
switch i := x.(type) {
case nil:
	printString("x is nil")                // type of i is type of x (interface{})
case int:
	printInt(i)                            // type of i is int
case float64:
	printFloat64(i)                        // type of i is float64
case func(int) float64:
	printFunction(i)                       // type of i is func(int) float64
case bool, string:
	printString("type is bool or string")  // type of i is type of x (interface{})
default:
	printString("don't know the type")     // type of i is type of x (interface{})
}
```

另外一种写法为：

```golang
v := x  // x is evaluated exactly once
if v == nil {
	i := v                                 // type of i is type of x (interface{})
	printString("x is nil")
} else if i, isInt := v.(int); isInt {
	printInt(i)                            // type of i is int
} else if i, isFloat64 := v.(float64); isFloat64 {
	printFloat64(i)                        // type of i is float64
} else if i, isFunc := v.(func(int) float64); isFunc {
	printFunction(i)                       // type of i is func(int) float64
} else {
	_, isBool := v.(bool)
	_, isString := v.(string)
	if isBool || isString {
		i := v                         // type of i is type of x (interface{})
		printString("type is bool or string")
	} else {
		i := v                         // type of i is type of x (interface{})
		printString("don't know the type")
	}
}
```

在 type switch 中 fallthrough 语句是不允许使用的；

在 [hoisie/web](https://github.com/hoisie/web) 项目的 server.go 中有如下代码：

```golang
...
	switch handler.(type) {
	case http.Handler:
		s.routes = append(s.routes, route{r: r, cr: cr, method: method, httpHandler: handler.(http.Handler)})
	case reflect.Value:
		fv := handler.(reflect.Value)
		s.routes = append(s.routes, route{r: r, cr: cr, method: method, handler: fv})
	default:
		fv := reflect.ValueOf(handler)
		s.routes = append(s.routes, route{r: r, cr: cr, method: method, handler: fv})
	}
...
```







