

[toc]


# raft/raft.go

## Config 结构体

```golang
// 启动一个 raft node 所需的配置参数
type Config struct {
	// 本地 raft node 的 ID ；ID 值不允许为 0 ；
	ID uint64

	// peers 中包含当前所在 raft cluster 中所有 node 的 ID 值（包括自身）；
	// 只允许启动一个新 raft cluster 时设置该值；
	// 基于之前的，设置过 peers 值的配置信息重启 raft node ，会导致 panic 发生；
	// 应该认为 peers 值是私有的，当前仅允许用于测试目的；
	peers []uint64
	
	// ElectionTick 是指在两次选举发生之间，允许 Node.Tick 触发的次数；
	// 若在 ElectionTick 次数被达到时，follower 没有从当前 term 中的 leader 处接收到任何消息，则将以 candidate 身份启动新一轮选举；
	// ElectionTick 的值必须大于 HeartbeatTick 的值；
	// 建议 ElectionTick = 10 * HeartbeatTick 以避免不必要的 leader 切换的发生；
	ElectionTick int

    // HeartbeatTick 是指在两次心跳发生之间，允许 Node.Tick 触发的次数；
    // 即 leader 会在每隔 HeartbeatTick 数量的 tick 后，通过发送心跳消息维护其领导地位
	HeartbeatTick int

	// Storage 作为 raft node 的信息存储使用；
	// raft node 会生成需要保存到 Storage 中的 entry 和 state 信息；
	// raft node 会在需要时从 Storage 中读取已持久化的 entry 和 state 信息；
	// raft node 会在重启时从 Storage 中读取之前的 state 和配置信息；
	Storage Storage

	// Applied 是指被 apply 过的最后一个 index 值；
	// 仅允许在重启 raft node 时设置该值；
	// raft node 不会将小于等于该值的 entry 返回给应用；
	// 若该值在重启时被重置，raft node 可能会返回之前已 apply 过的 entry ；
	// 该配置和应用的具体使用情况比较相关；
	Applied uint64

	// MaxSizePerMsg 限制了每条 append 消息的最大大小；
	// 较小的值能够降低 raft node 的 recovery 成本（针对初始 probing 和常规操作中的消息丢失场景）；
	// 另一方面，MaxSizePerMsg 会影响到常规（日志）复制行为的吞吐量；
	// 注意：math.MaxUint64 代表 unlimited ；0 代表每条消息中最多只能有一个 entry ；
	MaxSizePerMsg uint64

	// MaxInflightMsgs 限制了处于乐观复制阶段中，允许存在的，处于 inflight 状态的 append 消息数量；
	// 应用的传输层通常具有其自身 TCP/UDP 层面的发送缓冲区，设置 MaxInflightMsgs 值能够有效避免该发送缓冲区的溢出；
	// TODO (xiangli): feedback to application to limit the proposal rate?
	MaxInflightMsgs int

	// CheckQuorum 标识 leader 是否需要检查 quorum 的活跃情况；
	// 在发生选举超时的情况下，若发现 quorum 处于 unactive 状态，leader 会降级为 follower ；
	CheckQuorum bool

	// Logger 用作 raft log 的日志记录器；
	// For multinode which can host multiple raft group, each raft group can have its own logger
	Logger Logger
}
```

## raft 结构体

```golang

type raft struct {
	id uint64

	Term uint64
	Vote uint64

	readState ReadState

	// the log
	raftLog *raftLog

	maxInflight int
	maxMsgSize  uint64
	prs         map[uint64]*Progress

	state StateType

	votes map[uint64]bool

	msgs []pb.Message

	// the leader id
	lead uint64

	// leadTransferee 表示 leader transfer target 的 ID ，只要该值不为 0
	// Follow the procedure defined in raft thesis 3.10.
	leadTransferee uint64
	// 若存在未被 apply 的配置，则新配置需要被忽略
	pendingConf bool

	// number of ticks since it reached last electionTimeout or received a
	// valid message from current leader when it is a follower.
	// 作为 leader 或 candidate 角色，自上一次发生选举超时后，已经经过的 tick 数目；
	// 作为 follower 角色，从当前 leader 接收到一条有效消息时，已经经过的 tick 数目；
	electionElapsed int

	// 自上一次发生心跳超时后，已经经过的 tick 数目；
	// 只有 leader 需要维护该值；
	heartbeatElapsed int

	checkQuorum bool

	heartbeatTimeout int
	electionTimeout  int
	
	// randomizedElectionTimeout 为一个值范围在 [electiontimeout, 2 * electiontimeout - 1] 之前的随机数；
	// 当 raft node 变更其 state 为 follower 或 candidate 时充值该值；	randomizedElectionTimeout int

	rand *rand.Rand
	tick func()
	step stepFunc

	logger Logger
}

```


