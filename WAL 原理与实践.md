


# Write-ahead logging

https://en.wikipedia.org/wiki/Write-ahead_logging

在计算机科学中，write-ahead logging (WAL) 是一种为数据库系统提供了原子性和持久性 （ACID 中的 A  和 D）的技术族；

在使用了 WAL 的系统中，所有修改操作在被 apply 前都会被写入日志文件；通常情况下，redo 和 undo 信息都会被保存在该日志中；

上述处理方式可以通过一个例子进行解释说明；
假定某程序在执行某操作的过程中，突然发生了断电异常；在重启后，程序最好能够知道之前执行的操作是否已成功完成，还是只成功了一半，亦或失败了；如果使用了 write-ahead log ，程序就能够通过检查对应的日志文件，确定出异常断电时执行操作的预期结果和实际结果的差异；基于比较结果，程序就能够确定需要 undo 掉哪些操作，完成哪些操作，或者只是维持原状；

WAL 能够以 in-place 方式对数据库进行更新；另外一种实现原子更新的方式是基于 shadow paging ，该方式并非 in-place ；以 in-place 方式进行更新的主要优势在于可以有效的减少 index 和 block list 的修改；

文件系统通常会使用一种 WAL 的变体，专用于文件系统元数据的保存，通常称为 journaling ；



# 产品中的使用

## kafka

kafka 使用磁盘文件保存收到的消息。它使用一种类似于 WAL（write ahead log）的机制来实现对磁盘的顺序读写，然后再定时的将消息批量写入磁盘。消息的读取基本也是顺序的。这正符合 MQ 的顺序读取和追加写特性；

## InfluxDB


## etcd

Write Ahead Log（预写式日志）是 etcd 用于持久化存储的日志格式。除了在内存中存有所有数据的状态以及节点的索引以外，etcd 还通过 WAL 进行持久化存储。在 WAL 中，所有的数据提交前都会事先记录日志。而 Snapshot 是 etcd 为了防止 WAL 文件中数据过多而创建的数据状态快照；Entry 表示存储的具体日志内容。

etcd 的存储分为**内存存储**和**持久化（硬盘）存储**两部分；内存存储除了顺序化的记录下所有用户对节点数据变更的记录外，还会对用户数据进行索引、建堆等方便查询的操作。而持久化存储则使用了预写式日志（WAL, Write Ahead Log）进行记录存储。

**在 WAL 的体系中，所有数据在提交之前都会进行日志记录。**在 etcd 的持久化存储目录中，有两个子目录。一个是 WAL ，存储着所有事务的变化记录；另一个则是 snapshot ，用于存储某一个时刻etcd 所有目录的数据。通过 WAL 和 snapshot 相结合的方式，etcd 可以有效的进行数据存储和节点故障恢复等操作。

**既然有了 WAL 实时存储了所有的变更，为什么还需要 snapshot 呢？**随着使用量的增加，WAL 存储的数据会暴增，为了防止磁盘很快就爆满，etcd 默认每 10000 条记录做一次 snapshot ，经过snapshot 以后的 WAL 文件就可以删除。而通过 API 可以查询的 etcd 历史操作默认为 1000 条。

**WAL（Write Ahead Log）最大的作用是记录了整个数据变化的全部历程。**在 etcd 中，所有数据的修改在提交前，都要先写入到 WAL 中。

使用 WAL 进行数据的存储使得 etcd 拥有两个重要功能：

- **故障快速恢复**： 当你的数据遭到破坏时，就可以通过执行所有 WAL 中记录的修改操作，快速从最原始的数据恢复到数据损坏前的状态；
- **数据回滚（undo）/重做（redo）**：因为所有的修改操作都被记录在 WAL 中，需要回滚或重做，只需要反向或正向执行日志中的操作即可；

WAL 与 snapshot 在 etcd 中的命名规则：
在 etcd 的数据目录中，WAL 文件以 `$seq-$index.wal` 的格式存储。最初始的 WAL 文件是0000000000000000-0000000000000000.wal，表示是所有 WAL 文件中的第 0 个，初始的 Raft 状态编号为 0 。运行一段时间后可能需要进行日志切分，把新的条目放到一个新的 WAL 文件中。

snapshot 的存储命名则比较容易理解，以 `$term-$index.wal` 格式进行命名存储。term 和 index就表示存储 snapshot 时数据所在的 raft 节点状态，当前的任期编号以及数据项位置信息。

从代码逻辑中可以看到，WAL 有两种模式：**读（read）模式**和**数据添加（append）模式**；两种模式不能同时成立。一个新创建的 WAL 文件处于 append 模式，并且不会进入到 read 模式。一个本来存在的 WAL 文件被打开的时候必然是 read 模式，并且只有在所有记录都被读完的时候，才能进入 append 模式，进入 append 模式后也不会再进入 read 模式。这样做有助于保证数据的完整与准确。



