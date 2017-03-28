# Erlang 调试之 crashdump

------


## 基于 crashdump 的调试

### [erlang的erl_crash.dump产生以及如何解读](http://mryufeng.iteye.com/blog/113854)

正常情况下，当 erlang 进程发生错误并且没有进行 catch 的时候，**emulator** 就会自动产生 `erl_crash.dump` 以提供 crash 时 emulator 的详细情况，类似于 unix 的 core dump ；

如下几个环境变量用于控制 dump 行为：

- **`ERL_CRASH_DUMP`**
如果 emulator 需要生成一个 crash dump 文件，则该变量的值将作为 crash dump 文件的文件名；如果该坏境变量未设置，crash dump 文件将以 `erl_crash.dump` 作为文件名生成在当前目录中；
- **`ERL_CRASH_DUMP_NICE`**
Unix systems: 如果 emulator 需要生成一个 crash dump 文件，则会将该变量的值作为进程的 nice 值，以便降低其优先级；该值的允许范围为 1 至 39（若设置超过 39 的值则会被替换为 39）；最大值 39 会将 process 设置为最低优先级；
- **`ERL_CRASH_DUMP_SECONDS`**
Unix systems: 该变量设置了当 emulator 写 crash dump 文件时允许耗费的时间（以秒为单位）；当指定时间过去后，emulator 将会被 SIGALRM 信号所终止；

除了被动产生 dump 以外, 用户还可以主动产生 dump 文件，方法有两种：

1. erlang 控制台执行 `CTRL+c,A`
2. `kill -s SIGUSR1 <pid>`

产生的 `erl_crash.dump` 是个纯文本，可能非常大，特别是当你有成千上万的 process 和 port 时；该信息对于系统调优有非常大的意义。

查看方式可参见文档 `erl5.5.5/erts-5.5.5/doc/html/crash_dump.html`；请注意：输出信息中的内存单位，有的是 byte，有的是 WORD ；

更方便的是用工具 webtool 来察看 web 界面，比较直观。

> 总结：
>> （本文存在的问题）
>
> - 通过 `CTRL+c,A` 生成 `erl_crash.dump` 的情况，无法在通过 `-remsh` 登录目标节点的场景中使用（因为得到的文件内容并非目标节点 crash 信息）；
> - 通过 `kill -s SIGUSR1 <pid>` 方式得到的 `erl_crash.dump` 文件内容虽正确，但会导致目标进程 crash ，不满足“无破坏性”要求；


----------


### [产生crashdump的三种方法](http://blog.yufeng.info/archives/2737)


crashdump 对于 erlang 系统来讲，如同 core 对于 c/c++ 程序一样宝贵：针对系统问题的修复提供了最详细的资料。当然 erlang 很贴心了提供了网页版的 crashdump_view 帮助用户解读数据，使用方法如下：

```
crashdump_viewer:start().
```

因为 crashdump 文本文件里面记录了大量系统相关的信息，这些信息对于分析系统的**性能**，**状态**，**排除问题**提供了不可替代的功能。所以很需要在系统正常运作的时候，得到 crashdump 文件。

除了坐等系统有问题自动产生 crashdump 以外，还有两种方法来手动产生 crashdump ；

方法如下：

1. `erlang:halt("abort").`
2. 在 erlang shell 下输入 `CTRL+c,A`


> 总结：
>> （本文存在的问题）
>
> - 通过 `erlang:halt("abort").` 方式得到的 `erl_crash.dump` 文件内容虽正确，但会导致目标进程 crash ，不满足“无破坏性”要求；


----------


### [erlang 在线生成crashdump](http://www.cnblogs.com/lulu/p/4149217.html)

一般来说，生成 dump 有 4 种方式：

1. `erlang:halt("abort").`
2. 在 erlang shell 下输入 `CTRL+c,A`
3. 等着进程自己崩溃产生 dump 文件
4. `kill -SIGUSR1 <pid>`（shell 无法进入时可以使用）

不过上述方式均会导致 node crash 掉；而通常我们的需求是：**只想得到一份进程状态的 snapshot 方便分析，而不产生任何破坏性**；

在 google groups 中找到一种解决方案，如下：