----------


# raft/node.go

## SoftState 结构体

```golang

// SoftState 维护了方便 logging 和 debugging 当状态信息；
// SoftState 是可变的，且不需要被持久化到 WAL 中；
type SoftState struct {
	Lead      uint64
	RaftState StateType
}

```

## Ready 结构体

```golang
// Ready 结构封装了针对以下内容的处理：
// 1.读就绪的 entry 和消息;
// 2.待保存到 stable storage 的 entry 和消息；
// 3.待 commit 或 待发送到 peers 的 entry 和消息；
// Ready 中的所有 fileld 都是只读的；
type Ready struct {
	// Node 当前具有的，可变的 state ；
	// 若不存在更新，则 SoftState 为 nil ；
	// 针对 SoftState 的 consume 和 store 行为不作要求；
	*SoftState

	// 在消息被发送前，将要保存到 stable storage 中的 Node 的当前 state ；
	// 若不存在更新，则 HardState 将等于 empty state ；
	// 该文件头部定义了 emptyState = pb.HardState{}
	pb.HardState

	// ReadState 用于 node 在本地提供读请求线性化处理的场景（当其已 apply 的 index 大于 ReadState 中的 index 时）；
	// 当 raft 收到 msgReadIndex 时，将会返回 readState ；
	// 该返回值只对请求读的请求有效；
	ReadState

	// 在消息被发送前，将要保存到 stable storage 中的 entries
	Entries []pb.Entry

	// Snapshot 指将要被保存到 stable storage 中的 snapshot ；
	Snapshot pb.Snapshot

	// CommittedEntries 指待 commit 到 store/状态机 中的 entries ；
	// 这些 entries 之前已经被 commit 到 stable storage 中了；
	CommittedEntries []pb.Entry

	// Messages 指在 Entries 被 commit 到 stable storage 后，需要被发送的 outbound 消息；
	// 若其中包含了 MsgSnap 消息，则应用必须在 snapshot 成功接收后，或调用 ReportSnapshot 失败时，报告回 raft ；
	Messages []pb.Message
}
```


## Node 接口

