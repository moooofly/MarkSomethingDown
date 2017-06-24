# [Active vs Passive Polling in Packet Processing](http://www.ntop.org/pf_ring/active-vs-passive-polling-in-packet-processing/)

时常有 PF_RING 用户问到：应该使用 **passive polling** 技术（即调用 `pfring_poll()`）还是使用 `active polling` 技术（即基本等价于实现了一个 **active loop** ，直到下一个待处理 packet 出现；学院派的人们可能会回答说 (passive) polling 就是正确答案；这个回答从很多方面来说都是有道理的，包括从节省 CPUs 耗电的角度；但不幸的是，在实际中故事会有所不同；

如果你打算在无事可做时，避免浪费 CPU 时钟周期（cycles），即没有 packet 待处理时，你应该调用 `pfring_poll()` 或 `poll`/`select` 以告知系统在发现有 packet 处理时唤醒程序；而如果你创建的是 **active polling loop** ，你可能就会实现如下代码逻辑

```c
while(<no packet available>) { usleep(1); }
```

这种实现能够（理论上）减少 CPU loop ，令其稍微休息一下（1 微秒）；这会是一种很好的实践，如果 `usleep()` 或 `nanosleep()` 真正持续的时间确实如你所指定的值时；不幸的是，结果并不是这样；这些函数均通过系统调用实现的 sleep 功能；我们知道，简单的系统调用的代价是非常低的（例如，你可以基于这里的[测试程序](https://github.com/tsuna/contextswitch/blob/master/timesyscall.c)进行相关测试），并且在通常情况下[低于 100 nsec/call](http://blog.tsunanet.net/2010/11/how-long-does-it-take-to-make-context.html)，即和 1 usec sleep 相比可以忽略不计；下面让我们基于下面的简单程序(sleep.c)来具体测量一下 `usleep()` 和 `nanosleep()` 的精确度到底是多少；

```c
#include <string.h>
#include <sys/time.h>
#include <stdio.h>
 
double delta_time_usec(struct timeval *now, struct timeval *before) {
  time_t delta_seconds;
  time_t delta_microseconds;
 
  delta_seconds      = now->tv_sec  - before->tv_sec;
 
  if(now->tv_usec > before->tv_usec)
     delta_microseconds = now->tv_usec - before->tv_usec;
  else
    delta_microseconds = now->tv_usec - before->tv_usec + 1000000;  /* 1e6 */
 
  return((double)(delta_seconds * 1000000) + (double)delta_microseconds);
}
 
int main(int argc, char* argv[]) {
  int i, n = argc > 1 ? atoi(argv[1]) : 100000;
  static struct timeval start, end;
  struct timespec req, rem;
  int how_many = 1;
 
  gettimeofday(&start, NULL);
  for (i = 0; i < n; i++)
    usleep(how_many);
  gettimeofday(&end, NULL);
 
  printf("usleep(1) took %f usecs\n", delta_time_usec(&end, &start) / n);
 
  gettimeofday(&start, NULL);
  for (i = 0; i < n; i++) {
    req.tv_sec = 0, req.tv_nsec = how_many;
    nanosleep(&req, &rem);
  }
  gettimeofday(&end, NULL);
 
  printf("nanosleep(1) took %f usecs\n", delta_time_usec(&end, &start) / n);
}
```

不同机器上的运行结果会些许差异，但基本在 60 usec 左右；

```shell
# ./sleep
usleep(1) took 56.248760 usecs
nanosleep(1) took 65.165280 usecs
```

这意味着当使用 `usleep()` 和 `nanosleep()` 进行 1 微秒 sleep 时，实际上 sleep 了大概 60 微秒；这个并不令人觉得奇怪，已被许多人发现了[[1](http://stackoverflow.com/questions/12823598/effect-of-usleep0-in-c-on-linux)] [[2](https://lists.freebsd.org/pipermail/freebsd-arch/2012-March/012417.html)]；

而这种结果在实际中意味着什么呢？我们知道，**若线速为 10G ，那么将会每隔 67 nsec 接收到一个 packet ，即 60 usec 的 sleep 将会接收大约 895 个 packets** ，因此你必须设计好缓冲（buffers）以便处理这种情况；这个情况还意味着，在要求严苛的场景中，你只能选择使用纯粹的 **active polling** 方式（即不能在你的 active packet poll 中调用任何形式的 usleep/nanosleep 函数），例如当你从两块网络适配器上接收 time-merging packets 时；

结论就是：**Active polling** 并不优雅，但在进行 packets 的高速处理时，这种方式可能是必须的，以便获得好的/准确的结果；PF_RING 应用同时支持 passive 和 active polling ；例如，在 `pfcount` 中你可以使用 `-a` 标识强制要求使用 active packet polling ，在增大 CPU 负载（100%）的前提下，（在一些情况下）增加 packet 捕获性能（即当你使用 active polling 时，请确保你将你的应用保定到一个 CPU 核心上；在 `pfcount` 中你可以通过 `-g` 达成）；