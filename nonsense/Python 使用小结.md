# Python 使用小结

- 只使用 python 进行逻辑构建，作为粘合剂使用；
- 安装特定版本时，需要安装配套的文档；在 python 中微小的变更都可能导致函数不可用；使用错误文档，容易抓瞎；
- Jupyter 玩起来；用于私人环境，用户能拿到所有权限，不建议给别人用，可以放在虚拟机中；
- IPython 必备；
-  在 python 交互环境中执行 python 文件中的程序；

```python
➜  Python python
Python 2.7.13 (default, Dec 17 2016, 23:03:43)
[GCC 4.2.1 Compatible Apple LLVM 8.0.0 (clang-800.0.42.1)] on darwin
Type "help", "copyright", "credits" or "license" for more information.
>>>
>>> import os
>>> os.system('python test.py')
1+2
0
>>>
>>>
>>> import test
1+2
>>>
```

- 使用三引号引用双引号和单引号

```python
>>> '''I'm "Dad"'''
'I\'m "Dad"'
>>>
```

- 缩进：python 建议缩进统一使用四个空格实现；
- 只有 if..elif..else 没有 switch ；
- for 的潜台词是“枚举”，本质上是一种语法糖；
- python 中的 lambda 主要用于延迟调用；
- python 中的限制：只允许定义一行的 lambda ；
- python 中不提倡使用 getter 和 setter ；


----------


## Python2 v.s. Python3 

### print 问题

- python2 中 `print` 为关键字；允许使用 `print 'xxx'` 关键字用法和 `print('xxx')` 函数用法；
- python3 中 `print` 为函数；只允许使用 `print('xxx')` 函数用法；

### input 和 raw_input

```python
# input 是 python3 写法，raw_input 是对应的 python2 写法。
try:
    input = raw_input
except NameError:
    pass
```

> python2 中的 input 有使用风险；

### 缩进问题

- P2 允许空格和 Tab 混用；
- P3 不能混合使用 Tab 和空格；

### 字符串问题（Unicode）

- python2 通过 `u'xxx'` 表示 Unicode ；通过 `'xxx'` 表示普通字符串；
- python3 默认认为所有字符串都为 Unicode；

### 比较问题

- python2 下允许数字和字符进行大小比较；
- python3 下禁止；

> 建议：不混合使用；

### 除法问题 /

- python2 中 "/" 是整数除法；
- python3 中 "/" 是浮点除法；