```
-module(custom_crashdump).
-compile(export_all).

crash_dump() ->
    Date = erlang:list_to_binary(rfc1123_local_date()),
    Header = binary:list_to_bin([<<"=erl_crash_dump:0.2\n">>,Date,<<"\nSystem version: ">>]),
    Ets = ets_info(),
    Report = binary:list_to_bin([Header,erlang:list_to_binary(erlang:system_info(system_version)),
                                 erlang:system_info(info),erlang:system_info(procs),Ets,erlang:system_info(dist),
                                 <<"=loaded_modules\n">>,binary:replace(erlang:system_info(loaded),
                                                                        <<"\n">>,<<"\n=mod:">>,[global])]),
    file:write_file("erl_crash.dump",Report).

ets_info() ->
    binary:list_to_bin([ets_table_info(T)||T<-ets:all()]).

ets_table_info(Table) ->
    Info = ets:info(Table),
    Owner = erlang:list_to_binary(erlang:pid_to_list(proplists:get_value(owner,Info))),
    TableN = erlang:list_to_binary(erlang:atom_to_list(proplists:get_value(name,Info))),
    Name = erlang:list_to_binary(erlang:atom_to_list(proplists:get_value(name,Info))),
    Objects = erlang:list_to_binary(erlang:integer_to_list(proplists:get_value(size,Info))),
    binary:list_to_bin([<<"=ets:">>,Owner,<<"\nTable: ">>,TableN,<<"\nName: ">>,Name,
                        <<"\nObjects: ">>,Objects,<<"\n">>]).

rfc1123_local_date() ->
    rfc1123_local_date(os:timestamp()).
rfc1123_local_date({A,B,C}) ->
    rfc1123_local_date(calendar:now_to_local_time({A,B,C}));
rfc1123_local_date({{YYYY,MM,DD},{Hour,Min,Sec}}) ->
    DayNumber = calendar:day_of_the_week({YYYY,MM,DD}),
    lists:flatten(
        io_lib:format("~s, ~2.2.0w ~3.s ~4.4.0w ~2.2.0w:~2.2.0w:~2.2.0w GMT",
                      [httpd_util:day(DayNumber),DD,httpd_util:month(MM),YYYY,Hour,Min,Sec]));
rfc1123_local_date(Epoch) when erlang:is_integer(Epoch) ->
    rfc1123_local_date(calendar:gregorian_seconds_to_datetime(Epoch+62167219200)).
```

erlang 自己的 crash dump 是一个关于进程详细信息的快照文本文件，而此方式是自己拼接一个类似的文件（内容相对简略了许多）；

> 总结：
>> （本文存在的问题）
>
> - 基于上述脚本确实能够在满足“无破坏性”要求的前提下，成功获取进程相关信息，但由于功能上有简化，因此需要自己按需调整呢；
> - 上述脚本将 `erl_crash.dump` 文件生成在该脚本所在目录，非目标进程可执行文件对应的目录；


----------


### [【原创】Erlang 之 erl_crash.dump 生成](https://my.oschina.net/moooofly/blog/630946)

结论一：

- 上文的一些结论是存在问题的，已经进行了标注；
- 问题在于上述试验中，我是通过 `-remsh` 方式登录到目标节点上的，即（remsh 的行为实现）在本地创建一个 Erlang 节点，同时在远端节点上启动初始 shell ，那么此时无论是使用 `Ctrl+c,Ctrl+c`，或 `Ctrl+c,a` ，还是 `Ctrl+c,A` ，终止的都是在远端启动的那个初始 shell ，因此并不会导致目标进程退出；而此时获取到的 `erl_crash.dump` 文件当然也就不是目标进程对应的崩溃文件；
- 通过 `-remsh` 登录后执行 `erlang:halt("abort").` 命令，会令目标进程（ERTS）暴力退出，并以 "abort" 字符串作为 Slogan 生成 `erl_crash.dump` 文件。因为是基于目标进程信息生成的崩溃文件，因此必然比上面终止初始 shell 进程时生成的崩溃文件内容大；
- 通过 `SIGUSR1` 令目标进程退出，并生成 `erl_crash.dump` 文件的方式也是可以的。

结论二：

- `Ctrl+c,Ctrl+c` 和 `Ctrl+c,a` 什么都不会生成，即使是基于 console 启动程序时；
- `Ctrl+c,A` 可以生成 `erl_crash.dump` 和 `core.xxx` （要放开 `ulimit -c`）；
- `erlang:halt("abort").` 只会生成 `erl_crash.dump` （即使放开 `unlimit -c`）；
- `erlang:halt(abort).` 只会生成 `core.xxx` （要放开 `unlimit -c`）；
- 通过 `SIGUSR1` 终止 erlang 进程，可以生成 `erl_crash.dump` 和 `core.xxx` （要放开 `ulimit -c`）；