```golang

// Node 表示 raft cluster 中的一个 node ；
type Node interface {
	// 实现 Node 所需的，对应一个 tick 的，内部逻辑时钟；
	// 选举超时和心跳超时都以 tick 为单位；
	Tick()
	
	// 令当前 Node 转变为 candidate 状态，并发起 leader 竞选；
	Campaign(ctx context.Context) error
	
	// 发起将指定 data 附加到 log 中的议案
	Propose(ctx context.Context, data []byte) error
	
    // 发起配置变更提案；
    // 至多允许一个配置变更提案在特定时间里被 consensus 处理；
    // 应用需要在 apply 过 EntryConfChange 类型的 entry 后，调用 ApplyConfChange ；
	ProposeConfChange(ctx context.Context, cc pb.ConfChange) error
	
	// 基于给定消息向前步进状态机 ctx.Err() will be returned, if any.
	Step(ctx context.Context, msg pb.Message) error

    // 获取一个能够返回当前时间点状态的 channel
    // 当前 Node 的用户必须在通过 Ready 成功获取到状态后，调用 Advance 函数
    // 注意：只有前一次通过 Ready 获取到的所有 committed entries 和 snapshots 都成功 apply 后才允许后一次通过 Ready 获取到的 committed entries 被 apply
	Ready() <-chan Ready

	// 通知 Node 当前应用已经成功将 progress 更新成最新 Ready 内容；
	// 允许当前 node 获取下一次的 Ready 内容；
    // 一般情况下，应用应该在 apply 掉最新 Ready 内容后调用 Advance ；
    // 然而，作为一种优化手段，应用可以在 apply 一些 commands 的同时调用 Advance ；
    // 例如，当最新 Ready 包含了一个 snapshot 数据，则应用可能需要花很长时间 apply 该 snapshot 数据；
    // 为了能够持续接收 Ready 内容而不阻塞 raft 处理过程，允许在 apply 最新 Ready 结束前调用 Advance ；
    // 为了确保该优化手段的安全性，当应用收到了一个 softState.RaftState 等于 Candidate 的 Ready 内容时，
    // 则必须确保（同步机制）先 apply 掉全部未决的配置变更处理（如果存在的话）；
	//
	// 下面是一个简单的，可以等待所有 pending entries 被 apply 的解决方案：
	// ```
	// ...
	// rd := <-n.Ready()
	// go apply(rd.CommittedEntries) // 以 FIFO 顺序执行异步的，最佳化 apply 操作
	// if rd.SoftState.RaftState == StateCandidate {
	//     waitAllApplied()
	// }
	// n.Advance()
	// ...
	//```
	Advance()
	
    // 将配置变更 apply 到本地 node 上；
    // 该函数回返回一个 ConfState protobuf 不透明数据，且必须被记录到 snapshots 中；
    // 该函数永远不回返回 nil ；
    // 该函数只会返回匹配 MemoryStorage.Compact 的指针
	ApplyConfChange(cc pb.ConfChange) *pb.ConfState
	
	// 返回 raft 状态机的当前状态
	Status() Status
	
	// 用于报告在最后一次 send 中指定 node 不可达情况
	ReportUnreachable(id uint64)
	
	// 用于报告已发送 snapshot 的状态
	ReportSnapshot(id uint64, status SnapshotStatus)
	
	// 针对当前 node 执行必要的终止操作
	Stop()
}
```


```golang
// node 完整实现了 Node 接口
type node struct {
	propc      chan pb.Message
	recvc      chan pb.Message
	confc      chan pb.ConfChange
	confstatec chan pb.ConfState
	readyc     chan Ready
	advancec   chan struct{}
	tickc      chan struct{}
	done       chan struct{}
	stop       chan struct{}
	status     chan chan Status

	logger Logger
}
```


----------


# raft/storage.go

## Storage 接口

```golang

