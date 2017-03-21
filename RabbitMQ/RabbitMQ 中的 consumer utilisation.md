# RabbitMQ 中的 consumer utilisation

标签（空格分隔）： RabbitMQ

---


RabbitMQ 文档中关于 `consumer_utilisation` 的说明如下：

> Fraction of the time (between 0.0 and 1.0) that the queue is able to immediately deliver messages to consumers. This can be less than 1.0 if consumers are limited by network congestion or prefetch count.

在 RabbitMQ 的 web 页面上同样可以看到类似信息：

> Fraction of the time that the queue is able to immediately deliver messages to consumers. If this number is less than 100% you may be able to deliver messages faster if:
>
> - There were more consumers or
> - The consumers were faster or
> - The consumers had a higher prefetch count


有此可知

- consumer_utilisation 是介于 0.0 和 1.0 之间的表示时间占比的一个数值（也可以理解成百分比）；
- consumer_utilisation 用于描述 queue 能否将 messages 立即投递给 consumers 的能力；
- 当出现网络拥塞或者 prefetch count 设置不合理时，该值可能低于 1.0 ；
- 保持该值为 1.0 即 100% 为合理的常态，否则你可能需要增加 consumer 数量，或令当前 consumer 更快，或调整 consumer 的 prefetch count 为更高值；

在 `rabbit_amqqueue_process.erl` 中

```erlang
i(consumer_utilisation, #q{consumers = Consumers}) ->
    %% 获取 consumer 总数
    case rabbit_queue_consumers:count() of
        0 -> '';
        %% 如果 consumer 数量不为 0 则计算利用率
        _ -> rabbit_queue_consumers:utilisation(Consumers)
    end;
```

在 `rabbit_queue_consumers.erl` 中 `count/0` 的实现如下

```erlang
%% These are held in our process dictionary
-record(cr, {ch_pid,    %% 当前 channel 的 pid
             ...
             consumer_count,  %% 当前 channel 上的 consumer 数量
             ...
             unsent_message_count}).
...
%% 获取所有 channel 记录信息
all_ch_record() -> [C || {{ch, _}, C} <- get()].
...
%% 获取每一个 channel 记录中的 consumer 数量并求总和
count() -> lists:sum([Count || #cr{consumer_count = Count} <- all_ch_record()]).
```

`utilisation/1` 的实现如下

```erlang
%% 获取当前的利用率平均值
%%
%% Since => active 的起始时间
%% Avg => 之前的利用率均值
utilisation(#state{use = {active, Since, Avg}}) ->
    use_avg(time_compat:monotonic_time(micro_seconds) - Since, 0, Avg);
%% Since => inactive 的起始时间
%% Active => active 的持续时间
%% Avg => 之前的利用率均值
utilisation(#state{use = {inactive, Since, Active, Avg}}) ->
    use_avg(Active, time_compat:monotonic_time(micro_seconds) - Since, Avg).
```

其中 `time_compat:monotonic_time/1` 用户获取单调递增的系统时间（用于计算时间差）；

而 `use_avg/3` 的实现如下

```erlang
...
%% Utilisation average calculations are all in μs.
-define(USE_AVG_HALF_LIFE, 1000000.0).
...
use_avg(0, 0, Avg) ->
    Avg;
use_avg(Active, Inactive, Avg) ->
    %% 活跃时间长度 ＋ 不活跃时间长度
    Time = Inactive + Active,
    %% 基于权重得到新的（利用率） avg 值
    rabbit_misc:moving_average(Time, ?USE_AVG_HALF_LIFE, Active / Time, Avg).
```

在 `rabbit_misc.erl` 中

```erlang
moving_average(_Time, _HalfLife, Next, undefined) ->
    Next;
%% We want the Weight to decrease as Time goes up (since Weight is the
%% weight for the current sample, not the new one), so that the moving
%% average decays at the same speed regardless of how long the time is
%% between samplings. So we want Weight = math:exp(Something), where
%% Something turns out to be negative.
%%
%% We want to determine Something here in terms of the Time taken
%% since the last measurement, and a HalfLife. So we want Weight =
%% math:exp(Time * Constant / HalfLife). What should Constant be? We
%% want Weight to be 0.5 when Time = HalfLife.
%%
%% Plug those numbers in and you get 0.5 = math:exp(Constant). Take
%% the log of each side and you get math:log(0.5) = Constant.
%%
%% 根据调用关系，可知
%% Time => Inactive + Active
%% HalfLife => ?USE_AVG_HALF_LIFE
%% Next => Active / Time
%% Current => Avg
%% 
moving_average(Time,  HalfLife,  Next, Current) ->
    %% 根据xx算法得出一个加权权重值
    Weight = math:exp(Time * math:log(0.5) / HalfLife),
    %% 基于权重得到新的 avg 值
    Next * (1 - Weight) + Current * Weight.
```


----------

和 inactive 状态相关的代码如下：

- 如果 priority_queue 为空，则认为当前 queue 上的 consumer 处于 inactive 状态

```erlang
inactive(#state{consumers = Consumers}) ->
    priority_queue:is_empty(Consumers).
```

- 当 priority_queue 为空，则更新当前 queue 上的 consumer 状态为 inactive

```erlang
deliver(FetchFun, QName, ConsumersChanged,
        State = #state{consumers = Consumers}) ->
    case priority_queue:out_p(Consumers) of
        {empty, _} -> %% 若 priority_queue 为空，则更新状态为 inactive
            {undelivered, ConsumersChanged,
             State#state{use = update_use(State#state.use, inactive)}};
    ...
```

和 active 状态相关的代码如下：

- 新建时，认为 consumer 处于 active 状态

```erlang
new() -> #state{consumers = priority_queue:new(),
                use       = {active,
                             time_compat:monotonic_time(micro_seconds),
                             1.0}}.
```

- 在从 block 状态恢复成 unblock 状态时，则更新状态为 active

```
unblock(C = #cr{blocked_consumers = BlockedQ, limiter = Limiter},
        State = #state{consumers = Consumers, use = Use}) ->
        ...
        {Blocked, Unblocked} ->
        ...
            {unblocked,
             State#state{consumers = priority_queue:join(Consumers, UnblockedQ),
                         use       = update_use(Use, active)}}
```