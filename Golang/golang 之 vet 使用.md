# golang 之 vet 使用

标签（空格分隔）： golang

---

## [Command vet](https://golang.org/cmd/vet/)

`Vet` 命令负责检测 Go 源码，并报告其中可疑的部分，例如，`Printf` 调用时参数和格式字串不对应；`Vet` 使用启发式算法，因此不能保证报告出的所有内容都对应到真正的问题上，但可以锁定非编译器造成的所有错误；

三种触发方式：

- 针对 package

```
go vet package/path/name
```

vets 指定路径下的 package ；

- 针对 files

```
go tool vet source/directory/*.go
```

vets 指定名字的文件，所有文件必须属于相同的 package ；

- 针对 directory

```
go tool vet source/directory
```

recursively descends the directory, vetting each package it finds.

若 Vet 的 exit code 是 2 ，则表示使用该工具有误；1 表示检测出了需要报告的问题；0 表示没有问题；需要注意的是，由于采用了非可靠的启发式算法，该工具并不能检测出所有可能的问题，因此，仅应将该工具作为辅助排错工具，而非全面、精准的错误排查工具；

默认情况下，`-all` flag 被使用，表示将运行全部 test 检查；如果显式设置任何 flags 为 true ，则只有设置内容对应的 test 才被运行；相反地，如果有 flag 被显式设置为 false ，则也只有对应的测试被取消掉；由此可知，`-printf=true` 会运行 printf 检查，`-printf=false` 会运行所有检查，除了 printf 检查；

可执行的检查内容有：

- Assembly declarations
    Flag: `-asmdecl`
    
    汇编文件和 Go 函数声明之间无法匹配；

- Useless assignments
    Flag: `-assign`
    
    检测无用的赋值语句；

- Atomic mistakes
    Flag: `-atomic`
    
    针对 `sync`/`atomic` package 的常见错误使用；

- Boolean conditions
    Flag: `-bool`
    
    涉及布尔运算符的错误；

- Build tags
    Flag: `-buildtags`
    
    与 `+build` tags 相关的、存在问题的格式或位置；

- Invalid uses of cgo
    Flag: `-cgocall`
    
    针对 cgo 指针传递规则违反情况的检测；

- Unkeyed composite literals
    Flag: `-composites`
    
    未使用 field-keyed 语法的复合 struct 字面量；

- Copying locks
    Flag: `-copylocks`
    
    检测 Locks 被错误的按值传递情况；

- HTTP responses used incorrectly
    Flag: `-httpresponse`
    
    Mistakes deferring a function call on an HTTP response before checking whether the error returned with the response was nil.

- 未成功调用 `WithCancel` 返回的 `cancelation` 函数 
    Flag: `-lostcancel`
    
    由 `context.WithCancel`，`WithTimeout` 和 `WithDeadline` 返回的 `cancelation` 函数必须得到调用，否则新建（派生）的 context 将会一直存在，直到其 parent context 被取消；（background context 永远无法被取消）

- Methods
    Flag: `-methods`
    
   针对 methods 使用了非标准形式的签名（造成签名相似问题），包括：
    
    > Format GobEncode GobDecode MarshalJSON MarshalXML
    Peek ReadByte ReadFrom ReadRune Scan Seek
    UnmarshalJSON UnreadByte UnreadRune WriteByte
    WriteTo

- Nil function comparison
    Flag: `-nilfunc`
    
    将 functions 和 nil 进行了比较；

- Printf family
    Flag: `-printf`
    
    针对 Printf family 中函数的可疑调用，包括任何带有如下名字的函数，不管大小写：
    
    > Print Printf Println
    Fprint Fprintf Fprintln
    Sprint Sprintf Sprintln
    Error Errorf
    Fatal Fatalf
    Log Logf
    Panic Panicf Panicln
    
    使用 `-printfuncs` flag 可以重定义上述列表；若函数命名以 'f' 结尾，则该函数被假定接受一个格式描述字串，形式同 `fmt.Printf` 中的定义；如果不符合要求，那么 vet 将会针对看起来像格式描述字串的参数发出“抱怨”；
    
    该选项还能够检测其他错误，例如使用 Writer 作为 Printf 的第一个参数的情况；

- Range loop variables
    Flag: `-rangeloops`
    
    在 closures 中不正确的使用 range loop 变量；

- Shadowed variables
    Flag: `-shadow=false` (experimental; must be set explicitly)
    
    可能被不小心 shadowed 掉的 Variables ；

- Shifts
    Flag: `-shift`
    
    Shifts equal to or longer than the variable's length.

- Struct tags
    Flag: `-structtags`
    
    Struct tags that do not follow the format understood by reflect.StructTag.Get. Well-known encoding struct tags (json, xml) used with unexported fields.

- Tests and documentation examples
    Flag: `-tests`
    
    Mistakes involving tests including functions with incorrect names or signatures and example tests that document identifiers not in the package.

- Unreachable code
    Flag: `-unreachable`
    
    报告永远不会被执行到的代码；

- Misuse of unsafe Pointers
    Flag: `-unsafeptr`
    
    Likely incorrect uses of `unsafe.Pointer` to convert integers to pointers. A conversion from `uintptr` to `unsafe.Pointer` is invalid if it implies that there is a uintptr-typed word in memory that holds a pointer value, because that word will be invisible to stack copying and to the garbage collector.

- Unused result of certain function calls
    Flag: `-unusedresult`
    
    Calls to well-known functions and methods that return a value that is discarded. By default, this includes functions like fmt.Errorf and fmt.Sprintf and methods like String and Error. The flags -unusedfuncs and -unusedstringmethods control the set.

### Other flags ¶

These flags configure the behavior of vet:

```
-all (default true)
	使能全部 non-experimental 检查；
-v
	Verbose mode
-printfuncs
	指定一份逗号分隔的、print-like 函数名列表，以补充标准列表内容；
	更多信息，详见关于 -printf flag 的讨论；
-shadowstrict
	是否针对 shadowing 情况进行严格检查；可能导致大量信息输出；
```

Vet 是一种简单的、针对 Go 源码中静态错误的检查器；可以在 `doc.go` 中查看更多信息；


----------


帮助文档查看：

```
go doc cmd/vet
go tool vet -h
```