## mongodb

按照 Mongodb 默认的配置，￼WiredTiger 的写操作会先写入 Cache ，并持久化到 WAL（Write ahead log），每 60s 或 log 文件达到 2GB 时会做一次 Checkpoint ，将当前的数据持久化，产生一个新的快照。Wiredtiger 连接初始化时，首先将数据恢复至最新的快照状态，然后根据 WAL 恢复数据，以保证存储可靠性。


## HBase


## SQLite

SQLite 在 3.7.0 开始引入 WAL 技术，全称叫 Write Ahead Log（预写日志）；

其原理是：修改并不直接写入到数据库文件中，而是写入到另外一个称为 WAL 的文件中；如果事务失败，WAL 中的记录会被忽略，撤销修改；如果事务成功，它将在随后的某个时间被写回到数据库文件中，提交修改；

WAL 使用**检查点**将修改写回数据库，默认情况下，当 WAL 文件发现有 1000 页修改时，将自动调用检查点。这个页数大小可以自行配置。

WAL技术带来以下优点：
- 读写操作不再互相阻塞，一定程度上解决了 SQLite 在处理高并发上的性能瓶颈；
- 大多数场景中，带来很大的性能提升；

以下为[官网信息](http://www.sqlite.org/wal.html)摘录：

SQLite 实现原子提交和回滚的方式是基于 rollback journal ；从 3.7.0 版本开始，一种新的 "Write-Ahead Log" 选项出现了（缩写为 "WAL"）；

使用 WAL 取代 rollback journal 在一些方面有优有劣；

优势为：
- WAL 在大多数场景下都是非常快的；
- WAL 提供了更好的并发性能，因为 readers 不会阻塞 writers ，并且 writer 也不会阻塞 readers ；读和写可以并发进行；
- 基于 WAL ，磁盘 I/O 操作更加倾向于顺序读写；
- WAL 使用了 many fewer fsync() 操作，因此更加不容易受系统中  fsync() 调用问题的影响；

劣势为：
- WAL 通常会要求 VFS 支持共享内存原语；(例外：[WAL without shared memory](http://www.sqlite.org/wal.html#noshm)) unix 和 windows 中内置的 VFS 支持该特性，但针对定制操作系统的、第三方扩展 VFS 可能就不支持；
- 所有使用数据库的进程必须位于相同的主机上，WAL 无法工作于跨网络文件系统中；
- 包含针对多个关联（ATTACHed）数据库的变更的事务，对于每个单独的数据库而言都是原子的，  但将所有数据库作为整来看时（跨数据库时），则不是原子的；
- 在读非常多，写非常少的应用场景中，WAL 方式可能会比传统 rollback-journal 方式稍慢一点点（大概 1% 或 2%）；

传统的 rollback journal 方式会将一份原始未经改变的数据库内容写入写入单独的 rollback journal 文件中，之后再将变更内容直接写到数据库文件中；在发生 crash 或 ROLLBACK 时，保留在 rollback journal 中的原始内容将被重放到数据库文件里，以便将数据库文件重置回原始状态；而 COMMIT 操作会在 rollback journal 文件删除时触发；

WAL 方式有所不同；原始内容被保留在数据库文件中，而变更被 append 到单独的 WAL 文件中；COMMIT 行为可以通过将某个标志 commit 的特殊记录值 append 到 WAL 文件的方式进行记录；因此，COMMIT 行为可能发生于未将任何内容写入原始数据库的情况下；这也就允许 readers 能够继续在最初未改变的数据库上继续操作，与此同时，发生的变更会同时被 commit 到 WAL 中；多个事务可以被 append 到单个 WAL 文件到尾部；

关于 Checkpoint 的说明：

很显然，每个人都会想将所有 append 到 WAL 文件中的事务回写到原始数据库中；将 WAL 文件中的事务回写到数据库的行为称为 "checkpoint" ；

另外一种判别 way to think about the difference between rollback and write-ahead log is that in the rollback-journal approach, there are two primitive operations, reading and writing, whereas with a write-ahead log there are now three primitive operations: reading, writing, and checkpointing.

By default, SQLite does a checkpoint automatically when the WAL file reaches a threshold size of 1000 pages. (The SQLITE_DEFAULT_WAL_AUTOCHECKPOINT compile-time option can be used to specify a different default.) Applications using WAL do not have to do anything in order to for these checkpoints to occur. But if they want to, applications can adjust the automatic checkpoint threshold. Or they can turn off the automatic checkpoints and run checkpoints during idle moments or in a separate thread or process.


## Berkeley DB

Write-ahead-logging is the term that describes the underlying implementation that Berkeley DB uses to ensure recoverability. What it means is that before any change is made to a database, information about the change is written to a database log.