// Storage 定义为接口；可被应用程序实现为从 storage 中获取 log entries 的用途；
//
// 只要任意一个 Storage 方法返回了错误，raft 实例就将进入不可操作状态，并拒绝参与竞选；
// 在这种情况下，需要应用程序自己负责状态的清理和恢复；
type Storage interface {
	// 返回已保存的 HardState 和 ConfState 信息；
	InitialState() (pb.HardState, pb.ConfState, error)

	// 返回 [lo,hi) 范围内的 log entries 切片；
	// MaxSize 参数限制了允许返回的 log entries 的总大小，但是要求至少返回一个
	// entry ，如果确实存在的话；
	Entries(lo, hi, maxSize uint64) ([]pb.Entry, error)

	// 返回 entry i 的 term 值，范围在 [FirstIndex()-1, LastIndex()] 之内；
	// 在 FirstIndex 之前的，指定 entry 的 term 会被保留用于匹配目的，
	// 即使 the rest of that entry may not be available.
	Term(i uint64) (uint64, error)

	// 返回 log 中最后一个 entry 的 index 值
	LastIndex() (uint64, error)

	// 返回第一个 log entry 的 index 值，该 entry 可能可以通过 Entries 方法获取到；
	// （更老的 entries 已经被合并到了最新的 snapshot 中；如果 storage 中仅包含了已经
	// 无用的 entry ，那么第一条 log entry 将无法获取到）
	FirstIndex() (uint64, error)

	// 返回最近的 snapshot ；
	// 若 snapshot 暂时尚不存在，则应该返回 ErrSnapshotTemporarilyUnavailable ，
	// 以便 raft 状态机能够得知 Storage 需要一些时间来准备 snapshot ，而之后就可以
	// 调用 Snapshot 方法了；
	Snapshot() (pb.Snapshot, error)
}
```

## MemoryStorage 结构体

```golang
// MemoryStorage 基于 in-memory array 实现了 Storage 接口
type MemoryStorage struct {
	// 保护针对所有 field 的访问；
	// MemoryStorage 的大多数方法都运行在 raft goroutine 中，但是 Append()
	// 方法被运行在应用的 goroutine 中；
	sync.Mutex

	hardState pb.HardState
	snapshot  pb.Snapshot
	// ents[i] 对应了 raft log 中的 i+snapshot.Metadata.Index 位置
	ents []pb.Entry
}
```


# wal/wal.go

## WAL 结构体

```golang
// WAL 是 stable storage 的一种逻辑表达；
// 在特定时刻，WAL 只能处于 read 模式或 append 模式中的一种；
// 新创建的 WAL 处于 append 模式，且已经准备好 append 记录；
// 而刚被打开的 WAL 会处于 read 模式，且已经准备好 read 记录；
// WAL 会在读取了之前的全部记录之后，进入 append 就绪状态；
type WAL struct {
	dir      string           // 底层文件所在目录
	metadata []byte           // metadata 会被记录在每一个 WAL 文件头部
	state    raftpb.HardState // hardstate 会被记录在每一个 WAL 文件头部

	start     walpb.Snapshot // snapshot to start reading
	decoder   *decoder       // 用于解码记录的 decoder
	readClose func() error   // closer for decode reader

	mu      sync.Mutex
	enti    uint64   // 保存到 wal 中的 last entry 对应的 index
	encoder *encoder // 用于编码记录的 encoder

	locks []*fileutil.LockedFile // the locked files the WAL holds (the name is increasing)
	fp    *filePipeline
}
```


# raft/progress.go

## Progress 结构体

[Progress 设计文档](https://github.com/coreos/etcd/blob/master/raft/design.md)

```golang
// Progress 表示从 leader 角度看到的 follower 的 progress ；
// leader 需要维护所有 follower 的 progress ，并根据 follower 所处的 progress 发送相应的 entry（即复制消息）；
// 复制消息就是带有 log entries 的 msgApp
type Progress struct {
    // Match 是 highest known matched entry 的 index ；
    // 如果 leader 对 follower 的复制状态一无所知，则 Match 被设置为 0 ；
    // Next 是将要复制给 follower 的 first entry 的 index ；
    // Leader 会将 Next 中的 entries 放入下一条复制消息中的最新 entry 里（翻译是否有错误） 
	Match, Next uint64
	
	// State 值决定了 leader 应该如何与 follower 进行交互；
	// 当 State 的值为 ProgressStateProbe 时，leader 在每个心跳周期内至多发送一条复制消息；另外，还会对 follower 的实际 progress 情况进行 probe 探测；
	// 
	// 当 State 的值为 ProgressStateReplicate 时，leader 会在发送了复制消息后，乐观
	// 的将被发送的，最新 entry（的 index）进行增长；
	// 这是一种针对快速复制 log entry 给 follower 的一种优化过的 state ；
	// 
	// 当 State 的值为 ProgressStateSnapshot 时，leader 应该已经完成了 snapshot 发送，
	// 同时正在停止发送任何复制消息；
	State ProgressStateType

	// Paused 被用在 ProgressStateProbe 状态中；
	// 当 Paused 为 true 时，raft 应该停止发送复制消息给对应的 peer ；
	Paused bool

	// PendingSnapshot 用在 ProgressStateSnapshot 状态中；
	// 如果存在 pending snapshot ，pendingSnapshot 会被设置成该 snapshot 的 index ；
	// 如果设置了 pendingSnapshot ，则当前 Progress 的复制过程将被停止；
	// raft 将不回重新发送 snapshot 直到收到当前这个 pending snapshot 的失败报告；
	PendingSnapshot uint64

	// 若 progress 当前处于活跃状态，则 RecentActive 为 true ；
	// 从相关的 follower 处收到任何消息都表明当前 progress 是活跃的；
	// RecentActive 可能在选举超时后被重置为 false ；
	RecentActive bool

	// inflights is a sliding window for the inflight messages.
	// When inflights is full, no more message should be sent.
	// When a leader sends out a message, the index of the last
	// entry should be added to inflights. The index MUST be added
	// into inflights in order.
	// When a leader receives a reply, the previous inflights should
	// be freed by calling inflights.freeTo.
	// inflights 定义了针对 inflight 消息的滑动窗口；
	// 当 inflights 满了时，将不会再有消息被发送；
	// 当 leader 发送了一条消息，最后一条 entry 对应的 index 应该被添加到
	// inflights 中；
	// 添加到 inflights 中的 index 必须遵照一定的顺序；
	// 当 leader 收到了一条应答，则之前的 inflights 应该通过调用 inflights.freeTo 
	// 进行释放；
	ins *inflights
}
```


# raft/status.go

## Status 结构体

```golang
type Status struct {
	ID uint64

	pb.HardState
	SoftState

	Applied  uint64
	Progress map[uint64]Progress
}
```


# raft/raftpb/raft.pb.go

## EntryType 自定义类型

```golang
type EntryType int32