----------


### [【原创】Erlang 之 erl_crash.dump 文件分析](https://my.oschina.net/moooofly/blog/632689)

1. 基于 crashdump_viewer 的 web 页面进行 erl_crash.dump 分析；
2. 基于 recon 的 erl_crashdump_analyzer.sh 分析脚本 进行 erl_crash.dump 分析；


----------

## 关于 halt 的说明


### [erlang 手册之 halt/0,1,2](http://erlang.org/doc/man/erlang.html#halt-0)


- **`halt() -> no_return()`**

The same as `halt(0, [])`. Example:

```
> halt().
os_prompt%
```

- **`halt(Status) -> no_return()`**

Types:

> Status = integer() >= 0 | **abort** | **string()**

The same as `halt(Status, [])`. Example:

```
> halt(17).
os_prompt% echo $?
17
os_prompt%
```

- **`halt(Status, Options) -> no_return()`**

Types:

> Status = integer() >= 0 | **abort** | **string()**
> Options = [Option]
> Option = **{flush, boolean()}**

`Status` must be a **non-negative integer**, a **string**, or the atom `abort`. **Halts the Erlang runtime system**. Has no return value. 

Depending on `Status`, the following occurs:

- **integer()**
  The runtime system exits with integer value `Status` as status code to the calling environment (OS).
- **string()**
  An **Erlang _crash dump_ is produced** with `Status` as **slogan**. Then the runtime system exits with `status` code 1. Note that only code points in the range 0-255 may be used and the string will be truncated if longer than 200 characters.
- **abort**
  The runtime system aborts **producing a _core dump_**, if that is enabled in the OS.

> Note
> 
> On many platforms, the OS supports only status codes 0-255. A too large status code is truncated by clearing the high bits.

For **integer** `Status`, the Erlang runtime system **closes all ports and allows async threads to finish their operations before exiting**. To exit without such flushing, use Option as `{flush,false}`.

For statuses **string()** and `abort`, option `flush` is **ignored** and flushing is **not done**.

> 要点：
> 
> - halt/0,1,2 的效果均为终止 Erlang 的运行时系统；
> - 若 status 为非负数字，则 ERTS 以对应数字作为退出码返回；
> - 若 status 为 string ，则 ERTS 以数字 1 作为退出码返回；并以 string 内容做为 slogan 生成 crashdump 文件；
> - 若 status 为 abort 原子，则 ERTS 终止并生成 core dump（在 OS 允许的情况下）；
> - 若 status 为非负数字，则 ERTS 为优雅退出（清理行为参考上面的描述）；
> - 若 status 为 string 或 abort 原子，则 ERTS 以非优雅方式退出（无 flush 动作）；


### [erlang:halt/0,1,2 问题](http://erlang.org/download/otp_src_R15B01.readme)

```
    OTP-9985  == erts stdlib ==

	      When an escript ends now all printout to standard output and
	      standard error gets out on the terminal. This bug has been
	      corrected by changing the behaviour of erlang:halt/0,1, which
	      should fix the same problem for other escript-like
	      applications, i.e that data stored in the output port driver
	      buffers got lost when printing on a TTY and exiting through
	      erlang:halt/0,1.

	      The BIF:s erlang:halt/0,1 has gotten improved semantics and
	      there is a new BIF erlang:halt/2 to accomplish something like
	      the old semantics. See the documentation.

	      Now erlang:halt/0 and erlang:halt/1 with an integer argument
	      will close all ports and allow all pending async threads
	      operations to finish before exiting the emulator. Previously
	      erlang:halt/0 and erlang:halt(0) would just wait for pending
	      async threads operations but not close ports. And
	      erlang:halt/1 with a non-zero integer argument would not even
	      wait for pending async threads operations.

	      To roughly the old behaviour, to not wait for ports and async
	      threads operations when you exit the emulator, you use
	      erlang:halt/2 with an integer first argument and an option
	      list containing {flush,false} as the second argument. Note
	      that now is flushing not dependant of the exit code, and you
	      can not only flush async threads operations which we deemed
	      as a strange behaviour anyway.

	      Also, erlang:halt/1,2 has gotten a new feature: If the first
	      argument is the atom 'abort' the emulator is aborted
	      producing a core dump, if the operating system so allows.
```