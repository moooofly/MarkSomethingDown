

场景：在通过 rabbit-stress 进行压力测试的时候，会创建出大量的 queue ，并在某些情况下，创建出的 queue 会残留下来，因此需要一种全部删除的方法；


# 方法一

It is possible to reset the state of the entire broker to factory settings by removing the database directory. That will remove queues, but alot more. 

> 通过 rabbitmqctl 的 reset 命令重置 mnesia 数据库，过于暴力，不建议；

# 方法二

You could download the broker definitions using the management plugin, remove the queues section, reset the broker database and then upload the 
modified broker definition. That should preserve all the broker settings but cause only queues to be removed. 

> 相当于针对 reset 操作进行了保护：导出了 fabic 配置信息备用；该行为有意义，但用于删除 queue 这个目的，有点繁琐；

# 方法三

Another option is to list all queues using "`rabbitmqctl list_queues -q name`" or "`rabbitmqadmin --format=bash list queues name`" and feed that 
information into a program that deletes individual queues, e.g. "`rabbitmqadmin delete queue name=...`". The advantage of this method is that the broker need not be shut down. 

> 适用于构成脚本进行批量处理；

# 方法四

Another option for deleting all queues (since 3.2.0) is to define a policy which gives every queue a very short expiry. 

```shell
$ rabbitmqctl set_policy deleter ".*" '{"expires":1}' --apply-to queues
```

Then - assuming your queues are unused, which I assume is why you're deleting them - this will expire them near-immediately. Remember to 
delete the policy again afterwards, or you'll struggle to declare anything... 

This line creates a policy named **deleter** which applies to all names (.*) and will be applied to queues. 

```shell
$ rabbitmqctl clear_policy deleter 
```

> 最方便的方法，没有之一；

# 方法五

You can also delete the vhost containing the queues - of course this deletes the exchanges in that vhost as well...

> 不推荐的方式；


参考：[这里](http://rabbitmq.1065348.n5.nabble.com/Deleting-all-queues-in-rabbitmq-td30933.html)

----------


# How to delete a single or multiple queues in RabbitMQ

![](https://www.cloudamqp.com/images/blog/header-bunnies-faq-questions.jpg)

Frequently Asked RabbitMQ Question: We created 1000+ queues by accident - how do we delete them? This article explains how to delete a single or multiple queues in RabbitMQ.

There are different options of how to delete queues in RabbitMQ. The web based UI can be used via the RabbitMQ Management Interface, a queue policy can be added or a scripts can be used via rabbitmqadmin or HTTP API curl. A script or a queue policy is recommended if you need to delete multiple queues.

Delete queues via:

- RabbitMQ Management Interface
- rabbitmqadmin
- HTTP API curl
- Queue Policy


## RabbitMQ Management Interface

A queue can be deleted from the RabbitMQ Management Interface. Enter the queue tab and go to the bottom of the page. You will find a dropdown "Delete / Purge". Press Delete to the left to delete the queue.

> 该方式一次只能删除一个 queue

## rabbitmqadmin

The management plugin ships with a command line tool [rabbitmqadmin](http://www.rabbitmq.com/management-cli.html) which can perform the same actions as the web-based UI (the RabbitMQ management interface).

### Delete one queue:

```shell
$ rabbitmqadmin delete queue name=name_of_queue
```

> 需要注意 rabbitmqadmin 使用的端口为 management 插件的监听端口，通过 `--port=PORT, -P PORT` 进行指定；

In CloudAMQP the management plugin is assigned port 443 and the ssl flag has to to be used, as shown below.

```shell
$ rabbitmqadmin --host=HOST --port=443 --ssl --vhost=VHOST --username=USERNAME --password=PASSWORD delete queue name=QUEUE_NAME
```

![Delete RabbitMQ message queue](https://www.cloudamqp.com/images/blog/delete-queue-in-rabbitmq.png)

### Delete multiple queues

The script below will:
- Add all queues into a file called q.txt. You can open the file and remove the queues from the file that you would like to keep.
- Loop the list of queues and for each queue delete it.

```shell
$ rabbitmqadmin -f tsv -q list queues name > q.txt
$ while read -r name; do rabbitmqadmin -q delete queue name="${name}"; done < q.txt
```

In CloudAMQP the management plugin is assigned port 443 and the ssl flag has to to be used, as shown below.

```shell
$ rabbitmqadmin --host=HOST --port=443 --ssl --vhost=VHOST --username=USERNAME --password=PASSWORD -f tsv -q list queues name > q.txt
$ while read -r name; do rabbitmqadmin -q --host=HOST --port=443 --ssl --vhost=VHOST --username=USERNAME --password=PASSWORD delete queue name="${name}"; done < q.txt
```

![delete multiple queue rabbitmq](https://www.cloudamqp.com/images/blog/delete-multiple-queue-rabbitmq.png)

> 删 N 个 queue 的方式和删 1 个 queue 的方法是一样的；
> 每删除一个 queue 都会进行一次如下 HTTP 交互，并且是串行删除
>> -----> DELETE /api/queues/%2F/queueName HTTP/1.1
>> <---- HTTP/1.1 204 No Content



## Policy

Add a policy that matches the queue names with an auto expire rule. A policy can be added by entering the **Management Interface** and then pressing the admin tab.

> 也可以通过 `rabbitmqctl` 命令进行设置；

Note that this will only work for unused queues, and don't forget to delete the policy after it has been applied.


Policy to auto delete queues in RabbitMQ Policy to auto delete queues in RabbitMQ

![](https://www.cloudamqp.com/images/blog/policy-delete-all-queues.png)
![](https://www.cloudamqp.com/images/blog/policy-delete-overview.png)

> 效果非常赞，瞬间将匹配的 queue 全部清除；
> 使用中需要注意的点：
>> - 匹配规则错误将造成误删；
>> - 通过 policy 完成删除行为后，切记将该 policy 移除；


## HTTP API, curl

The RabbitMQ Management plugin provides an HTTP-based API for management and monitoring of your RabbitMQ server. In CloudAMQP the management plugin is assigned port 443 and SSL has to be used.

```shell
curl -i -XDELETE https://USERNAME:PASSWORD@HOST:PORT/api/queues/VHOST/QUEUE_NAME 
```

> 这种方式貌似一次也只能删除一个 queue


参考：[这里](https://www.cloudamqp.com/blog/2016-06-21-how-to-delete-queues-in-rabbitmq.html)