const (
	EntryNormal     EntryType = 0
	EntryConfChange EntryType = 1
)
```

## MessageType 自定义类型

```golang
type MessageType int32

const (
	MsgHup            MessageType = 0
	MsgBeat           MessageType = 1
	MsgProp           MessageType = 2
	MsgApp            MessageType = 3
	MsgAppResp        MessageType = 4
	MsgVote           MessageType = 5
	MsgVoteResp       MessageType = 6
	MsgSnap           MessageType = 7
	MsgHeartbeat      MessageType = 8
	MsgHeartbeatResp  MessageType = 9
	MsgUnreachable    MessageType = 10
	MsgSnapStatus     MessageType = 11
	MsgCheckQuorum    MessageType = 12
	MsgTransferLeader MessageType = 13
	MsgTimeoutNow     MessageType = 14
	MsgReadIndex      MessageType = 15
	MsgReadIndexResp  MessageType = 16
)
```

## ConfChangeType 自定义类型

```golang
type ConfChangeType int32

const (
	ConfChangeAddNode    ConfChangeType = 0
	ConfChangeRemoveNode ConfChangeType = 1
	ConfChangeUpdateNode ConfChangeType = 2
)
```


----------

## Entry 结构体
```golang
type Entry struct {
	Type             EntryType `protobuf:"varint,1,opt,name=Type,json=type,enum=raftpb.EntryType" json:"Type"`
	Term             uint64    `protobuf:"varint,2,opt,name=Term,json=term" json:"Term"`
	Index            uint64    `protobuf:"varint,3,opt,name=Index,json=index" json:"Index"`
	Data             []byte    `protobuf:"bytes,4,opt,name=Data,json=data" json:"Data,omitempty"`
	XXX_unrecognized []byte    `json:"-"`
}
```

## SnapshotMetadata 结构体
```golang
type SnapshotMetadata struct {
	ConfState        ConfState `protobuf:"bytes,1,opt,name=conf_state,json=confState" json:"conf_state"`
	Index            uint64    `protobuf:"varint,2,opt,name=index" json:"index"`
	Term             uint64    `protobuf:"varint,3,opt,name=term" json:"term"`
	XXX_unrecognized []byte    `json:"-"`
}
```

## Snapshot 结构体
```golang
type Snapshot struct {
	Data             []byte           `protobuf:"bytes,1,opt,name=data" json:"data,omitempty"`
	Metadata         SnapshotMetadata `protobuf:"bytes,2,opt,name=metadata" json:"metadata"`
	XXX_unrecognized []byte           `json:"-"`
}
```

## Message 结构体
```golang
type Message struct {
	Type             MessageType `protobuf:"varint,1,opt,name=type,enum=raftpb.MessageType" json:"type"`
	To               uint64      `protobuf:"varint,2,opt,name=to" json:"to"`
	From             uint64      `protobuf:"varint,3,opt,name=from" json:"from"`
	Term             uint64      `protobuf:"varint,4,opt,name=term" json:"term"`
	LogTerm          uint64      `protobuf:"varint,5,opt,name=logTerm" json:"logTerm"`
	Index            uint64      `protobuf:"varint,6,opt,name=index" json:"index"`
	Entries          []Entry     `protobuf:"bytes,7,rep,name=entries" json:"entries"`
	Commit           uint64      `protobuf:"varint,8,opt,name=commit" json:"commit"`
	Snapshot         Snapshot    `protobuf:"bytes,9,opt,name=snapshot" json:"snapshot"`
	Reject           bool        `protobuf:"varint,10,opt,name=reject" json:"reject"`
	RejectHint       uint64      `protobuf:"varint,11,opt,name=rejectHint" json:"rejectHint"`
	XXX_unrecognized []byte      `json:"-"`
}
```

## HardState 结构体
```golang
type HardState struct {
	Term             uint64 `protobuf:"varint,1,opt,name=term" json:"term"`
	Vote             uint64 `protobuf:"varint,2,opt,name=vote" json:"vote"`
	Commit           uint64 `protobuf:"varint,3,opt,name=commit" json:"commit"`
	XXX_unrecognized []byte `json:"-"`
}
```

## ConfState 结构体
```golang
type ConfState struct {
	Nodes            []uint64 `protobuf:"varint,1,rep,name=nodes" json:"nodes,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}
```

## ConfChange 结构体
```golang
type ConfChange struct {
	ID               uint64         `protobuf:"varint,1,opt,name=ID,json=iD" json:"ID"`
	Type             ConfChangeType `protobuf:"varint,2,opt,name=Type,json=type,enum=raftpb.ConfChangeType" json:"Type"`
	NodeID           uint64         `protobuf:"varint,3,opt,name=NodeID,json=nodeID" json:"NodeID"`
	Context          []byte         `protobuf:"bytes,4,opt,name=Context,json=context" json:"Context,omitempty"`
	XXX_unrecognized []byte         `json:"-"`
}
```


# rafthttp/transport.go

## Raft 接口

```golang
type Raft interface {
	Process(ctx context.Context, m raftpb.Message) error
	IsIDRemoved(id uint64) bool
	ReportUnreachable(id uint64)
	ReportSnapshot(id uint64, status raft.SnapshotStatus)
}
```

## Transporter 接口

```golang
type Transporter interface {
	// Start starts the given Transporter.
	// Start MUST be called before calling other functions in the interface.
	Start() error
	// Handler returns the HTTP handler of the transporter.
	// A transporter HTTP handler handles the HTTP requests
	// from remote peers.
	// The handler MUST be used to handle RaftPrefix(/raft)
	// endpoint.
	Handler() http.Handler
	// Send sends out the given messages to the remote peers.
	// Each message has a To field, which is an id that maps
	// to an existing peer in the transport.
	// If the id cannot be found in the transport, the message
	// will be ignored.
	Send(m []raftpb.Message)
	// SendSnapshot sends out the given snapshot message to a remote peer.
	// The behavior of SendSnapshot is similar to Send.
	SendSnapshot(m snap.Message)
	// AddRemote adds a remote with given peer urls into the transport.
	// A remote helps newly joined member to catch up the progress of cluster,
	// and will not be used after that.
	// It is the caller's responsibility to ensure the urls are all valid,
	// or it panics.
	AddRemote(id types.ID, urls []string)
	// AddPeer adds a peer with given peer urls into the transport.
	// It is the caller's responsibility to ensure the urls are all valid,
	// or it panics.
	// Peer urls are used to connect to the remote peer.
	AddPeer(id types.ID, urls []string)
	// RemovePeer removes the peer with given id.
	RemovePeer(id types.ID)
	// RemoveAllPeers removes all the existing peers in the transport.
	RemoveAllPeers()
	// UpdatePeer updates the peer urls of the peer with the given id.
	// It is the caller's responsibility to ensure the urls are all valid,
	// or it panics.
	UpdatePeer(id types.ID, urls []string)
	// ActiveSince returns the time that the connection with the peer
	// of the given id becomes active.
	// If the connection is active since peer was added, it returns the adding time.
	// If the connection is currently inactive, it returns zero time.
	ActiveSince(id types.ID) time.Time
	// Stop closes the connections and stops the transporter.
	Stop()
}
```

## Transport 结构体

```golang
// Transport implements Transporter interface. It provides the functionality
// to send raft messages to peers, and receive raft messages from peers.
// User should call Handler method to get a handler to serve requests
// received from peerURLs.
// User needs to call Start before calling other functions, and call
// Stop when the Transport is no longer used.
type Transport struct {
	DialTimeout time.Duration     // maximum duration before timing out dial of the request
	TLSInfo     transport.TLSInfo // TLS information used when creating connection

	ID          types.ID   // local member ID
	URLs        types.URLs // local peer URLs
	ClusterID   types.ID   // raft cluster ID for request validation
	Raft        Raft       // raft state machine, to which the Transport forwards received messages and reports status
	Snapshotter *snap.Snapshotter
	ServerStats *stats.ServerStats // used to record general transportation statistics
	// used to record transportation statistics with followers when
	// performing as leader in raft protocol
	LeaderStats *stats.LeaderStats
	// ErrorC is used to report detected critical errors, e.g.,
	// the member has been permanently removed from the cluster
	// When an error is received from ErrorC, user should stop raft state
	// machine and thus stop the Transport.
	ErrorC chan error

	streamRt   http.RoundTripper // roundTripper used by streams
	pipelineRt http.RoundTripper // roundTripper used by pipelines

	mu      sync.RWMutex         // protect the remote and peer map
	remotes map[types.ID]*remote // remotes map that helps newly joined member to catch up
	peers   map[types.ID]Peer    // peers map

	prober probing.Prober
}
```

# rafthttp/remote.go

## remote 结构体

```golang
type remote struct {
	id       types.ID
	status   *peerStatus
	pipeline *pipeline
}
```

# rafthttp/peer.go

## Peer 接口

```golang
type Peer interface {
	// send sends the message to the remote peer. The function is non-blocking
	// and has no promise that the message will be received by the remote.
	// When it fails to send message out, it will report the status to underlying
	// raft.
	send(m raftpb.Message)

	// sendSnap sends the merged snapshot message to the remote peer. Its behavior
	// is similar to send.
	sendSnap(m snap.Message)

	// update updates the urls of remote peer.
	update(urls types.URLs)

	// attachOutgoingConn attaches the outgoing connection to the peer for
	// stream usage. After the call, the ownership of the outgoing
	// connection hands over to the peer. The peer will close the connection
	// when it is no longer used.
	attachOutgoingConn(conn *outgoingConn)
	// activeSince returns the time that the connection with the
	// peer becomes active.
	activeSince() time.Time
	// stop performs any necessary finalization and terminates the peer
	// elegantly.
	stop()
}
```

## peer 结构体

```golang
// peer is the representative of a remote raft node. Local raft node sends
// messages to the remote through peer.
// Each peer has two underlying mechanisms to send out a message: stream and
// pipeline.
// A stream is a receiver initialized long-polling connection, which
// is always open to transfer messages. Besides general stream, peer also has
// a optimized stream for sending msgApp since msgApp accounts for large part
// of all messages. Only raft leader uses the optimized stream to send msgApp
// to the remote follower node.
// A pipeline is a series of http clients that send http requests to the remote.
// It is only used when the stream has not been established.
type peer struct {
	// id of the remote raft peer node
	id types.ID
	r  Raft

	status *peerStatus

	picker *urlPicker

	msgAppV2Writer *streamWriter
	writer         *streamWriter
	pipeline       *pipeline
	snapSender     *snapshotSender // snapshot sender to send v3 snapshot messages
	msgAppV2Reader *streamReader
	msgAppReader   *streamReader

	sendc chan raftpb.Message
	recvc chan raftpb.Message
	propc chan raftpb.Message

	mu     sync.Mutex
	paused bool

	cancel context.CancelFunc // cancel pending works in go routine created by peer.
	stopc  chan struct{}
}
```